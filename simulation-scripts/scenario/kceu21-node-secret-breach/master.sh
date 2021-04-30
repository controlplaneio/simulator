#!/bin/bash

# creates squawk container that's sharing its SA token on port 80, RBAC is cluster admin. Get the RBAC, start a priv pod to escalate to the host, find the flag on the host

set -Eeuxo pipefail

NS_PRIMARY="squawk"
SA_USER="hashjack"
USER_ROLE="hashjack-role"
SA_SQUAWK="indomitable-sasquawksh"
SQUAWK_CLUSTER_ROLE="dangerous-cluster-role"


# ===

touch /tmp/.done
