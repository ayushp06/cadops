APP_NAME := cadops
GO ?= go
BIN_DIR := bin
DIST_DIR := dist

ifeq ($(OS),Windows_NT)
EXE := .exe
endif

BINARY := $(BIN_DIR)/$(APP_NAME)$(EXE)

.PHONY: build install test fmt clean build-all

build:
	mkdir -p $(BIN_DIR)
	$(GO) build -o $(BINARY) ./cmd/cadops

install:
	$(GO) install ./cmd/cadops

test:
	$(GO) test ./...

fmt:
	$(GO) fmt ./...

clean:
	rm -rf $(BIN_DIR) $(DIST_DIR)

build-all:
	mkdir -p $(DIST_DIR)
	mkdir -p $(DIST_DIR)/cadops-windows-amd64
	mkdir -p $(DIST_DIR)/cadops-windows-arm64
	mkdir -p $(DIST_DIR)/cadops-linux-amd64
	mkdir -p $(DIST_DIR)/cadops-linux-arm64
	mkdir -p $(DIST_DIR)/cadops-darwin-amd64
	mkdir -p $(DIST_DIR)/cadops-darwin-arm64
	GOOS=windows GOARCH=amd64 $(GO) build -o $(DIST_DIR)/cadops-windows-amd64/cadops.exe ./cmd/cadops
	GOOS=windows GOARCH=arm64 $(GO) build -o $(DIST_DIR)/cadops-windows-arm64/cadops.exe ./cmd/cadops
	GOOS=linux GOARCH=amd64 $(GO) build -o $(DIST_DIR)/cadops-linux-amd64/cadops ./cmd/cadops
	GOOS=linux GOARCH=arm64 $(GO) build -o $(DIST_DIR)/cadops-linux-arm64/cadops ./cmd/cadops
	GOOS=darwin GOARCH=amd64 $(GO) build -o $(DIST_DIR)/cadops-darwin-amd64/cadops ./cmd/cadops
	GOOS=darwin GOARCH=arm64 $(GO) build -o $(DIST_DIR)/cadops-darwin-arm64/cadops ./cmd/cadops
