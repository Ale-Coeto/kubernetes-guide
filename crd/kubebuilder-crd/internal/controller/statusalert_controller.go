/*
Copyright 2025 Alejandra.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	examplev1 "github.com/Ale-Coeto/status-alerts/api/v1"
)

// StatusAlertReconciler reconciles a StatusAlert object
type StatusAlertReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder // Add for Kubernetes events
}

// +kubebuilder:rbac:groups=example.example.com,resources=statusalerts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=example.example.com,resources=statusalerts/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=example.example.com,resources=statusalerts/finalizers,verbs=update
// +kubebuilder:rbac:groups=example.com,resources=testobjects,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the StatusAlert object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.22.1/pkg/reconcile
func (r *StatusAlertReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// 1. Fetch the StatusAlert
	var statusAlert examplev1.StatusAlert
	if err := r.Get(ctx, req.NamespacedName, &statusAlert); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 2. Get all TestObjects in the specified namespace
	testObjects, err := r.getTestObjects(ctx, statusAlert)
	if err != nil {
		log.Error(err, "Failed to get TestObjects")
		return ctrl.Result{RequeueAfter: time.Minute}, nil
	}

	// 3. Check each TestObject for status changes and create events/logs
	eventsCreated := 0
	logsGenerated := 0
	fileLogsWritten := 0

	for _, testObj := range testObjects {
		if r.hasStatusChanged(testObj) {
			// Create events if enabled
			if statusAlert.Spec.EnableEvents {
				r.createStatusEvent(testObj, statusAlert)
				eventsCreated++
			}

			// Log to standard output if enabled
			if statusAlert.Spec.EnableLogging {
				r.logStatusChange(log, testObj, statusAlert)
				logsGenerated++
			}

			// Write to file if enabled
			if statusAlert.Spec.EnableFileLog {
				err := r.writeStatusToFile(testObj, statusAlert)
				if err != nil {
					log.Error(err, "Failed to write to file", "testObject", testObj.GetName())
				} else {
					fileLogsWritten++
				}
			}

			// Update the TestObject to track that we've processed this status change
			r.markStatusProcessed(ctx, testObj)
		}
	}

	// 4. Update StatusAlert status with counts
	statusChanged := false
	if eventsCreated > 0 {
		statusAlert.Status.EventsGenerated += int32(eventsCreated)
		statusChanged = true
	}
	if logsGenerated > 0 {
		statusAlert.Status.LogsGenerated += int32(logsGenerated)
		statusChanged = true
	}
	if fileLogsWritten > 0 {
		statusAlert.Status.FileLogsWritten += int32(fileLogsWritten)
		statusChanged = true
	}

	if statusChanged {
		statusAlert.Status.Status = "Active"
		statusAlert.Status.Message = fmt.Sprintf("Events: %d, Logs: %d, Files: %d", eventsCreated, logsGenerated, fileLogsWritten)

		if err := r.Status().Update(ctx, &statusAlert); err != nil {
			log.Error(err, "Failed to update StatusAlert status")
		}
	}

	return ctrl.Result{RequeueAfter: time.Minute * 5}, nil // Check again in 5 minutes
}

// getTestObjects retrieves TestObjects based on StatusAlert configuration
func (r *StatusAlertReconciler) getTestObjects(ctx context.Context, statusAlert examplev1.StatusAlert) ([]unstructured.Unstructured, error) {
	var testObjects unstructured.UnstructuredList
	testObjects.SetAPIVersion("example.com/v1")
	testObjects.SetKind("TestObjectList")

	// Set up list options
	listOptions := []client.ListOption{}
	if statusAlert.Spec.WatchNamespace != "" {
		listOptions = append(listOptions, client.InNamespace(statusAlert.Spec.WatchNamespace))
	}

	err := r.List(ctx, &testObjects, listOptions...)
	if err != nil {
		return nil, fmt.Errorf("failed to list TestObjects: %w", err)
	}

	return testObjects.Items, nil
}

// logStatusChange logs TestObject status changes to standard output
func (r *StatusAlertReconciler) logStatusChange(log logr.Logger, testObj unstructured.Unstructured, statusAlert examplev1.StatusAlert) {
	// Get current status from TestObject
	currentState, _, _ := unstructured.NestedString(testObj.Object, "status", "state")
	currentMessage, _, _ := unstructured.NestedString(testObj.Object, "status", "message")

	// Get previous state from annotations
	previousState := ""
	if annotations := testObj.GetAnnotations(); annotations != nil {
		previousState = annotations["status-alert/previous-state"]
	}

	log.Info("TestObject status changed",
		"statusAlert", statusAlert.Name,
		"testObject", testObj.GetName(),
		"namespace", testObj.GetNamespace(),
		"previousState", previousState,
		"currentState", currentState,
		"message", currentMessage,
		"timestamp", time.Now().Format(time.RFC3339),
	)
}

// hasStatusChanged checks if TestObject status has changed by comparing with annotations
func (r *StatusAlertReconciler) hasStatusChanged(testObj unstructured.Unstructured) bool {
	// Get current status
	currentState, _, _ := unstructured.NestedString(testObj.Object, "status", "state")

	// Get previous state from annotations
	previousState := ""
	if annotations := testObj.GetAnnotations(); annotations != nil {
		previousState = annotations["status-alert/previous-state"]
	}

	// Status changed if states are different
	return previousState != currentState
}

// writeStatusToFile writes TestObject status changes to a local file
func (r *StatusAlertReconciler) writeStatusToFile(testObj unstructured.Unstructured, statusAlert examplev1.StatusAlert) error {
	// Get current status from TestObject
	currentState, _, _ := unstructured.NestedString(testObj.Object, "status", "state")
	currentMessage, _, _ := unstructured.NestedString(testObj.Object, "status", "message")

	// Get previous state from annotations
	previousState := ""
	if annotations := testObj.GetAnnotations(); annotations != nil {
		previousState = annotations["status-alert/previous-state"]
	}

	// Create log entry
	logEntry := fmt.Sprintf("[%s] StatusAlert: %s | TestObject: %s/%s | %s -> %s | %s\n",
		time.Now().Format("2006-01-02 15:04:05"),
		statusAlert.Name,
		testObj.GetNamespace(), testObj.GetName(),
		previousState, currentState,
		currentMessage,
	)

	// Use default path if not specified
	filePath := statusAlert.Spec.LogFilePath
	if filePath == "" {
		filePath = "/var/log/status-alerts.log"
	}

	// Append to file
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file %s: %w", filePath, err)
	}
	defer file.Close()

	_, err = file.WriteString(logEntry)
	if err != nil {
		return fmt.Errorf("failed to write to log file: %w", err)
	}

	return nil
}

// createStatusEvent creates a Kubernetes event for status changes
func (r *StatusAlertReconciler) createStatusEvent(testObj unstructured.Unstructured, statusAlert examplev1.StatusAlert) {
	// Get current status from TestObject
	currentState, _, _ := unstructured.NestedString(testObj.Object, "status", "state")
	currentMessage, _, _ := unstructured.NestedString(testObj.Object, "status", "message")

	// Get previous state from annotations
	previousState := ""
	if annotations := testObj.GetAnnotations(); annotations != nil {
		previousState = annotations["status-alert/previous-state"]
	}

	// Create event message
	eventMessage := fmt.Sprintf("TestObject %s status changed from '%s' to '%s': %s",
		testObj.GetName(), previousState, currentState, currentMessage)

	// Determine event type based on current state
	eventType := "Normal"
	if currentState == "Failed" {
		eventType = "Warning"
	}

	// Create the event
	r.Recorder.Event(&statusAlert, eventType, "StatusChanged", eventMessage)
}

// markStatusProcessed updates the TestObject's annotations to track processed states
func (r *StatusAlertReconciler) markStatusProcessed(ctx context.Context, testObj unstructured.Unstructured) error {
	// Get current state
	currentState, _, _ := unstructured.NestedString(testObj.Object, "status", "state")

	// Update annotations
	annotations := testObj.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}
	annotations["status-alert/previous-state"] = currentState
	annotations["status-alert/last-processed"] = time.Now().Format(time.RFC3339)
	testObj.SetAnnotations(annotations)

	// Update the object
	return r.Update(ctx, &testObj)
}

// SetupWithManager sets up the controller with the Manager.
func (r *StatusAlertReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&examplev1.StatusAlert{}).
		Named("statusalert").
		Complete(r)
}
