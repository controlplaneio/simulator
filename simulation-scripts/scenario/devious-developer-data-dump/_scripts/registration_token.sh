#!/usr/bin/env bash

set -Eeuo pipefail

USER="ctf_admin"
PASS="admin123"
BASEURL="localhost:8080"


#capture CSRF token
CSRF=$(curl -sS -c cookie.jar "$BASEURL" | awk -F' ' '/csrfToken/ {print $2}' | tr -d "',")
#echo "Got CSRF: $CSRF"

curl -sSL -b cookie.jar -c cookie.jar -XPOST "$BASEURL/user/login" -d "user_name=$USER&password=$PASS&_csrf=$CSRF"

TOKEN=$(curl -sSL -b cookie.jar "$BASEURL/admin/runners" | grep -A3 'Registration Token' | awk -F'"' '/value/ {print $4}')
echo "$TOKEN"
