.PHONY: run build test lint

run:
	PINGOU_API_KEY=dev-key go run ./cmd/pingou/...

build:
	go build -o bin/pingou ./cmd/pingou/...

test:
	go test ./...

lint:
	golangci-lint run