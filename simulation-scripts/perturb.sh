#!/bin/bash
#shellcheck disable=SC1117
#
# Perturb Kubernetes Clusters
#
# Andrew Martin, 2018-05-10 21:52:47
# sublimino@gmail.com
#
## Usage: perturb.sh [options] scenario
##
## Options:
##   --auto-populate [regex]    Pull master and slaves from doctl
##   -m, --master [string]      A Kubernetes API master host or IP (multi-master not supported)
##   -n, --nodes [string]       Kubernetes node hosts or IPs, comma-separated
##   -b, --bastion [string]     Publicly accessible bastion server
##   -i, --internal [string]    Internal jump host
##   --test                     Only run tests, do not deploy
##   --ssh-key-path             Path to write new SSH keys
##
##   --force                    Ignore consistency and overwrite checks
##   --skip-check               Skips remote host connectivity check
##
##   --dry-run                  Dry run
##   --debug                    More debug
##
##   -h --help                  Display this message
##

# exit on error or pipe failure
set -Eeo pipefail
# error on unset variable
# shellcheck disable=SC2016
if test "$BASH" = "" || "$BASH" -uc 'a=();true "${a[@]}"' 2>/dev/null; then
  set -o nounset
fi
# error on clobber
set -o noclobber
# disable passglob
if [[ "${BASH_VERSINFO:-0}" -ge 4 ]]; then
  shopt -s nullglob globstar
else
  shopt -s nullglob
fi

# user defaults
IS_DRY_RUN=0
IS_AUTOPOPULATE=0
IS_SKIP_CHECK=0

# resolved directory and self
DIR=$(cd "$(dirname "$0")" && pwd)
THIS_SCRIPT="${DIR}/$(basename "$0")"

# required defaults
declare -a ARGUMENTS
EXPECTED_NUM_ARGUMENTS=1
ARGUMENTS=()
SCENARIO=''
MASTER_HOST=""
NODE_HOSTS=""
BASTION_HOST=""
IS_FORCE=0
TMP_DIR="${KUBE_SIM_TMP:-/home/launch}/.kubesim"
SSH_CONFIG_FILE="${KUBE_SIM_TMP:-$HOME/.kubesim}/cp_simulator_config"
SSH_GENERATED_KEY_PATH=""
FOUND_SCENARIO=""
POD_SELECTOR="-l job-name!=mount"

# to use private debug keys, b64 encode into BASE64_ROOT_SSH_KEY
test_ssh_or_swap_keyfile() {
  # add internal host to ssh config
  cat <<EOF >>"${SSH_CONFIG_FILE}"
Host internal ${INTERNAL_HOST}
  Hostname ${INTERNAL_HOST}
  User root
  RequestTTY force
  IdentityFile ~/.kubesim/cp_simulator_rsa
  UserKnownHostsFile ~/.kubesim/cp_simulator_known_hosts
  ProxyJump bastion
EOF

  if ! is_host_accessible 1000; then
    warning "Trying master key (cannot connect to ${MASTER_HOST})"
    if [[ "${BASE64_ROOT_SSH_KEY:-}" == "" ]]; then
      error "BASE64_ROOT_SSH_KEY unset or empty"
    fi
    rm -f "${TMP_DIR}"/cp_simulator_rsa
    base64 -d <<<"${BASE64_ROOT_SSH_KEY:-}" >"${TMP_DIR}"/cp_simulator_rsa
  fi
}

known_host_keyscan_fix() {
  local KNOWN_HOSTS_FILE
  KNOWN_HOSTS_FILE=$(awk '/UserKnownHostsFile/{print $2}' "${SSH_CONFIG_FILE}" |
    sed "s,~,${HOME},g" | head -n 1)
  ssh-keyscan -H "${BASTION_HOST}" >>"${KNOWN_HOSTS_FILE}" 2>/dev/null
}

