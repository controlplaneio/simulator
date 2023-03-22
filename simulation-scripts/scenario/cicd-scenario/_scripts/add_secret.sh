#!/usr/bin/env bash

set -Eeuo pipefail

GITEA_SERVER_USER="ctf_admin"
GITEA_SERVER_PASSWORD="ahXeehohsoo2suej4tee0ol5xeeteM1w"
BASEURL="localhost:8080"
ORG="rescue-drop"
REPO="production-image-build"
SECRET_NAME="DEPLOY_USER"
SECRET_VALUE="storeimage"


CSRF=$(curl -sS -c cookie.jar "$BASEURL" | awk -F' ' '/csrfToken/ {print $2}' | tr -d "',")
curl -sSL -b cookie.jar -c cookie.jar -XPOST "$BASEURL/user/login" -d "user_name=$GITEA_SERVER_USER&password=$GITEA_SERVER_PASSWORD&_csrf=$CSRF"
echo "Got CSRF: $CSRF"
CSRF=$(curl -sS -c cookie.jar -b cookie.jar "$BASEURL/$ORG/$REPO/settings/secrets" | awk -F' ' '/csrfToken/ {print $2}' | tr -d "',")
echo "Got CSRF: $CSRF"
RESP=$(curl -sSL -b cookie.jar -c cookie.jar -XPOST "$BASEURL/$ORG/$REPO/settings/secrets" -d "_csrf=$CSRF&title=${SECRET_NAME}&content=${SECRET_VALUE}" -v 2>&1)
SUCCESS=$(echo "$RESP" | grep -F "macaron_flash" | perl -0777 -pe 's/.*flash=success%3DThe%2Bsecret(.*)has%2Bbeen%2Badded.*/\1/msg')
if [[ "$SUCCESS" == "" ]]; then
    echo "Adding CI secret ${SECRET_NAME} failed"
    exit 1
fi
echo "$SUCCESS"
