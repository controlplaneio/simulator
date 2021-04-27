#!/bin/bash

set -Eeuxo pipefail

IMAGE="localhost:5000/excelsior:0x5F3759DF"

DIR=$(mktemp -d)
cd "${DIR}"
cleanup() {
  rm -rf "${DIR}"
}
trap cleanup EXIT

main() {
  prep_file
  build_image
  drop_flag_in_image
  host_image
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
}

build_image() {
  cat <<EOF | docker build . -f - --tag "${IMAGE}"
FROM debian
COPY .bash_history /root/.bash_history
RUN bash /root/.bash_history
CMD sleep "\$((6000 + (RANDOM % 300)))"
EOF
}

drop_flag_in_image() {
  cat <<'EOF' | docker run -i --cidfile=cidfile \
    "${IMAGE}" bash /dev/stdin

DIR=$(find /usr/share/.never/gonna | sort -R | head -n 1) && mkdir -p ${DIR}/\, && echo 'flag_ctf{aea91cbef9d2caaa}' > $_/\,
EOF

  docker commit "$(cat cidfile)" "${IMAGE}"
}


main
