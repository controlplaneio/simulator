FROM docker.io/library/node:18-bullseye AS scenario-tools

ENV NODE_ENV production

WORKDIR /app/scenario-tools
COPY --chown=1000 ./tools/scenario-tools/ .

RUN npm install --production \
 && npm prune                \
 && npm link

# Attack
FROM debian:bullseye-slim

# We're using sh not bash at this point
# hadolint ignore=DL4006
RUN apt-get update                                                                \
 && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends     \
    bash                                                                          \
    bash-completion                                                               \
    binutils                                                                      \
    ca-certificates                                                               \
    curl                                                                          \
    gettext-base                                                                  \
    hping3                                                                        \
    iproute2                                                                      \
    jq                                                                            \
    man-db                                                                        \
    nmap                                                                          \
    openssh-client                                                                \
 && rm -rf /var/lib/apt/lists/*

ARG YQ_VERSION=3.4.1
ARG GOSS_VERSION=v0.3.21

## Install Goss
RUN curl -sSOL https://github.com/aelsabbahy/goss/releases/download/${GOSS_VERSION}/goss-linux-amd64 \
 && install -Dm755 -s goss-linux-amd64 /usr/local/bin/goss

# Install YQ
RUN curl -sL https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64 \
    -o /usr/local/bin/yq                                                                  \
 && install -Dm755 -s /usr/local/bin/yq

RUN mkdir /scenario-tools
WORKDIR /scenario-tools
COPY --chown=1000 --from=scenario-tools /app/scenario-tools/*.json /scenario-tools/
COPY --chown=1000 --from=scenario-tools /app/scenario-tools/bin /scenario-tools/bin
COPY --chown=1000 --from=scenario-tools /app/scenario-tools/lib /scenario-tools/lib
COPY --chown=1000 --from=scenario-tools /app/scenario-tools/node_modules /scenario-tools/node_modules

RUN mkdir /attack
WORKDIR /attack
COPY ./attack/ .

# Add login message
COPY ./attack/scripts/motd /usr/local/bin/motd
COPY ./attack/scripts/bashrc /etc/bash.bashrc

COPY --chown=1000 ./launch-files/bash_aliases /root/.bash_aliases
COPY --chown=1000 ./launch-files/inputrc /root/.inputrc

RUN printf "                    \n\
if [ -f ~/.bash_aliases ]; then \n\
    . ~/.bash_aliases           \n\
fi                              \n\
" >> /root/.bashrc

# This container is root intentionally

CMD ["./scripts/entrypoint.sh"]
