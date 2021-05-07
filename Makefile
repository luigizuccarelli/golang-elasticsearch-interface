.PHONY: all test build clean

all: test build

build: 
	mkdir -p build
	go build -o build -tags real ./...

build-dev:
	mkdir -p build
	GOOS=linux go build -ldflags="-s -w" -o build ./...
	chmod 755 build/microservice
	chmod 755 build/uid_entrypoint.sh

test:
	go test -v -coverprofile=tests/results/cover.out -tags fake ./...

cover:
	go tool cover -html=tests/results/cover.out -o tests/results/cover.html

clean:
	rm -rf build/microservice
	go clean ./...

container:
	podman build -t quay.io/luigizuccarelli/golang-elasticsearch-interface:1.16.3 .

push:
	podman push --authfile=/home/lzuccarelli/config.json quay.io/luigizuccarelli/golang-elasticsearch-interface:1.16.3 
