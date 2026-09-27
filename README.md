# slack-backbone

A multi-team Slack socket-mode application written in Go. Exposes slash commands, notifies channels about events, and exposes an **MCP (Model Context Protocol) server** so agents can programmatically post messages, send notifications, check health, trigger deployments, and more — all without a public DNS or HTTP server required.

## Features

- **Socket mode** via `Asafrose/bolt-go` — connects to Slack over WebSocket
- **Multi-team support** — one binary, N concurrent workspaces
- **Slash commands** with Block Kit rich layouts
- **MCP server** (6 tools + 3 resources + 3 prompts) for agent integration
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

# 3. Run locally (CLI mode — slash commands)
go run . --config teams.yaml

# 4. Build a binary
go build -o slack-backbone .
```

## MCP Mode (Agent Integration)

Run the app as an MCP server so external agents can call Slack operations:

```bash
# Stdio transport (default, for local agents)
slack-backbone --mode mcp --config teams.yaml

# Streamable HTTP transport (for remote clients)
slack-backbone --mode mcp --http-port 8080 --config teams.yaml
```

### MCP Tools (6)

| Tool | Description |
|------|-------------|
| `slack_post` | Send a message to a Slack channel or DM a user (`team`, `channel`, `text`) |
| `slack_notify` | Send an ephemeral private message to a specific user (`team`, `user_id`, `text`) |
| `slack_status` | Check health / running state of all teams (optional `team` param) |
| `slack_deploy` | Trigger a deployment to a target environment (`team`, `env`) |
| `slack_confirm` | Start an interactive yes/no confirmation flow (`team`, `question`) |
| `slack_list_teams` | List all configured workspaces |

### MCP Resources (3)

| Resource | Description |
|----------|-------------|
| `channels/{team}` | List of accessible public and private channels for a team |
| `users/{team}` | Directory of Slack users for a team |
| `teams` | Metadata about all configured Slack workspaces |

### MCP Prompts (3)

| Prompt | Description |
|--------|-------------|
| `summarize_channel` | Summarize the last N messages from a Slack channel |
| `alert_on_keyword` | Watch a channel for a keyword and notify when found |
| `daily_status_report` | Generate and post a daily status report to a channel |

## Configuration

Create a `teams.yaml` file:

```yaml
teams:
  - name: "prod-workspace"
    bot_token: "xoxb-..."
    app_token: "xapp-..."
    commands: all
    restricted_commands: []

  - name: "staging-workspace"
    bot_token: "xoxb-..."
    app_token: "xapp-..."
    commands: all
    restricted_commands: [status]
```

## CLI Commands

| Command | Scope | Description |
|---------|-------|-------------|
| `/slack-backbone health` | Global | Health check across all teams |
| `/slack-backbone status` | Global | Show running status per team |
| `/slack-backbone deploy --env <env>` | Scoped | Trigger a deployment |
| `/slack-backbone confirm` | Scoped | Interactive confirmation flow |

## CLI Flags

| Flag | Description |
|------|-------------|
| `--config <path>` | Path to `teams.yaml` config file |
| `-t, --team <name>` | Target a specific team (default: all) |
| `--log-level <level>` | Log level: `debug`, `info`, `warn`, `error` (default: `info`) |
| `--mode <cli\|mcp>` | Operation mode (default: `cli`) |
| `--http-port <n>` | HTTP port for MCP streamable transport (MCP mode only; default: stdio) |

## Project Structure

```
├── main.go           # Entry point + signal handling
├── cmd/              # Cobra CLI commands (root, mcp subcommand)
├── handlers/         # Command handler implementations
├── mcp/              # MCP server: tools, resources, prompts, tests
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
# Local development (CLI mode)
go run . --log-level debug

# Production binary (CLI mode)
go build -o slack-backbone .
./slack-backbone --config teams.yaml --team prod-workspace

# MCP mode via stdio
./slack-backbone --mode mcp --config teams.yaml

# MCP mode via streamable HTTP
./slack-backbone --mode mcp --http-port 8080 --config teams.yaml

# Docker (CI builds and tags with commit SHA)
docker build -t slack-backbone:${{ github.sha }} .
docker run --rm slack-backbone:${{ github.sha }} --log-level debug
```

## Acknowledgments

This project was built with the help of several tools and services:

- **[Berd](https://github.com/block/berd/)** — the desktop app where this conversation took place, providing a rich interface for working with agents.
- **[Goose](https://github.com/goose)**, the coding agent used to design, implement, and review all changes in this project.
- **Unsloth Studio** ([unslothai/unsloth](https://github.com/unslothai/unsloth)) — the inference engine powering the LLM used during development.
- **llama.cpp** ([ggml-org/llama.cpp](https://github.com/ggml-org/llama.cpp)) — the inference engine that powers Unsloth Studio's runtime.
- **Qwen-3.8-35B-A3B** ([Empero AI](https://huggingface.co/empero-ai/Qwen3.8-35B-A3B-Distill-GGUF)) — the model powering this session, running on top of Unsloth Studio's inference engine.

## License

This project is licensed under the [MIT License](LICENSE).
