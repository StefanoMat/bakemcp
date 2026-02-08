.PHONY: build test fmt lint install snapshot

build:
	go build -o bakemcp ./cmd/bakemcp

install: build
	cp bakemcp /usr/local/bin/bakemcp

test:
	go test ./...

fmt:
	gofmt -s -w .
	go mod tidy

lint: fmt
	go vet ./...

snapshot:
	goreleaser release --snapshot --clean