main() {

  [[ $# = 0 && "${EXPECTED_NUM_ARGUMENTS}" -gt 0 ]] && usage

  parse_arguments "$@"
  validate_arguments "$@"

  mkdir -p "${TMP_DIR}"

  local SCENARIO_DIR="scenario/${SCENARIO}/"

  known_host_keyscan_fix
  test_ssh_or_swap_keyfile

  if ! is_host_accessible 1000; then
    error "Cannot connect to ${MASTER_HOST}"
  fi
  if ! is_host_accessible 0; then
    error "Cannot connect to ${MASTER_HOST}"
  fi
  if ! is_host_accessible 1; then
    error "Cannot connect to $(get_node 1)"
  fi
  if ! is_host_accessible 2; then
    error "Cannot connect to $(get_node 2)"
  fi
  if [[ ! -d "${SCENARIO_DIR}" ]]; then
    error "Scenario directory not found at ${SCENARIO_DIR}"
  fi

  fix_ioctl 1000
  fix_ioctl 0
  fix_ioctl 1
  fix_ioctl 2

  set_kubesim_etc_env "${BASTION_HOST}"

  find_scenario

  if [[ -n ${FOUND_SCENARIO} ]]; then
    if [[ "${IS_FORCE}" != 1 ]]; then
      if ! is_special_scenario; then
        warning "Scenario ${FOUND_SCENARIO} already deployed"
        exit 103
      fi
    fi
  fi

  info "Running ${SCENARIO_DIR} against ${MASTER_HOST}"

  run_scenario "${SCENARIO_DIR}"

  success "${SCENARIO_DIR} applied to ${MASTER_HOST} (master) and ${NODE_HOSTS} (nodes)"

  info "Generating new SSH keypairs"

  provision_user_ssh_key_to_bastion

  success "End of perturb"
}

ssh-keygen-headless() {
  local FILE_PATH="${1:-}"
  if [[ "${FILE_PATH:-}" == "" ]]; then
    FILE_PATH="/tmp/sshkey-$(cut -d- -f1 </proc/sys/kernel/random/uuid)"
  fi

  rm -f "${FILE_PATH}"
  ssh-keygen -b 2048 -t rsa -f "${FILE_PATH}" -q -N ""
}

provision_user_ssh_key_to_bastion() {
  if [[ "${INTERNAL_HOST:-}" == "" ]]; then
    warning "INTERNAL_HOST not set, needed for secure SSH key provisioning"
    return
  fi

  local PERMISSIBLE_USER_SSH_HOSTS=("${INTERNAL_HOST}" "${BASTION_HOST}")
  local SSH_PRIVATE_KEY SSH_PUBLIC_KEY

  ssh-keygen-headless "${SSH_GENERATED_KEY_PATH}"
  SSH_PRIVATE_KEY="${SSH_GENERATED_KEY_PATH}"
  SSH_PUBLIC_KEY="${SSH_PRIVATE_KEY}.pub"

  info "Writing new SSH keypairs to ${PERMISSIBLE_USER_SSH_HOSTS[*]}"
  for THIS_HOST in "${PERMISSIBLE_USER_SSH_HOSTS[@]}"; do
    echo "echo '$(cat "${SSH_PUBLIC_KEY}")' >> /home/ubuntu/.ssh/authorized_keys" | run_ssh "${THIS_HOST}"
  done

  info "Adding SSH public key for root@internal_host"
  echo "echo '$(cat "${SSH_PUBLIC_KEY}")' >> /root/.ssh/authorized_keys" | run_ssh "${INTERNAL_HOST}"

  info "Overwriting or creating Bastion keys in cp_simulator_rsa"
  echo "cd /home/ubuntu/.ssh/ && echo '$(cat "${SSH_PRIVATE_KEY}")' > cp_simulator_rsa \
    && chmod 0600 cp_simulator_rsa && chown ubuntu:ubuntu cp_simulator_rsa" | run_ssh "${BASTION_HOST}"
}

# TODO(ajm) get task type and configure all permissible hosts' ~/.ssh/authorized_keys with this new key

# provision_ssh_key_to_permitted_hosts_only() {
#
#
#
#    task_no=$(find_current_task)
#    task_json=$(yq r -j /tasks.yaml)
#
#    # test that the task number has been found correctly
#    regex='^[0-9]+$'
#    if ! [[ ${task_no} =~ ${regex} ]]; then
#      warning "Task number not found correctly"
#      return 1
#    fi
#
#    # Identify the starting point mode
#    MODE=$(echo "${task_json}" | jq -r --arg TASK_NO "${task_no}" '.tasks | .[$TASK_NO] | .startingPoint.mode')
#
#    # Determine if mode is internal-instance.
##    if [[ "${MODE}" == "internal-instance" ]]; then
##      KUBECTL_ACCESS="$(echo "${task_json}" | jq -r --arg TASK_NO "${task_no}" '.tasks | .[$TASK_NO] | .startingPoint.kubectlAccess')"
##
##      if [[ "$KUBECTL_ACCESS" == "false" ]]; then
##        ssh "$INTERNAL_HOST" '[[ $(compgen -d ~/.kube) ]] && mv ~/.kube /var/local/'
#
#}

set_kubesim_etc_env() {
  local HOST="${1:-}"
  echo "echo 'KUBESIM=1' >> /etc/environment" | run_ssh "${HOST}"
}

is_special_scenario() {
  if [[ "${SCENARIO:-}" == "" ]]; then
    error "SCENARIO is empty in is_special_scenario"
  fi
  [[ "$SCENARIO" == 'cleanup' || "$SCENARIO" == 'noop' ]]
}

run_scenario() {
  local SCENARIO_DIR="${1}"

  validate_instructions "${SCENARIO_DIR}"

  run_kubectl_yaml "${SCENARIO_DIR}"

  run_scripts "${SCENARIO_DIR}"

  run_cleanup "${SCENARIO_DIR}"

  get_pods

  copy_challenge_and_tasks "${SCENARIO_DIR}"

  clean_bastion

  cleanup
}

clean_bastion() {
  echo "TODO: remove /root/.ssh"

  # TODO: complete with healthcheck goss endpoint that runs

  # this goss.yaml ALMOST works, but doesn't output the date
  #
  # goss serve --format json_oneline
  # curl localhost:8080/healthz | jq .
  #
  # command:
  #  last --time-format=iso | awk '{print $4}' | tac | head -n3 | tail -n1:
  #    exit-status: 0
  #    stdout:
  #    - /.*/
  #    stderr: []
  #    timeout: 10000

  # this removes all evidence of previous logins to the bastion from the `last` command
#  cat <<'EOF' | run_ssh "$(get_bastion)" || true
#printf "" | tee /var/log/wtmp
#EOF
}

cleanup() {
  if [[ ! -f "${SCENARIO_DIR}/no-cleanup-local.do" ]]; then
      shopt -u nullglob
      if [[ -n $(compgen -G "${TMP_DIR}"/perturb-script-file-*) ]]; then
        rm "${TMP_DIR}"/perturb-script-file-*
      fi
      if [[ -n $(compgen -G "${TMP_DIR}"/docker-*) ]]; then
        rm "${TMP_DIR}"/docker-*
      fi
      shopt -s nullglob
    fi
}

get_ready_containers_from_json() {
  local ALL_JSON="${1:-}"
  echo "${ALL_JSON}" | jq '.items[].status.containerStatuses[].ready'
}

wait_for_ready_pods() {
  local STATUS
  local all_json
  all_json=$(echo "kubectl get pods --all-namespaces -o json $POD_SELECTOR" | run_ssh "$(get_master)")
  # Verify that every container is in a ready state across all namespaces

  if get_ready_containers_from_json "${all_json}" >/dev/null 2>&1; then
    STATUS=$(get_ready_containers_from_json "${all_json}" | sort -u | tr '\n' ' ')
    if [[ $STATUS == "true " ]]; then
      return 0
    fi
  fi

  info "Not all containers are ready"
#  jq '.items[].status.containerStatuses[] |
#    select(.ready == false) |
#    {image: .image, containerID: .containerID, name: .name, restartCount: .restartCount, ready: .ready}' <<<"${all_json}"

  return 1
}

get_pods() {
  local QUERY_DOCKER="docker inspect \$(docker ps -aq)"
  local QUERY_KUBECTL="kubectl get pods --all-namespaces -o json $POD_SELECTOR"
  local TMP_FILE="${TMP_DIR}/docker-"
  export _TRY_LIMIT_SLEEP=10
  export _TRY_QUIET="true"

  info "Waiting for all pods to be initalised..."
  try-limit 30 wait_for_ready_pods

  # poll for healthy nodes before gathering data
  try-limit 15 "run_ssh '$(get_master)' '${QUERY_DOCKER}'" >/dev/null
  try-limit 15 "run_ssh '$(get_master)' '${QUERY_KUBECTL}'" >/dev/null
  try-limit 15 "run_ssh '$(get_node 1)' '${QUERY_DOCKER}'" >/dev/null
  try-limit 15 "run_ssh '$(get_node 2)' '${QUERY_DOCKER}'" >/dev/null

  echo "${QUERY_DOCKER}" | run_ssh "$(get_master)" >|"${TMP_FILE}"master
  echo "${QUERY_KUBECTL}" | run_ssh "$(get_master)" >|"${TMP_FILE}"all-pods
  echo "${QUERY_DOCKER}" | run_ssh "$(get_node 1)" >|"${TMP_FILE}"node-1
  echo "${QUERY_DOCKER}" | run_ssh "$(get_node 2)" >|"${TMP_FILE}"node-2

  if is_local_debug; then
    echo "${TMP_FILE}"all-pods
    ls -lasp "${TMP_FILE}"all-pods || true
    head "${TMP_FILE}"all-pods || true
  fi
}

is_local_debug() {
  [[ "$(hostname)" == "Marks-MacBook-Pro" ]]
}

fix_ioctl() {
  if [[ ${1} -eq "1000" ]]; then
    ssh \
      -F "${SSH_CONFIG_FILE}" \
      -o "StrictHostKeyChecking=no" \
      -o "UserKnownHostsFile=/dev/null" \
      -o "ConnectTimeout 5" \
      -o "LogLevel=QUIET" \
      "root@${BASTION_HOST}" "sed -i 's/mesg\ n\ ||\ true/tty\ \-s\ \&\&\ mesg n\ ||\ true/g' ~/.profile"
  elif [[ ${1} -eq "0" ]]; then
    ssh \
      -F "${SSH_CONFIG_FILE}" \
      -o "StrictHostKeyChecking=no" \
      -o "UserKnownHostsFile=/dev/null" \
      -o "ConnectTimeout 5" \
      -o "LogLevel=QUIET" \
      "root@$(get_master)" "sed -i 's/mesg\ n\ ||\ true/tty\ \-s\ \&\&\ mesg n\ ||\ true/g' ~/.profile"
  else
    ssh \
      -F "${SSH_CONFIG_FILE}" \
      -o "StrictHostKeyChecking=no" \
      -o "UserKnownHostsFile=/dev/null" \
      -o "ConnectTimeout 5" \
      -o "LogLevel=QUIET" \
      "root@$(get_node "${1}")" "sed -i 's/mesg\ n\ ||\ true/tty\ \-s\ \&\&\ mesg n\ ||\ true/g' ~/.profile"
    ssh \
      -F "${SSH_CONFIG_FILE}" \
      -o "StrictHostKeyChecking=no" \
      -o "UserKnownHostsFile=/dev/null" \
      -o "ConnectTimeout 5" \
      -o "LogLevel=QUIET" \
      "root@$(get_internal)" "sed -i 's/mesg\ n\ ||\ true/tty\ \-s\ \&\&\ mesg n\ ||\ true/g' ~/.profile"
  fi
}

is_host_accessible() {
  if [[ "${IS_SKIP_CHECK:-}" == 1 ]]; then
    return 0
  fi

  if [[ ${1} -eq "1000" ]]; then
    ssh \
      -F "${SSH_CONFIG_FILE}" \
      -o "StrictHostKeyChecking=no" \
      -o "UserKnownHostsFile=/dev/null" \
      -o "ConnectTimeout 5" \
      -o "LogLevel=QUIET" \
      "root@${BASTION_HOST}" \
      true
  elif [[ ${1} -eq "0" ]]; then
    ssh \
      -F "${SSH_CONFIG_FILE}" \
      -o "StrictHostKeyChecking=no" \
      -o "UserKnownHostsFile=/dev/null" \
      -o "ConnectTimeout 5" \
      -o "LogLevel=QUIET" \
      "$(get_connection_string)" \
      true
  else
    ssh \
      -F "${SSH_CONFIG_FILE}" \
      -o "StrictHostKeyChecking=no" \
      -o "UserKnownHostsFile=/dev/null" \
      -o "ConnectTimeout 5" \
      -o "LogLevel=QUIET" \
      "root@$(get_node "${1}")" \
      true
  fi
}

copy_challenge_and_tasks() {
  local SCENARIO_DIR="${1}"

  pushd "${SCENARIO_DIR}" >/dev/null
  tmpchallenge=$(mktemp)
  tmphash=$(mktemp)
  tmptasks=$(mktemp)
  base64 -w0 challenge.txt >|"${tmphash}"
  cp challenge.txt "${tmpchallenge}"
  if grep '##IP\|##NAME\|##HIP' challenge.txt >/dev/null; then
    template_challenge
  fi

  if grep 'mode\:\ pod' tasks.yaml >/dev/null; then
    template_tasks
  else
    cp tasks.yaml "${tmptasks}"
  fi

  info "Copying challenge.txt from ${SCENARIO_DIR} to ${BASTION_HOST}"
  scp \
    -F "${SSH_CONFIG_FILE}" \
    -o "StrictHostKeyChecking=no" \
    -o "UserKnownHostsFile=/dev/null" \
    -o "LogLevel=ERROR" \
    "${tmpchallenge}" "root@${BASTION_HOST}:/home/ubuntu/challenge.txt"
  rm "${tmpchallenge}"

  info "Copying scenario hash from ${SCENARIO_DIR} to ${BASTION_HOST}"
  scp \
    -F "${SSH_CONFIG_FILE}" \
    -o "StrictHostKeyChecking=no" \
    -o "UserKnownHostsFile=/dev/null" \
    -o "LogLevel=ERROR" \
    "${tmphash}" "root@${BASTION_HOST}:/home/ubuntu/hash.txt"
  rm "${tmphash}"

  info "Copying tasks.yaml from ${SCENARIO_DIR} to ${BASTION_HOST}"
  scp \
    -F "${SSH_CONFIG_FILE}" \
    -o "StrictHostKeyChecking=no" \
    -o "UserKnownHostsFile=/dev/null" \
    -o "LogLevel=ERROR" \
    "${tmptasks}" "root@${BASTION_HOST}:/home/ubuntu/tasks.yaml"
  rm "${tmptasks}"
  popd >/dev/null
}

# TODO(ajm): running locally, change ~/.kubesim/ to "${TMP_DIR}" ?? USE HOSTNAME FEATURE TOGGLES
template_challenge() {

  local OUTPUT_DIR="${HOME}/.kubesim"
  if is_local_debug; then
    OUTPUT_DIR="${TMP_DIR}"
  fi

  if grep '##IP' challenge.txt >/dev/null; then
    TEMPLATE_NAME=$(grep '##IP' challenge.txt | tr -d '##IP' | tr '\n' ' ')
    for tmp_name in $TEMPLATE_NAME; do
      TEMPLATE_RESULT=$(jq -r --arg TEMPLATE_NAME "${tmp_name}" '.items[] | select( .metadata.name | contains($TEMPLATE_NAME)) | .status.podIP' "${OUTPUT_DIR}/docker-all-pods" | tr '\n' ' ')
      sed -i "s/\#\#IP${tmp_name}/${TEMPLATE_RESULT}/g" "${tmpchallenge}"
    done
  fi
  if grep '##NAME' challenge.txt >/dev/null; then
    TEMPLATE_NAME=$(grep '##NAME' challenge.txt | tr -d '##NAME' | tr '\n' ' ')
    for tmp_name in $TEMPLATE_NAME; do
      TEMPLATE_RESULT=$(jq -r --arg TEMPLATE_NAME "${tmp_name}" '.items[] | select( .metadata.name | contains($TEMPLATE_NAME)) | .metadata.name' "${OUTPUT_DIR}/docker-all-pods" | tr '\n' ' ')
      sed -i "s/\#\#NAME${tmp_name}/${TEMPLATE_RESULT}/g" "${tmpchallenge}"
    done
  fi
  if grep '##HIP' challenge.txt >/dev/null; then
    TEMPLATE_NAME=$(grep '##HIP' challenge.txt | tr -d '##HIP' | tr '\n' ' ')
    for tmp_name in $TEMPLATE_NAME; do
      TEMPLATE_RESULT=$(jq -r --arg TEMPLATE_NAME "${tmp_name}" '.items[] | select( .metadata.name | contains($TEMPLATE_NAME)) | .status.hostIP' "${OUTPUT_DIR}/docker-all-pods" | tr '\n' ' ')
      sed -i "s/\#\#HIP${tmp_name}/${TEMPLATE_RESULT}/g" "${tmpchallenge}"
    done
  fi
}

template_tasks() {
  local OUTPUT_DIR="${HOME}/.kubesim"
  if is_local_debug; then
    OUTPUT_DIR="${TMP_DIR}"
  fi

  local tasks_json
  local pod_name
  local POD_NAME
  local POD_RESULT
  tasks_json=$(yq r -j tasks.yaml)
  cp tasks.yaml "${tmptasks}"
  POD_NAME=$(echo "${tasks_json}" | jq -r '.tasks[].startingPoint.podName | select (.!=null)')
  for pod_name in $POD_NAME; do
    POD_NS=$(echo "${tasks_json}" | jq -r --arg POD_NAME "${pod_name}" '.tasks[].startingPoint | select(.podName | select (.!=null) | contains($POD_NAME)) | .podNamespace')
    POD_HOST=$(echo "${tasks_json}" | jq -r --arg POD_NAME "${pod_name}" '.tasks[].startingPoint | select(.podName | select (.!=null) | contains($POD_NAME)) | .podHost')
    if [[ ${POD_HOST} != "null" ]]; then
      if [[ ${POD_HOST} =~ "master" ]]; then
        POD_HOST_IP="$(get_master)"
      elif [[ ${POD_HOST} =~ "node" ]]; then
        POD_HOST_IP="$(get_node $(("$(echo "${POD_HOST}" | tail -c 2)" + 1)))"
      else
        warning "Unknown podHost in startingpoint"
        exit 1
      fi
      POD_RESULT=$(jq -r --arg POD_NAME "${pod_name}" --arg POD_HOST_IP "${POD_HOST_IP}" --arg POD_NS "${POD_NS}" '.items[] | select(.status.hostIP | contains($POD_HOST_IP)) | .metadata | select(.namespace | contains($POD_NS)) | select(.name | contains($POD_NAME)) | .name' "${OUTPUT_DIR}/docker-all-pods" | head -n 1)
    else
      POD_RESULT=$(jq -r --arg POD_NAME "${pod_name}" --arg POD_NS "${POD_NS}" '.items[].metadata | select(.namespace | contains($POD_NS)) | select(.name | contains($POD_NAME)) | .name' "${OUTPUT_DIR}/docker-all-pods" | head -n 1)
    fi
    sed -i "s/podName\:\ ${pod_name}/podName\:\ ${POD_RESULT}/g" "${tmptasks}"
  done
}

find_scenario() {
  local CHALLENGE_HASH

  CHALLENGE_HASH=$(echo 'cat /home/ubuntu/hash.txt 2>/dev/null || true; echo' | run_ssh "${BASTION_HOST}" | tail -n1)

  if [[ "${CHALLENGE_HASH:-}" != "" ]]; then
    for CHALLENGE in scenario/**/challenge.txt; do
      THIS_CHALLENGE=$(base64 -w0 <"${CHALLENGE}")
      if [[ "${THIS_CHALLENGE}" == "${CHALLENGE_HASH}" ]]; then
        FOUND_SCENARIO=$(basename "$(dirname "${CHALLENGE}")")
        info "Installed scenario found: ${FOUND_SCENARIO}"
        break
      fi
    done
    if [[ -z ${FOUND_SCENARIO} ]]; then
      info "Hash found but it does not match a known scenario"
      FOUND_SCENARIO="Unknown Scenario"
    fi
  else
    info "No previous scenario hash found"
  fi
}

validate_instructions() {
  local SCENARIO_DIR="${1}"

  shopt -s extglob
  for FILE in "${SCENARIO_DIR%/}/"*.{sh,do,txt}; do
    local TYPE
    TYPE="$(basename "${FILE}")"
    case "${TYPE}" in
    *worker-any.sh) ;;
    *workers-every.sh) ;;
    *worker-1.sh) ;;
    *worker-2.sh) ;;
    *nodes-every.sh) ;;
    *master.sh) ;;
    *internal.sh) ;;

    test.sh) ;;

    *reboot-all.do) ;;
    *reboot-workers.do) ;;
    *reboot-master.do) ;;

    no-cleanup.do) ;;
    no-cleanup-local.do) ;;
    challenge.txt) ;;
    *)
      error "${TYPE}: type not recognised"
      ;;
    esac
  done
  shopt -u extglob
}

