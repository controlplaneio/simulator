FROM debian:buster-slim

RUN apt-get update && apt-get install -y golang build-essential git ca-certificates

RUN mkdir -p /go/src/github.com/controlplaneio/simulator-standalone/cli
WORKDIR /go/src/github.com/controlplaneio/simulator-standalone/cli

RUN addgroup --quiet build
RUN adduser --quiet --disabled-password --gecos "" --ingroup build build
RUN chown -R build:build ./

ENV GO111MODULE=on

USER build

CMD ["make"]

