#!/bin/bash

set -euo pipefail

export KUBECONFIG=/etc/kubernetes/admin.conf

for ns in $(kubectl get ns -o jsonpath='{.items[*].metadata.name}'|tr ' ' '\n'); do
    kubectl label --overwrite ns "$ns" "kubernetes.io/metadata.name=$ns"
done

echo 'flag_ctf{AreYouWatchingClosely}' > /root/flag.txt
chmod 400 /root/flag.txt
