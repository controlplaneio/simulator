ARG GOLANG_IMAGE=golang:1.21.5-alpine3.19@sha256:55f716237933c85cee01748700755b4ac8736fb1ca974c9aed051691b68d6dc2
ARG GOLANGCI_LINT_IMAGE=golangci/golangci-lint:latest@sha256:fb70c9b2e6d0763141f057abcafde7f88d5e4bb3b5882d6b14bc79382f04481c
ARG PACKER_IMAGE=hashicorp/packer:1.10@sha256:a10638519af09f5ecad52b6eb4eab489377e4e89f30ea46832f1f401a234d783
ARG TERRAFORM_IMAGE=hashicorp/terraform:1.6@sha256:d593c353357a3db5a795c2ba0b998580cf12bad9125807bd877092c2e813279b
ARG UBUNTU_IMAGE=ubuntu:mantic@sha256:8d093e0651575a6437cc4a3d561f892a345d263aeac6156ef378fe6a4ccabd4c

FROM ${GOLANGCI_LINT_IMAGE}

WORKDIR /app

COPY . .

RUN /usr/bin/golangci-lint run -v -c .golangci.yml

FROM ${GOLANG_IMAGE} as BUILDER

WORKDIR /build

COPY go.* ./
RUN go mod download

COPY . ./

RUN go build -v -o /simulator cmd/container/main.go

FROM ${PACKER_IMAGE} as PACKER
FROM ${TERRAFORM_IMAGE} as TERRAFORM
FROM ${UBUNTU_IMAGE}

WORKDIR simulator

COPY --from=PACKER /bin/packer /usr/local/bin/packer
COPY --from=TERRAFORM /bin/terraform /usr/local/bin/terraform

RUN apt update && \
    apt install -y ca-certificates openssh-client ansible-core && \
    rm -rf /var/lib/apt/lists/* && \
    ansible-galaxy collection install kubernetes.core

COPY --from=BUILDER /simulator /usr/local/bin/simulator

USER ubuntu

ENTRYPOINT ["/usr/local/bin/simulator"]
