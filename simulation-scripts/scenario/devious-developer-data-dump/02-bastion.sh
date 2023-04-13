#!/bin/bash

set -Eeuo pipefail

eval $(grep -E 'export (NODE|MASTER)_IP_ADDR' /opt/bash_login_script)

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

## Configure Gitea
#CSRF=$(curl -sS -c cookie.jar "$BASEURL" --retry 3 --retry-connrefused --retry-delay 5 | awk -F' ' '/csrfToken/ {print $2}' | tr -d "',")
#curl -sSL -b cookie.jar -c cookie.jar -XPOST "$BASEURL/user/login" -d "user_name=$USER&password=$PASS&_csrf=$CSRF"
#TOKEN=$(curl -sSL -b cookie.jar "$BASEURL/admin/runners" | grep -A3 'Registration Token' | awk -F'"' '/value/ {print $4}')
