# Agent: Slack Backbone

## Identity

A Go-based multi-team Slack application that exposes slash commands, notifies channels about events, and provides an **MCP (Model Context Protocol) server** for agent integration. Runs in socket mode (no public URL needed).

## Capabilities

| Capability | Description |
|------------|-------------|
| Health checks | `/slack-backbone health` — verify all connected teams are responsive |
| Status reporting | `/slack-backbone status` — show running state per team |
| Deployments | `/slack-backbone deploy --env <env>` — trigger deployments (stub) |
| Confirmations | `/slack-backbone confirm` — interactive yes/no button flow |
| MCP Tools (6) | `slack_post`, `slack_notify`, `slack_status`, `slack_deploy`, `slack_confirm`, `slack_list_teams` |
| MCP Resources (3) | `channels/{team}`, `users/{team}`, `teams` — read-only data agents can query |
| MCP Prompts (3) | `summarize_channel`, `alert_on_keyword`, `daily_status_report` — reusable agent workflows |

## Configuration

Reads a `teams.yaml` file at startup:

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

## Commands & Scopes

| Command | Scope | Available On |
|---------|-------|--------------|
| `/slack-backbone health` | Global | All teams |
| `/slack-backbone status` | Global | All teams |
| `/slack-backbone deploy --env <env>` | Scoped | Teams where not restricted |
| `/slack-backbone confirm` | Scoped | Interactive button responses |

## MCP Tool Interface

Agents interact via the MCP server (stdio or streamable HTTP transport):

### Tools
- `slack_post` — Send a message to a channel or DM (`team`, `channel`, `text`)
- `slack_notify` — Send an ephemeral private message to a user (`team`, `user_id`, `text`)
- `slack_status` — Check health of all teams (optional `team` param)
- `slack_deploy` — Trigger a deployment (`team`, `env`)
- `slack_confirm` — Start interactive confirmation flow (`team`, `question`)
- `slack_list_teams` — List all configured workspaces

### Resources
- `channels/{team}` — Accessible public/private channels
- `users/{team}` — Slack user directory
- `teams` — Workspace metadata

### Prompts
- `summarize_channel` — Summarize last N messages from a channel
- `alert_on_keyword` — Watch a channel for a keyword, notify on match
- `daily_status_report` — Generate/post a daily status report

## Response Types

- **Immediate reply** — fast commands return a Block Kit message directly.
- **Ephemeral message** — sensitive info shown only to the triggering user.
- **Async callback** — long-running ops post a result later to the originating channel.

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `SLACK_BOT_TOKEN` | Yes | Bot token (`xoxb-...`) per team |
| `SLACK_APP_TOKEN` | Yes | App-level token (`xapp-...`) per team |
| `SLACK_BACKBONE_LOG_LEVEL` | No | `debug`, `info`, `warn`, `error` (default: `info`) |

## CLI Flags

| Flag | Description |
|------|-------------|
| `--config <path>` | Path to `teams.yaml` config file |
| `-t, --team <name>` | Target a specific team (default: all) |
| `--log-level <level>` | Override log verbosity (`debug`, `info`, `warn`, `error`) |
| `--mode <cli\|mcp>` | Operation mode (default: `cli`) |
| `--http-port <n>` | HTTP port for MCP streamable transport (MCP mode only) |

## Build & Run

```bash
go mod tidy
go run . --config teams.yaml
# or
go build -o slack-backbone . && ./slack-backbone --config teams.yaml

# MCP mode via stdio
./slack-backbone --mode mcp --config teams.yaml

# MCP mode via streamable HTTP
./slack-backbone --mode mcp --http-port 8080 --config teams.yaml
```

## Docker

```bash
docker build -t slack-backbone .
docker run -e SLACK_BOT_TOKEN=xoxb-... -e SLACK_APP_TOKEN=xapp-... slack-backbone
# MCP mode via env vars:
docker run -e MODE=mcp -e HTTP_PORT=8080 -e SLACK_BOT_TOKEN=xoxb-... slack-backbone
```
