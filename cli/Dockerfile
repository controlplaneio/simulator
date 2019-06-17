FROM golang:1.12 AS builder

RUN mkdir -p /go/src/github.com/controlplaneio/simulator-standalone/cli

WORKDIR /go/src/github.com/controlplaneio/simulator-standalone/cli

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

COPY . .

RUN dep ensure -v -vendor-only \
  && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o simulator ./cmd/simulator/*

# ===

FROM debian:buster-slim

RUN addgroup --quiet app
RUN adduser --quiet --disabled-password --gecos "" --no-create-home --ingroup app app

WORKDIR /home/app

COPY --from=builder /go/src/github.com/controlplaneio/simulator-standalone/cli/dist/simulator /usr/local/bin/
RUN chown -R app:app ./

USER app

CMD ["simulator"]

