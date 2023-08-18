# syntax=docker/dockerfile:1

FROM controlplaneoffsec/kubectl

WORKDIR /
RUN apt update
RUN useradd -ms /bin/bash swashbyter
USER swashbyter

COPY --chown=swashbyter:swashbyter diary.md /home/swashbyter/diary.md
COPY --chown=swashbyter:swashbyter treasure_map_1 /home/swashbyter/treasure_map_1
COPY --chown=swashbyter:swashbyter treasure_map_7 /home/swashbyter/treasure_map_7