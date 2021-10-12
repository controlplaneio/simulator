#!/bin/bash

set -eo pipefail

event=${1:-"$(< /dev/stdin)"}

level=$(echo "$event" | awk '{print $2}')

wall "Bye bye, better luck next time !!!" 2>/dev/null
sleep 1

PIDS=$(ps aux | grep -v awk | awk '/sshd:/ {print $2}' | tr '\n' ' ')
if [[ "$PIDS" != "" ]]; then
    kill -9 "$PIDS"
fi

TOKEN=$(cat /var/lib/kubelet/pods/*/volumes/kubernetes.io~secret/falco-token-*/token)
CA="/etc/kubernetes/pki/ca.crt"
SERVER=$(awk '/server:/ {print $2}' /etc/kubernetes/kubelet.conf)

kubectl --server="$SERVER" --token="$TOKEN" --certificate-authority="$CA" delete pods,ds -n default --field-selector='metadata.name!=start'
kubectl --server="$SERVER" --token="$TOKEN" --certificate-authority="$CA" delete pods,ds -n the-quiet-place --field-selector='metadata.name!=start'

tcpkill port 10250 &>/dev/null & sleep 10; killall tcpkill
