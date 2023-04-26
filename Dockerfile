#---------------------------------------#
# JS unit tests                         #
#---------------------------------------#
FROM docker.io/library/node:18-bullseye AS scenario-tools

WORKDIR /app/scenario-tools
COPY --chown=1000 ./tools/scenario-tools/ .

# Run javascript linting and unit tests
RUN npm install \
 && npm test

#---------------------------------------#
# Dependencies & Linting                #
#---------------------------------------#
FROM debian:bullseye-slim AS dependencies
# We're using sh not bash at this point
# hadolint ignore=DL4006
RUN apt-get update                                                                    \
    && DEBIAN_FRONTEND=noninteractive apt-get install  -y --no-install-recommends     \
    binutils                                                                          \
    curl                                                                              \
    build-essential                                                                   \
    ca-certificates                                                                   \
    git                                                                               \
    shellcheck                                                                        \
    unzip                                                                             \
    jq

# Default configuration for dep
ARG YQ_VERSION=3.4.1
ARG GOSS_VERSION=v0.3.21
ARG HADOLINT_VERSION=v2.12.0
ARG TERRAFORM_VERSION=1.4.1
ARG lint_user=lint

## Install YQ
RUN curl -sSOL https://github.com/mikefarah/yq/releases/download/${YQ_VERSION}/yq_linux_amd64 \
    && install -Dm755 -s yq_linux_amd64 /usr/local/bin/yq \
    && yq --version

## Install Goss
RUN curl -sSOL https://github.com/aelsabbahy/goss/releases/download/${GOSS_VERSION}/goss-linux-amd64 \
    && install -Dm755 -s goss-linux-amd64 /usr/local/bin/goss \
    && goss --version

# Install Hadolint and setup non-root lint user
RUN curl -sSOL https://github.com/hadolint/hadolint/releases/download/${HADOLINT_VERSION}/hadolint-Linux-x86_64 \
    && install -Dm755 hadolint-Linux-x86_64 /usr/local/bin/hadolint \
    && hadolint --version

# Install terraform
RUN curl -sSL --fail "https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip" \
    -o terraform.zip \
 && unzip terraform.zip \
 && install -Dm755 -s terraform /usr/local/bin/

RUN useradd -ms /bin/bash ${lint_user} \
    && mkdir -p /app

WORKDIR /app

# Copy Dockerfiles, hadolint config and scripts
COPY --chown=1000 scripts/ /app/scripts/
COPY --chown=1000 attack/ /app/attack/
COPY --chown=1000 simulation-scripts/ /app/simulation-scripts/
COPY --chown=1000 kubesim /app/kubesim
COPY --chown=1000 Dockerfile .hadolint.yaml /app/
COPY --chown=1000 terraform/modules/AWS/CloudInitCommon/bashrc /app/CloudInitCommon/bashrc
COPY --chown=1000 launch-files/bashrc /app/launch-files/bashrc

USER ${lint_user}

