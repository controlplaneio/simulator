# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.18-bullseye AS build

WORKDIR /hashjacker

COPY /static /var/www/static
COPY /tmp /tmp

# Go Webapp
COPY go.mod ./
COPY *.go ./

RUN go build -o /var/www/hashjacker -ldflags='-w -s' ./main.go

##
## Deploy
##
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /var/www /var/www
COPY --from=build /tmp /tmp
COPY --from=build /go/bin /usr/bin

EXPOSE 8080

CMD ["/var/www/hashjacker"]
