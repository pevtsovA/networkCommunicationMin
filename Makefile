BINARY_NAME=bin

all: build

build: build-proxy build-server

build-proxy:
	@if not exist bin mkdir bin
	go build -o $(BINARY_NAME) ./proxy

build-server:
	@if not exist bin mkdir bin
	go build -o $(BINARY_NAME) ./server

ifeq ($(OS),Windows_NT)
clean:
	@del /Q bin\*
else
clean:
	@rm -rf bin/*
endif

fmt:
	go fmt ./...

test:
	@echo "run tests..."
	go test ./...
	@echo "tests completed!"

generate-mocks:
	@echo "Generate interfaces mocks..."
	go generate ./...
	@echo "Generation done!"