# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o slack-backbone .

# Runtime stage
FROM alpine:3.20
RUN apk --no-cache add ca-certificates tini

# MCP transport configuration
# Default mode is "cli" (slash commands). Set to "mcp" for agent integration.
# When mode=mcp, optionally set HTTP_PORT for streamable HTTP transport.
ENV MODE="cli"
ENV HTTP_PORT=""

COPY --from=builder /app/slack-backbone /usr/local/bin/
COPY scripts/docker-entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD slack-backbone --help >/dev/null 2>&1 || exit 1

# tini handles signal forwarding; entrypoint.sh expands env vars and passes $@ through.
ENTRYPOINT ["/entrypoint.sh"]
CMD []
