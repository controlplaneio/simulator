cat <<EOF | run_ssh "$(get_master)"
set -euxo pipefail

chmod -x /etc/update-motd.d/*

false || true
echo TESTED GOOD
EOF

