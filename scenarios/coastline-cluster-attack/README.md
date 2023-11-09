# Capture the Flag - Coastline Cluster Attacker

## Overview

The player starts within a jumpbox (pod) with a DMZ namespace. From here the player must:

* obtain access to the cluster
* compromise the worker node
* obtain cluster admin to access the master node

### Learning Objectives

* Learn how ephemeral containers work
* Use common enumeration kubernetes techniques
* Understand different methods for obtaining credentials

## Notes

Deploying Elastic Search via Elastic Cloud on Kubernetes requires customisation for KubeSim. This includes:

* Restricting the Operator on the master-node (as it's highly permissive)
* Creating a Storage Class and PV (local) for the PVC to attach to
* Customising the Elastic Search resource with a volume claim template

Alongside this attaching a Fluentd daemonset requires access to the Elastic Cloud secret. This is located here:

```bash
kubectl get secret coastline-es-elastic-user -o go-template='{{.data.elastic | base64decode}}'
```

This secret is referenced as a environment variable in the daemonset


