DOCKER_COMPOSE_EXISTS := $(shell command -v docker-compose 2>/dev/null)

ifdef DOCKER_COMPOSE_EXISTS
	GO = docker-compose exec resource go
else
	GO = go
endif

build:
	$(GO) build -o bin/check ./cmd/check
	$(GO) build -o bin/in ./cmd/in

test: test-short test-integration

test-integration:
	$(GO) -v -parallel 1 -run Integration $$(go list ./...)

test-short:
	$(GO) -v -short $$(go list ./...)
