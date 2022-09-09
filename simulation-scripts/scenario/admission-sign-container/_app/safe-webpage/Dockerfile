# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.18-bullseye AS build

WORKDIR /safe-webpage

# Safe Website
COPY /keys /var/www/keys
COPY /static /var/www/static
COPY /safe-sign-keys /safe-sign-keys

# Go Webapp
COPY go.mod ./
COPY *.go ./

RUN go build -o /var/www/safe-secure -ldflags='-w -s' ./main.go

# Install Cosign
RUN go install github.com/sigstore/cosign/cmd/cosign@v1.11.1

##
## Deploy
##
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build --chown=nonroot:nonroot /var/www /var/www
COPY --from=build --chown=nonroot:nonroot /safe-sign-keys /safe-sign-keys
COPY --from=build /go/bin /usr/bin

ENV COSIGN_PASSWORD=""

USER nonroot:nonroot

#COPY --chown=nonroot:nonroot /safe-sign-keys /safe-sign-keys

EXPOSE 8090

CMD ["/var/www/safe-secure"]
