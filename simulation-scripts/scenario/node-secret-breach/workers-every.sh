#!/bin/bash

set -Eeuxo pipefail

mkdir -pv /node-secrets
echo "node-secret=flag_ctf{7e778526e78d7fdf}" > /node-secrets/secret.txt
