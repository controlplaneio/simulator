# syntax=docker/dockerfile:1

FROM controlplaneoffsec/kubectl

WORKDIR /
RUN apt update && apt install -y curl jq sudo vim
RUN useradd -ms /bin/bash swashbyter
RUN sed -i '/# User alias specification/a swashbyter ALL=(ALL) NOPASSWD:ALL' /etc/sudoers
USER swashbyter

COPY --chown=swashbyter:swashbyter diary.md /home/swashbyter/diary.md
COPY --chown=swashbyter:swashbyter treasure-map-1 /home/swashbyter/treasure-map-1
COPY --chown=swashbyter:swashbyter treasure-map-7 /home/swashbyter/treasure-map-7