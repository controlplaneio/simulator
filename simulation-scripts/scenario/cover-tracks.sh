#!/bin/bash

set -Eeuo pipefail

# remove cmd history
rm -f /root/.bash_history
rm -f /home/ubuntu/.bash_history

# clean syslogs
find /var/log/journal /run/log/ -name "*.journal" -delete 2>/dev/null
echo "" > /var/log/syslog

# rm cloud init logs
rm -f /var/log/cloud-init.log
rm -f /var/log/cloud-init-output.log
rm -rf /var/lib/cloud/
rm -rf /run/cloud-init/
rm -f /tmp/cloud-init.log

# TODO(JCP): clean kube events

# clean user logins
echo "" > /var/log/wtmp
echo "" > /var/log/lastlog
echo "" > /var/log/auth.log
