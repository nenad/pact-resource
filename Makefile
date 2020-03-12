build:
	go build -o bin/check ./cmd/check
	go build -o bin/in ./cmd/in

test-integration:
	go test $$(go list ./cmd/...)

test-short:
	go test -short $$(go list ./...)
