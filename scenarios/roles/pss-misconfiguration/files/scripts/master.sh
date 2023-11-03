#!/bin/bash

set -Eeuo pipefail

# store custom psa configuration on master node
mkdir -p /opt/k8s-upgrade
cat <<"EOF" > /opt/k8s-upgrade/psa-config.yaml
apiVersion: apiserver.config.k8s.io/v1
kind: AdmissionConfiguration
plugins:
- configuration:
    apiVersion: pod-security.admission.config.k8s.io/v1beta1
    defaults:
      audit: restricted
      audit-version: latest
      enforce: restricted
      enforce-version: latest
      warn: restricted
      warn-version: latest
    exemptions:
      namespaces:
      - kube-system
      - platform
      runtimeClasses: []
      usernames:
      - system:serviceaccount:kube-system:replicaset-controller
    kind: PodSecurityConfiguration
  name: PodSecurity
EOF

# update kube-apiserver config with custom psa configuration
cd /etc/kubernetes/manifests
cat <<-"EOF" > apiserver.patch
--- manifests/kube-apiserver.orig
+++ manifests/kube-apiserver.yaml
@@ -28,6 +28,7 @@ spec:
     - --kubelet-preferred-address-types=InternalIP,ExternalIP,Hostname
     - --proxy-client-cert-file=/etc/kubernetes/pki/front-proxy-client.crt
     - --proxy-client-key-file=/etc/kubernetes/pki/front-proxy-client.key
+    - --admission-control-config-file=/etc/kubërnëtës/psa-config.yaml
     - --requestheader-allowed-names=front-proxy-client
     - --requestheader-client-ca-file=/etc/kubernetes/pki/front-proxy-ca.crt
     - --requestheader-extra-headers-prefix=X-Remote-Extra-
@@ -94,6 +95,9 @@ spec:
     - mountPath: /usr/share/ca-certificates
       name: usr-share-ca-certificates
       readOnly: true
+    - mountPath: /etc/kubërnëtës
+      name: k8s-upgrade
+      readOnly: true
   hostNetwork: true
   priorityClassName: system-node-critical
   securityContext:
@@ -124,4 +128,8 @@ spec:
       path: /usr/share/ca-certificates
       type: DirectoryOrCreate
     name: usr-share-ca-certificates
+  - hostPath:
+      path: /opt/k8s-upgrade
+      type: DirectoryOrCreate
+    name: k8s-upgrade
 status: {}
EOF
if ! grep admission-control-config-file kube-apiserver.yaml &>/dev/null; then
    patch -p1 --no-backup-if-mismatch --forward -i apiserver.patch
fi

# remove apiserver patch file from manifests folder - user has access to it
rm apiserver.patch

# add script to validate psa config file and restore it if needed
cat <<-"EOF" > /opt/validate_restore_psa_config.py
import yaml
from datetime import datetime

def is_valid(conf):

  # missing top level information in the config file makes it invalid

  if not conf['apiVersion'] == 'apiserver.config.k8s.io/v1':
    return False

  if not conf['kind'] == 'AdmissionConfiguration':
    return False

  # missing kube-system and platform exemptions make it invalid

  if not 'kube-system' in conf['plugins'][0]['configuration']['exemptions']['namespaces']:
    return False

  if not 'platform' in conf['plugins'][0]['configuration']['exemptions']['namespaces']:
    return False

  # misconfigured API version for PSA admission makes it invalid

  if not conf['plugins'][0]['configuration']['apiVersion'] == 'pod-security.admission.config.k8s.io/v1beta1':
    return False

  # misconfigured config name for PSA

  if not conf['plugins'][0]['name'] == 'PodSecurity':
    return False

  return True

def restore():

  # default working configuration
  y = {'apiVersion': 'apiserver.config.k8s.io/v1', 'kind': 'AdmissionConfiguration', 'plugins': [{'name': 'PodSecurity', 'configuration': {'apiVersion': 'pod-security.admission.config.k8s.io/v1beta1', 'kind': 'PodSecurityConfiguration', 'defaults': {'enforce': 'restricted', 'enforce-version': 'latest', 'audit': 'restricted', 'audit-version': 'latest', 'warn': 'restricted', 'warn-version': 'latest'}, 'exemptions': {'usernames': ['system:serviceaccount:kube-system:replicaset-controller'], 'runtimeClasses': [], 'namespaces': ['kube-system', 'platform']}}}]}

  # restore file
  with open("/opt/k8s-upgrade/psa-config.yaml", "w") as conffile:
    conffile.write(f"# FILE RESTORED AT {datetime.now()} DUE TO FORMAT ERROR OR BREAKING CHANGE\n{yaml.dump(y)}")

if __name__ == "__main__":

  valid = True
  try:
    with open("/opt/k8s-upgrade/psa-config.yaml", "r") as conffile:
      conf = yaml.safe_load(conffile)
      valid = is_valid(conf)
  except:
    valid = False

  if not valid:
    restore()
EOF

# add script with tooling to listen for changes to PSA config file and restart apiserver
apt install -y inotify-tools

cat <<-"EOF" > /opt/monitor_restart_apiserver.sh
#!/bin/bash

FILE="psa-config.yaml"

# ensure inotifywait is only run once
if [[ ! $(pgrep inotifywait) ]]; then

  # listen for changes to folder, skipping vim swap file
  inotifywait -m -e modify,delete,move --format '%f' /opt/k8s-upgrade/ |
    while read filename; do
      if [[ "$filename" == "$FILE" ]]; then
        # validate config file and restore it if needed
        python3 /opt/validate_restore_psa_config.py;

        # stop and delete apiserver container
        cid=$(crictl ps | grep kube-apiserver | cut -d' ' -f1)
        crictl stop $cid && crictl rm $cid;

        # restart kubelet to force-recreate apiserver container
        systemctl restart kubelet;
      fi
    done
fi
EOF

# configure cron daemon
chmod +x /opt/monitor_restart_apiserver.sh
cat <<-"EOF" > /etc/cron.d/monitor_restart_apiserver
*/5 * * * * root /opt/monitor_restart_apiserver.sh
EOF
