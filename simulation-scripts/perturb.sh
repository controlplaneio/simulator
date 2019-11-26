#!/bin/bash
#
# Perturb Kubernetes Clusters
#
# Created by: Andrew Martin, 2018-05-10 21:52:47 sublimino@gmail.com
#
# Refactored for kubesim by: Jon Kent 2019-11-21
#
## Usage: perturb.sh [options] scenario
##
## Options:
##   -m, --master [string]      A Kubernetes API master host or IP (multi-master not supported)
##   -s, --slaves [string]      Kubernetes slave hosts or IPs, comma-separated
##   -b, --bastion [string]     Publicly accessible bastion server
##
##   -h --help                  Display this message
##

# exit on error or pipe failure
set -eo pipefail

# error on clobber
set -o noclobber

if [[ "${BASH_VERSINFO:-0}" -ge 4 ]];  then
  shopt -s nullglob globstar
  set -o nounset
else
  shopt -s nullglob
fi

# user defaults
SSH_CONFIG_FILE="$HOME/.ssh/cp_simulator_config"

# resolved directory and self
DIR=$(cd "$(dirname "$0")" && pwd)
THIS_SCRIPT="${DIR}/$(basename "$0")"

# required defaults
declare -a ARGUMENTS
ARGUMENTS=()
SCENARIO=''
MASTER_HOST=""
SLAVE_HOSTS=""
BASTION_HOST=""

main() {

  parse_arguments "$@"

  [[ "${#ARGUMENTS[@]}" -gt 1 ]] && error "Only one scenario accepted"

  [[ ${MASTER_HOST:-} ]] || error "--master must be an IP or hostname, or comma-delimited list"
  [[ ${SLAVE_HOSTS:-} ]] || error "--slaves must be an IP or hostname, or comma-delimited list"

  SCENARIO="${ARGUMENTS[0]:-}" || true

  local SCENARIO_DIR="scenario/${SCENARIO}/"

  if [ ! -d "${SCENARIO_DIR}" ]; then
    error "${SCENARIO_DIR} not found"
  fi

  info "Running ${SCENARIO_DIR} against ${MASTER_HOST}"

  if ! is_master_accessible ;then
    error "Cannot connect to ${MASTER_HOST}"
  elif [[ ! -d "${SCENARIO_DIR}" ]]; then
    error "Scenario directory not found at ${SCENARIO_DIR}"
  fi

  warning "Instructions in scenario:"
  ls -lasp "${SCENARIO_DIR}"

  validate_instructions "${SCENARIO_DIR}"

  copy_challenge_and_tasks "${SCENARIO_DIR}"

  run_kubectl_yaml "${SCENARIO_DIR}"

  run_scripts "${SCENARIO_DIR}"

  success "${SCENARIO_DIR} applied to ${MASTER_HOST} (master) and ${SLAVE_HOSTS} (slaves)"

  success "End of perturb"
}

parse_arguments() {
  while [ $# -gt 0 ]; do
    case $1 in
      -h | --help) usage ;;
      -m | --master)
        shift
        not_empty_or_usage "${1:-}"
        MASTER_HOST="${1}"
        ;;
      -s | --slaves)
        shift
        not_empty_or_usage "${1:-}"
        SLAVE_HOSTS="${1}"
        ;;
      -b | --bastion)
        shift
        not_empty_or_usage "${1:-}"
        BASTION_HOST="${1}"
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

is_master_accessible() {
  ssh \
    -F "${SSH_CONFIG_FILE}"  \
    -o "StrictHostKeyChecking=no" \
    -o "UserKnownHostsFile=/dev/null" \
    -o "ConnectTimeout 3" \
    "$(get_connection_string)" \
    true
}

