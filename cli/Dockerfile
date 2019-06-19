FROM golang:1.12 AS builder

RUN mkdir -p /go/src/github.com/controlplaneio/simulator-standalone/cli

WORKDIR /go/src/github.com/controlplaneio/simulator-standalone/cli

COPY . .

ENV GO111MODULE=on

RUN go install
RUN GOPATH= CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./dist/simulator

# ===

FROM debian:buster-slim

RUN addgroup --quiet app
RUN adduser --quiet --disabled-password --gecos "" --no-create-home --ingroup app app

WORKDIR /home/app

COPY --from=builder /go/src/github.com/controlplaneio/simulator-standalone/cli/dist/simulator /usr/local/bin/
RUN chown -R app:app ./

USER app

CMD ["simulator"]

