.PHONY: clean update

build/bundler: go.mod go.sum $(shell find . -type f -name '*.go') | build
	go build -v -o build/bundler ./cmd/bundler

build:
	mkdir -p build

clean:
	rm -rf build

update:
	go get -u ./...
	go mod tidy
