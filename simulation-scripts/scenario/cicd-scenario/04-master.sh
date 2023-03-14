#!/usr/bin/env bash

set -Eeuo pipefail
shopt -s expand_aliases

# Configure Gitea
curl -sSL https://dl.gitea.com/tea/0.9.2/tea-0.9.2-linux-amd64 -o tea
install -s tea /usr/local/bin

GITEA_SERVER_USER="ctf_admin"
GITEA_SERVER_PASSWORD="ahXeehohsoo2suej4tee0ol5xeeteM1w"
BASEURL="http://localhost:30080"

# shellcheck disable=SC2139
alias curl="curl -u \"$GITEA_SERVER_USER:$GITEA_SERVER_PASSWORD\" -H 'Content-type: application/json'"

#CSRF=$(curl -sS -c cookie.jar "$BASEURL" --retry 3 --retry-connrefused --retry-delay 5 | awk -F' ' '/csrfToken/ {print $2}' | tr -d "',")
#curl -sSL -b cookie.jar -c cookie.jar -XPOST "$BASEURL/user/login" -d "user_name=$GITEA_SERVER_USER&password=$GITEA_SERVER_PASSWORD&_csrf=$CSRF"
#CSRF=$(awk '/_csrf/ {print $7}' cookie.jar)
#RESP=$(curl -sSL -b cookie.jar -c cookie.jar -XPOST "$BASEURL/user/settings/applications" -d "_csrf=$CSRF&name=cli&scope=repo&scope=admin%3Aorg&scope=admin%3Apublic_key&scope=admin%3Arepo_hook&scope=admin%3Aorg_hook&scope=notification&scope=user&scope=delete_repo&scope=package&scope=admin%3Aapplication&scope=sudo" -v 2>&1)
#TOKEN=$(echo "$RESP" | grep -F "macaron_flash" | perl -0777 -pe 's/.*info%3D(.*)%26success%3DYour%2Bnew%2Btoken.*/\1/msg')

set -x
JSON=$(curl "$BASEURL/api/v1/users/ctf_admin/tokens" -XPOST --retry 3 --retry-connrefused --retry-delay 5 -d "{\"name\":\"cli_token\",\"scopes\":[\"all\"]}")
TOKEN=$(echo "$JSON" | jq .sha1 -r)

if [[ "$TOKEN" == "null" ]]; then
    echo "Failed to get cli TOKEN"
    exit 1
fi
tea login add -n ctf -u http://localhost:30080 -t "$TOKEN"
tea login default ctf

# Create developer user
USER="developer"
DATA="{\"email\":\"$USER@localhost.ctf\",\"username\":\"$USER\",\"password\":\"$(uuidgen)\"}"
JSON=$(curl "$BASEURL/api/v1/admin/users" -XPOST -d "$DATA")
SUCCESS=$(echo "$JSON" | jq .active -r)
if [[ "$SUCCESS" != "true" ]]; then
    echo "Creating user $USER failed"
    exit 1
fi

# Set Admin user to private
DATA="{\"login_name\":\"$GITEA_SERVER_USER\",\"visibility\":\"private\"}"
JSON=$(curl "$BASEURL/api/v1/admin/users/$GITEA_SERVER_USER" -XPATCH -d "$DATA")
SUCCESS=$(echo "$JSON" | jq .visibility -r)
if [[ "$SUCCESS" != "private" ]]; then
    echo "Editing user $USER failed"
    exit 1
fi

# Create org structure
ORG="supersecureorg"
tea org create "$ORG"
for NAME in developers devops reviewers; do
    DATA="{\"name\":\"$NAME\",\"permission\":\"write\",\"can_create_org_repo\":false,\"units\":[\"repo.code\",\"repo.pulls\"],\"includes_all_repositories\": true}"
    JSON=$(curl "$BASEURL/api/v1/orgs/$ORG/teams" -XPOST -d "$DATA")
    SUCCESS=$(echo "$JSON" | jq .id -r)
    if [[ "$SUCCESS" == "null" ]]; then
        echo "Creating team $NAME failed"
        exit 1
    fi
done

# Add reviewers
ID=$(curl "$BASEURL/api/v1/orgs/$ORG/teams/search?q=reviewers" | jq '.data[0].id' -r )
for r in developer1 dev2; do
    curl "$BASEURL/api/v1/team/$ID/members/$r" -XPUT
done
# Create repo
REPO="production-image-build"
tea repo create --owner "$ORG" --name "$REPO" --init

# Set repo features
# Requires patch
DATA='{"has_actions":true,"has_wiki":false,"has_releases":false,"has_projects":false}'
JSON=$(curl "$BASEURL/api/v1/repos/$ORG/$REPO" -XPATCH -d "$DATA")
SUCCESS=$(echo "$JSON" | jq .has_actions -r)
if [[ "$SUCCESS" != "true" ]]; then
    echo "Editing repo $REPO failed"
    exit 1
fi

# create branch protection
DATA='{"rule_name":"main","required_approvals":99,"enable_approvals_whitelist":true,"approvals_whitelist_teams":["reviewers"]}'
JSON=$(curl "$BASEURL/api/v1/repos/$ORG/$REPO/branch_protections" -XPOST -d "$DATA")
SUCCESS=$(echo "$JSON" | jq .created_at -r)
if [[ "$SUCCESS" == "null" ]]; then
    echo "Adding branch protection for $ORG/$REPO failed"
    exit 1
fi
