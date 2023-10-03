#!/bin/bash

SECRET_STORE="http://secret-service.private-services.svc.cluster.local:5050/api/v1/users"
DEX="http://dex.dex.svc.cluster.local:5556/dex/token"

if ID_TOKEN=""; then
  ID_TOKEN=$(curl -s -X POST "$DEX" \
  --data-urlencode "grant_type=password" \
  --data-urlencode "client_id=pod-checker" \
  --data-urlencode "client_secret=cG9kY2hlY2tlcmF1dGgK" \
  --data-urlencode "username=admin@pod-checker.local" \
  --data-urlencode "password=the-keys-to-the-kingdom" \
  --data-urlencode "scope=openid profile email" \
  | jq -r '.id_token')
fi

curl "${SECRET_STORE}" \
     -H "Authorization: Bearer ${ID_TOKEN}"