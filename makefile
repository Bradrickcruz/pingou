APP=pingou
CMD=./cmd/$(APP)
BIN=bin/$(APP)

VERSION?=dev
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildDate=$(BUILD_DATE)"

.PHONY: all run fmt build clean

all: fmt build run

run: fmt
		. ./.env.local && go run $(CMD)/...

fmt:
	gofumpt -w .

build:
	mkdir -p bin
	go build $(LDFLAGS) -o $(BIN) $(CMD)/...

clean:
	rm -rf bin