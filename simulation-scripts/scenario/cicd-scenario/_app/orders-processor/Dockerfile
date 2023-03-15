# syntax=docker/dockerfile:1

FROM python:3.10-alpine

WORKDIR /

COPY requirements.txt requirements.txt
RUN pip3 install -r requirements.txt

COPY process-orders.py process-orders.py

ENTRYPOINT [ "python3", "process-orders.py" ]