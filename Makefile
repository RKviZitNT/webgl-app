.PHONY: all build build-client build-server /
		build-all-systems /
		build-linux build-linux-arm /
		build-windows build-windows-arm /
		build-mac build-mac-arm /
		run clean help

GO := go

BUILD_DIR := build
SERVER_DIR := $(BUILD_DIR)
CLIENT_DIR := $(BUILD_DIR)/static

build: prepare-client build-client prepare-server build-server

prepare-server:
	@echo "Preparing server files..."
	@mkdir -p $(SERVER_DIR)

prepare-client:
	@echo "Preparing client files..."
	@mkdir -p $(CLIENT_DIR)
	@cp -v "$(shell $(GO) env GOROOT)/misc/wasm/wasm_exec.js" $(CLIENT_DIR)
	@cp -v templates/* $(CLIENT_DIR)
	@cp -r shaders $(CLIENT_DIR)
	@cp -r assets $(CLIENT_DIR)

build-client:
	@echo "Building WebAssembly client..."
	@GOOS=js GOARCH=wasm $(GO) build -o $(CLIENT_DIR)/main.wasm cmd/client/main.go

build-server:
	@echo "Building server for current OS..."
	@$(GO) build -o $(SERVER_DIR)/server cmd/server/main.go

build-all-systems: prepare-client build-client prepare-server build-linux build-linux-arm build-windows build-windows-arm build-mac build-mac-arm

build-linux:
	@echo "Building server for Linux (amd64)..."
	@GOOS=linux GOARCH=amd64 $(GO) build -o $(SERVER_DIR)/server-linux-amd64 cmd/server/main.go

build-windows:
	@echo "Building server for Windows (amd64)..."
	@GOOS=windows GOARCH=amd64 $(GO) build -o $(SERVER_DIR)/server-windows-amd64.exe cmd/server/main.go

build-mac:
	@echo "Building server for MacOS (amd64)..."
	@GOOS=darwin GOARCH=amd64 $(GO) build -o $(SERVER_DIR)/server-mac-amd64 cmd/server/main.go

build-linux-arm:
	@echo "Building server for Linux (arm64)..."
	@GOOS=linux GOARCH=arm64 $(GO) build -o $(SERVER_DIR)/server-linux-arm cmd/server/main.go

build-windows-arm:
	@echo "Building server for Windows (arm64)..."
	@GOOS=windows GOARCH=arm64 $(GO) build -o $(SERVER_DIR)/server-windows-arm.exe cmd/server/main.go

build-mac-arm:
	@echo "Building server for MacOS (arm64)..."
	@GOOS=darwin GOARCH=arm64 $(GO) build -o $(SERVER_DIR)/server-mac-arm cmd/server/main.go

run:
	@echo "Starting server..."
	@./$(SERVER_DIR)

clean:
	@echo "Cleaning build directory..."
	@rm -rf $(BUILD_DIR)