FROM golang:1.22.1-alpine3.18 AS builder

RUN go version

COPY ./ /0xterminal-helper
WORKDIR /0xterminal-helper

RUN go mod download
RUN go build -v -o ./.bin/terminal ./cmd/terminal/main.go

# Lightweight docker container with binary files
FROM alpine:latest

WORKDIR /app

COPY --from=builder /0xterminal-helper/.bin ./.bin
COPY --from=builder /0xterminal-helper/config ./config

CMD ["./.bin/terminal"]
