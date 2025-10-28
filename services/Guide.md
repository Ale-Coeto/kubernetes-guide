#

### Create a cluster

```sh 
kind create cluster --name testcluster
```

### Write services

Add backend or frontend services with their dockerfiles.

### Build docker images

In services/backend

```sh
docker build -t backend .
```

In services/frontend

```sh
docker build -t frontend .
```

### Load images to kind

This will allow the kind docker containers to access the images.
This images are copied, so if the build changes, it is necessary to re-run the command.

```
kind load docker-image backend:latest --name testcluster
kind load docker-image frontend:latest --name testcluster
```

### Add manifests
The manifests will define the deployments and services.
The deployments will then run pods that for example run the backend container.
The services will act as load balancer to call healthy pods.

After adding the manifests in k8s, it is necessary to apply them:

In services/k8s

```sh
kubectl apply -f .
```

### Verify 

```sh
kubectl get pods
kubectl get svc
kubectl get deployments
```

If running in wsl, to access port, you can forward with:

```sh
kubectl port-forward svc/frontend-service 8080:80
```