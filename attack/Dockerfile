FROM debian:buster-slim

RUN apt-get update                                                               \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
    bash                                                                         \
    ca-certificates                                                              \
    openssh-client                                                               \
    hping3                                                                       \
    nmap                                                                         \
    curl                                                                         \
 && rm -rf /var/lib/apt/lists/*

## Install Goss
ENV GOSS_VERSION v0.3.7
RUN curl -sL https://github.com/aelsabbahy/goss/releases/download/${GOSS_VERSION}/goss-linux-amd64 \
         -o /usr/local/bin/goss                                                                    \
    && chmod +rx /usr/local/bin/goss

# Add login message
COPY ./scripts/attack-motd /usr/local/bin/attack-motd
COPY ./scripts/bashrc /etc/bash.bashrc

RUN mkdir /attack
WORKDIR /attack
COPY . .

# This container is root intentionally

CMD ["./scripts/entrypoint.sh"]