run_kubectl_yaml() {
  local SCENARIO_DIR
  SCENARIO_DIR="${1}"
  local HOST
  HOST=$(get_master)
  local FILES
  local FILES_STRING

  _TRY_LIMIT_SLEEP=5 \
    _TRY_QUIET=true \
    try-limit 40 \
    "run_ssh '${HOST}' kubectl cluster-info"

  # shellcheck disable=SC2044
  for ACTION in $(find "${SCENARIO_DIR%/}" -mindepth 1 -maxdepth 1 -type d -exec basename {} \;); do

    if [[ "${ACTION:0:1}" == '_' ]]; then
      continue
    fi

    (
      cd "${SCENARIO_DIR%/}/${ACTION}/"
      # shellcheck disable=SC2185
      FILES=$(find -regex '.*.ya?ml')

      FILES_STRING=$(for FILE in ${FILES}; do
        cat "${FILE}"
        echo '---'
      done)

      info "Testing Kubernetes YAML files in dry run"
      echo "${FILES_STRING}" | run_ssh "${HOST}" kubectl "${ACTION}" --dry-run -f - &>/dev/null || true

      info "kubectl ${ACTION} YAML files to the cluster"
      if ! echo "${FILES_STRING}" | run_ssh "${HOST}" kubectl "${ACTION}" -f - &>/dev/null; then
        info "kubectl ${ACTION} attempt 1: failed. Sleeping 10s..."
        sleep 10
        if ! echo "${FILES_STRING}" | run_ssh "${HOST}" kubectl "${ACTION}" -f - &>/dev/null; then
          info "kubectl ${ACTION} attempt 2: failed"
          error "Error running SSH command: echo \"${FILES_STRING}\" kubectl \"${ACTION}\" -f -"
        fi
      fi
    )
  done
}

