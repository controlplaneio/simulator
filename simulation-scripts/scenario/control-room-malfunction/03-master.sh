#!/bin/bash

set -Eeuo pipefail

kubectl rollout restart -n control-room deployment control-panel
