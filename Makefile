.PHONY: build test test-all server client1 client2 client3

build:
	@go build ./...

test:
	@go test ./... -race -short

test-all:
	@go test ./... -race

PROTOCOL = tcp

server:
	@go run ./cmd/server/*.go -protocol $(PROTOCOL)

client1:
	@go run ./cmd/client/*.go -protocol $(PROTOCOL) -user '{"user_id":1, "friends": [2, 3, 4]}'

client2:
	@go run ./cmd/client/*.go -protocol $(PROTOCOL) -user '{"user_id":2, "friends": [1]}'

client3:
	@go run ./cmd/client/*.go -protocol $(PROTOCOL) -user '{"user_id":3, "friends": [1, 2]}'