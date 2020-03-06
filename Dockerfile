FROM golang:alpine as builder

COPY . /resource
WORKDIR /resource

ENV CGO_ENABLED 0
RUN go build -o /assets/check ./check
RUN go build -o /assets/in ./in

FROM alpine:edge AS resource
RUN apk add --no-cache bash tzdata ca-certificates unzip zip gzip tar
COPY --from=builder /assets/ /opt/resource/
RUN chmod +x /opt/resource/*

# Test binaries exist
RUN stat /opt/resource/check /opt/resource/in
