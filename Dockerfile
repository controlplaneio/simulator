#--------------------------#
# Dependencies and Linting #
#--------------------------#
FROM debian:buster-slim AS dependencies

RUN apt-get update                                                               \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
    build-essential                                                              \
    ca-certificates                                                              \
    curl                                                                         \
    golang                                                                       \
    git                                                                          \
    shellcheck                                                                   \
    unzip

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

# Install JQ
ARG JQ_VERSION=1.6
RUN curl -sL https://github.com/stedolan/jq/releases/download/jq-${JQ_VERSION}/jq-linux64 \
      -o /usr/local/bin/jq                                                                \
    && chmod +x /usr/local/bin/jq

## Install YQ
ARG YQ_VERSION=2.7.2
RUN curl -sL https://github.com/mikefarah/yq/releases/download/${YQ_VERSION}/yq_linux_amd64 \
      -o /usr/local/bin/yq                                                                  \
    && chmod +x /usr/local/bin/yq

## Install Goss
ARG GOSS_VERSION=v0.3.7
RUN curl -sL https://github.com/aelsabbahy/goss/releases/download/${GOSS_VERSION}/goss-linux-amd64 \
         -o /usr/local/bin/goss                                                                    \
    && chmod +rx /usr/local/bin/goss

# Install Hadolint
ARG HADOLINT_VERSION=v1.16.3
RUN curl -sL https://github.com/hadolint/hadolint/releases/download/${HADOLINT_VERSION}/hadolint-Linux-x86_64 \
        -o /usr/local/bin/hadolint                                                                            \
    && chmod +x /usr/local/bin/hadolint

# Setup non-root lint user
ARG lint_user=lint
RUN useradd -ms /bin/bash ${lint_user} \
    && mkdir /app

WORKDIR /app

# Copy Dockerfiles, hadolint config and scripts
COPY --chown=1000 Dockerfile .hadolint.yaml ./
COPY --chown=1000 scripts/ ./scripts/
COPY --chown=1000 attack/ ./attack/
COPY --chown=1000 kubesim ./kubesim

USER ${lint_user}

# Lint Dockerfiles
RUN hadolint Dockerfile            \
    &&  hadolint attack/Dockerfile \
# Lint shell scripts
    && shellcheck scripts/*        \
    && shellcheck attack/scripts/* \
    && shellcheck kubesim

#-----------------------#
# Golang Build and Test #
#-----------------------#
FROM debian:buster-slim AS build-and-test

RUN apt-get update                                                               \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
    awscli                                                                       \
    build-essential                                                              \
    ca-certificates                                                              \
    git                                                                          \
    golang                                                                       \
    openssh-client                                                               \
    unzip

COPY --from=dependencies /terraform-bundle/* /usr/local/bin/

# Setup non-root build user
ARG build_user=build
RUN useradd -ms /bin/bash ${build_user}

# Create golang src directory
RUN mkdir -p /go/src/github.com/controlplaneio/simulator-standalone

# Create an empty public ssh key file for the tests
RUN mkdir -p /home/${build_user}/.ssh && echo  "ssh-rsa FOR TESTING" > /home/${build_user}/.ssh/cp_simulator_rsa.pub \
# Create module cache and copy manifest files
    &&  mkdir -p /home/${build_user}/go/pkg/mod
COPY ./go.* /go/src/github.com/controlplaneio/simulator-standalone/

# Give ownership of module cache and src tree to build user
RUN chown -R ${build_user}:${build_user} /go/src/github.com/controlplaneio/simulator-standalone \
    && chown -R ${build_user}:${build_user} /home/${build_user}

# Run all build and test steps as build user
USER ${build_user}

# Install golang module dependencies before copying source to cache them in their own layer
WORKDIR /go/src/github.com/controlplaneio/simulator-standalone
RUN go mod download

# Add the full source tree
COPY --chown=1000 .  /go/src/github.com/controlplaneio/simulator-standalone/
WORKDIR /go/src/github.com/controlplaneio/simulator-standalone/

# TODO: (rem) why is this owned by root after the earlier chmod?
USER root
RUN chown -R ${build_user}:${build_user} /go/src/github.com/controlplaneio/simulator-standalone/

USER ${build_user}

# Golang build and test
WORKDIR /go/src/github.com/controlplaneio/simulator-standalone
ENV GO111MODULE=on
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
    golang                                                                       \
    lsb-release                                                                  \
    make                                                                         \
    openssh-client                                                               \
    tcl                                                                          \
    tcl-expect                                                                   \
 && rm -rf /var/lib/apt/lists/*

# Add login message
COPY --from=build-and-test /go/src/github.com/controlplaneio/simulator-standalone/scripts/launch-motd /usr/local/bin/launch-motd
RUN echo '[ ! -z "$TERM" ] && source /usr/local/bin/launch-motd' >> /etc/bash.bashrc

# Use 3rd party dependencies from build
COPY --from=dependencies /usr/local/bin/jq /usr/local/bin/jq
COPY --from=dependencies /usr/local/bin/yq /usr/local/bin/yq
COPY --from=dependencies /usr/local/bin/goss /usr/local/bin/goss
COPY --from=dependencies /terraform-bundle/* /usr/local/bin/

# Copy statically linked simulator binary
COPY --from=build-and-test /go/src/github.com/controlplaneio/simulator-standalone/dist/simulator /usr/local/bin/simulator

# Setup non-root launch user
ARG launch_user=launch
RUN useradd -ms /bin/bash ${launch_user} \
    && mkdir /app                        \
    && chown -R ${launch_user}:${launch_user} /app \
    && mkdir -p /home/${launch_user}/.kubesim \
    && chown -R ${launch_user}:${launch_user} /home/${launch_user}/.kubesim

# Copy acceptance and smoke tests
COPY --chown=1000 --from=build-and-test /go/src/github.com/controlplaneio/simulator-standalone/test/ /app/test/


WORKDIR /app

# Add terraform and perturb/scenario scripts to the image and goss.yaml to verify the container
ARG config_file=./simulator.yaml
COPY --chown=1000 ./terraform/ ./terraform/
COPY --chown=1000 ./simulation-scripts/ ./simulation-scripts/
COPY --chown=1000 ./goss.yaml ./launch-entrypoint.sh ./acceptance.sh ./
COPY --chown=1000 ${config_file} /home/launch/.kubesim/

ENV SIMULATOR_SCENARIOS_DIR=/app/simulation-scripts/ \
    SIMULATOR_TF_DIR=/app/terraform/deployments/AWS

USER ${launch_user}

STOPSIGNAL SIGTERM

ENTRYPOINT [ "./launch-entrypoint.sh" ]
