.PHONY: all build build-client build-server build-linux build-windows run clean help

GO := go

BUILD_DIR := build
SERVER_DIR := $(BUILD_DIR)/server
CLIENT_DIR := $(BUILD_DIR)/client

build: prepare-client prepare-server build-client build-server

prepare-server:
	@echo "Preparing server files..."
	@mkdir -p $(SERVER_DIR)

prepare-client:
	@echo "Preparing client files..."
	@mkdir -p $(CLIENT_DIR)
	@cp -v "$(shell $(GO) env GOROOT)/misc/wasm/wasm_exec.js" $(CLIENT_DIR)
	@cp -v shaders-manifest.json $(CLIENT_DIR)
	@cp -v assets-manifest.json $(CLIENT_DIR)
	@cp -v templates/* $(CLIENT_DIR)
	@cp -r shaders $(CLIENT_DIR)
	@cp -r assets $(CLIENT_DIR)

build-client:
	@echo "Building WebAssembly client..."
	@GOOS=js GOARCH=wasm $(GO) build -o $(CLIENT_DIR)/main.wasm cmd/client/main.go

build-server:
	@echo "Building server for current OS..."
	@$(GO) build -o $(SERVER_DIR)/server cmd/server/main.go

build-linux:
	@echo "Building server for Linux (amd64)..."
	@GOOS=linux GOARCH=amd64 $(GO) build -o $(SERVER_DIR)/server-linux-amd64 cmd/server/main.go

build-windows:
	@echo "Building server for Windows x64..."
	@GOOS=windows GOARCH=amd64 $(GO) build -o $(SERVER_DIR)/server-windows-amd64.exe cmd/server/main.go

run:
	@echo "Starting server..."
	@./$(SERVER_DIR)/server

clean:
	@echo "Cleaning build directory..."
	@rm -rf $(BUILD_DIR)

help:
	@echo ""
	@echo "Targets:"
	@echo ""
	@echo "  build                	- Build server for current OS and client"
	@echo "  build-client         	- Build WebAssembly client only"
	@echo "  build-server         	- Build server for current OS"
	@echo "  build-linux 			- Build Linux x64 server"
	@echo "  build-windows 			- Build Windows x64 server"
	@echo ""
	@echo "  run                  	- Run server for current OS"
	@echo "  clean                	- Remove all build artifacts"
	@echo "  help                 	- Show this help message"
	@echo ""