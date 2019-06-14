FROM docker:18.09.2 AS static-docker-source

FROM debian:buster-slim

ENV TERRAFORM_VERSION 0.12.2
ENV JQ_VERSION 1.6
ENV YQ_VERSION 2.7.2
ENV GOSS_VERSION v0.3.7

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
      unzip \
 && rm -rf /var/lib/apt/lists/*

#
# 3rd party dependencies
#

# jq
RUN curl https://github.com/stedolan/jq/releases/download/jq-${JQ_VERSION}/jq-linux64 -Lo /usr/local/bin/jq
RUN chmod +x /usr/local/bin/jq
## yq
RUN curl https://github.com/mikefarah/yq/releases/download/${YQ_VERSION}/yq_linux_amd64 -Lo /usr/local/bin/yq
RUN chmod +x /usr/local/bin/yq
## Terraform
RUN curl -sLO https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip
RUN unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip
RUN mv terraform /usr/local/bin/
## goss
RUN curl -L https://github.com/aelsabbahy/goss/releases/download/${GOSS_VERSION}/goss-linux-amd64 -o /usr/local/bin/goss
RUN chmod +rx /usr/local/bin/goss
## docker
COPY --from=static-docker-source /usr/local/bin/docker /usr/local/bin/docker

RUN useradd -ms /bin/bash launch-user
RUN mkdir /app
ADD . /app
RUN chown -R launch-user:launch-user /app

USER launch-user
RUN ssh-keygen -f /home/launch-user/.ssh/id_rsa -t rsa -N ''
WORKDIR app

#
ENV TF_VAR_shared_credentials_file /app/credentials
CMD [ "/bin/bash" ]
