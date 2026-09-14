# Tutorial: multi-region failover with a merged ServiceExport

This tutorial walks through one of the most common Multi-Cluster Services API
use cases: running the same stateless service in clusters in different
regions, and letting clients fail over between regions automatically without
any client-side logic.

This tutorial is vendor agnostic. It describes the Kubernetes objects you
create and what they mean, not how any specific implementation wires up the
underlying networking. Consult your MCS API implementation's documentation
for installation and network connectivity requirements between clusters.

## Scenario

You run a `serve` Deployment and Service in two clusters, `cluster-east` and
`cluster-west`, each in its own region. You want clients in either region to
reach `serve` through a single name, and to keep working automatically if one
region's `serve` instances become unavailable.

## Prerequisites

- Two or more clusters that are part of the same ClusterSet, with an MCS API
  implementation installed and connected between them.
- `kubectl` access to each cluster.

## Step 1: deploy the same Service to every cluster

In each cluster, deploy your workload and a regular Kubernetes Service using
the same name and namespace:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: serve
  namespace: demo
spec:
  ports:
  - port: 80
    targetPort: 8080
  selector:
    app: serve
```

At this point each `serve` Service is only reachable from within its own
cluster, the same as any other Kubernetes Service.

## Step 2: export the Service from every cluster

In each cluster, create a `ServiceExport` with the same name and namespace as
the Service:

```yaml
apiVersion: multicluster.x-k8s.io/v1beta1
kind: ServiceExport
metadata:
  name: serve
  namespace: demo
```

A `ServiceExport` declares that the Service with the same name and namespace
should be consumable from other clusters in the ClusterSet. It has no
required spec fields, creating it is enough to opt the Service in.

Because `cluster-east` and `cluster-west` both export a Service named `serve`
in the `demo` namespace, your MCS API implementation merges them into a
single logical multi-cluster service. This is the "merged ServiceExport"
case: same name, same namespace, multiple clusters, treated as one service.

## Step 3: consume the merged service

Once both exports exist, your implementation creates a `ServiceImport`
named `serve` in the `demo` namespace of every cluster in the ClusterSet.
You do not create this object yourself, check for it to confirm the export
succeeded:

```console
$ kubectl get serviceimport -n demo serve
```

Pods in any cluster in the ClusterSet can now reach the merged service at:

```
serve.demo.svc.clusterset.local
```

Requests to that name are load balanced across the healthy endpoints behind
`serve` in every cluster that currently exports it, including both
`cluster-east` and `cluster-west`.

## Step 4: observe the failover

If every `serve` pod in `cluster-east` becomes unhealthy or that cluster
becomes unreachable, its endpoints drop out of the merged service and
traffic to `serve.demo.svc.clusterset.local` is served entirely from
`cluster-west`, with no client-side retry or region-aware logic required.
When `cluster-east` recovers, its endpoints rejoin the merged service
automatically.

## Cleanup

Delete the `ServiceExport` in a cluster to stop exporting `serve` from it.
The Service itself is unaffected and remains reachable within that cluster.
