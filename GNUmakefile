default: fmt lint install generate

build:
	go build -v ./...
generate:
	cd tools && go generate ./...
install: build
	go install -v ./...
lint:
	golangci-lint run
fmt:
	gofmt -s -w -e .
	
testacc:
	TF_ACC=1 VERSION=1.4.4 go test -v -cover -timeout 120m ./...

.PHONY: fmt lint test testacc build install generate