FROM golang:alpine as builder

COPY . /resource
WORKDIR /resource

ENV CGO_ENABLED 0
RUN apk update && apk add git

RUN go get github.com/cespare/reflex

ENTRYPOINT /go/bin/reflex -s -r '\.go$' -- go test $(go list ./...)
