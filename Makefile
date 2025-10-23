.PHONY: build clean test install run-cli run-build run-inspector run-server run-http run-http-custom run-sse run-sse-custom

BINARY_NAME=gqai
BUILD_DIR=dist

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) .

clean:
	rm -rf $(BUILD_DIR)
	go clean

test:
	go test -v ./...

install:
	go install
