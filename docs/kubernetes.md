---
layout: page
title: Kubernetes
hide_hero: true
show_sidebar: false
menubar: docs-menu
---

# Kubernetes

The following snippet runs the Fed Server on [Kubernetes](https://kubernetes.io/docs/tutorials/kubernetes-basics/) in the `apps` namespace. You can reach the fed instance at the following URL from inside the cluster.

```
# Needs to be ran from inside the cluster
$ curl http://fed.apps.svc.cluster.local:8086/ping
PONG
```

Kubernetes manifest - save in a file (`fed.yaml`) and apply with `kubectl apply -f fed.yaml`.