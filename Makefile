BINARY_NAME=content_foundry
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X github.com/leo/content-foundry-cli/cmd.Version=$(VERSION)"

.PHONY: build install clean run lint test

build:
	go build $(LDFLAGS) -o $(BINARY_NAME) .

install: build
	install -m 755 $(BINARY_NAME) $(shell brew --prefix content_foundry)/bin/$(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME)
	go clean

run:
	go run $(LDFLAGS) .

lint:
	golangci-lint run

test:
	go test ./...
