#!/bin/bash

set -Eeuo pipefail

## manual cleanup
rm -f /root/.bash_history
journalctl --vacuum-time 1s --quiet 2>/dev/null
