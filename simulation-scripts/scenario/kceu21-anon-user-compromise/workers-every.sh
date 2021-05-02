#!/bin/bash

set -Eeuxo pipefail

BASE_IMAGE="debian:latest"
IMAGE="registry.41RG4P.controlplane.io/excalibur:0x5F3759DF"
IMAGE_DECOY_1="registry.41RG4P.controlplane.io/merlin:0x5F3759DF"
IMAGE_DECOY_2="registry.41RG4P.controlplane.io/guinevere:0x5F3759DF"

cleanup() {
  rm -rf "${DIR}" || true
}

main() {
  DIR=$(mktemp -d)
  trap cleanup EXIT

  cd "${DIR}"

  prep_file
  build_image
  drop_flag_in_image
}

prep_file() {

  cat <<EOF >.bash_history
#.._# :(){ :|:& };: #_..#
# Y29udHJvbC1wbGFuZS5pbyBhcmUgaGlyaW5nCg==
#
mkdir -p /usr/share/.never/gonna/{give/you/up,let/you/down}
mkdir -p /usr/share/.never/gonna/{make/you/cry,run/around/and/desert/you}
mkdir -p /usr/share/.never/gonna/{say/goodbye,tell/a/lie/and/hurt/you}
#
# pirates can be foolhardy // Hλ$ħ𝔍Ⱥ¢k / securi.fyi
EOF

  cat .bash_history
}

build_image() {
  cat <<EOF | docker build . -f - --tag "${IMAGE}"
FROM ${BASE_IMAGE}
COPY .bash_history /root/.bash_history
RUN bash /root/.bash_history
CMD sleep "\$((6000 + (RANDOM % 300)))"
EOF
}

drop_flag_in_image() {
  if [[ -f cidfile ]]; then
    rm -f cidfile
  fi
  cat <<'EOF' | docker run -i --cidfile=cidfile "${IMAGE}" bash /dev/stdin
DIR=$(find /usr/share/.never/gonna | sort -R | head -n 1) && mkdir -p ${DIR}/\, && echo 'flag_ctf{aea91cbef9d2caaa}' > $_/\,
EOF


  docker commit "$(cat cidfile)" "${IMAGE}"
  docker tag "${BASE_IMAGE}" "${IMAGE_DECOY_1}"
  docker tag "${BASE_IMAGE}" "${IMAGE_DECOY_2}"


  docker run -t "${IMAGE}" find /usr/share || true

  # test with:
  # kubectl  run -it xxx-${RANDOM} --image=debian:latest bash --overrides '{"spec":{"imagePullPolicy": "Never"}}'
}

#main |& tee /tmp/worker-log
main
