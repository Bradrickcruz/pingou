APP=pingou
CMD=./cmd/$(APP)
BIN=bin/$(APP)
GO?=$(shell which go)

VERSION?=$(shell git describe --tags --always --dirty)
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildDate=$(BUILD_DATE)"

.PHONY: all dev fmt build test clean build-web docker-build docker-up docker-down docker-size release

all: fmt build dev

dev: fmt
	. ./.env && $(GO) run $(CMD)/...

fmt:
	gofumpt -w .

build: build-web
	mkdir -p bin
	CGO_ENABLED=1 $(GO) build $(LDFLAGS) -o $(BIN) $(CMD)/...

test:
	$(GO) test ./...

clean:
	rm -rf bin

build-web:
	cd web && npm install && npm run build

# ── docker ─────────────────────────────────────────────────────
docker-build:
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		-t $(APP):$(VERSION) \
		-t $(APP):latest .

docker-up:
	docker compose up --build

docker-down:
	docker compose down

docker-size:
	@echo "📦 Building image for size check..."
	@docker build \
		--quiet \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		-t $(APP):size-check . > /dev/null
	@SIZE=$$(docker images $(APP):size-check --format "{{.Size}}"); \
	echo ""; \
	echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"; \
	echo "📊 Tamanho da imagem final ($(APP):size-check)"; \
	echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"; \
	docker images $(APP):size-check --format "  Repositório: {{.Repository}}"; \
	docker images $(APP):size-check --format "  Imagem:      {{.Repository}}:{{.Tag}}"; \
	docker images $(APP):size-check --format "  Tamanho:     {{.Size}}"; \
	docker images $(APP):size-check --format "  Criado:      {{.CreatedAt}}"; \
	echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"; \
	echo ""
	@docker rmi $(APP):size-check > /dev/null
	@echo "✓ Imagem de verificação removida"

# ── release ────────────────────────────────────────────────────
release: build
	@echo "✓ Release $(VERSION) pronto em $(BIN)"