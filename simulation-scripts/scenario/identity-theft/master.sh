#!/bin/bash

set -Eeuxo pipefail

function GETDEXIP {
  (kubectl get pods -n dex -ojson | jq -r '.items[].status.podIP')
}

until [[ -z $(GETDEXIP) ]];
do
    sleep 5
done


function GETSSIP {
  (kubectl get pods secret-store -n private-services -ojson | jq -r '.status.podIP')
}

until [[ -z $(GETSSIP) ]];
do
    sleep 5
done

DEXIP=$(kubectl get pods -n dex -ojson | jq -r '.items[].status.podIP')
SSIP=$(kubectl get pods secret-store -n private-services -ojson | jq -r '.status.podIP')

SECRET_STORE="http://$SSIP:5050/api/v1/users"
DEX="http://$DEXIP:5556/dex/token"

USER1="{\"email\":\"admin@pod-checker.local\",\"firstName\":\"ultra\",\"lastName\":\"violet\",\"password\":\"the-keys-to-the-kingdom\", \"secret\": \"podcheckerauth\"}"
USER2="{\"email\":\"flag@pod-checker.local\",\"firstName\":\"ctf\",\"lastName\":\"flag\",\"password\":\"ctf_flag{BRAVO_BONUS_POINTS_FOR_CRACKING_THIS}\", \"secret\": \"ctf_flag{O_I_DC_WHAT_YOU_HAVE_DONE_THERE}\"}"
USER3="{\"email\":\"db@pod-checker.local\",\"firstName\":\"infra\",\"lastName\":\"red\",\"password\":\"administer-the-secret-store-247!\", \"secret\": \"dbsecretstoreauth\"}"
USER4="{\"email\":\"wakeward@pod-checker.local\",\"firstName\":\"kevin\",\"lastName\":\"ward\",\"password\":\"i-am-the-author-of-this!\", \"secret\": \"ihopeyouenjoyedthis!\"}"

USERS=("$USER1" "$USER2" "$USER3" "$USER4")

ID_TOKEN=$(curl -s -X POST "$DEX" \
  --data-urlencode "grant_type=password" \
  --data-urlencode "client_id=pod-checker" \
  --data-urlencode "client_secret=cG9kY2hlY2tlcmF1dGgK" \
  --data-urlencode "username=admin@pod-checker.local" \
  --data-urlencode "password=the-keys-to-the-kingdom" \
  --data-urlencode "scope=openid profile email" \
  | jq -r '.id_token')

# Create all users
for USER in "${USERS[@]}"
do
    curl -X POST "$SECRET_STORE" \
     -H "content-type: application/json" \
     -H "Authorization: Bearer $ID_TOKEN" \
     -d "$USER"
done