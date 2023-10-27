# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.18-bullseye AS build

WORKDIR /chest

# Go Webapp
COPY go.mod ./
COPY *.go ./

RUN go build -o /var/www/chest -ldflags='-w -s' ./main.go

##
## Deploy
##
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build --chown=nonroot:nonroot /var/www /var/www
COPY --from=build /go/bin /usr/bin

USER nonroot:nonroot

EXPOSE 8080

CMD ["/var/www/chest"]
