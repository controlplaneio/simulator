# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.20-bullseye AS build

WORKDIR /secret-store

COPY go.mod .
COPY go.sum .
COPY *.go ./

COPY . .
RUN go build -o secret-store -ldflags='-w -s' .

##
## Deploy
##
FROM gcr.io/distroless/base-debian10

WORKDIR /app

COPY --from=build /secret-store/secret-store .
COPY --from=build /secret-store/.env .

EXPOSE 5050

CMD ["./secret-store"]