# syntax=docker/dockerfile:1

## Build
FROM golang:1.20-buster AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY ./src ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /talos-os-config-generator ./api/

## Deploy
FROM gcr.io/distroless/static-debian11
WORKDIR /
COPY --from=build /talos-os-config-generator /talos-os-config-generator
EXPOSE 8080
USER nonroot:nonroot

ENTRYPOINT ["/talos-os-config-generator"]