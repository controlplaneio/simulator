# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.20-bullseye AS build

WORKDIR /path

COPY /static /var/www/static

# Go Webapp
COPY go.mod ./
COPY *.go ./

RUN go build -o /var/www/path -ldflags='-w -s' ./main.go

##
## Deploy
##
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /var/www /var/www

EXPOSE 8080

CMD ["/var/www/path"]
