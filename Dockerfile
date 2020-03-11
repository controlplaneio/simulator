#---------------------------------------#
# Dependencies, Linting & JS unit tests #
#---------------------------------------#
FROM debian:buster-slim AS dependencies

# We're using sh not bash at this point
# hadolint ignore=DL4006
RUN apt-get update                                                                    \
    && DEBIAN_FRONTEND=noninteractive apt-get install  -y --no-install-recommends     \
    curl                                                                              \
    software-properties-common                                                        \
    && curl -sL https://deb.nodesource.com/setup_13.x | bash -                        \
    && apt-get update                                                                 \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends      \
    build-essential                                                                   \
    ca-certificates                                                                   \
    git                                                                               \
    nodejs                                                                            \
    shellcheck                                                                        \
    unzip

# Download and save golang latest for use in other layers and install
ARG GO_INSTALL_VERSION=1.13.7
# hadolint ignore=DL3003,DL3010
RUN mkdir /downloads                                                  \
    && cd /downloads                                                  \
    && curl -sLO https://dl.google.com/go/go${GO_INSTALL_VERSION}.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go${GO_INSTALL_VERSION}.linux-amd64.tar.gz

ENV PATH $PATH:/usr/local/go/bin

# Install terraform
ENV GOPATH /go
ENV PATH $PATH:/go/bin
RUN mkdir -p /go/ && \
    chdir /go     && \
    go get -d -v github.com/hashicorp/terraform && \
    go install ./src/github.com/hashicorp/terraform/tools/terraform-bundle
COPY ./terraform/deployments/AWS/terraform-bundle.hcl .
RUN terraform-bundle package terraform-bundle.hcl && \
    mkdir -p terraform-bundle                     && \
    unzip -d terraform-bundle terraform_*.zip

# Default configuration for dep
ARG JQ_VERSION=1.6
ARG YQ_VERSION=2.7.2
ARG GOSS_VERSION=v0.3.7
ARG HADOLINT_VERSION=v1.16.3
ARG lint_user=lint

# Install JQ
RUN curl -sL https://github.com/stedolan/jq/releases/download/jq-${JQ_VERSION}/jq-linux64 \
      -o /usr/local/bin/jq                                                                \
    && chmod +x /usr/local/bin/jq

## Install YQ
RUN curl -sL https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64 \
      -o /usr/local/bin/yq                                                           \
    && chmod +x /usr/local/bin/yq

## Install Goss
RUN curl -sL https://github.com/aelsabbahy/goss/releases/download/${GOSS_VERSION}/goss-linux-amd64 \
         -o /usr/local/bin/goss                                                                    \
    && chmod +rx /usr/local/bin/goss

# Install Hadolint and setup non-root lint user
RUN curl -sL https://github.com/hadolint/hadolint/releases/download/${HADOLINT_VERSION}/hadolint-Linux-x86_64 \
        -o /usr/local/bin/hadolint                                                                            \
    && chmod +x /usr/local/bin/hadolint \
    && useradd -ms /bin/bash ${lint_user} \
    && mkdir /app

WORKDIR /app/scenario-tools

COPY --chown=1000 ./tools/scenario-tools/ /app/scenario-tools/

# Run javascript linting and unit tests
RUN npm install   \
    && npm test

WORKDIR /app

# Copy Dockerfiles, hadolint config and scripts
COPY --chown=1000 scripts/ /app/scripts/
COPY --chown=1000 attack/ /app/attack/
COPY --chown=1000 simulation-scripts/ /app/simulation-scripts/
COPY --chown=1000 kubesim /app/kubesim
COPY --chown=1000 Dockerfile .hadolint.yaml /app/
COPY --chown=1000 terraform/modules/AWS/Bastion/bashrc /app/Bastion/bashrc
COPY --chown=1000 terraform/modules/AWS/InternalHost/bashrc /app/InternalHost/bashrc
COPY --chown=1000 terraform/modules/AWS/Kubernetes/bashrc /app/Kubernetes/bashrc
COPY --chown=1000 launch-files/bashrc /app/launch-files/bashrc

USER ${lint_user}

# Lint Dockerfiles
RUN hadolint Dockerfile                         \
    &&  hadolint attack/Dockerfile              \