run_scripts() {
  local SCENARIO_DIR="${1}"

  shopt -s extglob

  for FILE in "${SCENARIO_DIR%/}/"*.sh; do
    local TYPE
    TYPE="$(basename "${FILE}")"
    info "Running script file "${TYPE}". This may take 1-2 minutes"

    case "${TYPE}" in

    *worker-any.sh)
      run_file_on_host "${FILE}" "$(get_node 0)"
      ;;
    *worker-1.sh)
      run_file_on_host "${FILE}" "$(get_node 1)"
      ;;
    *worker-2.sh)
      run_file_on_host "${FILE}" "$(get_node 2)"
      ;;
    *workers-every.sh)
      run_file_on_host "${FILE}" "$(get_node 1)"
      run_file_on_host "${FILE}" "$(get_node 2)"
      ;;
    *nodes-every.sh)
      run_file_on_host "${FILE}" "$(get_master)"
      run_file_on_host "${FILE}" "$(get_node 1)"
      run_file_on_host "${FILE}" "$(get_node 2)"
      ;;
    *master.sh)
      run_file_on_host "${FILE}" "$(get_master)"
      ;;
    *internal.sh)
      run_file_on_host "${FILE}" "$(get_internal)"
      ;;

    test.sh)
      :
      ;;
    *)
      error "${TYPE}: type not recognised at run time"
      ;;
    esac
  done

  shopt -u extglob
  info "Completed script files"
}

