FROM golang:alpine as builder

RUN apk update && apk add make

COPY . /resource
WORKDIR /resource

ENV CGO_ENABLED 0
RUN make build

FROM alpine:edge AS resource
RUN apk add --no-cache bash tzdata ca-certificates unzip zip gzip tar
COPY --from=builder /resource/bin/ /opt/resource/
RUN chmod +x /opt/resource/*

# Test binaries exist
RUN stat /opt/resource/check /opt/resource/in
