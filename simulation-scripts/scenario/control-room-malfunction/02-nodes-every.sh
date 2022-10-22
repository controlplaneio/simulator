#!/bin/bash

set -Eeuo pipefail

sed -i -e 's/disable_apparmor.*/disable_apparmor = true/g' /etc/containerd/config.toml
systemctl restart containerd
