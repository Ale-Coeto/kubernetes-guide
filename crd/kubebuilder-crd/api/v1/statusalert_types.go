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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// StatusAlertSpec defines the desired state of StatusAlert
type StatusAlertSpec struct {
	// Configuration for what to watch
	WatchKind      string `json:"watchKind"`      // e.g., "TestObject"
	WatchNamespace string `json:"watchNamespace"` // e.g., "default" or "" for all namespaces

	// Alert method toggles
	EnableEvents  bool `json:"enableEvents"`  // Kubernetes events
	EnableLogging bool `json:"enableLogging"` // Standard output logs
	EnableFileLog bool `json:"enableFileLog"` // Write to local file

	// Configuration for file logging (only used if EnableFileLog is true)
	LogFilePath string `json:"logFilePath,omitempty"` // e.g., "/var/log/status-changes.log"
}

// StatusAlertStatus defines the observed state of StatusAlert.
type StatusAlertStatus struct {
	// Track what objects we're currently watching
	WatchedObjects int32 `json:"watchedObjects,omitempty"`

	// Count of alerts sent
	EventsGenerated int32 `json:"eventsGenerated,omitempty"`
	LogsGenerated   int32 `json:"logsGenerated,omitempty"`
	FileLogsWritten int32 `json:"fileLogsWritten,omitempty"`

	// Current status
	Status  string `json:"status,omitempty"`  // "Active", "Error", "Stopped"
	Message string `json:"message,omitempty"` // Human readable status
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// StatusAlert is the Schema for the statusalerts API
type StatusAlert struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of StatusAlert
	// +required
	Spec StatusAlertSpec `json:"spec"`

	// status defines the observed state of StatusAlert
	// +optional
	Status StatusAlertStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// StatusAlertList contains a list of StatusAlert
type StatusAlertList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []StatusAlert `json:"items"`
}

func init() {
	SchemeBuilder.Register(&StatusAlert{}, &StatusAlertList{})
}
