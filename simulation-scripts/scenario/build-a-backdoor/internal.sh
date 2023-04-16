#!/bin/bash

# store custom psa configuration on master node
cat <<"EOF" > /root/hashjacker.yaml
apiVersion: v1
kind: Pod
metadata:
  name: ii-management-services
  namespace: ii-pord
  labels:
    app: ii
spec:
  containers:
  - name: hashjacker
    image: controlplaneoffsec/build-a-backdoor:hashjacker
    ports:
    - containerPort: 8080
    env:
    - name: P
      value: "ha5hjackk3y4d3crypt1ngs3cret1mg!"
  restartPolicy: Always
EOF

# Checks whether backdoor is accessible externally
cat <<- "EOF" > /root/backdoor-checker.sh
#!/bin/bash

set -Eeuo pipefail

function ING_IP {
    (kubectl get ing -n ii-prod -ojson | jq -r .items[].status.loadBalancer.ingress[].ip)
}

function NODE_PORT {
    (kubectl get svc ingress-nginx-controller -n ingress-nginx -ojson | jq '.spec.ports[] | select(.port==80)' | jq .nodePort)
}

function ING_PATH {
    (kubectl get ing ingress-ii-mgmt -n ii-prod -ojson | jq -r '.spec.rules[].http.paths[] | select(.backend.service.port.number==5724)' | jq -r .path)
}

function SECRET {
    (kubectl get secret flag -n ii-prod)
}

function POD {
    (kubectl get pods ii-management-services -n ii-pord)
}

SERVICE_MSG="Vulnerable Operational Port is Unavailable"

if [[ $(curl -L -u system:3f89a11b3fee0f6918c731702edb7244 http://$(ING_IP):$(NODE_PORT)/$(ING_PATH)) == "SYSTEM OPERATIONS" ]];
then
    if [[ $(SECRET) == 'Error from server (NotFound): secrets "flag" not found' && $(POD) == 'Error from server (NotFound): pods "ii-management-services" not found' ]];
    then
        echo "Secret and Pod exist"
        exit 0
    else
      (kubectl create secret generic flag --from-literal=flag=flag_ctf{BACKDOOR_OPEN_HASHJACKER_DEPLOYED_IN_PORD} --namespace=ii-prod)
      (kubectl apply -f /root/hashjacker.yaml)
    fi
else
    echo "$SERVICE_MSG"
fi
EOF

# configure cron daemon
chmod +x /root/backdoor-checker.sh
cat <<-"EOF" > /etc/cron.d/backdoor-checker
*/1 * * * * root /root/backdoor-checker.sh
EOF