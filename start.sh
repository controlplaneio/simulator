#!/bin/bash

if [[ -z ${AWS_SHARED_CREDENTIALS_FILE} ]]; then
  SIMULATOR_AWS_CREDS_PATH=${HOME}/.aws/
else
  SIMULATOR_AWS_CREDS_PATH=$(dirname ${AWS_SHARED_CREDENTIALS_FILE})
fi
CONTAINER_NAME=s1rvine/kubesim:test
SSH_CONFIG_PATH=${HOME}/.ssh/
KUBE_SIM_TMP=${HOME}/.kubesim/
SIMULATOR_CONFIG_FILE=${KUBE_SIM_TMP}/simulator.yaml
touch ${SIMULATOR_CONFIG_FILE}

curl https://raw.githubusercontent.com/controlplaneio/simulator-standalone/master/scripts/validate-requirements | bash

docker pull ${CONTAINER_NAME}

docker run                                                          \
  -h launch                                                         \
  -v ${SIMULATOR_CONFIG_FILE}:/app/simulator.yaml                   \
  -v ${SIMULATOR_AWS_CREDS_PATH}:/home/launch/.aws                  \
  -v ${SSH_CONFIG_PATH}:/home/launch/.ssh                           \
  -v ${KUBE_SIM_TMP}:/home/launch/.kubesim                          \
  --rm --init -it ${CONTAINER_NAME}