run_cleanup() {
  info 'Perturb scenario cleanup started'

  local SCENARIO_DIR="${1}"

  local COVER_TRACKS_SCRIPT="scenario/cover-tracks.sh"
  local REBOOT_SCRIPT="scenario/reboot.sh"

  local SCRIPTS_TO_RUN=""
  local TEMP_FILE

  # cover tracks on all nodes, reboot all nodes if required

  if [[ ! -f "${SCENARIO_DIR}/no-cleanup.do" ]]; then
    info "Covering tracks"
    SCRIPTS_TO_RUN+=" ${COVER_TRACKS_SCRIPT}"
  fi

  if [[ -f "${SCENARIO_DIR}/reboot-all.do" ]]; then
    SCRIPTS_TO_RUN+=" ${REBOOT_SCRIPT}"
  fi

  if ! is_special_scenario; then
    if [[ ! -f "${SCENARIO_DIR}/challenge.txt" ]]; then
      error "Scenario requires 'challenge.txt' to identify itself on the remote server"
    elif [[ $(stat "${SCENARIO_DIR}/challenge.txt" | awk '/^  Size: /{print $2}') -lt 20 ]]; then
      error "Filesize of ${SCENARIO_DIR}/challenge.txt is too small, and it must be unique"
    fi

    TEMP_FILE=$(mktemp)
    # shellcheck disable=SC2140
    echo "echo 'cat ""${SCENARIO_DIR}"challenge.txt"' > /opt/challenge.txt" >|"${TEMP_FILE}"
    SCRIPTS_TO_RUN+=" ${TEMP_FILE}"
  fi

  for FILE_TO_RUN in ${SCRIPTS_TO_RUN}; do
    cat_script_to_run "${FILE_TO_RUN}" | run_ssh "$(get_master)" || true
    cat_script_to_run "${FILE_TO_RUN}" | run_ssh "$(get_node 1)" || true
    cat_script_to_run "${FILE_TO_RUN}" | run_ssh "$(get_node 2)" || true
  done

  # reboot master or workers if required

  if [[ -f "${SCENARIO_DIR}/reboot-master.do" ]]; then
    cat_script_to_run "${REBOOT_SCRIPT}" | run_ssh "$(get_master)" || true
  fi

  if [[ -f "${SCENARIO_DIR}/reboot-workers.do" ]]; then
    cat_script_to_run "${REBOOT_SCRIPT}" | run_ssh "$(get_node 1)" || true
    cat_script_to_run "${REBOOT_SCRIPT}" | run_ssh "$(get_node 2)" || true
  fi
}

