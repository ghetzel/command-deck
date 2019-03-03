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

build-all:
	GOOS=linux   go build -o ~/lib/apps/cdeck/linux/amd64/cdeck $(LOCALS)
	GOOS=freebsd go build -o ~/lib/apps/cdeck/freebsd/amd64/cdeck $(LOCALS)
	GOOS=darwin  go build -o ~/lib/apps/cdeck/darwin/amd64/cdeck $(LOCALS)

demo:
	@./bin/cdeck -e 'default:blue:Command' 'default:166:Deck' 'default:240::' 'default:green:A' 'default:56:simple' 'default:124:prompt' 'default:6:generator.'
	@echo
