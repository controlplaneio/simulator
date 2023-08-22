# syntax=docker/dockerfile:1

FROM controlplaneoffsec/kubectl

WORKDIR /
RUN apt update && apt install -y curl jq
RUN useradd -ms /bin/bash swashbyter
USER swashbyter

# PUT USER IN SUDO GROUP

COPY --chown=swashbyter:swashbyter diary.md /home/swashbyter/diary.md
COPY --chown=swashbyter:swashbyter treasure-map-1 /home/swashbyter/treasure-map-1
COPY --chown=swashbyter:swashbyter treasure-map-7 /home/swashbyter/treasure-map-7