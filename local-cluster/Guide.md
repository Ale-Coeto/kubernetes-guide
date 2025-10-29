# Creating a cluster locally

Short guide to make a cluster locally to run frontend and backend services using kind.

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

In `services/backend`

```sh
docker build -t backend .
```

In `services/frontend`

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

In `services/k8s`:

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