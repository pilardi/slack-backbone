package cmd

import (
	"context"
	"log/slog"

	"github.com/pilardi/slack-backbone/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "slack-backbone",
	Short: "Multi-team Slack socket-mode app",
	Long: `A Go application that connects to multiple Slack workspaces
via socket mode (Bolt), exposes slash commands, and notifies
channels about events.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return run(cmd.Context())
	},
}

func Execute(ctx context.Context, logger *slog.Logger) error {
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.Flags().String("config", "", "Path to teams.yaml config file")
	rootCmd.Flags().StringP("team", "t", "", "Target a specific team (default: all)")
	rootCmd.Flags().String("log-level", "info", "Log level: debug|info|warn|error")

	config.BindFlags(rootCmd)
}

func run(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	slog.Info("starting slack-backbone", "teams", len(cfg.Teams))

	// TODO: initialize multi-team manager and start listening
	_ = ctx

	return nil
}
