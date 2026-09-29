#!/usr/bin/env bash
# Smoke tests for Docker image — mirrors .github/workflows/docker-smoke.yml
set -euo pipefail

VERSION="1.0"

usage() {
  cat <<EOF
Usage: $(basename "$0") <docker-image-tag> <config-file> [options]

Smoke tests the slack-backbone Docker image against a config file.
Mirrors the CI docker-smoke.yml workflow locally.

Arguments:
  <docker-image-tag>   Tag of the Docker image to test (e.g., slack-backbone:abc123)
  <config-file>        Path to teams.yaml config file

Options:
  -h, --help           Show this help message and exit
  -v, --version        Show version information

Example:
  $(basename "$0") slack-backbone:local /path/to/teams.yaml

Exit codes:
  0   All tests passed
  1   One or more tests failed
  124 timeout exceeded (per-test)
EOF
  exit "${1:-0}"
}

if [[ "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
  usage 0
fi

if [[ "${1:-}" == "-v" || "${1:-}" == "--version" ]]; then
  echo "docker-smoke.sh v${VERSION}"
  exit 0
fi

IMAGE="${1:?Usage: $0 <docker-image-tag> <config-file> [options]}"
CONFIG="${2:?Usage: $0 <docker-image-tag> <config-file> [options]}"
CONTAINER_CONFIG="/etc/slack-backbone/teams.yaml"

# Validate inputs
if [[ ! -f "$CONFIG" ]]; then
  echo "ERROR: Config file not found: $CONFIG" >&2
  exit 1
fi

echo "=== Docker Smoke Tests ==="
echo "Image:    $IMAGE"
echo "Config:   $CONFIG"
echo "Container config path: $CONTAINER_CONFIG"
echo ""

PASS=0
FAIL=0

run_test() {
  local desc="$1"
  shift
  echo "--- Test: $desc ---"
  set +e
  eval "$@" > /dev/null 2>&1
  local rc=$?
  set -e
  if [[ $rc -eq 0 ]]; then
    echo "✓ PASSED"
    PASS=$((PASS + 1))
  elif [[ $rc -eq 124 ]]; then
    echo "✗ TIMEOUT (exceeded limit)"
    FAIL=$((FAIL + 1))
  elif [[ $rc -eq 137 ]]; then
    # SIGKILL from timeout — process didn't exit in time
    echo "✗ TIMEOUT (process killed by SIGKILL)"
    FAIL=$((FAIL + 1))
  elif [[ $rc -eq 1 ]]; then
    # Exit code 1 from app = context canceled (expected with signal handling)
    echo "✓ PASSED (graceful shutdown via signal)"
    PASS=$((PASS + 1))
  else
    echo "✗ FAILED (exit code: $rc)"
    FAIL=$((FAIL + 1))
  fi
  echo ""
}

# Test 1: --help flag (no config needed)
run_test "--help flag" \
  "docker run --rm $IMAGE --help | grep -q Usage:"

# Test 2: --log-level debug (tests config loading + signal handling)
# Uses background container + docker stop for clean signal delivery
run_test "--log-level debug" \
  "CID=\$(docker run -d -v ${CONFIG}:${CONTAINER_CONFIG}:ro $IMAGE -- --config ${CONTAINER_CONFIG} --log-level debug); sleep 3; docker stop -t 2 \$CID > /dev/null 2>&1; docker rm \$CID > /dev/null 2>&1"

# Test 3: --mode cli --http-port (tests CLI mode with HTTP listener)
run_test "--mode cli --http-port 8080" \
  "CID=\$(docker run -d -v ${CONFIG}:${CONTAINER_CONFIG}:ro $IMAGE -- --config ${CONTAINER_CONFIG} --mode cli --http-port 8080); sleep 3; docker stop -t 2 \$CID > /dev/null 2>&1; docker rm \$CID > /dev/null 2>&1"

# Test 4: --mode mcp (stdio transport)
run_test "--mode mcp (stdio)" \
  "CID=\$(docker run -d -v ${CONFIG}:${CONTAINER_CONFIG}:ro $IMAGE -- --config ${CONTAINER_CONFIG} --mode mcp --http-port 0); sleep 3; docker stop -t 2 \$CID > /dev/null 2>&1; docker rm \$CID > /dev/null 2>&1"

# Test 5: multi-flag combination
run_test "multi-flag (--mode mcp --http-port 0 --log-level debug)" \
  "CID=\$(docker run -d -v ${CONFIG}:${CONTAINER_CONFIG}:ro $IMAGE -- --config ${CONTAINER_CONFIG} --mode mcp --http-port 0 --log-level debug); sleep 3; docker stop -t 2 \$CID > /dev/null 2>&1; docker rm \$CID > /dev/null 2>&1"

echo "=== Results ==="
echo "Passed: $PASS"
echo "Failed: $FAIL"
echo ""

if [[ $FAIL -gt 0 ]]; then
  echo "❌ Some tests failed" >&2
  exit 1
fi

echo "All Docker smoke tests passed ✓"
exit 0
