#!/bin/bash

# creates squawk container that's sharing its SA token on port 80, RBAC is cluster admin. Get the RBAC, start a priv pod to escalate to the host, find the flag on the host

set -Eeuxo pipefail

source "../../include/util.sh"

NS_PRIMARY="squawk"
NS_SECONDARY="none"

# ns

kubectl delete namespace --grace-period=1 "${NS_PRIMARY}" || true
kubectl create namespace "${NS_PRIMARY}"

# ===
# monitor and user SAs and RBAC
# ===

# user sa and rbac
NS_PRIMARY="squawk"

SA_USER="hashjack"
USER_ROLE="hashjack-role"
kubectl delete role "${USER_ROLE}" -n "${NS_PRIMARY}" || true

cat <<EOF | kubectl apply -n "${NS_PRIMARY}" -f -
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: ${USER_ROLE}
rules:
- apiGroups: [""]
  resources: ["services"]
  verbs: ["get", "watch", "list"]
EOF

kubectl delete serviceaccount "${SA_USER}" -n "${NS_PRIMARY}" || true
kubectl create serviceaccount "${SA_USER}" -n "${NS_PRIMARY}"

# cluster role bindings
kubectl delete rolebinding "${SA_USER}" -n "${NS_PRIMARY}" || true
kubectl create rolebinding "${SA_USER}" -n "${NS_PRIMARY}" \
  --role="${USER_ROLE}" \
  --serviceaccount="${NS_PRIMARY}":"${SA_USER}"

# ===

# monitor sa and rbac

SA_SQUAWK="indomitable-sasquawksh"
SQUAWK_CLUSTER_ROLE="dangerous-cluster-role"
kubectl delete clusterrole "${SQUAWK_CLUSTER_ROLE}" || true

cat <<EOF | kubectl create -f -
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: ${SQUAWK_CLUSTER_ROLE}
rules:
  - apiGroups:
      - '*'
    resources:
      - '*'
    verbs:
      - '*'
  - nonResourceURLs:
      - '*'
    verbs:
      - '*'
EOF

kubectl delete serviceaccount "${SA_SQUAWK}" -n "${NS_PRIMARY}" || true
kubectl create serviceaccount "${SA_SQUAWK}" -n "${NS_PRIMARY}"

# cluster role bindings
kubectl delete clusterrolebinding "${SA_SQUAWK}" -n "${NS_PRIMARY}" || true
kubectl create clusterrolebinding "${SA_SQUAWK}" -n "${NS_PRIMARY}" \
  --clusterrole="${SQUAWK_CLUSTER_ROLE}" \
  --serviceaccount="${NS_PRIMARY}":"${SA_SQUAWK}"

# ===


# squawking serviceaccount pod

NS_PRIMARY="squawk"
kubectl run --generator=deployment/apps.v1 \
  squawk \
  --namespace "${NS_PRIMARY}" \
  --image=busybox \
  --serviceaccount="${SA_SQUAWK}" \
  --expose \
  --port 80 \
  -- /bin/sh -ec "TOKEN=\`cat /var/run/secrets/kubernetes.io/serviceaccount/token\`; while true; do { printf 'HTTP/1.1 200 OK\r\n\n%s\n' \$TOKEN; } | nc -l -p 80; done"


#  -- /bin/bash -ec "while true; do { printf 'HTTP/1.1 200 OK\r\n\n# Hλ$ħ𝔍Ⱥ¢k's Indomitable SASquawkch\n\n\$(cat /etc/passwd)\n'; } | nc -l -p 80; done"
#  --serviceaccount="${SA_SQUAWK}" \


# user starting pod

kubectl run --generator=deployment/apps.v1 \
  hashjack \
  --namespace "${NS_PRIMARY}" \
  --image=bitnami/kubectl:latest \
  --command=true \
  --serviceaccount="${SA_USER}" \
  --expose \
  --port 80 \
  -- /bin/sh -xc "while true; do sleep 10; done"

# ===
