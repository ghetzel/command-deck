.PHONY: deps fmt test build package
.EXPORT_ALL_VARIABLES:

GO111MODULE ?= on
LOCALS      := $(shell find . -type f -name '*.go')

all: deps fmt build

deps:
	go get ./...
	go generate -x
	-go mod tidy

fmt:
	gofmt -w $(LOCALS)
	go vet ./...

test:
	go test -race ./...

build: fmt
	go build -o bin/cdeck $(LOCALS)