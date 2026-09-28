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

### 1. `handlers/deploy.go` — Stub
Returns a static `"🚀 Deploying to **production**..."` message. No real deployment logic yet.

**Suggested approach:** Accept an `env` argument, log the intent, and return a structured response indicating which environment would be targeted (e.g., `"deploying to staging"`). A full integration could call a CI/CD webhook or API endpoint.

### 2. `main.go` — Minimal
Entry point delegates to `cmd.Execute()`. Could benefit from:
- Graceful shutdown handling (`SIGINT`/`SIGTERM`)
- Structured logging setup with custom handler (JSON output)

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
