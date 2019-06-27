FROM debian:buster-slim

RUN apt-get update && apt-get install -y golang build-essential git ca-certificates curl unzip

RUN mkdir -p /go/src/github.com/controlplaneio/simulator-standalone/cli
WORKDIR /go/src/github.com/controlplaneio/simulator-standalone/cli

RUN addgroup --quiet build
RUN adduser --quiet --disabled-password --gecos "" --ingroup build build
RUN chown -R build:build ./

ENV GO111MODULE=on

## Terraform
ENV TERRAFORM_VERSION 0.12.3
RUN curl -sLO https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip
RUN unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip
RUN mv terraform /usr/local/bin/

USER build

CMD ["make"]
