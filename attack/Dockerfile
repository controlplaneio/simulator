#FROM kalilinux/kali-linux-docker:latest AS attack
FROM debian:buster-slim

RUN apt-get update                                                               \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
    bash                                                                         \
    hping3                                                                       \
    nmap                                                                         \
    curl                                                                         \
 && rm -rf /var/lib/apt/lists/*

# Add login message
COPY ./scripts/attack-motd /usr/local/bin/attack-motd
RUN echo '[ ! -z "$TERM" ] && attack-motd' >> /etc/bash.bashrc

# This container is root intentionally

CMD ["/bin/bash"]
