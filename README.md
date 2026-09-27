# slack-backbone

A multi-team Slack socket-mode application written in Go. Exposes slash commands and notifies channels about events — no public DNS or HTTP server required.

## Features

- **Socket mode** via `Asafrose/bolt-go` — connects to Slack over WebSocket
- **Multi-team support** — one binary, N concurrent workspaces
- **Slash commands** with Block Kit rich layouts
- **Configurable per-team**: allowed/restricted commands, default channels
- **Structured JSON logging** via `slog.NewJSONHandler`
- **Go 1.25+** only

## Quick Start

```bash
# 1. Install dependencies
go mod tidy

# 2. Copy the example config and edit tokens
cp teams.yaml.example teams.yaml
# Edit teams.yaml with your actual bot_token and app_token per team

# 3. Run locally
go run . --config teams.yaml

# 4. Build a binary
go build -o slack-backbone .
```

## Configuration

Create a `teams.yaml` file:

```yaml
teams:
  - name: "prod-workspace"
    bot_token: "xoxb-..."
    app_token: "xapp-..."
    default_channel: "#ops-notifications"
    commands: all
    restricted_commands: []

  - name: "staging-workspace"
    bot_token: "xoxb-..."
    app_token: "xapp-..."
    default_channel: "#deployments"
    commands: all
    restricted_commands: [status]
```

## Commands

| Command | Scope | Description |
|---------|-------|-------------|
| `/slack-backbone health` | Global | Health check across all teams |
| `/slack-backbone status` | Global | Show running status per team |
| `/slack-backbone deploy --env <env>` | Scoped | Trigger a deployment |
| `/slack-backbone confirm` | Scoped | Interactive confirmation flow |

## Project Structure

```
├── main.go           # Entry point + signal handling
├── cmd/              # Cobra CLI commands
├── handlers/         # Command handler implementations
├── slack/            # Multi-team Bolt manager + API client
├── config/           # Config loading (YAML/env)
├── .github/workflows/ci.yml  # CI: build, test, lint, format, mod-tidy, docker-build
├── Dockerfile        # Multi-stage build (golang:1.25-alpine → alpine:3.20)
├── .dockerignore     # Excludes .git, .github, *.md, .env, logs from build context
├── teams.yaml.example # Example config for local testing
├── go.mod / go.sum   # Go module dependencies
└── README.md
```

## Building & Running

```bash
# Local development
go run . --log-level debug

# Production binary
go build -o slack-backbone .
./slack-backbone --config teams.yaml --team prod-workspace

# Docker (CI builds and tags with commit SHA)
docker build -t slack-backbone:${{ github.sha }} .
docker run --rm slack-backbone:${{ github.sha }} --log-level debug
```

## Acknowledgments

This project was built with the help of several tools and services:

- **[Berd](https://github.com/block/berd/)** — the desktop app where this conversation took place, providing a rich interface for working with agents.
- **[Goose](https://github.com/goose),** the coding agent used to design, implement, and review all changes in this project.
- **Unsloth Studio** ([unslothai/unsloth](https://github.com/unslothai/unsloth)) — the inference engine powering the LLM used during development.
- **llama.cpp** ([ggml-org/llama.cpp](https://github.com/ggml-org/llama.cpp)) — the inference engine that powers Unsloth Studio's runtime.
- **Qwen-3.8-35B-A3B** ([Empero AI](https://huggingface.co/empero-ai/Qwen3.8-35B-A3B-Distill-GGUF)) — the model powering this session, running on top of Unsloth Studio's inference engine.

## License

MIT
