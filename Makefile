# -------
# USAGE:
# make build        # compile both node and server -> bin/
# make run-node     # run the node directly (no compile step)
# make run-server   # run the server directly (no compile)
# make test         # run all tests
# make test-race    # run tests with the race detector on
# make lint         # format + vet all code
# make clean        # remove bin/
# -------

BINARY_DIR := bin
NODE_BIN   := $(BINARY_DIR)/node
SERVER_BIN := $(BINARY_DIR)/server

.PHONY: all build build-node build-server run-node run-server test clean fmt vet

all: build

## Build both binaries
build: build-node build-server

build-node:
	@mkdir -p $(BINARY_DIR)
	go build -o $(NODE_BIN) ./cmd/node

build-server:
	@mkdir -p $(BINARY_DIR)
	go build -o $(SERVER_BIN) ./cmd/server

## Run
run-node:
	go run ./cmd/node

run-server:
	go run ./cmd/server

## Test
test:
	go test ./...

test-race:
	go test -race ./...

## Code quality
fmt:
	go fmt ./...

vet:
	go vet ./...

lint: fmt vet

## Clean build artifacts
clean:
	rm -rf $(BINARY_DIR)

