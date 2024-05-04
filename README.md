[![Tests](https://github.com/qba73/inspector/actions/workflows/test.yml/badge.svg)](https://github.com/qba73/inspector/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/qba73/inspector)](https://goreportcard.com/report/github.com/qba73/inspector)


# inspector

`inspector` is a CLI tool and a Kubernetes plugin for running cluster diagnostics, collecting cluster and Ingress Controller logs and diagnostics, and generating reports.

Here's how to install it and use as a CLI tool:

```shell
go install github.com/qba73/inspector/cmd/inspector@latest
```

`inspector` requires `kubectl` and the config file `${HOME}/.kube/config`.

To run it:

```shell
inspector
```

```shell
Usage:
   inspector [-v] namespace

Collect K8s and NIC diagnostics in the given namespace

In verbose mode (-v), prints out progess, steps and all data points to stdout.
```

```shell
inspector default
```

```shell
=== Cluster Info ===
Version: v1.29.2
ClusterID: f66852d1-6d39-40a9-b4c7-c05e39d22332
Nodes: 1
Platform: kind
=== Pods ===
&PodList{ListMeta:{ 11582  <nil>},Items:[]Pod{},}
=== Pod logs ===
...
```

## How it works

The program checks and collects K8s cluster and [Ingress Controller](https://kubernetes.io/docs/concepts/services-networking/ingress/) diagnostics data. It prints out data to the stdout. This allows the output to be piped to other tools for further parsing and processing.

## Collected data points

Currently collected data:

- K8s version
- K8s cluster id
- Number of nodes in the cluster
- K8s platform name
- Pods
- Logs from pods

Planned:

- Events
- ConfigMaps
- Services
- Deployments
- StatefulSets
- ReplicaSets
- Leases
- CRDs
- IngressClasses
- Ingresses
- IngressAnnotations
- Nodes metrics
- Ingress Controller stats, options and configuration

Future releases will add support for collecting [K8s Gateway API](https://kubernetes.io/docs/concepts/services-networking/gateway/) diagnostics.
