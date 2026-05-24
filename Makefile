# -------
# USAGE:
# make build        # compile both node and server -> bin/
# make run-node     # run the node directly (no compile step)
# make run-server   # run the server directly (no compile)
# make test         # run all tests
# make test-race    # run tests with the race detector on
# make proto 				# regenerate .proto Go bindings
# make lint         # format + vet all code
# make clean        # remove bin/
# -------

BINARY_DIR := bin
PROTO_DIR := proto
GEN_DIR := gen/testfs

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

## Proto
proto:
	protoc \
		--go_out=$(GEN_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(GEN_DIR) --go-grpc_opt=paths=source_relative \
		--proto_path=$(PROTO_DIR) \
		$(PROTO_DIR)/testfs.proto

## Code quality
fmt:
	go fmt ./...

vet:
	go vet ./...

lint: fmt vet

## Clean build artifacts
clean:
	rm -rf $(BINARY_DIR)

