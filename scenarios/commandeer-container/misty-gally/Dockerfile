# syntax=docker/dockerfile:1

FROM controlplaneoffsec/kubectl

WORKDIR /
RUN apt update && apt install -y curl jq
RUN useradd -ms /bin/bash stowaway
USER stowaway