package mcp

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Prompt definitions for reusable agent workflows.

// promptSummarizeChannel summarizes recent messages from a channel.
func promptSummarizeChannel() *mcp.Prompt {
	return &mcp.Prompt{
		Name:        "summarize_channel",
		Description: "Summarize the last N messages from a Slack channel.",
	}
}

// promptAlertOnKeyword watches for a keyword and notifies.
func promptAlertOnKeyword() *mcp.Prompt {
	return &mcp.Prompt{
		Name:        "alert_on_keyword",
		Description: "Watch a channel for a keyword and notify the user when found.",
	}
}

// promptDailyStatus sends a daily status report.
func promptDailyStatus() *mcp.Prompt {
	return &mcp.Prompt{
		Name:        "daily_status_report",
		Description: "Generate and post a daily status report to a channel.",
	}
}

// RegisterPrompts adds prompt definitions to the MCP server.
func RegisterPrompts(s *Server) {
	s.server.AddPrompt(promptSummarizeChannel(), makeSummarizeHandler())
	s.server.AddPrompt(promptAlertOnKeyword(), makeAlertHandler())
	s.server.AddPrompt(promptDailyStatus(), makeStatusHandler())
}

// ---- Prompt Handlers ----

func makeSummarizeHandler() mcp.PromptHandler {
	return func(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		args := req.Params.Arguments
		team := args["team"]
		channel := args["channel"]
		countStr := args["count"]

		if team == "" || channel == "" {
			return nil, fmt.Errorf("team and channel are required")
		}

		count := 10
		if n, err := fmt.Sscanf(countStr, "%d", &count); err != nil || n == 0 {
			count = 10
		}
		if count <= 0 {
			count = 10
		}

		return &mcp.GetPromptResult{
			Messages: []*mcp.PromptMessage{
				{
					Role: "user",
					Content: &mcp.TextContent{Text: fmt.Sprintf(
						"Please summarize the last %d messages from #%s in team %q.", count, channel, team,
					)},
				},
			},
		}, nil
	}
}

func makeAlertHandler() mcp.PromptHandler {
	return func(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		args := req.Params.Arguments
		team := args["team"]
		channel := args["channel"]
		keyword := args["keyword"]

		if team == "" || channel == "" || keyword == "" {
			return nil, fmt.Errorf("team, channel, and keyword are required")
		}

		return &mcp.GetPromptResult{
			Messages: []*mcp.PromptMessage{
				{
					Role: "user",
					Content: &mcp.TextContent{Text: fmt.Sprintf(
						"Watch for the keyword '%s' in #%s (team %q). When found, notify the user with a brief excerpt.",
						keyword, channel, team,
					)},
				},
			},
		}, nil
	}
}

func makeStatusHandler() mcp.PromptHandler {
	return func(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		args := req.Params.Arguments
		team := args["team"]
		channel := args["channel"]

		if team == "" || channel == "" {
			return nil, fmt.Errorf("team and channel are required")
		}

		return &mcp.GetPromptResult{
			Messages: []*mcp.PromptMessage{
				{
					Role: "user",
					Content: &mcp.TextContent{Text: fmt.Sprintf(
						"Generate a concise daily status report for team %q and post it to #%s.", team, channel,
					)},
				},
			},
		}, nil
	}
}
