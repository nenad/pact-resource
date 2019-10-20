FROM golang:alpine as builder

COPY . /resource
WORKDIR /resource

ENV CGO_ENABLED 0
RUN go build -o /assets/check ./check
RUN go build -o /assets/in ./in

RUN mkdir /tests && set -e; for pkg in $(go list ./...); do \
		go test -o "/tests/$(basename $pkg).test" -c $pkg; \
	done

FROM alpine:edge AS resource
RUN apk add --no-cache bash tzdata ca-certificates unzip zip gzip tar
COPY --from=builder assets/ /opt/resource/
RUN chmod +x /opt/resource/*

FROM resource AS tests
ARG BROKER_URL
COPY --from=builder /tests /resource-tests
WORKDIR /resource-tests
RUN set -e; for test in /resource-tests/*.test; do \
		$test; \
	done

FROM resource
