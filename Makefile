APP_NAME := cadops

.PHONY: build test fmt

build:
	go build ./cmd/cadops

test:
	go test ./...

fmt:
	go fmt ./...
