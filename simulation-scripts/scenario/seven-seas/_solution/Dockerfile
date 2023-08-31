# syntax=docker/dockerfile:1

FROM controlplaneoffsec/kubectl

WORKDIR /
RUN apt update && apt install -y curl jq sudo vim openssh-client
RUN useradd -ms /bin/bash swashbyter
RUN sed -i '/# User alias specification/a swashbyter ALL=(ALL) NOPASSWD:ALL' /etc/sudoers
USER swashbyter