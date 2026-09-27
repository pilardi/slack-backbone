package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/pilardi/slack-backbone/cmd"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger := slog.Default()

	if err := cmd.Execute(ctx, logger); err != nil {
		slog.Error("shutting down", "error", err)
		os.Exit(1)
	}
}
