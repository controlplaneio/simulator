#!/usr/bin/env bash
#
# Perturb Kubernetes Clusters
#
# Andrew Martin, 2018-05-10 21:52:47
# sublimino@gmail.com
#
## Usage: %SCRIPT_NAME% [options] scenario
##
## Options:
##   --auto-populate [regex]    Pull master and slaves from doctl
##   -m, --master [string]      A Kubernetes API master host or IP (multi-master not supported)
##   -s, --slaves [string]      Kubernetes slave hosts or IPs, comma-separated
##   --test                     Only run tests, do not deploy
##
##   --force                    Ignore consistency and overwrite checks
##   --skip-check               Skips remote host connectivity check
##   --add-keys [user,user...]  GitHub user keys to add (comma delimited)
##
##   --dry-run                  Dry run
##   --debug                    More debug
##
##   -h --help                  Display this message
##

# exit on error or pipe failure
set -eo pipefail
# error on unset variable
if test "$BASH" = "" || "$BASH" -uc 'a=();true "${a[@]}"' 2>/dev/null; then
  set -o nounset
fi
# error on clobber
set -o noclobber
# disable passglob
shopt -s nullglob globstar

# user defaults
DESCRIPTION="Perturb Kubernetes Clusters"
DEBUG=0
IS_DRY_RUN=0
IS_AUTOPOPULATE=0
IS_SKIP_CHECK=0

# resolved directory and self
declare -r DIR=$(cd "$(dirname "$0")" && pwd)
declare -r THIS_SCRIPT="${DIR}/$(basename "$0")"

# required defaults
declare -a ARGUMENTS
EXPECTED_NUM_ARGUMENTS=1
ARGUMENTS=()
SCENARIO=''
MASTER_HOST=""
SLAVE_HOSTS=""
IS_TEST_ONLY=0
IS_FORCE=0

main() {
  handle_arguments "$@"

  local SCENARIO_DIR="scenario/${SCENARIO}/"

  if [[ "${IS_TEST_ONLY:-}" == 1 ]]; then
    success "In test mode, skipping deployment of ${SCENARIO}"
  else

    source test-func.sh

    local FOUND_SCENARIO
    if FOUND_SCENARIO=$(find_scenario); then
      if [[ "${IS_FORCE}" != 1 && "${FOUND_SCENARIO}" == "${SCENARIO}" ]]; then
        if ! is_special_scenario; then
          error "Scenario ${SCENARIO} already deployed, reset deployment with 'cleanup' first"
        fi
      fi
      info "Found scenario ${FOUND_SCENARIO}"
    fi

    info "Running ${SCENARIO_DIR} against ${MASTER_HOST}"
  fi

  if ! is_master_accessible; then
    error "Cannot connect to ${MASTER_HOST}"
  elif ! is_scenario_dir_accessible; then
    error "Scenario directory not found at ${SCENARIO_DIR}"
  fi

  if [[ "${IS_TEST_ONLY:-}" != 1 ]]; then
    run_scenario "${SCENARIO_DIR}"

    success "${SCENARIO_DIR} applied to ${MASTER_HOST} (master) and ${SLAVE_HOSTS} (slaves)"
  fi

  if [[ "${IS_TEST_ONLY:-}" == 1 ]]; then
    # force use of this scenario
    run_test_loop "${SCENARIO}"
  else
    # autodetect scenario
    run_test_loop
  fi

  success "End of perturb"
}

is_special_scenario() {
  if [[ "${SCENARIO:-}" == "" ]]; then
    error "SCENARIO is empty in is_special_scenario"
  fi
  [[ "$SCENARIO" == 'cleanup' || "$SCENARIO" == 'noop' ]]
}

run_test_loop() {

  local IS_SUCCESS=0
  local TEST_RESULT=0
  local LOCAL_SCENARIO="${1:-}"

  while true; do

    source test-func.sh

    if test_scenario "${LOCAL_SCENARIO:-}"; then
      IS_SUCCESS=1
      break
    else
      TEST_RESULT=$?
      if [[ "${TEST_RESULT}" -eq 99 ]]; then
        warning "No tests found"
        IS_SUCCESS=0
        break
      fi
    fi

    echo "Test again? [n/q/Y]"
    read_prompt

  done

  if [[ "${IS_SUCCESS:-}" == 1 ]]; then
    success "finished"
  else
    warning "tests failed"
  fi
}

