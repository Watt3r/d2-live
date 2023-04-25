FROM golang:1.20 AS builder

ENV GOPATH=/root/go

ARG VERSION=dev

WORKDIR /app

COPY go.* ./

RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY *.go ./
COPY ./internal ./internal

RUN --mount=type=cache,mode=0755,target=/root/.cache/go-build --mount=type=cache,mode=0755,target=/root/go \
  go build \
  -ldflags="-w -s -X 'main.Version=${VERSION}'" \
  -o /http-server


EXPOSE 8090

ENTRYPOINT ["/http-server"]