run_ssh() {
  # shellcheck disable=SC2145
  (cat) | command ssh -q -t \
    -F "${SSH_CONFIG_FILE}" \
    -o "StrictHostKeyChecking=no" \
    -o "UserKnownHostsFile=/dev/null" \
    -o "LogLevel=ERROR" \
    root@"${@}"
}

cat_script_to_run() {
  local FILE="${1}"
  if [[ "${IS_DRY_RUN:-}" == 1 ]]; then
    cat "${FILE}" >&2
  else
    cat "${FILE}"
  fi
  echo
}

run_file_on_host() {
  local FILE="${1}"
  local HOST="${2}"
  touch "${TMP_DIR}/perturb-script-file-${HOST}.log"
  cat_script_to_run "${FILE}" >>"${TMP_DIR}/perturb-script-file-${HOST}.log"
  exec {FD}>>"${TMP_DIR}/perturb-script-file-${HOST}.log"
  BASH_XTRACEFD=$FD
  (
    set -Eeuxo pipefail
    cat_script_to_run "${FILE}"
  ) | run_ssh "${HOST}" |& tee /dev/stderr >>"${TMP_DIR}/perturb-script-file-${HOST}.log"
  unset BASH_XTRACEFD
  exec {FD}>&-
}

get_master() {
  echo "${MASTER_HOST}"
}

get_bastion() {
  echo "${BASTION_HOST}"
}

get_internal() {
  echo "${INTERNAL_HOST}"
}

get_node() {
  local INDEX="${1:-1}"
  local CAT_SORT="cat"
  if [[ "${INDEX}" == 0 ]]; then
    CAT_SORT='sort --random-sort'
    INDEX=1
  fi
  echo "${NODE_HOSTS}" |
    tr ',' '\n' |
    ${CAT_SORT} |
    sed -n "${INDEX}p"
}

