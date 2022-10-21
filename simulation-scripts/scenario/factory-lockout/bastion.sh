#!/bin/bash

set -Eeuo pipefail

eval $(grep -E 'export (NODE|MASTER)_IP_ADDR' ~ubuntu/.bash_login_script)

cat <<-EOF >/etc/systemd/system/socat.service
[Unit]
Description=Socat Forward

[Service]
Type=simple
ExecStartPre=-killall -9 -q socat
ExecStart=socat TCP4-LISTEN:8080,reuseaddr,fork TCP4:$MASTER_IP_ADDRESSES:30080
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

systemctl enable --now socat
