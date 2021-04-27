#!/bin/bash

# creates cryptominer deployment, and a monitor process that creates a secret when they're all deleted. RBAC in hashjack pod allows deletion of cryptominers

set -Eeuxo pipefail

source "../../include/util.sh"

NS_PRIMARY="avalon"
DEPLOYMENT="avalon"

refresh_ns

# run once so it crashes out
kubectl run --generator=deployment/apps.v1 \
  "${DEPLOYMENT}" \
  --namespace "${NS_PRIMARY}" \
  --image=control-plane.io/excelsior:0x5F3759DF \
  --image-pull-policy=Never \
  --restart=Never

### user

# ===
# monitor and user SAs
# ===

# user
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
  resources: ["pods", "pods/log"]
  verbs: ["get", "watch", "list", "create"]
- apiGroups: ["apps"]
  resources: ["deployments"]
  verbs: ["create"]
EOF

kubectl delete serviceaccount "${SA_USER}" -n "${NS_PRIMARY}" || true
kubectl create serviceaccount "${SA_USER}" -n "${NS_PRIMARY}"

# cluster role bindings
kubectl delete rolebinding "${SA_USER}" -n "${NS_PRIMARY}" || true
kubectl create rolebinding "${SA_USER}" -n "${NS_PRIMARY}" \
  --role="${USER_ROLE}" \
  --serviceaccount="${NS_PRIMARY}":"${SA_USER}"

# ===

pod_hashjack

# === solution

# kubectl run -n avalon  test-${RANDOM} --image-pull-policy=IfNotPresent    --image=control-plane.io/excelsior:0x5F3759DF     --command=true     --  /bin/bash -xc "id; cat /etc/passwd; sleep 5"
