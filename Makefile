
all: clean format build

format:
	go fmt ./...

build: deps
	@mkdir -p bin/
	go build --ldflags '-extldflags "-lz"' -o bin/photosort ./...

buildstatic: deps
	@mkdir -p bin/
	go build --ldflags '-extldflags "-static -lz"' -o bin/photosort ./...

deps:
	@go get -d -v ./...

clean:
	@rm -rf bin/

.PHONY: all clean format build
