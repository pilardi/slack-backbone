# Slack-Backbone — TODOs & Pending Work

*Date: 2026-09-29*

## 🎯 Active Todos (Unmerged)

| # | Task | Priority | Effort | Status |
|---|------|----------|--------|--------|

## 🆕 New Tasks

| # | Task | Priority | Effort | Status |
|---|------|----------|--------|--------|
| 22 | Update GitHub Actions CI to use `make` targets (build, test, lint, format, mod-tidy, docker-build) | P3 | Medium | 🔄 In Progress (PR pending review) |

## 📋 Completed & Merged (on main)

| PR | Title | Status |
|----|-------|--------|
| #22 | feat(ci): introduce Makefile build process (build, test, lint, format, docker) + docker-smoke | ✅ Done (PR #22) |
| #1 | fix(ci): update to Asafrose/bolt-go API and fix type mismatches | ✅ |
| #2 | chore(ci): update actions to Node.js 24-compatible versions + pin ubuntu-24.04 | ✅ |
| #3 | feat: implement real health check with Slack API connectivity verification | ✅ |
| #4 | feat: implement real status handler with --team flag support | ✅ |
| #5 | feat: wire up run() to register handlers and start Bolt apps (#5) | ✅ |
| #6 | feat: implement config.Load() + add teams.yaml.example | ✅ |
| #7 | feat: wire confirm button callbacks via app.Action() middleware | ✅ |
| #8 | docs: regenerate REVIEW.md with all completed work | ✅ |
| #9 | feat: add structured JSON logging via slog.NewJSONHandler | ✅ |
| #10 | docs: mark Dockerfile as done, remove from remaining tasks | ✅ |
| #11 | docs: add Acknowledgments section | ✅ |
| #18 | Add code coverage badge to README | ✅ Done (PR #18) |
| #20 | Add `.github/workflows/docker-smoke.yml` for real-token validation | ✅ Done (PR #20) |
| #12 | refactor: remove unused default_channel                | ✅   |
| #13 | feat(mcp): add MCP server for agent Slack communication| ✅   |       |                      |
| #14 | docs: convert REVIEW.md → TODOs.md with actionable items | ✅ |
| #15 | feat(mcp): wire slack_deploy to delegate to handlers.DeployHandler | ✅ |
| #16 | feat(handlers): implement real deploy logic with env validation + integration tests | ✅ |
| #8  | ci: replace explicit actions/cache with setup-go cache:true (simplified 4 cache steps) | ✅ Done (PR #21) |

## 🔍 Resolved Gaps (no longer applicable)

### ~~`handlers/deploy.go` — CLI Stub~~ ✅ Done (PR #16)
Now implements: env validation (`production/staging/development/qa`), `--channel` targeting, structured logging via `slog.InfoContext`, simulated deploy ID.

### ~~`mcp/tools.go::handleDeploy` — Separate Stub~~ ✅ Done (PR #15)
Now delegates to `handlers.DeployHandler{}.Run()` for consistency with CLI behavior.

### ~~`main.go` — Minimal entry point~~ ✅ Done
Already had graceful shutdown (`signal.NotifyContext`) and structured JSON logging (`slog.NewJSONHandler`). No further work needed here.

## 📊 Code Stats (as of 2026-09-28)

| Metric | Value |
|--------|-------|
| Total Go files | 14 (+ 5 test files across handlers, mcp, integration) |
| Lines of code (excl. tests) | ~970 |
| Unit test coverage (`handlers/`) | ~46% |
| Unit test coverage (`mcp/`) | ~72% |
| Combined coverage (all packages) | ~61% |
| Integration test files | 3 (`deploy_test.go`, `full_test.go`, plus handler tests) |
| Handlers remaining as stubs | **none** ✅ |

## 📝 Notes

- The `.env.example` file is kept for onboarding reference; tokens are primarily sourced from `teams.yaml`, but `.env`/`.env.local` are still read by godotenv.
- Config precedence: **CLI flags > env vars > .env.local > .env > teams.yaml > defaults**
- CI runs 7 jobs: `build`, `test`, `integration`, `coverage`, `lint`, `format`, `mod-tidy`, `docker-build`
- PR #17 (unmerged) adds the `integration` and `coverage` CI jobs
- Task 22 will replace hardcoded shell commands in `.github/workflows/ci.yml` with `make` targets, reducing duplication between local dev and CI
