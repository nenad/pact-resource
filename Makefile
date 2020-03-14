DOCKER_COMPOSE_EXISTS := $(shell command -v docker-compose 2>/dev/null)

ifdef DOCKER_COMPOSE_EXISTS
	GOTEST = docker-compose exec resource go test
else
	GOTEST = go test
endif

build:
	go build -o bin/check ./cmd/check
	go build -o bin/in ./cmd/in

test: test-short test-integration

test-integration:
	$(GOTEST) -v -parallel 1 -run Integration $$(go list ./...)

test-short:
	$(GOTEST) -v -short $$(go list ./...)
