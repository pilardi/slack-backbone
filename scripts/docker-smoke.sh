#!/usr/bin/env bash
# Smoke tests for Docker image — mirrors .github/workflows/docker-smoke.yml
set -euo pipefail

IMAGE="${1:?Usage: $0 <docker-image-tag>}"
CONFIG="${2:?Usage: $0 <image-tag> <config-file>}"
CONTAINER_CONFIG="/etc/slack-backbone/teams.yaml"

echo "--- Test 1: --help flag ---"
timeout 5s docker run --rm "$IMAGE" --help | grep -q "Usage:"
echo "✓ --help works"

echo ""
echo "--- Test 2: --log-level debug ---"
timeout 5s docker run --rm -v "${CONFIG}:${CONTAINER_CONFIG}:ro" "$IMAGE" -- --config "${CONTAINER_CONFIG}" --log-level debug || [ $? -eq 124 ]
echo "✓ --log-level debug works"

echo ""
echo "--- Test 3: --mode cli --http-port ---"
timeout 10s docker run --rm -v "${CONFIG}:${CONTAINER_CONFIG}:ro" "$IMAGE" -- --config "${CONTAINER_CONFIG}" --mode cli --http-port 8080 || [ $? -eq 124 ]
echo "✓ --mode cli works"

echo ""
echo "--- Test 4: --mode mcp (stdio transport) ---"
timeout 5s docker run --rm -v "${CONFIG}:${CONTAINER_CONFIG}:ro" "$IMAGE" -- --config "${CONTAINER_CONFIG}" --mode mcp --http-port 0 || [ $? -eq 124 ]
echo "✓ --mode mcp works"

echo ""
echo "--- Test 5: multi-flag combination ---"
timeout 5s docker run --rm -v "${CONFIG}:${CONTAINER_CONFIG}:ro" "$IMAGE" -- --config "${CONTAINER_CONFIG}" --mode mcp --http-port 0 --log-level debug || [ $? -eq 124 ]
echo "✓ multi-flag combination works"

echo ""
echo "All Docker smoke tests passed ✓"
