APP=pingou
CMD=./cmd/$(APP)
BIN=bin/$(APP)
GO?=$(shell which go)

VERSION?=$(shell git describe --tags --always --dirty)
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildDate=$(BUILD_DATE)"

.PHONY: all run fmt build test clean

all: fmt build run

run: fmt
		. ./.env.local && $(GO) run $(CMD)/...

fmt:
	gofumpt -w .

build:
	mkdir -p bin
	$(GO) build $(LDFLAGS) -o $(BIN) $(CMD)/...

test:
	$(GO) test ./...

clean:
	rm -rf bin