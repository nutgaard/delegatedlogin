BINARY_NAME=delegated-login

.PHONY: all build run clean test

all: test build

build:
	go build -o bin/${BINARY_NAME} cmd/main/main.go

run:
	go run cmd/local/main.go

clean:
	go clean
	go mod tidy
	rm -rf -v bin/*

test:
	go test -count=1 ./...