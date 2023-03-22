#!/usr/bin/env bash

set -Eeuo pipefail
shopt -s expand_aliases

GITEA_SERVER_USER="ctf_admin"
GITEA_SERVER_PASSWORD="ahXeehohsoo2suej4tee0ol5xeeteM1w"
BASEURL="http://localhost:8080"

# shellcheck disable=SC2139
alias curl="curl -sS -u \"$GITEA_SERVER_USER:$GITEA_SERVER_PASSWORD\" -H 'Content-type: application/json'"

USERS=$(curl "$BASEURL/api/v1/admin/users" | jq '.[].login' -r | tr '\n' ' ')

for USR in $USERS; do
    if [[ "$USR" == "$GITEA_SERVER_USER" ]]; then
        continue
    fi
    SUCCESS=$(curl "$BASEURL/api/v1/admin/users/$USR" -XDELETE -o /dev/null -w "%{http_code}" )
    if [[ "$SUCCESS" != "204" ]]; then
        echo "Deleting $USR failed with $SUCCESS"
    fi
done

ORG="rescue-drop"
REPO="production-image-build"
SUCCESS=$(curl "$BASEURL/api/v1/repos/$ORG/$REPO" -XDELETE -o /dev/null -w "%{http_code}" )
if [[ "$SUCCESS" != "204" ]]; then
    echo "Deleting $REPO failed with $SUCCESS"
fi

SUCCESS=$(curl "$BASEURL/api/v1/orgs/$ORG" -XDELETE -o /dev/null -w "%{http_code}" )
if [[ "$SUCCESS" != "204" ]]; then
    echo "Deleting $ORG failed with $SUCCESS"
fi

SUCCESS=$(curl "$BASEURL/api/v1/users/$GITEA_SERVER_USER/tokens/cli_token" -XDELETE -o /dev/null -w "%{http_code}" )
if [[ "$SUCCESS" != "204" ]]; then
    echo "Deleting user token failed with $SUCCESS"
fi
