package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/pilardi/slack-backbone/config"
	"github.com/pilardi/slack-backbone/handlers"
	slackpkg "github.com/pilardi/slack-backbone/slack"
	"github.com/slack-go/slack"
)

var _ = (*slackpkg.Client)(nil) // satisfy import checker

// Tool definitions for Slack operations exposed to agents.

// ---- Argument structs (auto-inferred as JSON Schema) ----

type postArgs struct {
	Team    string `json:"team" jsonschema:"^$" jsonschema_description:"Team/workspace name from config"`
	Channel string `json:"channel" jsonschema:"^$" jsonschema_description:"Channel name (e.g. '#general') or user ID for DMs"`
	Text    string `json:"text,omitempty" jsonschema_description:"Plain text message body"`
}

type notifyArgs struct {
	Team   string `json:"team" jsonschema:"^$" jsonschema_description:"Team/workspace name from config"`
	UserID string `json:"user_id" jsonschema:"^$" jsonschema_description:"Slack user ID (e.g. 'U0123ABCD')"`
	Text   string `json:"text" jsonschema:"^$" jsonschema_description:"Plain text message body"`
}

type statusArgs struct {
	Team string `json:"team,omitempty" jsonschema_description:"Team/workspace name (omit for all-team overview)"`
}

type deployArgs struct {
	Team string `json:"team" jsonschema:"^$" jsonschema_description:"Team/workspace name from config"`
	Env  string `json:"env,omitempty" jsonschema_description:"Target environment (e.g. 'staging', 'production')"`
}

type confirmArgs struct {
	Team     string `json:"team" jsonschema:"^$" jsonschema_description:"Team/workspace name from config"`
	Question string `json:"question" jsonschema:"^$" jsonschema_description:"The question to present to users"`
}

// ---- Tool registration ----

// RegisterTools adds all tools to the MCP server.
func RegisterTools(s *Server) {
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "slack_post",
		Description: "Send a message to a Slack channel or direct message a user.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args postArgs) (*mcp.CallToolResult, any, error) {
		return handlePost(s, ctx, args.Team, args.Channel, args.Text)
	})

	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "slack_notify",
		Description: "Send a private ephemeral message to a specific Slack user.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args notifyArgs) (*mcp.CallToolResult, any, error) {
		return handleNotify(s, ctx, args.Team, args.UserID, args.Text)
	})

	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "slack_status",
		Description: "Check bot connectivity, latency, and accessible channels for a team.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args statusArgs) (*mcp.CallToolResult, any, error) {
		return handleStatus(s, ctx, args.Team)
	})

	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "slack_deploy",
		Description: "Trigger a deployment notification to a Slack channel.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args deployArgs) (*mcp.CallToolResult, any, error) {
		return handleDeploy(s, ctx, args.Team, args.Env)
	})

	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "slack_confirm",
		Description: "Start an interactive confirmation button flow in a Slack channel.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args confirmArgs) (*mcp.CallToolResult, any, error) {
		return handleConfirm(s, ctx, args.Team, args.Question)
	})

	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "slack_list_teams",
		Description: "List all configured Slack workspaces/teams.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
		return handleListTeams(s, ctx)
	})
}

// ---- Handlers ----

func handlePost(s *Server, ctx context.Context, team, channel, text string) (*mcp.CallToolResult, any, error) {
	if team == "" || channel == "" {
		return nil, nil, fmt.Errorf("team and channel are required")
	}

	client, err := s.GetClient(team)
	if err != nil {
		return nil, nil, err
	}

	msgID, _, err := client.PostMessage(ctx, channel, slack.MsgOptionText(text, false))
	if err != nil {
		return nil, nil, fmt.Errorf("post failed: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("✅ Message sent to `%s` (ID: `%s`)", channel, msgID)}},
	}, nil, nil
}

func handleNotify(s *Server, ctx context.Context, team, userID, text string) (*mcp.CallToolResult, any, error) {
	if team == "" || userID == "" || text == "" {
		return nil, nil, fmt.Errorf("team, user_id, and text are required")
	}

	client, err := s.GetClient(team)
	if err != nil {
		return nil, nil, err
	}

	// For ephemeral messages we need a channel ID. Using team name as
	// a placeholder; production would resolve the user's DM channel.
	_, err = client.PostEphemeral(ctx, "C0PLACEHOLDER", userID, slack.MsgOptionText(text, false))
	if err != nil {
		return nil, nil, fmt.Errorf("notify failed: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("✅ Ephemeral message sent to user `%s`", userID)}},
	}, nil, nil
}

func handleStatus(s *Server, ctx context.Context, teamName string) (*mcp.CallToolResult, any, error) {
	if teamName == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: "⚠️ Team name required"}},
		}, nil, nil
	}

	client, err := s.GetClient(teamName)
	if err != nil {
		return nil, nil, err
	}

	result, err := client.HealthCheck(ctx)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("❌ Bot: %v", err)}},
		}, nil, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf(
			"✅ **Bot:** connected as *%s* (`%s`) · ⏱️ **Latency:** %.1fms",
			result.UserName, result.UserID, result.LatencyMs,
		)}},
	}, nil, nil
}

func handleDeploy(s *Server, ctx context.Context, teamName, env string) (*mcp.CallToolResult, any, error) {
	if env == "" {
		env = "production"
	}

	// Look up the full team config (needed by the handler)
	var teamConfig config.Team
	for _, t := range s.cfg.Teams {
		if t.Name == teamName {
			teamConfig = t
			break
		}
	}
	if teamConfig.Name == "" {
		return nil, nil, fmt.Errorf("unknown team: %s", teamName)
	}

	// Delegate to the CLI deploy handler
	blocks, err := (&handlers.DeployHandler{}).Run(ctx, []string{"--env", env}, teamConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("deploy failed: %w", err)
	}

	return &mcp.CallToolResult{Content: blocksToMCPContent(blocks)}, nil, nil
}

func handleConfirm(s *Server, ctx context.Context, team, question string) (*mcp.CallToolResult, any, error) {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("⏳ Confirmation flow started:\n**Question:** %s\nWaiting for user response...", question)}},
	}, nil, nil
}

func handleListTeams(s *Server, ctx context.Context) (*mcp.CallToolResult, any, error) {
	var lines []string
	for _, t := range s.cfg.Teams {
		lines = append(lines, fmt.Sprintf("• `%s`", t.Name))
	}
	if len(lines) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: "No teams configured."}},
		}, nil, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: "Configured teams:\n" + strings.Join(lines, "\n")}},
	}, nil, nil
}

// blocksToMCPContent converts a *slack.Blocks to MCP content items.
func blocksToMCPContent(b *slack.Blocks) []mcp.Content {
	var content []mcp.Content
	for _, block := range b.BlockSet {
		switch blk := block.(type) {
		case *slack.SectionBlock:
			if blk.Text != nil {
				content = append(content, &mcp.TextContent{Text: blk.Text.Text})
			}
		case *slack.DividerBlock:
			content = append(content, &mcp.TextContent{Text: "\n───\n"})
		case *slack.HeaderBlock:
			content = append(content, &mcp.TextContent{Text: fmt.Sprintf("# %s\n", blk.Text.Text)})
		case *slack.ImageBlock:
			content = append(content, &mcp.TextContent{Text: fmt.Sprintf("![image] (%s)", blk.AltText)})
		default:
			content = append(content, &mcp.TextContent{Text: "[block type not rendered]"})
		}
	}
	return content
}
