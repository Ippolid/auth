##запускать из корня
#FROM golang:1.23.3-alpine AS builder
#
#COPY . /github.com/Ippolid/auth/source/
#WORKDIR /github.com/Ippolid/auth/source/
#
#RUN go mod download
#RUN go build -o ./bin/server cmd/grpc_server/main.go
#
#FROM alpine:latest
#
#WORKDIR /root/
#COPY --from=builder /github.com/Ippolid/auth/source/bin/server .
#COPY ./deploy/.env .env
#
#CMD ["./server"]

# syntax=docker/dockerfile:1.2
FROM golang:1.23.3-alpine AS builder

COPY . /github.com/Ippolid/auth/source/
WORKDIR /github.com/Ippolid/auth/source/

# кешируем папку модулей между сборками
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

RUN go build -o ./bin/server cmd/grpc_server/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/Ippolid/auth/source/bin/server .
COPY ./deploy/.env .env

CMD ["./server"]
