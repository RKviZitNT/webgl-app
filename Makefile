all: build run

build: build-server build-client
	mkdir -p build/server
	mkdir -p build/client
	cp -v "$(shell go env GOROOT)/misc/wasm/wasm_exec.js" build/client
	cp -v config_shaders.json build/client
	cp -v config_assets.json build/client
	cp -v templates/* build/client
	cp -r shaders build/client
	cp -r assets build/client

build-client:
	@echo "Building WebAssembly client..."
	GOOS=js GOARCH=wasm go build -o build/client/main.wasm cmd/client/main.go

build-server:
	@echo "Building server..."
	go build -o build/server/server cmd/server/main.go

run:
	@echo "Starting server..."
	./build/server/server

clean:
	rm -rf build