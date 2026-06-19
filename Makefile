BINARY_NAME=hermes
CMD_PATH=./cmd/hermes
DIST_DIR=dist

.PHONY: build run test clean build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64 build-all

build: clean
	go build -o $(BINARY_NAME) $(CMD_PATH)

run: build
	./$(BINARY_NAME)

test:
	go test ./...

build-linux-amd64:
	mkdir -p $(DIST_DIR)
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 $(CMD_PATH)

build-linux-arm64:
	mkdir -p $(DIST_DIR)
	CGO_ENABLED=1 GOOS=linux GOARCH=arm64 CC=aarch64-linux-gnu-gcc go build -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 $(CMD_PATH)

build-darwin-amd64:
	mkdir -p $(DIST_DIR)
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 $(CMD_PATH)

build-darwin-arm64:
	mkdir -p $(DIST_DIR)
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 $(CMD_PATH)

build-all: build-linux-amd64 build-linux-arm64

clean:
	rm -f $(BINARY_NAME)
	rm -rf $(DIST_DIR)
