# Slack-Backbone — Remaining Work Review

*Date: 2026-09-27*

## ✅ Already Implemented (Merged)

| Component | Status | Notes |
|-----------|--------|-------|
| `handlers/health.go` | ✅ Complete | Real API calls: auth.test + conversations.list pagination |
| `handlers/status.go` | ✅ Complete | Real API calls, `--team` flag support |
| `slack/client.go` | ✅ Complete | HealthCheck, GetAccessibleChannels, PostMessage, PostEphemeral |
| `slack/manager.go` | ✅ Complete | Multi-team Bolt app lifecycle (Register, StartAll) |
| `handlers/interface.go` | ✅ Complete | Handler contract |
| `handlers/scope.go` | ✅ Complete | Global/Scoped scoping |
| `slack/blocks.go` | ✅ Complete | Block Kit helpers (NewBlocks, SectionBlock, DividerBlock, ActionBlockWithButton) |
| `.github/workflows/ci.yml` | ✅ Complete | 5 jobs: build, test, lint, format, mod-tidy |
| `cmd/root.go::run()` | ✅ **Done** (PR #5) | Full wiring: creates Manager per team, registers all handlers on Bolt apps, calls StartAll(), blocks forever |

## ⚠️ Stubbed / Incomplete

### 1. `config/config.go::Load()` — **Critical**
```go
func Load() (*Config, error) {
    v := viper.New()
    v.SetDefault("log_level", "info")
    v.AutomaticEnv()
    v.SetEnvPrefix("slack_backbone")
    v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
    _ = v.BindEnv("config", "CONFIG_FILE")
    return nil, nil  // ← always returns nil!
}
```
**What's needed:** Actually read the `--config` file via Viper (`v.ReadInConfig()`, `v.Unmarshal()`).

### 2. `handlers/deploy.go` — **Stub**
Returns a static `"🚀 Deploying to **production**..."` message. No actual deployment logic.

### 3. `handlers/confirm.go` — **Partial**
Has button UI (`✅ Confirm` / `❌ Cancel`) but no callback handler for button interactions. The buttons are rendered but never wired to a response action.

### 4. `main.go` — **Minimal**
Entry point delegates to `cmd.Execute()`. Could benefit from structured logging setup (custom handler with JSON output).

## 📋 Summary of Remaining Tasks

| Priority | Task | Effort |
|----------|------|--------|
| P0 | Fix `config.Load()` to actually parse the YAML config file | Small |
| P1 | Implement real deploy logic (or at least a more realistic response) | Small |
| P1 | Wire confirm button callbacks via Bolt's `ViewSubmission` middleware | Medium |
| P2 | Add structured JSON logging (slog handler with custom format) | Small |
| P3 | Add a `teams.yaml.example` for local testing | Small |
| P4 | Add Dockerfile / multi-stage build | Medium |

## 📊 Code Stats

- **Total Go files:** 15
- **Lines of code (excl. tests):** ~500
- **Test coverage (handlers):** 10/10 tests passing
- **Files with TODOs:** `config/config.go`
- **Handlers remaining as stubs:** deploy, confirm

## 🔍 Notable Gaps

1. **✅ Resolved:** Command dispatching now works — all 4 handlers are registered on each team's Bolt app via `app.Command()`.
2. **No `.env` file reading** — Viper's `v.ReadInConfig()` is missing, so the config file path from `--config` flag is never used.
3. **No logging configuration** — `slog.Default()` uses console output; no JSON structured logs for production use.
