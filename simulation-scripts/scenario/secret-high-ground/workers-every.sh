#!/bin/bash

docker build -t jenkins/jenkins:lts-alpine - &> /dev/null <<EOF
FROM jenkins/jenkins:lts-alpine

USER root

RUN mkdir -pv /secrets && touch /secrets/aws-creds && echo "AWS_ACCESS_KEY_ID=EXAMPLEKEY12345" >> /secrets/aws-creds && echo "AWS_SECRET_ACCESS_KEY=SECRETACCESSKEY9876!*" >> /secrets/aws-creds

USER jenkins
EOF
