# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.20-bullseye AS build

WORKDIR /pod-checker

COPY go.mod .
COPY go.sum .
COPY *.go ./

COPY . .
RUN go build -o pod-checker -ldflags='-w -s' .

##
## Deploy
##
FROM gcr.io/distroless/base-debian10

WORKDIR /app

COPY --from=build /pod-checker/pod-checker .
COPY --from=build /pod-checker/.env .
COPY --from=build /pod-checker/frontend/static/ frontend/static/
COPY --from=build /pod-checker/frontend/template/ frontend/template/

EXPOSE 8080

CMD ["./pod-checker"]