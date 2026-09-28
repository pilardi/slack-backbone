# Slack-Backbone — TODOs & Pending Work

*Date: 2026-09-27*

## 🎯 Active Todos (Unmerged)

| # | Task | Priority | Effort | Status |
|---|------|----------|--------|--------|
| 1 | Implement real deploy logic in `handlers/deploy.go` (replace static stub) | P1 | Small | ⬜ Open |
| 2 | Wire `slack_deploy` MCP tool to call the deploy handler | P1 | Small | ⬜ Open |
| 3 | Add `.env.example` documentation to README about env var precedence | P2 | Small | ✅ Done (PR #14) |

## 📋 Completed & Merged

These items were addressed in PRs and are now on `main`:

- ✅ PR #5 — Full wiring of `cmd/root.go::run()` (Manager per team, all handlers registered on Bolt apps, StartAll(), blocks forever)
- ✅ PR #6 — Config loading via viper + godotenv (`.env`/`.env.local` support)
- ✅ PR #7 — `handlers/confirm.go` with buttons + action callbacks
- ✅ PR #9 — Structured JSON logging via `slog.NewJSONHandler(os.Stdout, nil)`
- ✅ PR #10 — Multi-stage Dockerfile + `.dockerignore` + smoke test in CI
- ✅ PR #12 — Removed unused `default_channel`
- ✅ PR #13 — MCP server (6 tools, 3 resources, 3 prompts) with stdio + streamable HTTP transports

## 🔍 Notable Gaps

### 1. `handlers/deploy.go` — CLI Stub
Parses `--env <value>` but returns a static `"🚀 Deploying to **%s**..."` message. No real deployment logic yet.

**Suggested approach:**
- Accept an `env` argument (already parsed), log the intent
- Return a structured response indicating which environment would be targeted
- Optionally call a CI/CD webhook or API endpoint for real deployments

### 2. `mcp/tools.go::handleDeploy` — MCP Stub (separate from CLI)
The MCP tool `slack_deploy` has its own stub handler (`handleDeploy`) that returns a static message. It should ideally delegate to `handlers/deploy.go::Run()` for consistency with CLI behavior.

### 3. `main.go` — Already has graceful shutdown + structured logging ✅
Entry point already handles:
- Graceful shutdown via `signal.NotifyContext(SIGINT, SIGTERM)`
- Structured JSON logging via `slog.NewJSONHandler(os.Stdout, nil)`

**Remaining:** Could add a custom `slog.Handler` with contextual fields (team, version) for richer log output.

## 📊 Code Stats

| Metric | Value |
|--------|-------|
| Total Go files | 14 (+ 2 test files) |
| Lines of code (excl. tests) | ~914 |
| Test coverage (handlers) | 10/10 tests passing |
| Files with TODOs | none (resolved in PR #6) |
| Handlers remaining as stubs | deploy |

## 📝 Notes

- The `.env.example` file is kept for onboarding reference; tokens are primarily sourced from `teams.yaml`, but `.env`/`.env.local` are still read by godotenv.
- Config precedence: **CLI flags > env vars > .env.local > .env > teams.yaml > defaults**
