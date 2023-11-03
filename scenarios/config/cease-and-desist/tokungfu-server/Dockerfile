FROM golang:1.18-bullseye AS build

WORKDIR /

# Go Webapp
COPY go.mod ./
COPY go.sum ./
COPY *.go ./

RUN go build -o /tokungfu -ldflags='-w -s' ./main.go

##
## Deploy
##
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build --chown=nonroot:nonroot /tokungfu /app/tokungfu

USER nonroot:nonroot

EXPOSE 8080

CMD ["/app/tokungfu"]
