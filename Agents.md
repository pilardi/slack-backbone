# Agent: Slack Backbone

## Identity
A Go-based multi-team Slack application that exposes slash commands and notifies channels about events. Runs in socket mode (no public URL needed).

## Capabilities

| Capability | Description |
|------------|-------------|
| Health checks | `/slack-backbone health` — verify all connected teams are responsive |
| Status reporting | `/slack-backbone status` — show running state per team |
| Deployments | `/slack-backbone deploy --env <env>` — trigger deployments (stub) |
| Confirmations | `/slack-backbone confirm` — interactive yes/no button flow |

## Configuration

Reads a `teams.yaml` file at startup:

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

## Commands & Scopes

| Command | Scope | Available On |
|---------|-------|--------------|
| `/slack-backbone health` | Global | All teams |
| `/slack-backbone status` | Global | All teams |
| `/slack-backbone deploy --env <env>` | Scoped | Teams where not restricted |
| `/slack-backbone confirm` | Scoped | Interactive button responses |

## Response Types

- **Immediate reply** — fast commands return a Block Kit message directly.
- **Ephemeral message** — sensitive info shown only to the triggering user.
- **Async callback** — long-running ops post a result later to `default_channel`.

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `SLACK_BOT_TOKEN` | Yes | Bot token (`xoxb-...`) per team |
| `SLACK_APP_TOKEN` | Yes | App-level token (`xapp-...`) per team |
| `SLACK_BACKBONE_LOG_LEVEL` | No | `debug`, `info`, `warn`, `error` (default: `info`) |

## CLI Flags

| Flag | Description |
|------|-------------|
| `--config <path>` | Path to `teams.yaml` |
| `-t, --team <name>` | Target a specific team (default: all) |
| `--log-level <level>` | Override log verbosity |

## Build & Run

```bash
go mod tidy
go run . --config teams.yaml
# or
go build -o slack-backbone . && ./slack-backbone --config teams.yaml
```

## Docker

```bash
docker build -t slack-backbone .
docker run -e SLACK_BOT_TOKEN=xoxb-... -e SLACK_APP_TOKEN=xapp-... slack-backbone
```
