# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o slack-backbone .

# Runtime stage
FROM alpine:3.20
RUN apk --no-cache add ca-certificates

# MCP transport configuration
# Default mode is "cli" (slash commands). Set to "mcp" for agent integration.
# When mode=mcp, optionally set HTTP_PORT for streamable HTTP transport.
ENV MODE="cli"
ENV HTTP_PORT=""

COPY --from=builder /app/slack-backbone /usr/local/bin/

ENTRYPOINT ["slack-backbone"]
CMD ["--mode", "${MODE}", "--http-port", "${HTTP_PORT}"]
