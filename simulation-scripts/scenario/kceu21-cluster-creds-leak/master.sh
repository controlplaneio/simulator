#!/bin/bash

# creates cryptominer deployment, and a monitor process that creates a secret when they're all deleted. RBAC in hashjack pod allows deletion of cryptominers

set -Eeuxo pipefail

NS_PRIMARY="crypto-roller"
NS_SECONDARY="crypto-monitor"
MINER_NAMES="pseudocash sigzstash doshcoin moolahro"
SA_USER="hashjack"
USER_ROLE="hashjack-role"

# ns

#kubectl delete namespace --grace-period=1 "${NS_SECONDARY}" || true &
#kubectl delete namespace --grace-period=1 "${NS_PRIMARY}" || true
#kubectl create namespace "${NS_PRIMARY}"
#wait
#kubectl create namespace "${NS_SECONDARY}"

# ===
# monitor and user SAs
# ===

# user
#kubectl delete role "${USER_ROLE}" -n "${NS_PRIMARY}" || true

#cat <<EOF | kubectl apply -n "${NS_PRIMARY}" -f -
#apiVersion: rbac.authorization.k8s.io/v1
#kind: Role
#metadata:
#  name: ${USER_ROLE}
#rules:
#- apiGroups: [""]
#  resources: ["pods", "pods/log", "secrets"]
#  verbs: ["get", "watch", "list"]
#- apiGroups: [""]
#  resources: ["pods", "secrets"]
#  verbs: ["delete"]
#- apiGroups: [""]
#  resources: ["pods/exec"]
#  verbs: ["create"]
#- apiGroups: ["apps"]
#  resources: ["deployments"]
#  verbs: ["delete", "deletecollection"]
#EOF

#kubectl delete serviceaccount "${SA_USER}" -n "${NS_PRIMARY}" || true
#kubectl create serviceaccount "${SA_USER}" -n "${NS_PRIMARY}"

# cluster role bindings
#kubectl delete rolebinding "${SA_USER}" -n "${NS_PRIMARY}" || true
#kubectl create rolebinding "${SA_USER}" -n "${NS_PRIMARY}" \
#  --role="${USER_ROLE}" \
#  --serviceaccount="${NS_PRIMARY}":"${SA_USER}"

# ===

# monitor

#SA_MONITOR="monitor"
#MONITOR_CLUSTER_ROLE="dangerous-cluster-role"
#kubectl delete clusterrole "${MONITOR_CLUSTER_ROLE}" || true
#
#cat <<EOF | kubectl create -f -
#apiVersion: rbac.authorization.k8s.io/v1
#kind: ClusterRole
#metadata:
#  name: ${MONITOR_CLUSTER_ROLE}
#rules:
#  - apiGroups:
#      - '*'
#    resources:
#      - '*'
#    verbs:
#      - '*'
#  - nonResourceURLs:
#      - '*'
#    verbs:
#      - '*'
#EOF

#kubectl delete serviceaccount "${SA_MONITOR}" -n "${NS_SECONDARY}" || true
#kubectl create serviceaccount "${SA_MONITOR}" -n "${NS_SECONDARY}"

# cluster role bindings
#kubectl delete clusterrolebinding "${SA_MONITOR}" -n "${NS_SECONDARY}" || true
#kubectl create clusterrolebinding "${SA_MONITOR}" -n "${NS_SECONDARY}" \
#  --clusterrole="${MONITOR_CLUSTER_ROLE}" \
#  --serviceaccount="${NS_SECONDARY}":"${SA_MONITOR}"

# ===


# cryptominers


#for DEPLOYMENT_TYPE in ${MINER_NAMES}; do
#
#  DEPLOYMENT="${DEPLOYMENT_TYPE}"
#
#  # TODO: add "Hλ$ħ𝔍Ⱥ¢k" to network output?
#  kubectl run --generator=deployment/apps.v1 \
#    "${DEPLOYMENT}" \
#    --namespace "${NS_PRIMARY}" \
#    --image=busybox \
#    --labels=app=web,role="${DEPLOYMENT_TYPE}" \
#    --requests='cpu=10m,memory=32Mi' \
#    --restart='Never' \
#    --expose \
#    --port 80 \
#    --dry-run -o yaml \
#    -- sh -ec "while true; do { printf 'HTTP/1.1 200 OK\r\n\nMinerBot: v7.0.7-${DEPLOYMENT_TYPE}\n'; } | nc -l -p 80; done"
#
#  kubectl --namespace "${NS_PRIMARY}" \
#    scale deployment "${DEPLOYMENT}" --replicas=1
#done

# user starting pod

#kubectl run --generator=deployment/apps.v1 \
#  hashjack \
#  --namespace "${NS_PRIMARY}" \
#  --image=bitnami/kubectl:latest \
#  --command=true \
#  --serviceaccount="${SA_USER}" \
#  --expose \
#  --port 80 \
#    --dry-run -o yaml \
#  -- /bin/sh -xc "while true; do sleep 10; done"

# ===

# roller monitor

#DEPLOYMENT="monitor"
#kubectl delete deployment "${DEPLOYMENT}" --namespace "${NS_SECONDARY}" \
#  --grace-period=1 || true

#kubectl run --generator=deployment/apps.v1 \
#  "${DEPLOYMENT}" \
#  --namespace "${NS_SECONDARY}" \
#  --image=bitnami/kubectl:latest \
#  --command=true \
#  --serviceaccount="${SA_MONITOR}" \
#    --dry-run -o yaml \
#  -- /bin/sh -xc "while true; do \
#      kubectl get deployments -n ${NS_PRIMARY}; \
#      if kubectl get deployment -n ${NS_PRIMARY} ${MINER_NAMES}; then \
#        if kubectl get secret flag -n ${NS_PRIMARY}; then kubectl delete secret flag -n ${NS_PRIMARY}; fi; \
#      else \
#        kubectl create secret -n ${NS_PRIMARY} generic flag  "${DEPLOYMENT}"
#--from-literal='flag_ctf=flag_ctf{d43dc977ca951c52}'; \
#      fi; \
#      sleep 10; \
#    done"
#kubectl --namespace "${NS_SECONDARY}" \
#  scale deployment "${DEPLOYMENT}" --replicas=1


#touch /tmp/master-done

#kubectl apply -f resource/net-pol/web-deny-all.yaml -f resource/net-pol/test-services-allow.yaml;

## Start the nginx container using a different command and custom arguments.
#kubectl run nginx --image=nginx --command -- <cmd> <arg1> ... <argN>
#
## Start the perl container to compute π to 2000 places and print it out.
#kubectl run pi --image=perl --restart=OnFailure -- perl -Mbignum=bpi -wle 'print bpi(2000)'
#
## Start the cron job to compute π to 2000 places and print it out every 5 minutes.
#kubectl run pi --schedule="0/5 * * * ?" --image=perl --restart=OnFailure -- perl -Mbignum=bpi -wle 'print bpi(2000)'
