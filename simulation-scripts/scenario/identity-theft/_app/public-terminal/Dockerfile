# syntax=docker/dockerfile:1

FROM controlplaneoffsec/kubectl

WORKDIR /
RUN apt update && apt install -y curl jq sudo vim ncat
RUN useradd -ms /bin/bash elitehexor
RUN sed -i '/# User alias specification/a elitehexor ALL=(ALL) NOPASSWD:ALL' /etc/sudoers
USER elitehexor