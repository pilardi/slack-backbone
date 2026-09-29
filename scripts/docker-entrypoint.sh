#!/bin/sh
# slack-backbone Docker entrypoint
# Expands environment variables and execs the binary via tini for signal handling.

exec slack-backbone --mode "${MODE:-cli}" --http-port "${HTTP_PORT:-0}" "$@"
