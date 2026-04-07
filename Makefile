BINARY_NAME=content_foundry
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X github.com/leo/content-foundry-cli/cmd.Version=$(VERSION)"

.PHONY: build install clean run lint test

build:
	go build $(LDFLAGS) -o $(BINARY_NAME) .

install: build
	cp $(BINARY_NAME) $(shell go env GOPATH)/bin/$(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME)
	go clean

run:
	go run $(LDFLAGS) .

lint:
	go vet ./...

test:
	go test ./...
