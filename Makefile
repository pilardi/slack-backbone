# Slack-Backbone — Local Build & CI Replication
#
# Run `make help` to see available targets.
# All targets can be run locally to mirror what happens in .github/workflows/ci.yml

GO ?= go
COVER ?= unit-coverage.out
COMBINED_COVER ?= combined.out

.PHONY: help all build test integration coverage lint format fmt-fix mod-tidy docker-build clean

help: ## Show available targets and their descriptions
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | \
		awk 'BEGIN{FS=":.*?## "}{printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

all: build test integration coverage lint format mod-tidy ## Run all CI-like checks

# ─── Build ───────────────────────────────────────────────────────────────────

build: ## Build all packages (non-test code)
	$(GO) mod download
	$(GO) build ./...

# ─── Tests ──────────────────────────────────────────────────────────────────

test: clean ## Run unit tests with coverage profile
	@echo "=== Running unit tests ==="
	$(GO) test -v -coverprofile=$(COVER) ./handlers/... ./mcp/...
	@echo ""
	@echo "=== Unit Test Coverage ==="
	$(GO) tool cover -func=$(COVER) | tail -1

integration: ## Run integration tests
	@echo "=== Running integration tests ==="
	$(GO) test -v ./integration/...

coverage: ## Run all tests with combined coverage profile
	@echo "=== Running all tests (combined coverage) ==="
	$(GO) test -v -coverprofile=$(COMBINED_COVER) ./handlers/... ./mcp/... ./integration/...
	@echo ""
	@echo "=== Combined Test Coverage ==="
	$(GO) tool cover -func=$(COMBINED_COVER) | tail -5

# ─── Linting ────────────────────────────────────────────────────────────────

lint: ## Run golangci-lint with --timeout=5m
	@echo "=== Running linter ==="
	$(GO) run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.14.0 run --timeout=5m

# ─── Formatting ──────────────────────────────────────────────────────────────

format: ## Check go fmt drift (fails if any file needs formatting)
	@echo "=== Checking formatting ==="
	go fmt ./...
	if [ -n "$$(git status --porcelain)" ]; then \
		echo "Files are not properly formatted:"; \
		git diff; \
		exit 1; \
	fi
	@echo "gofmt: clean ✓"

fmt-fix: ## Reformat every file in place with go fmt
	@echo "=== Fixing formatting ==="
	$(GO) fmt ./...

# ─── Module ──────────────────────────────────────────────────────────────────

mod-tidy: ## Verify go.mod/go.sum are tidy (fails if diff detected)
	@echo "=== Checking go mod tidy ==="
	$(GO) mod tidy
	if [ -n "$$(git status --porcelain)" ]; then \
		echo "go.mod/go.sum out of sync:"; \
		git diff; \
		exit 1; \
	fi

# ─── Docker ──────────────────────────────────────────────────────────────────

docker-build: ## Build Docker image (no push, load to local daemon)
	@echo "=== Building Docker image ==="
	docker build -t slack-backbone:$(shell git rev-parse HEAD 2>/dev/null || echo local) .
	@echo ""
	@echo "Image tagged as slack-backbone:$(shell git rev-parse HEAD 2>/dev/null || echo local)"

docker-smoke: docker-build ## Run runtime smoke tests against the Docker image
	@echo "=== Running Docker smoke tests ==="
	IMAGE="slack-backbone:$(shell git rev-parse HEAD 2>/dev/null || echo local)" && \
	bash $(abspath scripts/docker-smoke.sh) "$$IMAGE" "$(abspath smoke-teams.yaml)"



# ─── Cleanup ────────────────────────────────────────────────────────────────

clean: ## Remove coverage profiles
	@rm -f $(COVER) $(COMBINED_COVER) 2>/dev/null || true
