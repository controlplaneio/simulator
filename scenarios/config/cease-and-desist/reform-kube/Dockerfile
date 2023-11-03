FROM golang:1.20-bullseye AS build

WORKDIR /

# Go Webapp
COPY go.mod .
COPY *.go .

RUN go build -o reform-kube-licensing-server -ldflags='-w -s' main.go

##
## Deploy
##
FROM controlplaneoffsec/kubectl

RUN apt update && apt install -y --no-install-recommends \
    ca-certificates \
    && apt clean all

WORKDIR /

RUN useradd -ms /bin/bash licenser
USER licenser

COPY --from=build --chown=licenser:licenser /reform-kube-licensing-server /reform-kube/reform-kube-licensing-server