package cmd

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/pilardi/slack-backbone/config"
	"github.com/pilardi/slack-backbone/mcp"
	"github.com/spf13/cobra"
)

// mcpCmd is the subcommand for running in MCP mode.
var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Run as an MCP server",
	Long: `Start slack-backbone as an MCP (Model Context Protocol) server.
Agents can connect via stdio or streamable HTTP to invoke Slack tools,
query resources, and use prompt templates.

Usage:
  slack-backbone --mode mcp [--config teams.yaml] [--http-port 8080]`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runMCP(cmd.Context())
	},
}

func init() {
	mcpCmd.Flags().String("config", "", "Path to teams.yaml config file")
	mcpCmd.Flags().Int("http-port", 0, "HTTP port for streamable transport (default: stdio only)")
}

func runMCP(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	logger := slog.Default()

	srv, err := mcp.NewServer(cfg, logger)
	if err != nil {
		return fmt.Errorf("failed to create MCP server: %w", err)
	}

	srv.InitClients()
	mcp.RegisterTools(srv)

	httpPort, _ := rootCmd.Flags().GetInt("http-port")

	if httpPort > 0 {
		logger.Info("starting MCP server over HTTP", "port", httpPort)
		return srv.StartHTTP(ctx, httpPort)
	}

	logger.Info("starting MCP server over stdio")
	return srv.Start(ctx)
}
