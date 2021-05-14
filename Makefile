ROOT_DIR := $(shell pwd)

default: clean build

build:
	go mod tidy
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o "$(ROOT_DIR)/bin/lambda" "$(ROOT_DIR)/cmd/lambda/main.go"
	zip -j "$(ROOT_DIR)/bin/lambda.zip" "$(ROOT_DIR)/bin/lambda"

clean:
	rm -rf "$(ROOT_DIR)/bin"
