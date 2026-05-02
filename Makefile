APP=pingou
CMD=./cmd/$(APP)
BIN=bin/$(APP)
GO?=$(shell which go)

VERSION?=$(shell git describe --tags --always --dirty)
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildDate=$(BUILD_DATE)"

.PHONY: all dev fmt build test clean build-web docker-build docker-up docker-down docker-size docker-startup-test release

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

docker-startup-test:
	@echo "🚀 Validando startup time do docker compose..."
	@echo ""
	@docker compose down > /dev/null 2>&1 || true
	@docker system prune -f > /dev/null 2>&1 || true
	@echo "⏱️  Iniciando container..."
	@START=$$(date +%s%N); \
	docker compose up -d > /dev/null 2>&1; \
	\
	for i in {1..60}; do \
		if docker ps 2>/dev/null | grep -q "healthy"; then \
			END=$$(date +%s%N); \
			ELAPSED_NS=$$(((END - START) / 1000000)); \
			ELAPSED_S=$$(( ELAPSED_NS / 1000 )); \
			ELAPSED_MS=$$(( ELAPSED_NS % 1000 )); \
			echo ""; \
			echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"; \
			echo "✓ Container saudável"; \
			echo "⏱️  Tempo de startup: $${ELAPSED_S}.$${ELAPSED_MS}s"; \
			if [ $${ELAPSED_S} -lt 30 ]; then \
				echo "✅ PASSOU: < 30s"; \
			else \
				echo "❌ FALHOU: >= 30s"; \
			fi; \
			echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"; \
			echo ""; \
			exit 0; \
		fi; \
		sleep 1; \
	done; \
	echo ""; \
	echo "❌ TIMEOUT: Container não ficou healthy em 60s"; \
	echo ""; \
	docker compose logs; \
	exit 1

# ── release ────────────────────────────────────────────────────
release: build
	@echo "✓ Release $(VERSION) pronto em $(BIN)"