# syntax=docker/dockerfile:1

## Build
FROM golang:1.19-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./src/api/*.go ./api/
COPY ./src/generator/*.go ./generator/

RUN go build -o /talos-os-config-generator ./api/

## Deploy
FROM gcr.io/distroless/base-debian10
WORKDIR /
COPY --from=build /talos-os-config-generator /talos-os-config-generator
EXPOSE 8080
USER nonroot:nonroot

ENTRYPOINT ["/talos-os-config-generator"]