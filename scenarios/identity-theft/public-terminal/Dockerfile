# syntax=docker/dockerfile:1

FROM controlplaneoffsec/kubectl

WORKDIR /
RUN apt update && apt install -y curl jq sudo vim ncat
RUN useradd -ms /bin/bash elit3hex0r
RUN sed -i '/# User alias specification/a elit3hex0r ALL=(ALL) NOPASSWD:ALL' /etc/sudoers
USER elit3hex0r

COPY --chown=elit3hex0r:elit3hex0r hacker-notes.md /home/elit3hex0r/hacker-notes.md