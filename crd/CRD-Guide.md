# Custom Resource Definition

We can create our own custom resources in kubernetes, which can be considered as objects that users can then use.

## Creating a CRD

1.  Define the CRD

Create a file like mycrd.yaml and apply the changes.

2.  Create a custom resource
For example, adding the widget.yaml and then applying the changes.

3. Verify it was created

```sh
kubectl get widgets
```

## Adding functionality

1. 

```sh
kubebuilder init --domain example.com --owner Alejandra --repo github.com/Ale-Coeto/widget-operator
kubebuilder create api --group example --version v2 --kind Widget
```

2. Update code
Added size type in `api/v1/widget_types.go`.

```sh
make generate
make manifests
```

Added counting logic in `internal/controller/widget_controller.go`.

3. Deploy locally (build and add to kind)

```sh
docker build -t widget-controller:latest .
kind load docker-image widget-controller:latest --name testcluster

```

4. Apply manifests from widget operator

In widget-operator/
```sh
make deploy IMG=widget-controller:latest
kubectl apply -f config/crd/bases/
```

5. Test it

Open logs

```sh
kubectl get pods -n widget-system
kubectl logs -n widget-system widget-operator-controller-manager-5c7bfb97b-fr77l -f
```

Create widgets

```sh
kubectl apply -f - <<EOF
apiVersion: example.com/v1
kind: Widget
metadata:
  name: widget-two
spec:
  size: "large"
EOF
```

You should see sth like
```sh
2025-10-29T17:20:49Z    INFO    Widget reconciled       {"controller": "widget", "controllerGroup": "example.example.com", "controllerKind": "Widget", "Widget": {"name":"widget-two","namespace":"default"}, "namespace": "default", "name": "widget-two", "reconcileID": "c8ff6889-cd02-4050-a52f-e21cba4ee9bc", "name": "widget-two", "size": "large", "totalWidgets": 2}
```