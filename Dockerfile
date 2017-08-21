FROM golang:alpine

RUN apk add --no-cache \
    bash \
    git \
    gcc \
    musl-dev \
    curl \
    openssh \
    mysql-client

ENV GOPATH="/go"
WORKDIR /go/src/github.com/bombsimon/laundry
