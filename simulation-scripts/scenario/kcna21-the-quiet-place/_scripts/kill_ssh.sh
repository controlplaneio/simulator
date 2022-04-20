#!/bin/bash

set -eo pipefail

## remember this has to finish relatively quickly as otherwise it struggles to keep up with events

event=${1:-"$(< /dev/stdin)"}

level=$(echo "$event" | awk '{print $2}')

if ! grep -iE '(emergency|alert|critical|error|warning)' <(echo -n "$level") &>/dev/null; then
        echo "Aborting action, event only $level"
        exit
fi

#wall "Bye bye, better luck next time !!!" 2>/dev/null
#sleep 1

PIDS=$(ps aux | grep -F -v '/usr/sbin/sshd' | awk '/[s]shd:/ {print $2}' | tr '\n' ' ')
if [[ "$PIDS" != "" ]]; then
    kill -9 "$PIDS"
fi

if [[ $(hostname) == "k8s-master-0" ]]; then
  KOPTS="--kubeconfig=/etc/kubernetes/admin.conf"
else
  TOKEN=$(cat /var/lib/kubelet/pods/*/volumes/kubernetes.io*/falco-token-*/token)
  CA="/etc/kubernetes/pki/ca.crt"
  SERVER=$(awk '/server:/ {print $2}' /etc/kubernetes/kubelet.conf)
  KOPTS="--server=$SERVER --token=$TOKEN --certificate-authority=$CA"
fi


# delete everything useful from non-allowed ns
ns=$(kubectl "$KOPTS" get ns -o jsonpath='{.items[*].metadata.name}'|tr ' ' '\n'|grep -v -E '(kube-system|falco|the-quiet-place)')
for ns in $ns; do
    kubectl "$KOPTS" delete pods,ds,deploy --all --force --grace-period=0
done

# removed any 'extras' from trusted ns
kubectl "$KOPTS" delete pods,ds,deploy -n the-quiet-place --field-selector='metadata.name!=start' --force --grace-period=0
kubectl "$KOPTS" delete pods,ds,deploy -n falco --selector='app!=falco' --force --grace-period=0
kubectl "$KOPTS" delete pods,deploy -n kube-system -l 'k8s-app notin (calico-node,kube-dns,kube-proxy,calico-kube-controllers), tier!=control-plane'