run_scenario() {
  local SCENARIO_DIR="${1}"

  warning "Instructions in scenario:"
  ls -lasp "${SCENARIO_DIR}"

  validate_instructions "${SCENARIO_DIR}"

  run_kubectl_yaml "${SCENARIO_DIR}"

  run_scripts "${SCENARIO_DIR}"

  run_cleanup "${SCENARIO_DIR}"

  add_github_keys
}

add_github_keys() {
  if [[ "${GITHUB_KEY_USERS:-}" != "" ]]; then
    local FILE
    FILE="scenario/authorized_keys.sh"
    (cat "${FILE}"; echo "write_keys ${GITHUB_KEY_USERS//,/ }") | run_ssh "$(get_master)"
    (cat "${FILE}"; echo "write_keys ${GITHUB_KEY_USERS//,/ }") | run_ssh "$(get_slave 1)"
    (cat "${FILE}"; echo "write_keys ${GITHUB_KEY_USERS//,/ }") | run_ssh "$(get_slave 2)"
  fi
}

is_master_accessible() {
  if [[ "${IS_SKIP_CHECK:-}" == 1 ]]; then
    return 0
  fi
  ssh \
    -o "StrictHostKeyChecking=no" \
    -o "UserKnownHostsFile=/dev/null" \
    -o "ConnectTimeout 3" \
    $(get_connection_string) \
    true
}

is_scenario_dir_accessible() {
  [[ -d "${SCENARIO_DIR}" ]]
}

validate_instructions() {
  local SCENARIO_DIR="${1}"

  shopt -s extglob
  for FILE in "${SCENARIO_DIR%/}/"*.{sh,do}; do
    echo "${FILE}"
    local TYPE="$(basename "${FILE}")"
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
      *)
        error "${TYPE}: type not recognised"
        ;;
    esac
  done
  shopt -u extglob
}

run_kubectl_yaml() {
  local SCENARIO_DIR="${1}"
  local HOST=$(get_master)
  local FILES
  local FILES_STRING

  for ACTION in $(find "${SCENARIO_DIR%/}" -mindepth 1 -maxdepth 1 -type d -exec basename {} \;); do

    if [[ "${ACTION:0:1}" == '_' ]]; then
      continue
    fi

    (
      cd "${SCENARIO_DIR%/}/${ACTION}/"
      FILES=$(find -regex '.*.ya?ml')

      info "running remotely: kubectl ${ACTION} -f ${FILES}"

      FILES_STRING=$(for FILE in ${FILES}; do
        cat "${FILE}"
        echo '---'
      done)

      echo "${FILES_STRING}" | run_ssh ${HOST} kubectl "${ACTION}" --dry-run -f -
      echo "${FILES_STRING}" | run_ssh ${HOST} kubectl "${ACTION}" -f -
    )
  done
}

