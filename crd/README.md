# Custom Resource Definition

We can create our own custom resources in kubernetes, which can be considered as objects that users can then use.

In this guide, there are two approaches. The first one uses plain yaml files and an sh script to simulate a controller. This is recommended to first understand the basics of CRDs. The second way uses `KubeBuilder`, which is a tool that allows to create a more complex Go controller. This is recommended for real applications.

### Plain CRD

This CRD is meant to create `TestObjects` that simulate doing a task. This resource is also necessary to do the Kubebuilder example, so make sure to complete this guide first.

- [Plain CRD - TestObject](plain-crd/README.md)

### Kubebuilder CRD

This CRD will check for changes in a `TestObject` status and write to a changes file. It will also add to logs and events.

- [Kubebuilder CRD - Status Alerts](kubebuilder-crd/README.md)
