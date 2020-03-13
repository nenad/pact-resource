build:
	go build -o bin/check ./cmd/check
	go build -o bin/in ./cmd/in

test: test-short test-integration

test-integration:
	go test -run Integration $$(go list ./...)

test-short:
	go test -short $$(go list ./...)