copy_challenge_and_tasks() {
  local SCENARIO_DIR="${1}"

  pushd "${SCENARIO_DIR}"
  info "Copying challenge.txt from ${SCENARIO_DIR} to ${BASTION_HOST}"
  scp \
    -F "${SSH_CONFIG_FILE}"  \
    -o "StrictHostKeyChecking=no" \
    -o "UserKnownHostsFile=/dev/null" \
    challenge.txt root@${BASTION_HOST}:/home/ubuntu/challenge.txt
  info "Copying tasks.yaml from ${SCENARIO_DIR} to ${BASTION_HOST}"
  scp \
    -F "${SSH_CONFIG_FILE}"  \
    -o "StrictHostKeyChecking=no" \
    -o "UserKnownHostsFile=/dev/null" \
    tasks.yaml root@${BASTION_HOST}:/home/ubuntu/tasks.yaml
  popd
}

validate_instructions() {
  local SCENARIO_DIR="${1}"

  shopt -s extglob
  for FILE in "${SCENARIO_DIR%/}/"*.{sh,do,txt}; do
    echo "${FILE}"
    local TYPE; TYPE="$(basename "${FILE}")"
    case "${TYPE}" in
      *worker-any.sh) ;;
      *workers-every.sh) ;;
      *worker-1.sh) ;;
      *worker-2.sh) ;;
      *nodes-every.sh) ;;
      *master.sh) ;;

      test.sh) ;;

      *reboot-all.do) ;;
      *reboot-workers.do) ;;
      *reboot-master.do) ;;

      no-cleanup.do) ;;
      challenge.txt) ;;
      *)
        error "${TYPE}: type not recognised"
        ;;
    esac
  done
  shopt -u extglob
}

run_kubectl_yaml() {
  local SCENARIO_DIR; SCENARIO_DIR="${1}"
  local HOST; HOST=$(get_master)
  local FILES
  local FILES_STRING

  # shellcheck disable=SC2044
  for ACTION in $(find "${SCENARIO_DIR%/}" -mindepth 1 -maxdepth 1 -type d -exec basename {} \;); do

    if [[ "${ACTION:0:1}" == '_' ]]; then
      continue
    fi

    (
      cd "${SCENARIO_DIR%/}/${ACTION}/"
      # shellcheck disable=SC2185
      FILES=$(find -regex '.*.ya?ml')

      info "running remotely: kubectl ${ACTION} -f ${FILES}"

      FILES_STRING=$(for FILE in ${FILES}; do
        cat "${FILE}"
        echo '---'
      done)

      echo "${FILES_STRING}" | run_ssh "${HOST}" kubectl "${ACTION}" --dry-run -f -
      echo "${FILES_STRING}" | run_ssh "${HOST}" kubectl "${ACTION}" -f -
    )
  done
}

run_scripts() {
  local SCENARIO_DIR="${1}"

  shopt -s extglob

  for FILE in "${SCENARIO_DIR%/}/"*.sh; do
    echo "${FILE}"
    local TYPE; TYPE="$(basename "${FILE}")"
    case "${TYPE}" in

      *worker-any.sh)
        run_file_on_host "${FILE}" "$(get_slave 0)"
        ;;
      *worker-1.sh)
        run_file_on_host "${FILE}" "$(get_slave 1)"
        ;;
      *worker-2.sh)
        run_file_on_host "${FILE}" "$(get_slave 2)"
        ;;
      *workers-every.sh)
        run_file_on_host "${FILE}" "$(get_slave 1)"
        run_file_on_host "${FILE}" "$(get_slave 2)"
        ;;
      *nodes-every.sh)
        run_file_on_host "${FILE}" "$(get_master)"
        run_file_on_host "${FILE}" "$(get_slave 1)"
        run_file_on_host "${FILE}" "$(get_slave 2)"
        ;;
      *master.sh)
        run_file_on_host "${FILE}" "$(get_master)"
        ;;

       test.sh)
        : ;;
      *)
        error "${TYPE}: type not recognised at run time"
        ;;
    esac
  done

  shopt -u extglob

}


