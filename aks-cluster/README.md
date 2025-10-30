# Creating a cluster through Azure Kubernetes Service

Short guide to make a cluster on AKS and deploy a tiny page, then set a static public IP using Azure.

## Prerequisites

This guide assumes that the user has installed and configured:

- **Azure CLI**: For managing Azure resources and AKS clusters
- **kubectl**: Kubernetes command-line tool for cluster management
- **Docker**: For building and pushing container images
- **Azure subscription**: With permissions to create AKS clusters and resource groups

### Installation links:
- [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli)
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [Docker](https://docs.docker.com/get-docker/)

### Azure CLI setup:
```sh
# Login to Azure
az login

# Set your subscription (if you have multiple)
az account set --subscription "your-subscription-name"

# Verify your account
az account show
```

## Creating a cluster

### 1. Create an AKS cluster
Through the Azure portal, navigate to Kubernetes Services (AKS) and create a new cluster.
Once it is ready, load it with the following command:

```sh
az aks get-credentials --resource-group <resource-group-name> --name <cluster-name>
kubectl config get-contexts # Verify it was added
kubectl config use-context <cluster-name> # To use it
```

## 2. Add frontend

Add an application or service (in this case check the [frontend/](frontend) files).
Push the docker image to dockerhub.

In `frontend/`
```sh
docker buildx build --platform linux/amd64 -t your-docker-username/small-page:v1 . # Ensure it is build for the same architecture in the aks cluster
docker push your-docker-username/small-page:v1
```

## 3. Add manifests

In this case, we'll add a manifest to create a namespace and a manifest for the frontend deployment and service:

- [Namespace manifest](k8s/namespace.yaml)
- [Frontend manifest](k8s/frontend.yaml) 
  - Make sure to update the image to match your username.
  - In this case, we're using load balancer

Then apply them:

```sh
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/frontend.yaml
```

## 4. View it

To quickly view it you can use the port-forward command:

```sh
kubectl port-forward deployment/frontend 8080:80 -n small-page
```

<br/>

To view it from the cluster's external ip, identify the external service IP:

```sh
kubectl get svc -n small-page
```

Then you can access the page: http://IP:80 or just http://IP

## 5. Optional: use a static IP

Using the default config means that the IP is dynamic and will change if the service changes. Therefore, we could use a static IP to assign a domain name or do other things.

First create a public ip using azure. It is important to set the resource group to the node cluster group.

```sh
# Get the node resource group
az aks show --resource-group <cluster resource group> --name <ClusterName> --query nodeResourceGroup -o tsv
```

With this resource group, create the public IP either using az cli or in the Azure portal (Public IP adresses).

```sh
# Create the IP with that resource group
az network public-ip create \
  --resource-group <node resource-group> \
  --name myStaticIP \
  --sku Standard \
  --allocation-method Static \
  --location eastus \
  --dns-name <your-dns-label>  # Optional: Creates yourapp.eastus.cloudapp.azure.com 
  # Make sure the location matches cluster location

# Show the IP
az network public-ip show \
  --resource-group <node-resource-group> \
  --name myStaticIP \
  --query ipAddress \
  --output tsv
```

Finally, set the ip to the load balancer in the manifest (uncomment and set your ip):

```yaml
type: LoadBalancer 
  loadBalancerIP: 20.118.70.76
```

With this we can now view the page through the actual public static IP or the dns if configured.
