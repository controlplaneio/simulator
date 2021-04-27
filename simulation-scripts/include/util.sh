#!/bin/bash

refresh_ns() {
  kubectl delete namespace --grace-period=1 "${NS_PRIMARY}" || true
  kubectl create namespace "${NS_PRIMARY}"
}

pod_hashjack() {
  # user starting pod

  [[ "${NS_PRIMARY}" != "" ]]

  kubectl run --generator=deployment/apps.v1 \
    hashjack \
    --namespace "${NS_PRIMARY}" \
    --image=bitnami/kubectl:latest \
    --command=true \
    --serviceaccount="${SA_USER}" \
    --expose \
    --port 80 \
    -- /bin/bash -xc "while sleep \$(echo 'H\316\273$\304\247\360\235\224\215\310\272\302\242k // control-plane.io' | base64 | tr -dc 0-9); do :; done"

#     -- /bin/sh -xc "while sleep \$(echo 'Hλ$ħ𝔍Ⱥ¢k // control-plane.io' | base64 | tr -dc 0-9); do :; done"
}

