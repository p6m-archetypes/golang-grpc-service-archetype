.PHONY: build test proto run

proto:
	buf generate

setup:
	go mod tidy

build:
	go build ./...

test:
	go test ./...

run:
	go run ./cmd/server