run_scripts() {
  local SCENARIO_DIR="${1}"

  shopt -s extglob

  for FILE in "${SCENARIO_DIR%/}/"*.sh; do
    echo "${FILE}"
    local TYPE="$(basename "${FILE}")"
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

run_cleanup() {
  warning 'Cleanup started'

  local SCENARIO_DIR="${1}"

  local COVER_TRACKS_SCRIPT="scenario/cover-tracks.sh"
  local REBOOT_SCRIPT="scenario/reboot.sh"

  local SCRIPTS_TO_RUN=""
  local TEMP_FILE

  # cover tracks on all nodes, reboot all nodes if required

  if [[ ! -f "${SCENARIO_DIR}/no-cleanup.do" ]]; then
    info "Covering tracks..."
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
    echo "echo '$(cat "${SCENARIO_DIR}/challenge.txt")' > /opt/challenge.txt" | tee "${TEMP_FILE}"
    SCRIPTS_TO_RUN+=" ${TEMP_FILE}"
  fi

  if [[ -f "${SCENARIO_DIR}/flag.txt" ]]; then
    TEMP_FILE=$(mktemp)
    echo "echo '$(cat "${SCENARIO_DIR}/flag.txt" | base64 -w0)' | base64 -d > /root/flag.txt" | tee "${TEMP_FILE}"
    warning "Flag"
    cat "${TEMP_FILE}"
    SCRIPTS_TO_RUN+=" ${TEMP_FILE}"
  fi

  warning "Running script '${SCRIPTS_TO_RUN}'"

  for FILE_TO_RUN in ${SCRIPTS_TO_RUN}; do
    get_file_to_run "${FILE_TO_RUN}" | run_ssh "$(get_master)" || true
    get_file_to_run "${FILE_TO_RUN}" | run_ssh "$(get_slave 1)" || true
    get_file_to_run "${FILE_TO_RUN}" | run_ssh "$(get_slave 2)" || true
  done

  # reboot master or workers if required

  if [[ -f "${SCENARIO_DIR}/reboot-master.do" ]]; then
    get_file_to_run "${REBOOT_SCRIPT}" | run_ssh "$(get_master)" || true
  fi

  if [[ -f "${SCENARIO_DIR}/reboot-workers.do" ]]; then
    get_file_to_run "${REBOOT_SCRIPT}" | run_ssh "$(get_slave 1)" || true
    get_file_to_run "${REBOOT_SCRIPT}" | run_ssh "$(get_slave 2)" || true
  fi
}

get_flag() {
  local SALT=$(get_salt)
  FLAG=$({ get_flag_raw; echo "${SALT}"; } | sha256sum)
  echo "${FLAG}"
}

get_flag_raw() {
  if [[ ! -f "${SCENARIO_DIR}/flag.txt" ]]; then
    error "${SCENARIO_DIR}/flag.txt not found"
  fi
  cat "${SCENARIO_DIR}/flag.txt"
}

get_salt() {
  date +%Y-%m%Y-%m%Y-%m%Y-%m%Y-%m%Y-%m
}

run_ssh() {
# TODO: test everything with `set -ex`
# (printf "%s\n\n" 'set -ex;' ; cat) | command ssh -q -t \
  (cat) | command ssh -q -t \
    -o "StrictHostKeyChecking=no" \
    -o "UserKnownHostsFile=/dev/null" \
    root@"${@}"
}

get_file_to_run() {
  local FILE="${1}"
  echo "set -x"
  if [[ "${IS_DRY_RUN:-}" == 1 ]]; then
    cat "${FILE}" >&2
  else
    cat "${FILE}"
  fi
  echo "echo PERTURBANCE COMPLETE ON \$(hostname) AT \$(date)"
  echo
}

run_file_on_host() {
  local FILE="${1}"
  local HOST="${2}"
  (
    set -x
    get_file_to_run "${FILE}"
  ) | run_ssh ${HOST}
  hr
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

get_connection_string() {
  echo "root@${MASTER_HOST}"
}

handle_arguments() {
  [[ $# = 0 && "${EXPECTED_NUM_ARGUMENTS}" -gt 0 ]] && usage

  parse_arguments "$@"
  validate_arguments "$@"

  if [[ "${IS_AUTOPOPULATE:-}" == 1 ]]; then
    if ! command doctl >/dev/null; then
      error "Please install doctl from https://github.com/digitalocean/doctl"
    fi

    if [[ "${DIGITALOCEAN_ACCESS_TOKEN:-}" == "" ]]; then
      warning "Please export DIGITALOCEAN_ACCESS_TOKEN. For example:"
      error "export DIGITALOCEAN_ACCESS_TOKEN=xxx"
    fi
  fi
}

parse_arguments() {
  while [ $# -gt 0 ]; do
    case $1 in
      -h | --help) usage ;;
      --dry-run)
        IS_DRY_RUN=1
        ;;
      --debug)
        DEBUG=1
        set -xe
        ;;
      --auto-populate)
        shift
        not_empty_or_usage "${1:-}"
        IS_AUTOPOPULATE=1
        AUTOPOPULATE_REGEX="${1}"
        ;;
      --test)
        IS_TEST_ONLY=1
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
      -s | --slaves)
        shift
        not_empty_or_usage "${1:-}"
        SLAVE_HOSTS="${1}"
        ;;
      --add-keys)
        shift
        not_empty_or_usage "${1:-}"
        GITHUB_KEY_USERS="${1}"
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

