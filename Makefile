BINARY := license-scanner

.PHONY: build test lint scan-self

build:
	go build -o bin/$(BINARY) ./cmd/scanner

test:
	go test ./...

lint:
	go vet ./...

scan-self:
	go run ./cmd/scanner scan . --format text
