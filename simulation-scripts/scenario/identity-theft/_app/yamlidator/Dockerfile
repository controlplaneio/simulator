# syntax=docker/dockerfile:1

FROM debian:testing as dl

WORKDIR /dl
RUN apt update && apt install -y --no-install-recommends \
    curl ca-certificates tar binutils

ENV VERSION="1.1.0"
RUN curl -sfSOL https://github.com/seredot/kubectl-example/releases/download/v$VERSION/kubectl-example_1.1.0_Linux_x86_64.tar.gz \
 && tar xvf kubectl-example_1.1.0_Linux_x86_64.tar.gz kubectl-example

ENV VERSION="1.24.0"
RUN curl -sfSOL "https://dl.k8s.io/release/v$VERSION/bin/linux/amd64/kubectl"


FROM ubuntu:20.04

RUN apt update && apt install -y --no-install-recommends \
    ca-certificates python3 python3-pip\
 && apt clean all

COPY --from=dl /dl/kubectl-example /usr/local/bin/
COPY --from=dl /dl/kubectl /usr/bin/

RUN chmod +x /usr/local/bin/kubectl-example
RUN chmod +x /usr/bin/kubectl

WORKDIR /

COPY requirements.txt requirements.txt
RUN pip3 install -r requirements.txt

RUN mkdir app
COPY app.py app.py
COPY pod-schema.py pod-schema.py

CMD [ "python3", "-m", "flask", "run", "--host=0.0.0.0" ]