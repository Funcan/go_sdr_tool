default: build

test:
	go test ./...

build: test
	go build