run_ssh() {
  # shellcheck disable=SC2145
  (cat) | command ssh -q -t \
    -F "${SSH_CONFIG_FILE}" \
    -o "StrictHostKeyChecking=no" \
    -o "UserKnownHostsFile=/dev/null" \
    root@"${@}"
}

run_file_on_host() {
  local FILE="${1}"
  local HOST="${2}"

  info "Setup scenario to be perturbed"

  scp \
    -F "${SSH_CONFIG_FILE}"  \
    -o "StrictHostKeyChecking=no" \
    -o "UserKnownHostsFile=/dev/null" \
    /app/simulation-scripts/"${FILE}" root@"${HOST}":/tmp/setup.sh

  ssh -q -t \
    -F "${SSH_CONFIG_FILE}" \
    -o "StrictHostKeyChecking=no" \
    -o "UserKnownHostsFile=/dev/null" \
    root@"${HOST}" "chmod +x /root/setup.sh && /tmp/setup.sh"

  info "Cleanup setup"
  ssh -q -t \
    -F "${SSH_CONFIG_FILE}" \
    -o "StrictHostKeyChecking=no" \
    -o "UserKnownHostsFile=/dev/null" \
    root@"${HOST}" "rm /tmp/setup.sh"

  echo "PERTURBANCE COMPLETE FOR /simulation-scripts/${FILE} ON ${HOST}"
}

get_master() {
  echo "${MASTER_HOST}"
}

get_slave() {
  local INDEX="${1:-1}"
  local CAT_SORT="cat"
  if [[ "${INDEX}" == 0 ]]; then
    CAT_SORT='sort --random-sort'
    INDEX=1
  fi
  echo "${SLAVE_HOSTS}" \
    | tr ',' '\n' \
    | ${CAT_SORT} \
    | sed -n "${INDEX}p"
}

# helper functions

usage() {
  [ "$*" ] && echo "${THIS_SCRIPT}: ${COLOUR_RED}$*${COLOUR_RESET}" && echo
  sed -n '/^##/p' "${THIS_SCRIPT}"
  exit 2
} 2>/dev/null

success() {
  [ "${*:-}" ] && RESPONSE="$*" || RESPONSE="Unknown Success"
  printf "%s\n" "$(log_message_prefix)${COLOUR_GREEN}${RESPONSE}${COLOUR_RESET}"
} 1>&2

info() {
  [ "${*:-}" ] && INFO="$*" || INFO="Unknown Info"
  printf "%s\n" "$(log_message_prefix)${COLOUR_WHITE}${INFO}${COLOUR_RESET}"
} 1>&2

warning() {
  [ "${*:-}" ] && ERROR="$*" || ERROR="Unknown Warning"
  printf "%s\n" "$(log_message_prefix)${COLOUR_RED}${ERROR}${COLOUR_RESET}"
} 1>&2

error() {
  [ "${*:-}" ] && ERROR="$*" || ERROR="Unknown Error"
  printf "%s\n" "$(log_message_prefix)${COLOUR_RED}${ERROR}${COLOUR_RESET}"
  exit 3
} 1>&2

log_message_prefix() {
  TIMESTAMP="[$(date +'%Y-%m-%dT%H:%M:%S%z')]"
  THIS_SCRIPT_SHORT=${THIS_SCRIPT/$DIR/.}
  tput bold 2>/dev/null
  echo -n "${TIMESTAMP} ${THIS_SCRIPT_SHORT}: "
}

is_empty() {
  [[ -z ${1-} ]] && return 0 || return 1
}

not_empty_or_usage() {
  is_empty "${1-}" && usage "Non-empty value required" || return 0
}

###########################################
#
# main section
#
###########################################

TERM="xterm-color"
COLOUR_RED=$(tput setaf 1 :-"" 2>/dev/null)
COLOUR_GREEN=$(tput setaf 2 :-"" 2>/dev/null)
COLOUR_WHITE=$(tput setaf 7 :-"" 2>/dev/null)
COLOUR_RESET=$(tput sgr0 :-"" 2>/dev/null)

main "$@"
