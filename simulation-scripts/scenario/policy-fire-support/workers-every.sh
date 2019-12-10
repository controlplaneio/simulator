#!/bin/bash

mkdir -pv /root/.ssh && ssh-keygen -C root -b 4096 -t rsa -N "" -f /root/.ssh/id_rsa
docker build -t jenkins:policy-fire-support - <<EOF

FROM jenkins

USER root
RUN apt update && apt install -y sudo
RUN echo 'ALL            ALL = (ALL) NOPASSWD: ALL' >> /etc/sudoers
USER jenkins
EOF
