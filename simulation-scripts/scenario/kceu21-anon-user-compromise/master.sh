#!/bin/bash

# creates cryptominer deployment, and a monitor process that creates a secret when they're all deleted. RBAC in hashjack pod allows deletion of cryptominers

set -Eeuxo pipefail

NS_PRIMARY="avalon"
DEPLOYMENT="avalon"

BASE_IMAGE="debian:latest"
IMAGE="registry.41RG4P.controlplane.io/excalibur:0x5F3759DF"
IMAGE_DECOY_1="registry.41RG4P.controlplane.io/merlin:0x5F3759DF"
IMAGE_DECOY_2="registry.41RG4P.controlplane.io/guinevere:0x5F3759DF"

main() {

  # allow time for works to complete
#  sleep 30
#
#cat <<EOF | kubectl apply -f - || true
#apiVersion: apps/v1
#kind: DaemonSet
#metadata:
#  namespace: "${NS_PRIMARY}"
#  labels:
#    run: avalon
#  name: avalon
#spec:
#  selector:
#    matchLabels:
#      run: avalon
#  template:
#    metadata:
#      labels:
#        run: avalon
#    spec:
#      containers:
#      - image: debian:latest
#        imagePullPolicy: Always
#        name: avalon
#        command:
#        - /bin/bash
#        - -xc
#        - while sleep 10; do :; done
#EOF


#  # run once so it crashes out
#  kubectl run --generator=deployment/apps.v1 \
#    "${DEPLOYMENT}" \
#    --namespace "${NS_PRIMARY}" \
#    --image="${IMAGE}" \
#    --image-pull-policy=Never \
#    --restart=Never

  ### user

  # ===
  # monitor and user SAs
  # ===

  # user
  SA_USER="hashjack"
  USER_ROLE="hashjack-role"
#  kubectl delete role "${USER_ROLE}" -n "${NS_PRIMARY}" || true

#  cat <<EOF | kubectl apply -n "${NS_PRIMARY}" -f -
#apiVersion: rbac.authorization.k8s.io/v1
#kind: Role
#metadata:
#  name: ${USER_ROLE}
#  namespace: ${NS_PRIMARY}
#rules:
#- apiGroups: [""]
#  resources: ["pods", "pods/log"]
#  verbs: ["get", "watch", "list", "create"]
#- apiGroups: ["apps"]
#  resources: ["deployments"]
#  verbs: ["create"]
#EOF

#  kubectl delete serviceaccount "${SA_USER}" -n "${NS_PRIMARY}" || true
#  kubectl create serviceaccount "${SA_USER}" -n "${NS_PRIMARY}"

  # cluster role bindings
#  kubectl delete rolebinding "${SA_USER}" -n "${NS_PRIMARY}" || true
#  kubectl create rolebinding "${SA_USER}" -n "${NS_PRIMARY}" \
#    --role="${USER_ROLE}" \
#    --serviceaccount="${NS_PRIMARY}":"${SA_USER}"

  kubectl create ns "${NS_PRIMARY}" || true
  sleep 5

  if ! kubectl get serviceaccounts default -n "${NS_PRIMARY}"; then
    sleep 20
    if ! kubectl get serviceaccounts default -n "${NS_PRIMARY}"; then
      sleep 20
      kubectl get serviceaccounts default -n "${NS_PRIMARY}"
    fi
  fi

  # deploy mad pod
  kubectl delete pod -n ${NS_PRIMARY} privateer-1 privateer-2 privateer-3 || true

  COUNT=0
  for POD in $IMAGE_DECOY_1 $IMAGE_DECOY_2 $IMAGE; do
    COUNT=$((COUNT + 1))

    cat <<EOF | kubectl create -f -
  apiVersion: v1
  kind: Pod
  metadata:
    name: privateer-${COUNT}
    namespace: ${NS_PRIMARY}
  spec:
    containers:
      - command:
          - /bin/bash
          - -xc
          - while sleep 10; do :; done
        image: ${IMAGE}
        name: testing
        imagePullPolicy: Never
EOF
  done

  # ensure further pods aren't scheduled
  # taint node 0 (e.g. 1 of 2: k8s-node-0, k8s-node-1)
#  kubectl taint node k8s-node-0 dedicated=special-user:NoSchedule

  # === solution

  # kubectl run -n avalon  test-${RANDOM} --image-pull-policy=IfNotPresent    --image=control-plane.io/excelsior:0x5F3759DF     --command=true     --  /bin/bash -xc "id; cat /etc/passwd; sleep 5"
}

#main |& tee /tmp/log-main
main
