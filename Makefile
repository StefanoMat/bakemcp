.PHONY: build test fmt lint install

build:
	go build -o bakemcp ./cmd/openapi2mcp

install: build
	cp bakemcp /usr/local/bin/bakemcp

test:
	go test ./...

fmt:
	gofmt -s -w .
	go mod tidy

lint: fmt
	go vet ./...
