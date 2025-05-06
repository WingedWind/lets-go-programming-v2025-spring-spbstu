.PHONY: build clean proto test run

# Build variables
GO=go
PROTOC=protoc
BINARY_NAME=contacts-service
GO_FILES=$(shell find . -name "*.go" -type f -not -path "./pkg/*")

# Build the binary
build:
	$(GO) build -o $(BINARY_NAME) cmd/server/main.go

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -rf pkg/api

# Generate proto files
proto:
	$(PROTOC) -I./api/proto \
		--go_out=./pkg/api --go_opt=paths=source_relative \
		--go-grpc_out=./pkg/api --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=./pkg/api --grpc-gateway_opt=paths=source_relative \
		api/proto/contacts/v1/contacts.proto

# Run tests
test:
	$(GO) test ./...

# Run the server
run: build
	./$(BINARY_NAME)

# Install dependencies
deps:
	$(GO) mod tidy

# Install tools
tools:
	$(GO) install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1
	$(GO) install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0
	$(GO) install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.7.3