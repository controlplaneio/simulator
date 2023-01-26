# syntax=docker/dockerfile:1

FROM controlplaneoffsec/kubectl

WORKDIR /
RUN useradd -ms /bin/bash sre
USER sre
COPY --chown=sre:sre chat-archive.enc /home/sre/chat-archive.enc
COPY --chown=sre:sre secure-archive-chat.sh /home/sre/secure-archive-chat.sh
RUN chmod +x /home/sre/secure-archive-chat.sh