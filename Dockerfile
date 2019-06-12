FROM docker:18.09.2 AS static-docker-source

FROM debian:buster

RUN \
  DEBIAN_FRONTEND=noninteractive \
    apt update && apt install --assume-yes --no-install-recommends \
      apt-transport-https \
      awscli \
      bzip2 \
      ca-certificates \
      curl \
      gettext-base \
      golang \
      lsb-release \
      make \
      openssh-client \
      gnupg \
      unzip \
  && rm -rf /var/lib/apt/lists/* \
  && curl https://github.com/stedolan/jq/releases/download/jq-1.5/jq-linux64 -Lo /usr/local/bin/jq \
  && chmod +x /usr/local/bin/jq \
  && curl https://github.com/mikefarah/yq/releases/download/2.1.1/yq_linux_amd64 -Lo /usr/local/bin/yq \
  && chmod +x /usr/local/bin/yq \
  && ssh-keyscan -H github.com gitlab.com bitbucket.org >> /etc/ssh/ssh_known_hosts \
  && useradd -ms /bin/bash jenkins

RUN cd $(mktemp -d) \
  && curl -sLO https://releases.hashicorp.com/terraform/0.12.1/terraform_0.12.1_linux_amd64.zip \
  && unzip terraform_0.12.1_linux_amd64.zip \
  && mv terraform /usr/local/bin/ \
  && terraform version

COPY --from=static-docker-source /usr/local/bin/docker /usr/local/bin/docker

RUN mkdir app
WORKDIR app
ADD . /app
ENV TF_VAR_shared_credentials_file /app/credentials
RUN docker --version
CMD bash
