.PHONY: build run clean vet fmt lint

APP_NAME = server
CMD_DIR = ./cmd/server

build:
	go build -o bin/$(APP_NAME) $(CMD_DIR)

run:
	go run $(CMD_DIR)

vet:
	go vet ./...

fmt:
	go fmt ./...

lint: vet

clean:
	rm -rf bin/

all: fmt vet build