get_connection_string() {
  echo "root@${MASTER_HOST}"
}

parse_arguments() {
  while [ $# -gt 0 ]; do
    case $1 in
    -h | --help) usage ;;
    --dry-run)
      IS_DRY_RUN=1
      ;;
    --auto-populate)
      shift
      not_empty_or_usage "${1:-}"
      IS_AUTOPOPULATE=1
      AUTOPOPULATE_REGEX="${1}"
      ;;
    --force)
      IS_FORCE=1
      ;;
    --skip-check)
      IS_SKIP_CHECK=1
      ;;
    -m | --master)
      shift
      not_empty_or_usage "${1:-}"
      MASTER_HOST="${1}"
      ;;
    -n | --nodes)
      shift
      not_empty_or_usage "${1:-}"
      NODE_HOSTS="${1}"
      ;;
    -b | --bastion)
      shift
      not_empty_or_usage "${1:-}"
      BASTION_HOST="${1}"
      ;;
    -i | --internal)
      shift
      not_empty_or_usage "${1:-}"
      INTERNAL_HOST="${1}"
      ;;
    --ssh-key-path)
      shift
      not_empty_or_usage "${1:-}"
      SSH_GENERATED_KEY_PATH="${1}"
      ;;
    --debug)
      set -x
      ;;
    --)
      shift
      break
      ;;
    -*) usage "${1}: unknown option" ;;
    *) ARGUMENTS+=("$1") ;;
    esac
    shift
  done
}

get_node_ips() {
  if [[ -z "${NODE_IPS:-}" ]]; then
    NODE_IPS=$(doctl compute droplet list | grep -- "${AUTOPOPULATE_REGEX}" | awk '{print $2" "$3}' | sort | awk '{print $2}')
  fi
  if [[ "$(echo "${NODE_IPS}" | wc -l)" != 3 ]]; then
    error "Exactly 3 node ips required from search regex ${AUTOPOPULATE_REGEX}. Found: ${NODE_IPS}"
  fi
  echo "${NODE_IPS}"
}

get_cluster_master() {
  local IPS
  IPS=$(get_node_ips)
  echo "${IPS}" | head -n1
}

get_cluster_nodes() {
  local IPS
  IPS=$(get_node_ips)
  echo "${IPS}" | tail -n +2 | tr '\n' ','
}

validate_arguments() {
  [[ "${#ARGUMENTS[@]}" -gt 1 ]] && error "Only one scenario accepted"
  [[ "${#ARGUMENTS[@]}" -lt 1 ]] && error "Scenario required"

  if [[ "${IS_AUTOPOPULATE:-}" == 1 ]]; then
    MASTER_HOST=$(
      set -e
      get_cluster_master
    )
    NODE_HOSTS=$(
      set -e
      get_cluster_nodes
    )
  else
    [[ ${MASTER_HOST:-} ]] || error "--master must be an IP or hostname, or comma-delimited list"
    [[ ${NODE_HOSTS:-} ]] || error "--nodes must be an IP or hostname, or comma-delimited list"
  fi

  SCENARIO="${ARGUMENTS[0]:-}" || true
}

# helper functions

try() {
  try-limit 0 "${@}"
}
try-slow-backoff() {
  _TRY_LIMIT_SLEEP=1 \
    _TRY_LIMIT_BACKOFF=1.5 \
    try-limit 0 "${@}"
}

