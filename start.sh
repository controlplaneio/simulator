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

docker pull ${CONTAINER_NAME}

docker run                                                          \
  -h launch                                                         \
  -v ${SIMULATOR_AWS_CREDS_PATH}:/home/launch/.aws                  \
  -v ${SSH_CONFIG_PATH}:/home/launch/.ssh                           \
  -v ${KUBE_SIM_TMP}:/home/launch/.kubesim                          \
  -e "AWS_ACCESS_KEY_ID"                                            \
  -e "AWS_SECRET_ACCESS_KEY"                                        \
  -e "AWS_REGION"                                                   \
  -e "AWS_DEFAULT_REGION"                                           \
  -e "AWS_PROFILE"                                                  \
  -e "AWS_DEFAULT_PROFILE"                                          \
  --rm --init -it ${CONTAINER_NAME}
