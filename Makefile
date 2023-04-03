build/bundler: $(shell find . -type f -name '*.go') | build
	go build -v -o build/bundler ./cmd/bundler

build:
	mkdir -p build

