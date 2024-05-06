[![Go Reference](https://pkg.go.dev/badge/github.com/qba73/inspector.svg)](https://pkg.go.dev/github.com/qba73/inspector)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/qba73/inspector)
[![Go Report Card](https://goreportcard.com/badge/github.com/qba73/inspector)](https://goreportcard.com/report/github.com/qba73/inspector)
![GitHub](https://img.shields.io/github/license/qba73/inspector)
[![Tests](https://github.com/qba73/inspector/actions/workflows/test.yml/badge.svg)](https://github.com/qba73/inspector/actions/workflows/test.yml)

# inspector

Before using `inspector` you need to have [kubectl](https://kubernetes.io/docs/tasks/tools/) binary installed and configured (config file `${HOME}/.kube/config`).

`inspector` is a CLI tool and a Kubernetes plugin for running cluster diagnostics, collecting cluster and Ingress Controller logs and diagnostics, and generating reports.

Here's how to install it and use as a CLI tool:

```shell
go install github.com/qba73/inspector/cmd/inspector@latest
```

To run it:

```shell
inspector
```

```shell
Usage:
   inspector [-v] [-h] -n namespace

Collect K8s and NIC diagnostics in the given namespace

In verbose mode (-v), prints out progess, steps and all data points to stdout.
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

Planned:

- Nodes metrics
- Ingress Controller stats, options and configuration

Future releases will add support for collecting [K8s Gateway API](https://kubernetes.io/docs/concepts/services-networking/gateway/) diagnostics.
