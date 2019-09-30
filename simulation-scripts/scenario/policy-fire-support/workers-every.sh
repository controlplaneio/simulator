#!/bin/bash

mkdir -pv /root/.ssh && ssh-keygen -C root -b 4096 -t rsa -N "" -f /root/.ssh/id_rsa &> /dev/null
docker build -t jenkins:policy-fire-support - &> /dev/null <<EOF

FROM jenkins

USER root
RUN apt update && apt install -y sudo
RUN echo 'ALL            ALL = (ALL) NOPASSWD: ALL' >> /etc/sudoers
USER jenkins
EOF

kubectl apply -f - <<EOF
apiVersion: policy/v1beta1
kind: PodSecurityPolicy
metadata:
  name: admin
  annotations:
    seccomp.security.alpha.kubernetes.io/allowedProfileNames: '*'
spec:
  privileged: true
  allowPrivilegeEscalation: true
  allowedCapabilities:
  - '*'
  volumes:
  - '*'
  hostNetwork: true
  hostPorts:
  - min: 0
    max: 65535
  hostIPC: true
  hostPID: true
  runAsUser:
    rule: 'RunAsAny'
  seLinux:
    rule: 'RunAsAny'
  supplementalGroups:
    rule: 'RunAsAny'
  fsGroup:
    rule: 'RunAsAny'
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: admin-psp
rules:
- apiGroups: ['policy']
  resources: ['podsecuritypolicies']
  verbs:     ['use']
  resourceNames:
  - admin
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: admin-psp
  namespace: kube-system
roleRef:
  kind: ClusterRole
  name: admin-psp
  apiGroup: rbac.authorization.k8s.io
subjects:
  - kind: Group
  apiGroup: rbac.authorization.k8s.io
  name: system:serviceaccounts
EOF

sed -i '/    - --enable-admission-plugins=NodeRestriction/c\    - --enable-admission-plugins=NodeRestriction,PodSecurityPolicy' /etc/kubernetes/manifests/kube-apiserver.yaml
