# syntax=docker/dockerfile:1

FROM controlplaneoffsec/kubectl

WORKDIR /root
COPY todolist.md /root/todolist.md

RUN apt update && apt install -y curl jq
