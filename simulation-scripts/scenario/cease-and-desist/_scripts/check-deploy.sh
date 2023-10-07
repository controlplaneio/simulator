#!/bin/bash

set -Eeuo pipefail

LICENSE="valid"
LICENSE_KEY="2fc593b894ef1402987d2595487d9763"

POD="{\"apiVersion\":\"v1\",\"kind\":\"Pod\",\"metadata\":{\"annotations\":{},\"name\":\"tokungfu-server\",\"namespace\":\"production\"},\"spec\":{\"containers\":[{\"env\":[{\"name\":\"FLAG\",\"valueFrom\":{\"secretKeyRef\":{\"key\":\"flag\",\"name\":\"flag\"}}}],\"image\":\"ttl.sh/wakeward-321sa5-4324ff-pr:12h\",\"name\":\"tokungfu-server\",\"ports\":[{\"containerPort\":8080,\"name\":\"http\"}]}],\"restartPolicy\":\"Always\"}}"

function LICENSE_CHECK {
  (kubectl get pods rkls -n licensing -ojson | jq -r '.metadata.labels.license')
}

function LICENSE_KEY_CHECK {
  (kubectl get pods rkls -n licensing -ojson | jq -r '.metadata.labels.license_key')
}
# Waiting for licensing server to be labelled
while [[ $(LICENSE_CHECK) != "${LICENSE}" && $(LICENSE_KEY_CHECK) != "${LICENSE_KEY}" ]];
do
sleep 30
done
(echo "$POD" | kubectl create -f -)