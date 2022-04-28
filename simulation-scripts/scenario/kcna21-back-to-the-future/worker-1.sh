#!/bin/bash

set -Eeuo pipefail

export KUBECONFIG=/etc/kubernetes/kubelet.conf

mkdir -p /root/.kube
kubectl -n kube-public get cm player-config -ojson | jq '.data.kubeconfig' -r | base64 -d > /root/.kube/config

export KUBECONFIG=/root/.kube/config

timedatectl set-ntp off
date --set="25 OCT 1985 07:53:00"
