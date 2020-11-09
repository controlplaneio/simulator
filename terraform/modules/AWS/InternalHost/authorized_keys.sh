#!/bin/bash

set -x

authorized_keys_path="/root/.ssh/authorized_keys"

mkdir -p "/root/.ssh"
(
  for USER in "${@:-}"; do
    curl --fail --max-time 10 "https://github.com/${USER}.keys" || {
      sleep 2
      curl --fail --max-time 10 "https://github.com/${USER}.keys"
    } || true
  done
) >>"${authorized_keys_path}"

# dedupe
awk '!x[$0]++' "${authorized_keys_path}" \
  >"${authorized_keys_path}.new" \
  && mv "${authorized_keys_path}.new" "${authorized_keys_path}"