# Lint shell scripts
    && shellcheck scripts/*                     \
    && shellcheck attack/scripts/*              \
    && shellcheck simulation-scripts/perturb.sh \
    && shellcheck kubesim                       \
    && shellcheck Bastion/bashrc                \
    && shellcheck InternalHost/bashrc           \
    && shellcheck Kubernetes/bashrc             \
    && shellcheck launch-files/bashrc

WORKDIR /app/scenario-tools

#-----------------------#
# Golang Build and Test #
#-----------------------#
FROM debian:buster-slim AS build-and-test

RUN apt-get update                                                               \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
    awscli                                                                       \
    build-essential                                                              \
    ca-certificates                                                              \
    curl                                                                         \
    git                                                                          \
    openssh-client                                                               \
    unzip

# Install golang version downloaded in dependency stage
COPY --from=dependencies /terraform-bundle/* /usr/local/bin/
# hadolint ignore=DL3010
COPY --from=dependencies /downloads/go*.linux-amd64.tar.gz .
# We want to minimise layers to keep the build fast
# hadolint ignore=DL3010
RUN tar -C /usr/local -xzf go*.linux-amd64.tar.gz \
    && rm go*.linux-amd64.tar.gz
ENV PATH $PATH:/usr/local/go/bin

# Setup non-root build user
ARG build_user=build
RUN useradd -ms /bin/bash ${build_user} \
# Create golang src directory
    &&  mkdir -p /go/src/github.com/kubernetes-simulator/simulator \
# Create an empty public kubesim directory the tests
    && mkdir -p /home/${build_user}/.kubesim                            \
# Create module cache and copy manifest files
    &&  mkdir -p /home/${build_user}/go/pkg/mod

COPY ./go.* /go/src/github.com/kubernetes-simulator/simulator/

# Give ownership of module cache and src tree to build user
RUN chown -R ${build_user}:${build_user} /go/src/github.com/kubernetes-simulator/simulator \
    && chown -R ${build_user}:${build_user} /home/${build_user}

# Run all build and test steps as build user
USER ${build_user}

# Install golang module dependencies before copying source to cache them in their own layer
WORKDIR /go/src/github.com/kubernetes-simulator/simulator

# Add the full source tree
COPY --chown=1000 Makefile /go/src/github.com/kubernetes-simulator/simulator/
COPY --chown=1000 prelude.mk /go/src/github.com/kubernetes-simulator/simulator/
COPY --chown=1000 main.go /go/src/github.com/kubernetes-simulator/simulator/
COPY --chown=1000 pkg/  /go/src/github.com/kubernetes-simulator/simulator/pkg
COPY --chown=1000 cmd/  /go/src/github.com/kubernetes-simulator/simulator/cmd
COPY --chown=1000 test/  /go/src/github.com/kubernetes-simulator/simulator/test

WORKDIR /go/src/github.com/kubernetes-simulator/simulator/

# TODO: (rem) why is this owned by root after the earlier chmod?
USER root
# We're using sh not bash at this point
# hadolint ignore=DL4006
RUN chown -R ${build_user}:${build_user} /go/src/github.com/kubernetes-simulator/simulator/ \
    && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b /usr/local/bin v1.22.2

USER ${build_user}

# Golang build and test
WORKDIR /go/src/github.com/kubernetes-simulator/simulator
RUN make test-unit

#------------------#
# Launch Container #
#------------------#
FROM debian:buster-slim

RUN apt-get update                                                               \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
    awscli                                                                       \
    bash                                                                         \
    bash-completion                                                              \
    bzip2                                                                        \
    ca-certificates                                                              \
    curl                                                                         \
    file                                                                         \
    gettext-base                                                                 \
    gnupg                                                                        \
    lsb-release                                                                  \
    make                                                                         \
    openssh-client                                                               \
    tcl                                                                          \
    tcl-expect                                                                   \
 && rm -rf /var/lib/apt/lists/*

# Install golang version downloaded in dependency stage
COPY --from=dependencies /terraform-bundle/* /usr/local/bin/
# hadolint ignore=DL3010
COPY --from=dependencies /downloads/go*.linux-amd64.tar.gz .
RUN tar -C /usr/local -xzf go*.linux-amd64.tar.gz \
    && rm go*.linux-amd64.tar.gz
ENV PATH $PATH:/usr/local/go/bin

# Add login message
COPY ./scripts/launch-motd /usr/local/bin/launch-motd
RUN echo '[ ! -z "$TERM" ] && source /usr/local/bin/launch-motd' >> /etc/bash.bashrc

# Use 3rd party dependencies from build
COPY --from=dependencies /usr/local/bin/jq /usr/local/bin/jq
COPY --from=dependencies /usr/local/bin/yq /usr/local/bin/yq
COPY --from=dependencies /usr/local/bin/goss /usr/local/bin/goss
COPY --from=dependencies /terraform-bundle/* /usr/local/bin/

# Copy statically linked simulator binary
COPY --from=build-and-test /go/src/github.com/kubernetes-simulator/simulator/dist/simulator /usr/local/bin/simulator

# Setup non-root launch user
ARG launch_user=launch
RUN useradd -ms /bin/bash ${launch_user} \
    && mkdir /app                        \
    && chown -R ${launch_user}:${launch_user} /app \
    && mkdir -p /home/${launch_user}/.kubesim \
    && chown -R ${launch_user}:${launch_user} /home/${launch_user}/.kubesim

# Copy acceptance and smoke tests
COPY --chown=1000 --from=build-and-test /go/src/github.com/kubernetes-simulator/simulator/test/ /app/test/


WORKDIR /app

# Add terraform and perturb/scenario scripts to the image and goss.yaml to verify the container
ARG config_file="./launch-files/simulator.yaml"
COPY --chown=1000 ./terraform/ ./terraform/
COPY --chown=1000 ./simulation-scripts/ ./simulation-scripts/
COPY --chown=1000                     \
  ./launch-files/goss.yaml            \
  ./launch-files/launch-entrypoint.sh \
  ./launch-files/test-acceptance.sh   \
  ./
COPY --chown=1000 ./launch-files/bash_aliases /home/launch/.bash_aliases
COPY --chown=1000 ./launch-files/inputrc /home/launch/.inputrc
COPY --chown=1000 ${config_file} /home/launch/.kubesim/

COPY --chown=1000 launch-files/bashrc /home/launch/.bashrc

ENV SIMULATOR_SCENARIOS_DIR=/app/simulation-scripts/ \
    SIMULATOR_TF_DIR=/app/terraform/deployments/AWS

USER ${launch_user}

STOPSIGNAL SIGTERM

ENTRYPOINT [ "./launch-entrypoint.sh" ]
