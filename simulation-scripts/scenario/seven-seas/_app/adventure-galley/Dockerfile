# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.20-bullseye AS build

WORKDIR /adventure-galley

COPY /static /var/www/static

# Go Webapp
COPY go.mod ./
COPY *.go ./

RUN go build -o /var/www/target -ldflags='-w -s' ./main.go

##
## Deploy
##
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build --chown=nonroot:nonroot /var/www /var/www

USER nonroot:nonroot

EXPOSE 8080

CMD ["/var/www/target"]
