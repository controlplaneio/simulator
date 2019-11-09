FROM debian:buster-slim

# We're using sh not bash at this point
# hadolint ignore=DL4006
RUN apt-get update                                                                \
    && DEBIAN_FRONTEND=noninteractive apt-get install  -y --no-install-recommends \
    curl                                                                          \
    software-properties-common                                                    \
    && curl -sL https://deb.nodesource.com/setup_13.x | bash -                    \
    && apt-get update                                                             \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends  \
    bash                                                                          \
    bash-completion                                                               \
    ca-certificates                                                               \
    gettext-base                                                                  \
    hping3                                                                        \
    man-db                                                                        \
    nmap                                                                          \
    nodejs                                                                        \
    openssh-client                                                                \
 && rm -rf /var/lib/apt/lists/*

RUN echo "PS1='\u@\h # '" > /root/.bashrc

## Install Goss
ENV GOSS_VERSION v0.3.7
RUN curl -sL https://github.com/aelsabbahy/goss/releases/download/${GOSS_VERSION}/goss-linux-amd64 \
         -o /usr/local/bin/goss                                                                    \
    && chmod +rx /usr/local/bin/goss

# Add login message
COPY ./attack/scripts/attack-motd /usr/local/bin/attack-motd
COPY ./attack/scripts/bashrc /etc/bash.bashrc

ENV NODE_ENV production
RUN mkdir /scenario-tools
WORKDIR /scenario-tools
COPY ./tools/scenario-tools .
RUN npm install --production \
    && npm prune             \
    && npm link



RUN mkdir /attack
WORKDIR /attack
COPY ./attack/ .

# This container is root intentionally

CMD ["./scripts/entrypoint.sh"]
