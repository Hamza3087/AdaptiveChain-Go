.PHONY: all build test clean

all: build

build:
	go build -o bin/blockchain ./cmd/blockchain
	go build -o bin/blockchain-cli ./cmd/blockchain-cli
	go build -o bin/blockchain-api ./cmd/blockchain-api

test:
	go test -v ./pkg/...

clean:
	rm -rf bin/

run-node:
	./bin/blockchain -id node-1 -addr 127.0.0.1:8000

run-api:
	./bin/blockchain-api -api 127.0.0.1:8080

generate-key:
	./bin/blockchain-cli generate-key key.pem
