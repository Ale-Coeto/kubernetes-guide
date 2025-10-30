# Basic guide for Kubernetes

This repo contains some guides to do different things with kubernetes and contains some useful commands as well.

## Short guides

- **[Local cluster](local-cluster/)**: Making a local cluster with kind and deploying simple backend and frontend services.

- **[Custom Resource Definition](crd/)**: Creating a basic CRD using kubebuilder.

## Useful commands

### Cluster Management
```sh
# View cluster info
kubectl cluster-info
kubectl config get-contexts # Shows clusters
kubectl config use-context <context-name> # Use cluster

# Node information
kubectl get nodes
kubectl describe node <node-name>
```

### Namespaces
```sh
# List namespaces
kubectl get namespaces
kubectl get ns

# Create namespace
kubectl create namespace <namespace-name>

# Delete namespace
kubectl delete namespace <namespace-name>

# Set default namespace
kubectl config set-context --current --namespace=<namespace-name>
```

### Pods
```sh
# List pods
kubectl get pods # gets pods in default ns
kubectl get pods -o wide # gets more details
kubectl get pods -n <namespace> # gets pods in a specific namespace
kubectl get pods -A # gets all pods

# Pod details
kubectl describe pod <pod-name>
kubectl logs <pod-name>
kubectl logs -f <pod-name>  # Follow logs

# Execute commands in pods
kubectl exec -it <pod-name> -- /bin/bash
kubectl exec <pod-name> -- <command>
```

### Deployments
```sh
# List deployments
kubectl get deployments
kubectl get deploy

# Deployment details
kubectl describe deployment <deployment-name>
kubectl rollout status deployment <deployment-name> # gets rolloyt status

# Scale deployments
kubectl scale deployment <deployment-name> --replicas=3

# Rolling updates
kubectl rollout restart deployment <deployment-name>
kubectl rollout undo deployment <deployment-name> # rolls back
```

### Services
```sh
# List services
kubectl get services
kubectl get svc

# Service details
kubectl describe service <service-name>

# Port forwarding
kubectl port-forward service/<service-name> <external-port>:<internal-port>
kubectl port-forward pod/<pod-name> <local-port>:<pod-port>
```

### ConfigMaps & Secrets
```sh
# ConfigMaps
kubectl get configmaps
kubectl describe configmap <configmap-name>
kubectl create configmap <name> --from-file=<path>

# Secrets
kubectl get secrets
kubectl describe secret <secret-name>
kubectl create secret generic <name> --from-literal=<key>=<value>
```

### Apply & Delete Resources
```sh
# Apply manifests
kubectl apply -f <file.yaml>
kubectl apply -f <directory>/
kubectl apply -f <url>

# Delete resources
kubectl delete -f <file.yaml>
kubectl delete <resource-type> <resource-name>
kubectl delete pod <pod-name>
kubectl delete deployment <deployment-name>

# Force delete
kubectl delete pod <pod-name> --force --grace-period=0
```

### Debugging & Troubleshooting
```sh
# Events
kubectl get events
kubectl get events --sort-by=.metadata.creationTimestamp

# Resource usage
kubectl top nodes
kubectl top pods

# Describe resources
kubectl describe <resource-type> <resource-name>

# Get YAML/JSON output
kubectl get <resource-type> <resource-name> -o yaml
kubectl get <resource-type> <resource-name> -o json
```

### Kind-specific Commands
```sh
# Create cluster
kind create cluster --name <cluster-name>

# Delete cluster
kind delete cluster --name <cluster-name>

# Load Docker images
kind load docker-image <image-name> --name <cluster-name>

# List clusters
kind get clusters
```

