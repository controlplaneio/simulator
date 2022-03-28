FROM docker.io/library/debian:11

RUN apt update && apt install -y curl \
 && apt clean

RUN curl -LO "https://dl.k8s.io/release/v1.20.11/bin/linux/amd64/kubectl" \
 && install kubectl /usr/local/bin/
