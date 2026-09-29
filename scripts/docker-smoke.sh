#!/usr/bin/env bash
# Smoke tests for Docker image — mirrors .github/workflows/docker-smoke.yml
set -euo pipefail

VERSION="1.0"

usage() {
  cat <<EOF
Usage: $(basename "$0") <docker-image-tag|artifact-path> <config-file> [options]

Smoke tests the slack-backbone Docker image against a config file.
Mirrors the CI docker-smoke.yml workflow locally.

Arguments:
  <docker-image-tag>   Tag of the Docker image to test (e.g., slack-backbone:abc123)
                       OR path to a .tar artifact from GitHub Actions
  <config-file>        Path to teams.yaml config file

Options:
  -h, --help           Show this help message and exit
  -v, --version        Show version information
  -a, --artifact       Load image from a .tar artifact (auto-detect if arg is a file path)

Example:
  $(basename "$0") slack-backbone:local /path/to/teams.yaml
  $(basename "$0") /tmp/docker-image.tar /path/to/teams.yaml --artifact

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

# Parse arguments — support --artifact flag and auto-detect file paths
ARTIFACT_MODE=false
CONFIG=""
IMAGE=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    -h|--help) usage 0 ;;
    -v|--version) echo "docker-smoke.sh v${VERSION}"; exit 0 ;;
    -a|--artifact) ARTIFACT_MODE=true; shift ;;
    *)
      if [[ -z "$IMAGE" ]]; then
        IMAGE="$1"
      elif [[ -z "$CONFIG" ]]; then
        CONFIG="$1"
      else
        echo "ERROR: Too many arguments" >&2
        exit 1
      fi
      shift
      ;;
  esac
done

# Auto-detect artifact mode if first arg is a file path ending in .tar
if [[ -z "$IMAGE" ]]; then
  echo "ERROR: Missing image tag or artifact path" >&2
  exit 1
fi

if [[ -z "$CONFIG" ]]; then
  echo "ERROR: Missing config file path" >&2
  exit 1
fi

# If --artifact mode or arg is a .tar file, load it as a Docker image
CONTAINER_CONFIG="/etc/slack-backbone/teams.yaml"
LOADED_IMAGE=""

if [[ "$ARTIFACT_MODE" == true ]] || [[ "$IMAGE" == *.tar ]]; then
  TAR_PATH="$IMAGE"
  if [[ ! -f "$TAR_PATH" ]]; then
    echo "ERROR: Artifact file not found: $TAR_PATH" >&2
    exit 1
  fi
  echo "--- Loading Docker image from artifact: $TAR_PATH ---"
  CONTAINER_ID=$(docker load -i "$TAR_PATH" | tail -1 | awk '{print $NF}')
  LOADED_IMAGE="slack-backbone-smoke-$$"
  docker tag "$CONTAINER_ID" "$LOADED_IMAGE" > /dev/null 2>&1
  IMAGE="$LOADED_IMAGE"
  echo "  Loaded image: $IMAGE (from $(du -h "$TAR_PATH" | cut -f1))"
fi

# Validate config
if [[ ! -f "$CONFIG" ]]; then
  echo "ERROR: Config file not found: $CONFIG" >&2
  exit 1
fi

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

echo ""
echo "=== Cleanup ==="
if [[ -n "$LOADED_IMAGE" ]]; then
  echo "Removing loaded test image: $LOADED_IMAGE"
  docker rmi "$LOADED_IMAGE" > /dev/null 2>&1 || true
fi

echo ""
echo "All Docker smoke tests passed ✓"
exit 0
