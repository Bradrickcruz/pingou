# ─── Stage 1: React ───────────────────────────────────────────
FROM node:22-alpine AS web-builder

WORKDIR /app/web
COPY web/package*.json ./
RUN npm ci
COPY web/ ./
RUN npm run build

# ─── Stage 2: Go ──────────────────────────────────────────────
FROM golang:1.25-alpine AS go-builder

# gcc necessário para go-sqlite3 (CGO)
RUN apk add --no-cache gcc musl-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=web-builder /app/internal/handler/dist ./internal/handler/dist

ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_DATE=unknown

RUN CGO_ENABLED=1 GOOS=linux go build \
  -ldflags="-s -w \
  -X github.com/Bradrickcruz/pingou/cmd/pingou/commands.version=${VERSION} \
  -X github.com/Bradrickcruz/pingou/cmd/pingou/commands.commit=${COMMIT} \
  -X github.com/Bradrickcruz/pingou/cmd/pingou/commands.date=${BUILD_DATE}" \
  -o /bin/pingou ./cmd/pingou

# ─── Stage 3: Final ───────────────────────────────────────────
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=go-builder /bin/pingou .

EXPOSE 8080

# Health check via endpoint público /healthz (sem auth)
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget -qO- http://localhost:8080/healthz || exit 1

ENTRYPOINT ["./pingou"]