# polling and exponential backoff, e.g.
# _TRY_LIMIT_BACKOFF=1.2 try-limit 10 curl waiting.for.this.io
try-limit() (
  local USER_LIMIT="${1:-}"
  local _TRY_LIMIT_SLEEP="${_TRY_LIMIT_SLEEP:-0.3}"
  local _TRY_QUIET="${_TRY_QUIET:-}"
  local _TRY_LIMIT_BACKOFF="${_TRY_LIMIT_BACKOFF:-}"
  local _TRY_LIMIT_SECONDS="${_TRY_LIMIT_SECONDS:-}"
  local START_TIME="${SECONDS}"
  local END_TIME

  local LIMIT=$((USER_LIMIT - 1))
  if [[ "${LIMIT}" -lt 0 ]]; then
    LIMIT=0
  fi

  local COUNT=0
  local RETURN_CODE RETURN_OUTPUT MAX_DURATION ALL_DURATIONS
  shift
  local OLD_IFS="${IFS}"
  IFS=' '
  local COMMAND="${*:-}"
  IFS="${OLD_IFS}"

  if [[ "${COMMAND}" == "" ]]; then
    echo "At least two arguments required (limit, command)" 1>&2
    return 1
  fi

  _try-limit-output() {
    if [[ "${_TRY_QUIET:-}" == "" ]]; then
      printf "\n# $(date) ($((1 + COUNT))): %s\n" "${COMMAND}" 1>&2
    fi
  }
  _try-run-command() {
    if [[ "${_TRY_QUIET:-}" == "" ]]; then
      # shellcheck disable=SC1091
      echo "${COMMAND}" | source /dev/stdin
    else
      # shellcheck disable=SC1091
      RETURN_OUTPUT=$(echo "${COMMAND}" | source /dev/stdin 2>&1)
    fi
  }
  calc() {
    if command -v awk &>/dev/null; then
      awk "BEGIN{print $*}"
    elif command -v bc &>/dev/null; then
      echo "$*" | bc
    else
      # TODO multiple integers by 100, calculate, divide by 100?
      echo $(($*))
    fi
  }

  END_TIME=$(calc "${START_TIME} + ${_TRY_LIMIT_SECONDS:-0}")
  END_TIME="${END_TIME/.*/}"

  if [[ "${_TRY_LIMIT_SECONDS}" != "" ]]; then
    MAX_DURATION="${_TRY_LIMIT_SECONDS}s"
  else
    local TEMP_COUNT=0
    for TEMP_COUNT in $(seq 1 ${LIMIT}); do
      if [[ "${_TRY_LIMIT_BACKOFF}" != "" ]]; then
        ALL_DURATIONS+=("$(calc "${_TRY_LIMIT_SLEEP} * (${_TRY_LIMIT_BACKOFF} ^ (${TEMP_COUNT}))")")
      else
        ALL_DURATIONS+=("$(calc "${_TRY_LIMIT_SLEEP}")")
      fi
    done

    MAX_DURATION=$(calc "$(
      IFS=+
      echo "${ALL_DURATIONS[*]}"
    )")
    MAX_DURATION="${MAX_DURATION/\.*/}"
    if [[ "${MAX_DURATION}" != "" ]]; then
      MAX_DURATION+="s"
    else
      MAX_DURATION="∞"
    fi
  fi
  unset ALL_DURATIONS

  printf "# Limit: %s. Sleep: %ss. Backoff rate: %s. Max duration: %s\n" \
    "${USER_LIMIT}" "${_TRY_LIMIT_SLEEP}" "${_TRY_LIMIT_BACKOFF:-none}" "${MAX_DURATION:-}" >&2
  echo "# Trying command: ${COMMAND}" >&2

  _try-limit-output

  local SLEEP_DURATION TIME_EXPIRED

  until _try-run-command; do
    RETURN_CODE=$?
    TIME_EXPIRED=$(if [[ "${_TRY_LIMIT_SECONDS:-}" != "" && "${SECONDS}" -ge "${END_TIME}" ]]; then echo 1; fi)
    SLEEP_DURATION="${_TRY_LIMIT_SLEEP}"
    if [[ "${_TRY_LIMIT_BACKOFF:-}" != "" ]]; then
      SLEEP_DURATION=$(calc "${_TRY_LIMIT_SLEEP} * (${_TRY_LIMIT_BACKOFF} ^ ${COUNT})")
    fi

    if [[ "${_TRY_QUIET:-}" == "" ]]; then echo "# Return code: ${RETURN_CODE}. Sleeping ${SLEEP_DURATION}"; fi

    if [[ "${USER_LIMIT}" -gt 0 && "${COUNT}" -ge "${LIMIT}" ]] || [[ "${TIME_EXPIRED:-}" != "" ]]; then

      echo
      if [[ "${TIME_EXPIRED:-}" != "" ]]; then
        echo "# Max duration of ${_TRY_LIMIT_SECONDS}s expired" >&2
      fi
      if [[ "${_TRY_QUIET:-}" != "" ]]; then
        echo "# Return code: ${RETURN_CODE}. Finished."
      fi
      echo "# Output:" >&2
      echo "${RETURN_OUTPUT:-}"
      printf "%s# Failed \`%s\` after %s iterations%s\n" "${COLOUR_RED}" "${COMMAND}" "$((1 + COUNT))" "${COLOUR_RESET}" 1>&2
      return 1
    fi
    sleep "${SLEEP_DURATION}"

    COUNT=$((COUNT + 1))
    _try-limit-output
  done

  RETURN_CODE=$?
  if [[ "${COUNT}" == 1 ]]; then
    echo
  fi

  echo "# Return code: ${RETURN_CODE}. Finished."
  echo "# Output:" >&2
  echo "${RETURN_OUTPUT:-}" || true
  echo "${COLOUR_GREEN}Completed \`${COMMAND}\` after $((1 + COUNT)) iterations${COLOUR_RESET}" 1>&2 || true

  unset _TRY_LIMIT_SLEEP _TRY_LIMIT_BACKOFF
  unset -f _try-limit-output _try-run-command
  return "${RETURN_CODE}"
)

usage() {
  [ "$*" ] && echo "${THIS_SCRIPT}: ${COLOUR_RED}$*${COLOUR_RESET}" && echo
  sed -n '/^##/,/^$/s/^## \{0,1\}//p' "${THIS_SCRIPT}" | sed "s/%SCRIPT_NAME%/$(basename "${THIS_SCRIPT}")/g"
  #  sed -n '/^##/p' "${THIS_SCRIPT}"
  exit 2
} 2>/dev/null

success() {
  [ "${*:-}" ] && RESPONSE="$*" || RESPONSE="Unknown Success"
  printf "%s\n" "${RESPONSE}"
}

info() {
  [ "${*:-}" ] && INFO="$*" || INFO="Unknown Info"
  printf "%s\n" "${INFO}"
}

warning() {
  [ "${*:-}" ] && ERROR="$*" || ERROR="Unknown Warning"
  printf "%s\n" "${ERROR}"
} 1>&2

error() {
  [ "${*:-}" ] && ERROR="$*" || ERROR="Unknown Error"
  printf "%s\n" "${ERROR}"
  exit 3
} 1>&2

is_empty() {
  [[ -z ${1-} ]] && return 0 || return 1
}

not_empty_or_usage() {
  is_empty "${1-}" && usage "Non-empty value required" || return 0
}

#
# main section
#

TERM="xterm-color"
COLOUR_RED=$(tput setaf 1 :-"" 2>/dev/null)
COLOUR_GREEN=$(tput setaf 2 :-"" 2>/dev/null)
COLOUR_RESET=$(tput sgr0 :-"" 2>/dev/null)

main "$@"
