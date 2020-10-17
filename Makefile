.PHONY: all build test

build:
	mkdir -p bin
	CGO_ENABLED=0 go build -o bin/funds-loader

test:
	go test -race ./...
