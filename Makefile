.PHONY: all build run test lint tidy

tidy:
	go mod tidy

build: tidy
	go build -o bin/gateway main.go

run: build
	./bin/gateway api

test: tidy
	go test -v ./...

lint:
	go vet ./...
