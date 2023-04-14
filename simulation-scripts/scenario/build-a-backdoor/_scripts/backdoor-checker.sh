#!/bin/bash

set -Eeuo pipefail

#!/bin/bash

function SERVICE_CHECK {
    (kubectl get svc mgmt-service -n ii-prod -ojson | jq '.spec.ports[] | select(.targetPort==5724)' | jq .targetPort)
}

function NETPOL_CHECK {
    (kubectl get netpol ii-mgmt-np -n ii-prod -ojson | jq '.spec.ingress[].ports[] | select(.port==5724)' | jq .port)
}

function ING_IP {
    (kubectl get ing -n ii-prod -ojson | jq -r .items[].status.loadBalancer.ingress[].ip)
}

function NODE_PORT {
    (kubectl get svc ingress-nginx-controller -n ii-prod -ojson | jq '.spec.ports[] | select(.targetPort==5724)' | jq .nodePort)
}

SERVICE_MSG="Vulnerable Operational Port is Unavailable"

until [[ $(SERVICE_CHECK) == "5724" && $(NETPOL_CHECK) == "5724" ]];
do
    echo "$SERVICE_MSG"
    sleep 30
done

until [[ $(curl -L -u system:3f89a11b3fee0f6918c731702edb7244 http://$(ING_IP):$(NODE_PORT)/admin) == "SYSTEM OPERATIONS" ]];
do
   echo "Waiting for backdoor"
   sleep 10
done

# DEPLOY SECRET
(kubectl create secret generic flag --from-literal=flag=flag_ctf{BACKDOOR_OPEN_HASHJACKER_DEPLOYED} --namespace=ii-prod)

# DEPLOY CONTAINER
