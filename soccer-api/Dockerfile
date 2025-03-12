# DEV DOCKERFILE
FROM golang:1.23-alpine AS builder
WORKDIR /usr/src/app

ENV PROJECT_DIR=/app \
    GO111MODULE=on \
    CGO_ENABLED=0

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY . .

# HotReload - for development purpouse
RUN go get github.com/githubnemo/CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon
