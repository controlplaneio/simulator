#!/usr/bin/env bash

set -Eeuo pipefail

kubectl delete -n kube-system ds kube-proxy --now --ignore-not-found
kubectl delete -f https://raw.githubusercontent.com/projectcalico/calico/v3.24.1/manifests/calico.yaml --now --ignore-not-found
