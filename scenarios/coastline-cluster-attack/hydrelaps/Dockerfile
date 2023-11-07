# syntax=docker/dockerfile:1

FROM ubuntu:22.04

WORKDIR /

RUN apt update && apt install -y python3 python3-pip

COPY requirements.txt requirements.txt
RUN pip3 install -r requirements.txt

ENV PATH="${PATH}:/usr/local/bin/"

RUN mkdir app
COPY application.py app/application.py
COPY db-init.py app/db-init.py
COPY entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

ENTRYPOINT ["entrypoint.sh"]