get_cluster_slaves() {
  local IPS
  IPS=$(get_node_ips)
  echo "${IPS}" | tail -n +2 | tr '\n' ','
}

validate_arguments() {
  [[ "${#ARGUMENTS[@]}" -gt 1 ]] && error "Only one scenario accepted"
  [[ "${#ARGUMENTS[@]}" -lt 1 && "${IS_TEST_ONLY}" != 1 ]] && error "Scenario required"

  if [[ "${IS_AUTOPOPULATE:-}" == 1 ]]; then
    MASTER_HOST=$(
      set -e
      get_cluster_master
    )
    SLAVE_HOSTS=$(
      set -e
      get_cluster_slaves
    )
  else
    [[ ${MASTER_HOST:-} ]] || error "--master must be an IP or hostname, or comma-delimited list"
    [[ ${SLAVE_HOSTS:-} ]] || error "--slaves must be an IP or hostname, or comma-delimited list"
  fi

  SCENARIO="${ARGUMENTS[0]:-}" || true
}

# helper functions

usage() {
  [ "$*" ] && echo "${THIS_SCRIPT}: ${COLOUR_RED}$*${COLOUR_RESET}" && echo
  sed -n '/^##/,/^$/s/^## \{0,1\}//p' "${THIS_SCRIPT}" | sed "s/%SCRIPT_NAME%/$(basename "${THIS_SCRIPT}")/g"
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

error_env_var() {
  error "${1} environment variable required"
}

log_message_prefix() {
  local TIMESTAMP="[$(date +'%Y-%m-%dT%H:%M:%S%z')]"
  local THIS_SCRIPT_SHORT=${THIS_SCRIPT/$DIR/.}
  tput bold 2>/dev/null
  echo -n "${TIMESTAMP} ${THIS_SCRIPT_SHORT}: "
}

is_empty() {
  [[ -z ${1-} ]] && return 0 || return 1
}

not_empty_or_usage() {
  is_empty "${1-}" && usage "Non-empty value required" || return 0
}

check_number_of_expected_arguments() {
  [[ "${EXPECTED_NUM_ARGUMENTS}" != "${#ARGUMENTS[@]}" ]] && {
    ARGUMENTS_STRING="argument"
    [[ "${EXPECTED_NUM_ARGUMENTS}" -gt 1 ]] && ARGUMENTS_STRING="${ARGUMENTS_STRING}"s
    usage "${EXPECTED_NUM_ARGUMENTS} ${ARGUMENTS_STRING} expected, ${#ARGUMENTS[@]} found"
  }
  return 0
}

hr() {
  printf '=%.0s' $(seq $(tput cols))
  echo
}

wait_safe() {
  local PIDS="${1}"
  for JOB in ${PIDS}; do
    wait "${JOB}"
  done
}

export CLICOLOR=1
export TERM="xterm-color"
export COLOUR_BLACK=$(tput setaf 0 :-"" 2>/dev/null)
export COLOUR_RED=$(tput setaf 1 :-"" 2>/dev/null)
export COLOUR_GREEN=$(tput setaf 2 :-"" 2>/dev/null)
export COLOUR_YELLOW=$(tput setaf 3 :-"" 2>/dev/null)
export COLOUR_BLUE=$(tput setaf 4 :-"" 2>/dev/null)
export COLOUR_MAGENTA=$(tput setaf 5 :-"" 2>/dev/null)
export COLOUR_CYAN=$(tput setaf 6 :-"" 2>/dev/null)
export COLOUR_WHITE=$(tput setaf 7 :-"" 2>/dev/null)
export COLOUR_RESET=$(tput sgr0 :-"" 2>/dev/null)

main "$@"
