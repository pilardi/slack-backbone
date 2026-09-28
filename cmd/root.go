package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Asafrose/bolt-go"
	"github.com/pilardi/slack-backbone/config"
	"github.com/pilardi/slack-backbone/handlers"
	slackpkg "github.com/pilardi/slack-backbone/slack"
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
	rootCmd.Flags().String("mode", "cli", "Operation mode: cli | mcp")
	rootCmd.Flags().Int("http-port", 0, "HTTP port for MCP streamable transport (MCP mode only)")

	config.BindFlags(rootCmd)
	rootCmd.AddCommand(mcpCmd)
}

func run(ctx context.Context) error {
	cfg, err := config.Load(rootCmd)
	if err != nil {
		return err
	}

	slog.Info("starting slack-backbone", "teams", len(cfg.Teams))

	// Create multi-team manager
	manager := slackpkg.NewManager()

	// For each team: create a client and register handlers on the Bolt app
	for _, team := range cfg.Teams {
		// Create API client for this team's bot token
		client := slackpkg.NewClient(team.BotToken)

		// Register the team's Bolt app with the manager (returns the app)
		app, err := manager.Register(ctx, team)
		if err != nil {
			return fmt.Errorf("failed to register team %q: %w", team.Name, err)
		}

		// Register slash commands on this team's app
		registerHandlers(app, ctx, client, team)
	}

	// Start all apps listening in background goroutines
	if err := manager.StartAll(ctx); err != nil {
		return fmt.Errorf("failed to start managers: %w", err)
	}

	slog.Info("all teams listening", "teams", len(cfg.Teams))
	// Block forever — StartAll is blocking per goroutine, but we want the main
	// process to stay alive. In practice, a real app would wait on a signal channel.
	select {}
}

// registerHandlers wires all handlers onto a Bolt app for the given team.
func registerHandlers(app *bolt.App, ctx context.Context, client *slackpkg.Client, team config.Team) {
	app.Command("health", adaptHandler(ctx, handlers.NewHealthHandler(client), team))
	app.Command("status", adaptHandler(ctx, handlers.NewStatusHandler(client), team))
	app.Command("deploy", adaptHandler(ctx, &handlers.DeployHandler{}, team))
	app.Command("confirm", adaptHandler(ctx, &handlers.ConfirmHandler{}, team))

	// Wire up button callbacks for confirm handler
	app.Action(bolt.ActionConstraints{ActionID: "confirmed"}, confirmCallback(team.Name, true))
	app.Action(bolt.ActionConstraints{ActionID: "cancelled"}, confirmCallback(team.Name, false))
}

// adaptHandler converts a Handler into Bolt's Command middleware signature.
func adaptHandler(ctx context.Context, h handlers.Handler, team config.Team) func(args bolt.SlackCommandMiddlewareArgs) error {
	return func(args bolt.SlackCommandMiddlewareArgs) error {
		// Parse text args (strip command prefix if present)
		text := strings.TrimPrefix(args.Command.Text, "/slack-backbone ")
		argsList := strings.Fields(text)

		blocks, err := h.Run(ctx, argsList, team)
		if err != nil {
			slog.Error("handler error", "handler", h.Name(), "team", team.Name, "error", err)
			return args.Ack(&bolt.CommandResponse{
				Text: fmt.Sprintf("❌ **Error:** %v", err),
			})
		}

		return args.Ack(&bolt.CommandResponse{
			Blocks: blocks.BlockSet,
		})
	}
}

// confirmCallback returns a Bolt action middleware handler for confirm buttons.
func confirmCallback(teamName string, confirmed bool) func(args bolt.SlackActionMiddlewareArgs) error {
	return func(args bolt.SlackActionMiddlewareArgs) error {
		status := "confirmed"
		if !confirmed {
			status = "cancelled"
		}
		text := fmt.Sprintf("✅ **%s!** · Team: **%s**", status, teamName)
		resp := interface{}(text)
		return args.Ack(&resp)
	}
}
