FROM golang:alpine

RUN apk add --no-cache \
    bash \
    git \
    gcc \
    musl-dev \
    curl \
    openssh \
    mysql-client

ENV GOPATH="/go" \
    PATH="$PATH:$GOPATH/bin"

WORKDIR /go/src/github.com/bombsimon/laundry

COPY ./ /go/src/github.com/bombsimon/laundry/

RUN go get -u github.com/golang/dep/cmd/dep && dep ensure

CMD go build -o laundry cmd/laundry/*.go && ./laundry
