all: build

build: build-server build-client
	mkdir -p build/static
	mkdir -p build/bin
	cp -v "$(shell go env GOROOT)/misc/wasm/wasm_exec.js" build/static/
	cp -v templates/* build/static
	cp -r shaders build/static

build-client:
	@echo "Building WebAssembly client..."
	GOOS=js GOARCH=wasm go build -o build/static/main.wasm cmd/client/main.go

build-server:
	@echo "Building server..."
	go build -o build/bin/server cmd/server/main.go

run:
	@echo "Starting server..."
	./build/bin/server

clean:
	rm -rf build