# Lint Dockerfiles & shell scripts
RUN hadolint Dockerfile &&                       \
    hadolint attack/Dockerfile &&                \
    shellcheck scripts/* &&                      \
    shellcheck attack/scripts/* &&               \
    shellcheck simulation-scripts/perturb.sh &&  \
    shellcheck kubesim &&                        \
    shellcheck CloudInitCommon/bashrc &&         \
    shellcheck launch-files/bashrc

#-----------------------#
# Golang Build and Test #
#-----------------------#
FROM docker.io/library/golang:1.16-bullseye AS build-and-test

RUN apt-get update                                                               \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
    awscli                                                                       \
    binutils                                                                          \
    build-essential                                                              \
    ca-certificates                                                              \
    curl                                                                         \
    git                                                                          \
    openssh-client                                                               \
    unzip

WORKDIR /app

# Setup non-root build user
ARG build_user=build
RUN useradd -ms /bin/bash ${build_user} \
# Create an empty public kubesim directory for the tests
 && mkdir -p /home/${build_user}/.kubesim \
 && chown -R ${build_user}:${build_user} /home/${build_user}

# Give ownership of src tree to build user
RUN mkdir -p simulator \
 && chown -R ${build_user}:${build_user} simulator
WORKDIR /app/simulator

# Run all build and test steps as build user
USER ${build_user}
ENV GOPATH /home/${build_user}/go
ENV PATH $PATH:$GOPATH/bin

# Add the full source tree
COPY --chown=1000 ./go.* ./
COPY --chown=1000 Makefile ./
COPY --chown=1000 prelude.mk ./
COPY --chown=1000 main.go ./

COPY --chown=1000 pkg/  ./pkg
COPY --chown=1000 cmd/  ./cmd
COPY --chown=1000 test/  ./test

# Golang build and test
COPY --from=dependencies /usr/local/bin/terraform /usr/local/bin/
RUN make test-unit

#------------------#
# Launch Container #
#------------------#
FROM debian:bullseye-slim

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
    jq                                                                           \
    lsb-release                                                                  \
    openssh-client                                                               \
 && rm -rf /var/lib/apt/lists/*

# Init Terraform
COPY --from=dependencies /usr/local/bin/terraform /usr/local/bin/

# Add login message
COPY ./scripts/launch-motd /usr/local/bin/launch-motd
RUN echo '[ ! -z "$TERM" ] && source /usr/local/bin/launch-motd' >> /etc/bash.bashrc

# Use 3rd party dependencies from build
COPY --from=dependencies /usr/local/bin/yq /usr/local/bin/yq
COPY --from=dependencies /usr/local/bin/goss /usr/local/bin/goss

# Copy simulator binary
COPY --from=build-and-test /app/simulator/dist/simulator /usr/local/bin/simulator

# Setup non-root launch user
ARG launch_user=launch
RUN useradd -ms /bin/bash ${launch_user} \
    && mkdir -p /app                     \
    && chown -R ${launch_user}:${launch_user} /app \
    && mkdir -p /home/${launch_user}/.kubesim \
    && chown -R ${launch_user}:${launch_user} /home/${launch_user}/.kubesim

# Copy acceptance and smoke tests
#COPY --chown=1000 --from=build-and-test /app/simulator/test/ /app/test/

# Copy scenario-tools
COPY --chown=1000 --from=scenario-tools /app/scenario-tools/*.json /app/scenario-tools/
COPY --chown=1000 --from=scenario-tools /app/scenario-tools/bin /app/scenario-tools/bin
COPY --chown=1000 --from=scenario-tools /app/scenario-tools/lib /app/scenario-tools/lib
COPY --chown=1000 --from=scenario-tools /app/scenario-tools/node_modules /app/scenario-tools/node_modules

WORKDIR /app

# Add terraform and perturb/scenario scripts to the image and goss.yaml to verify the container
ARG config_file="./launch-files/simulator.yaml"
COPY --chown=1000 ./terraform/ ./terraform/
RUN terraform -chdir=terraform/deployments/AWS init --get=true --backend=false \
 && chown -R ${launch_user}:${launch_user} terraform/deployments/AWS

COPY --chown=1000 ./simulation-scripts/ ./simulation-scripts/
COPY --chown=1000                     \
  ./launch-files/goss.yaml            \
  ./launch-files/launch-entrypoint.sh \
  ./
COPY --chown=1000 ./launch-files/bash_aliases /home/launch/.bash_aliases
COPY --chown=1000 ./launch-files/inputrc /home/launch/.inputrc
COPY --chown=1000 ${config_file} /home/launch/.kubesim/

COPY --chown=1000 launch-files/bashrc /home/launch/.bashrc

RUN curl --compressed --connect-timeout 5 -LO \
  "https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl" \
 && install -Dm755 ./kubectl /usr/local/bin/ \
 && rm kubectl

ENV SIMULATOR_SCENARIOS_DIR=/app/simulation-scripts/ \
    SIMULATOR_TF_DIR=/app/terraform/deployments/AWS \
    DISABLE_CHECKPOINT=true

USER ${launch_user}

STOPSIGNAL SIGTERM

ENTRYPOINT [ "./launch-entrypoint.sh" ]
