#!/bin/bash

set -Eeuxo pipefail

IMAGE="control-plane.io/valiant:effort"

main() {
  setup
  create_user_and_k8s_secret
  build_image
  docker_run
}

docker_run() {
  docker run -d --restart=always \
    -p 5678:5678 \
    "${IMAGE}" \
    -text='☣☠☣ This is a valiant effort! ☣☠☣ Follow the 22 white rabbits ✰ Malodorous regards, Hλ$ħ𝔍Ⱥ¢k // https://securi.fyi ✰'
}

prep_file() {
  cat <<EOF >"flag_file"
flag_ctf{d241b8126e0adff0}
EOF
}

create_user_and_k8s_secret() {
  adduser hashjack --disabled-password --gecos ""
  adduser hashjack docker
  SSH_KEY_PATH="/home/hashjack/.ssh/id_rsa"
  rm -rf "${SSH_KEY_PATH}"
  su hashjack -c "id; ssh-keygen -f ${SSH_KEY_PATH} -N ''"
  cat /home/hashjack/.ssh/id_rsa.pub >> /home/hashjack/.ssh/authorized_keys
  chown -R hashjack:hashjack /home/hashjack/

  SECRET="bootstrap-token-9jn6gr"
  if kubectl get secret -n kube-system "${SECRET}"; then
    kubectl delete secret -n kube-system "${SECRET}"
  fi
  kubectl create secret -n kube-system  generic "${SECRET}" \
    --from-file=ssh-privatekey="${SSH_KEY_PATH}"

  kubectl create serviceaccount -n default irredeemable-villainy
  kubectl create serviceaccount -n default tritically-pseudoreminiscent
  kubectl create serviceaccount -n kube-system sanctimonious-ismaticalnesses
  kubectl create serviceaccount -n kube-system plenitudinous-quinqueliteralism
  kubectl create serviceaccount -n kube-public misinterpretable-technicality
  kubectl create serviceaccount -n kube-public subcompensatory-superaverageness

  kubectl create secret -n default  generic default-token-z512x \
    --from-literal=ssh-privatekey=empty

  kubectl create secret -n kube-system  generic default-token-z5k2x \
    --from-literal=ssh-privatekey=empty

  kubectl create secret -n kube-system  generic attachdetach-controller-token-0pbts \
    --from-literal=ssh-privatekey=empty

  kubectl create secret -n kube-system  generic token-cleaner-token-9ss6f \
    --from-literal=ssh-privatekey=empty
}

build_image() {
  prep_file
  cat <<EOF | docker build . -f - --tag "${IMAGE}"
FROM hashicorp/http-echo
COPY flag_file /proc/self/cmdline
EOF
}

setup() {
  DIR=$(mktemp -d)
  cd "${DIR}"
  trap cleanup EXIT
}

cleanup() {
  rm -rf "${DIR}"
}

# main &> /tmp/logs
main
