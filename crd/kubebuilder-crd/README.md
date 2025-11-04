# Kubebuilder CRD

In this example, we are going to make a resource using the `kubebuilder` tool. This CRD is meant to be a monitoring tool for `TestObjects` which we defined in the [basic crd](../plain-crd/README.md) guide. Summarized, this tool will do 3 main things:

When the status of a `TestObject` changes:
- Add kubernetes events to our new CRD
- Add logging
- Write to an output file

## Prerequisites

This guide assumes that the user has installed and configured:

- **Go**: Programming language (v1.21+) for building the controller
- **kubebuilder**: Framework for building Kubernetes APIs and controllers
- **kubectl**: Kubernetes command-line tool for cluster management
- **A Kubernetes cluster**: Either local (kind, minikube) or remote cluster access
- **make**: Build tool for running Kubebuilder commands
- **TestObject CRD**: Must complete the [basic CRD guide](../plain-crd/README.md) first

### Installation links:
- [Go](https://golang.org/doc/install) (v1.21 or later)
- [kubebuilder](https://book.kubebuilder.io/quick-start.html#installation)
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation) (for local testing)

### Cluster setup:
If you don't have a cluster running, you can create a local one using kind:
```bash
kind create cluster --name kubebuilder-test
kubectl config use-context kind-kubebuilder-test
```

### Required dependencies:
Before starting, make sure you have completed the basic CRD guide and have TestObjects available:
```bash
kubectl apply -f ../plain-crd/test-object-definition.yaml
kubectl get crds | grep testobjects
```

## Steps

### 1. Initiallize Kubebuilder project

First, we'll create a kubebuilder project:

```sh
# You can also create a new crd
mkdir kubebuilder-crd2
cd kubebuilder-crd2

# Initialize project
kubebuilder init --domain example.com --owner Alejandra --repo github.com/Ale-Coeto/status-alerts

# Create an API (yes to create resource and controller)
# In this case the resource kind will be called StatusAlert
kubebuilder create api --group example --version v1 --kind StatusAlert
```

### 2. Update CRD Design

Now, we need to define the `spec` and `status`. This can be modified in `api/v1/statusalert_types.go`. 

In the spec, we added:

- WatchKind: the object to observe
- WatchNamespace: the namespace to observe
- Boolean flags: to enable alert methods
- Log filepath: where our output file will be written

In the status, we added: 

- WatchedObjects: that we're watching
- Counters: for alerts sent
- Status: to see current object status
- Message: readable message

Next, run the following commands to generate the CRD YAML files from the Go structs:

```sh
make generate  # Generates deepcopy methods and other boilerplate
make manifests # Generates CRD YAML files from your Go structs
```

This creates the CRD yaml in `config/crd/bases/`.

### 3. Update controller logic

Now, we need to define the logic for the controller in `internal/controller/statusalert_controller.go`.
Summarized, the logic in the reconciler is the following:

1. Fetch StatusAlert resources
2. Get testObjects that we'll be monitorin
3. Check each testobject for status changes and depending on the spec flags, add events, logging or update the log file.
4. Update StatusAlert counters

All of these changes are set in the `Reconcile` function, but it is also necessary to add rbac annotations:

```go
// +kubebuilder:rbac:groups=example.com,resources=testobjects,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
```

And add the recorder to the reconciler struct:

```go
// StatusAlertReconciler reconciles a StatusAlert object
type StatusAlertReconciler struct {
	client.Client
	Scheme  *runtime.Scheme
	Recorder record.EventRecorder  // Add for Kubernetes events
}
```

Finally, make sure to include all necessary imports.

Additionally, we need to update `cmd/main.go` to include the recorder in the reconciler struct:

```go
if err := (&controller.StatusAlertReconciler{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("statusalert-controller"), // Add this
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "StatusAlert")
		os.Exit(1)
	}
```

With this changes, we need to generate the manifests (in `kubebuilder-crd/`) because we added RBAC annotations:

```sh
make manifests
```

### 4. Apply CRDs

Now, we need to apply the manifests generated. (The following commands are from the path `kubebuilder-crd/`)

```sh
kubectl apply -f config/crd/bases/ # Apply
kubectl get crds # Verify the crd was created
```

Build the go project:

```sh
go build ./...
```

Run the controller locally:

```sh
make run
```

On a different terminal, create a StatusAlert sample:

```sh
kubectl apply -f config/samples/example_v1_statusalert.yaml
```

### 5. Verify functionality

#### Verify events
To verify events, run the following command:

```sh
kubectl describe StatusAlert statusalert-sample
```

You should be able to see events when `TestObjects` change (Feel free to add or update testobjects).

#### Verify logs

On the controller terminal you should be able to see logs such as:

```
INFO    TestObject status changed    statusAlert=statusalert-sample testObject=test-object-sample namespace=default previousState= currentState=Pending message=Task started timestamp=2025-11-04T...
```

#### Verify log file

```sh
cat /tmp/status-alerts.log
```

### 6. Cleanup

When you're done testing, clean up the resources:

```sh
# Delete StatusAlert resources
kubectl delete statusalert statusalert-sample
kubectl delete -f config/samples/

# Delete TestObjects (if no longer needed)
kubectl delete -f ../plain-crd/test-objects.yaml

# Remove CRDs from cluster
kubectl delete -f config/crd/bases/
kubectl delete -f ../plain-crd/test-object-definition.yaml

# Clean up log files
rm -f /tmp/status-alerts.log

# Delete the entire cluster (if using kind)
kind delete cluster --name kubebuilder-test
```

## Summary

In this guide we used `Kubebuilder` to create a custom resource definition that monitors `TestObjects` and does some actions when its status changes. In more real applications, the functionality can of course be more complex, like posting to an http service or storing logs in a db. Either way, the process is relatively simple, define the `spec` and `status`, make the manifests for the CRD, then write logic for the controller.