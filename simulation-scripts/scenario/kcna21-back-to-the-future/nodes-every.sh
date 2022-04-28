#!/bin/bash

set -Eeuo pipefail

# enhanced cleanup
rm /var/log/cloud-init.log || true
rm /var/log/cloud-init-output.log || true
