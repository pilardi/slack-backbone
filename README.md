# slack-backbone

A multi-team Slack socket-mode application written in Go. Exposes slash commands and notifies channels about events — no public DNS or HTTP server required.

## Features

- **Socket mode** via `slack-go/bolt` — connects to Slack over WebSocket
- **Multi-team support** — one binary, N concurrent workspaces
- **Slash commands** with Block Kit rich layouts
- **Configurable per-team**: allowed/restricted commands, default channels
- **Structured logging** via `slog`
- **Go 1.24+** only

## Quick Start

```bash
# 1. Install dependencies
go mod tidy

# 2. Set up your tokens (see .env.example)
cp .env.example .env
# Edit .env with your actual bot and app tokens

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
├── .env.example      # Template for local dev
├── Dockerfile        # Container build
└── README.md
```

## Building & Running

```bash
# Local development
go run . --log-level debug

# Production binary
go build -o slack-backbone .
./slack-backbone --config teams.yaml --team prod-workspace

# Docker
docker build -t slack-backbone .
docker run -e SLACK_BOT_TOKEN=xoxb-... -e SLACK_APP_TOKEN=xapp-... slack-backbone
```

## License

MIT
