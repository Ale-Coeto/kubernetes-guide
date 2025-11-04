# Basic CRD

In this example, we are going to create a CRD using only plain YAML and a an sh script as the controller. This will be a resource meant for testing purposes as it will simulate working on a task and keep a status (pending, succeeded or failed).

## Prerequisites

This guide assumes that the user has installed and configured:

- **kubectl**: Kubernetes command-line tool for cluster management and CRD operations
- **A Kubernetes cluster**: Either local (kind, minikube) or remote cluster access
- **Bash shell**: For running the shell script controller

### Installation links:
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation) (for local testing)

### Cluster setup:
If you don't have a cluster running, you can create a local one using kind:
```bash
kind create cluster --name crd-test
kubectl config use-context kind-crd-test
```

## Steps

### 1.  Define the CRD

First, we need to write the definition of our custom resource. For this, a custom resource manifest was added: [Test Object Definition](test-object-definition.yaml). 

In this manifest, we define the TestObject resource with its spec and status.

The `spec` is like the input or information that the user will provide when creating a `TestObject` resource. 

The `status` will be managed by controllers so that it can achieve a desired state or do other actions.

Make sure to apply changes:

```sh
kubectl apply -f test-object-definition.yaml
```

### 2.  Create a custom resource
Using the custom resource object that we created, we can instantiate a new testObject through a manifest like this one: [test-objects](test-objects.yaml). This manifest creates 2 TestObjects that will simulate working on a task. As we can see in the manifest, we define the `task` and `successMessage`, which are part of the spec that we defined in the custom resource definition, but we don't need to set the status.


Apply changes:

```sh
kubectl apply -f test-objects.yaml
```

### 3. Verify they were created

Now you can check that the objects were created and see the spec details:

```sh
kubectl get TestObjects # Shows all testobjects 
kubectl describe testobject email-task # Describes a testobject
```

However, notice that there is no information about the status. This happens because the status information should only be set by a controller, which we currently don't have.

### 4. Adding a controller

Now, we will create a controller using a simple script to show how controllers typically work. In this case, the controller first gets all of the testobjects. Then it checks the status. If there is no status, it sets it to `Pending`, then, if the status is `Pending`, it "processes" or simulates doing a task and then randomly sets the state to `Succeeded` or `Failed`. This script simulates a controller, but it checks for updates every 2 seconds (while true with a sleep). When using a framework like kubebuilder, it uses Kubernetes Watch API that runs only when something changes (Reconciliation loop), which makes it event-driven instead of polling based. Check the controller here: [Simple Controller](simple-controller.sh).

To run the controller, use the following command:

```sh
./simple-controller.sh
```

You can also see how the status changes live with the following command:

```sh
kubectl get testobjects -o custom-columns="NAME:.metadata.name,STATE:.status.state,MESSAGE:.status.message" --watch
```

Or you can view the description after the controller made the updates:

```sh
kubectl describe testobject email-task
kubectl describe testobject backup-task
```

### 5. Cleanup

When you're done testing, clean up the resources if you won't use them anymore. (If you're continuing with the kubebuilder guide, it is recommended not to delete the resources yet):

```sh
# Delete TestObjects
kubectl delete -f test-objects.yaml

# Remove CRD from cluster
kubectl delete -f test-object-definition.yaml

# Delete the entire cluster (if using kind)
kind delete cluster --name crd-test
```

### Summary

In this guide, we created a custom resource definition for a `TestObject` that has a `SuccessMessage` and `Task` as spec (input) and also has an `message` and `state` as status, which is handled by the controller. Finally, we created a very basic controller that constantly checks the status of our `testobjects` and if they are `pending` it simulates doing a task and then sets it to `success` or `fail`.