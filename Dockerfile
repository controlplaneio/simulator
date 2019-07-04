#--------------------------#
# Dependencies and Linting #
#--------------------------#
FROM debian:buster-slim AS dependencies

RUN apt-get update                                 \
    && apt-get install -y  --no-install-recommends \
    ca-certificates                                \
    curl                                           \
    unzip

# Install terraform
# TODO: (rem) use `terraform-bundle`
ENV TERRAFORM_VERSION 0.12.3
RUN curl -sLO https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip
RUN unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip
RUN mv terraform /usr/local/bin/

# Install JQ
ENV JQ_VERSION 1.6
RUN curl https://github.com/stedolan/jq/releases/download/jq-${JQ_VERSION}/jq-linux64 -Lo /usr/local/bin/jq
RUN chmod +x /usr/local/bin/jq

## Install YQ
ENV YQ_VERSION 2.7.2
RUN curl https://github.com/mikefarah/yq/releases/download/${YQ_VERSION}/yq_linux_amd64 -Lo /usr/local/bin/yq
RUN chmod +x /usr/local/bin/yq

## Install Goss
ENV GOSS_VERSION v0.3.7
RUN curl -L https://github.com/aelsabbahy/goss/releases/download/${GOSS_VERSION}/goss-linux-amd64 -o /usr/local/bin/goss
RUN chmod +rx /usr/local/bin/goss

# Install Hadolint
ENV HADOLINT_VERSION v1.16.3
RUN curl https://github.com/hadolint/hadolint/releases/download/${HADOLINT_VERSION}/hadolint-Linux-x86_64 -Lo /usr/local/bin/hadolint
RUN chmod +x /usr/local/bin/hadolint

# Setup non-root lint user
ARG lint_user=lint
RUN useradd -ms /bin/bash ${lint_user}

# Copy dockerfiles and hadolint configs
RUN mkdir /app
COPY Dockerfile  /app
COPY .hadolint.yaml  /app
COPY attack/ /app/attack/
RUN chown -R ${lint_user}:${lint_user} /app

WORKDIR /app
USER ${lint_user}

RUN hadolint Dockerfile
RUN hadolint attack/Dockerfile

#-----------------------#
# Golang Build and Test #
#-----------------------#
FROM debian:buster-slim AS build-and-test

RUN apt-get update                                \
    && apt-get install -y --no-install-recommends \
    golang                                        \
    build-essential                               \
    git                                           \
    ca-certificates                               \
    unzip

COPY --from=dependencies /usr/local/bin/terraform /usr/local/bin/terraform

# Setup non-root build user
RUN addgroup --quiet build && adduser --quiet --disabled-password --gecos "" --ingroup build build

# Create golang src directory
RUN mkdir -p /go/src/github.com/controlplaneio/simulator-standalone

# Create an empty public ssh key file for the tests
RUN mkdir -p /home/build/.ssh && echo  "ssh-rsa FOR TESTING" > /home/build/.ssh/id_rsa.pub
# Create module cache and copy manifest files
RUN mkdir -p /home/build/go/pkg/mod
COPY ./go.mod /go/src/github.com/controlplaneio/simulator-standalone
COPY ./go.sum /go/src/github.com/controlplaneio/simulator-standalone

# Give ownership of module cache and src tree to build user
RUN chown -R build:build /go/src/github.com/controlplaneio/simulator-standalone
RUN chown -R build:build /home/build/go

# Run all build and test steps as build user
USER build

# Install golang module dependencies before copying source to cache them in their own layer
WORKDIR /go/src/github.com/controlplaneio/simulator-standalone
RUN go mod download

# Add the full source tree
COPY .  /go/src/github.com/controlplaneio/simulator-standalone/
WORKDIR /go/src/github.com/controlplaneio/simulator-standalone/

# TODO: (rem) why is this owned by root after the earlier chmod?
USER root
RUN chown -R build:build /go/src/github.com/controlplaneio/simulator-standalone/

USER build

# Golang build and test
WORKDIR /go/src/github.com/controlplaneio/simulator-standalone
ENV GO111MODULE=on
RUN make test

#------------------#
# Launch Container #
#------------------#
FROM debian:buster-slim

RUN \
  DEBIAN_FRONTEND=noninteractive \
    apt update && apt install --assume-yes --no-install-recommends \
      bash \
      awscli \
      bzip2 \
      file \
      ca-certificates \
      curl \
      gettext-base \
      golang \
      lsb-release \
      make \
      openssh-client \
      gnupg \
 && rm -rf /var/lib/apt/lists/*

# Add login message
COPY --from=build-and-test /go/src/github.com/controlplaneio/simulator-standalone/scripts/launch-motd /usr/local/bin/launch-motd
RUN echo '[ ! -z "$TERM" ] && launch-motd' >> /etc/bash.bashrc

# Use 3rd party dependencies from build
COPY --from=dependencies /usr/local/bin/jq /usr/local/bin/jq
COPY --from=dependencies /usr/local/bin/yq /usr/local/bin/yq
COPY --from=dependencies /usr/local/bin/goss /usr/local/bin/goss
COPY --from=dependencies /usr/local/bin/terraform /usr/local/bin/terraform

# Copy statically linked simulator binary
COPY --from=build-and-test /go/src/github.com/controlplaneio/simulator-standalone/dist/simulator /usr/local/bin/simulator

# Setup non-root launch user
ARG launch_user=launch
RUN useradd -ms /bin/bash ${launch_user}
RUN mkdir /app
RUN chown -R ${launch_user}:${launch_user} /app

# Add terraform and perturb/scenario scripts to the image
COPY ./terraform /app/terraform
COPY ./simulation-scripts /app/simulation-scripts

# Add goss.yaml to verify the container
COPY ./goss.yaml /app

# Add simulator.yaml config file
# The path to the config file can be supplied at build time to provide a custom conig from the host
#
# docker build --build-arg config_file=/path/to/simulator.yaml .
ARG config_file=./simulator.yaml
COPY ${config_file} /app

ENV SIMULATOR_SCENARIOS_DIR /app/simulation-scripts/
ENV SIMULATOR_TF_DIR /app/terraform/deployments/AwsSimulatorStandalone
ENV TF_VAR_shared_credentials_file /app/credentials

USER ${launch_user}
RUN ssh-keygen -f /home/${launch_user}/.ssh/id_rsa -t rsa -N ''
WORKDIR /app

CMD [ "/bin/bash" ]
