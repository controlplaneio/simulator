ARG GOLANG_IMAGE=golang:1.21.3-alpine3.18@sha256:27c76dcf886c5024320f4fa8ceb57d907494a3bb3d477d0aa7ac8385acd871ea
ARG GOLANGCI_LINT_IMAGE=golangci/golangci-lint:latest@sha256:c87d8a1a6521748fee124920c8e9302934ed26c9d3d48982449192b420a34686
ARG PACKER_IMAGE=hashicorp/packer:1.9@sha256:03808122fbfdd88e03be0d21cce9b3317778319b415c77e88efe1a98db82c76a
ARG TERRAFORM_IMAGE=hashicorp/terraform:1.5@sha256:c3bc74e7a2a8fab8216cbbedf12a9637db09288806a6aa537b6f397cba04dd93
ARG UBUNTU_IMAGE=ubuntu:mantic@sha256:13f233a16be210b57907b98b0d927ceff7571df390701e14fe1f3901b2c4a4d7

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
