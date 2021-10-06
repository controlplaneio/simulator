#!/bin/bash

set -euo pipefail

export KUBECONFIG=/etc/kubernetes/admin.conf

for ns in $(kubectl get ns -o jsonpath='{.items[*].metadata.name}'|tr ' ' '\n'); do
    kubectl label --overwrite ns $ns kubernetes.io/metadata.name=$ns
done
