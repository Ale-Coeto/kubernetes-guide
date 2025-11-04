# Creating a cluster locally

Short guide to make a cluster locally to run frontend and backend services using kind and deploy them.

## Prerequisites

This guide assumes that the user has installed:

- **Docker**: For building and running container images
- **kubectl**: Kubernetes command-line tool for cluster management
- **kind**: Kubernetes in Docker - for creating local Kubernetes clusters
- **Go**: Required for building the backend service (optional if using pre-built images)

### Installation links:
- [Docker](https://docs.docker.com/get-docker/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)
- [Go](https://golang.org/doc/install) (optional) 


## Deploying services in a local cluster

### 1. Create a cluster

```sh 
kind create cluster --name testcluster
```

To use the cluster, make sure kubectl is on the correct context

```sh
kubectl config get-contexts                 # Shows clusters
kubectl config use-context kind-testcluster # Use cluster
```

### 2. Write services

Add backend or frontend services with their dockerfiles.
In this case, the frontend has a simple html page and a dockerfile using nginx. The backend is a basic go server.

- [Frontend](frontend/index.html)
- [Backend](backend/main.go)

### 3. Build docker images

With the services ready, it is necessary to pre-build the docker images to later load them in kind.

In `backend/`

```sh
docker build -t backend .
```

In `frontend/`

```sh
docker build -t frontend .
```

### 4. Load images to kind

Since kind runs clusters through docker containers, it will not directly have access to the images we just built, so it is necessary to load them using the following commands.

```sh
kind load docker-image backend:latest --name testcluster
kind load docker-image frontend:latest --name testcluster
```

This images are copied, so if the original images change, it is necessary to re-run the command.

### 5. Add manifests
The manifests will define the deployments and services needed for the application.

The deployments will then run pods that internally run the containers. For example, if we create a deployment for the backend that requires 3 replicas, then there will be 3 nodes, each one running a backend container.

On the other hand, the services will act as load balancer to call healthy pods. This ensures that for example our request reaches one of the nodes.

- [Frontend manifest](k8s/frontend.yaml)
- [Backend manifest](k8s/backend.yaml)

After adding the manifests, it is necessary to apply them:

In `k8s/`:

```sh
kubectl apply -f . # This applies all the manifests in the current path
kubectl apply -f frontend.yaml # This would only apply the frontend manifest
```

### 6. Verify 

Finally, verify that the pods were created.

```sh
kubectl get pods
kubectl get svc
kubectl get deployments
```

If running in wsl, to access port, you can forward with:

```sh
kubectl port-forward svc/frontend-service 8080:80
```

### Optional - Namespace

Optionally, we can also create resources under a namespace. There is already a manifest that creates the namespace called `local-cluster`: [namespace.yaml](k8s/namespace.yaml), but to create the deployments and services under the namespace, you'll need to edit the frontend and backend manifests (uncomment those sections).

Alternatively, you can run the apply with a ns:

In `k8s/`
```sh
kubectl apply -f frontend.yaml -n local-cluster
kubectl apply -f backend.yaml -n local-cluster
```

Or apply all manifests to a specific namespace:
```sh
kubectl apply -f . -n my-namespace
```

### 7. Clean up

To delete the cluster when done:

```sh
kind delete cluster --name testcluster
```

Or delete the ns:

```sh
kubectl delete namespace local-cluster
```

Or delete the nodes, deployments and services:

```sh
kubectl delete -f .
kubectl delete deployment backend frontend
kubectl delete service backend-service frontend-service
```