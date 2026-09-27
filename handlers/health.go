package handlers

import (
	"context"
	"fmt"

	"github.com/pilardi/slack-backbone/config"
	slackpkg "github.com/pilardi/slack-backbone/slack"
	"github.com/slack-go/slack"
)

// HealthHandler handles the /slack-backbone health command.
type HealthHandler struct {
	client *slackpkg.Client
}

// NewHealthHandler creates a new health handler with the given Slack client.
func NewHealthHandler(client *slackpkg.Client) *HealthHandler {
	return &HealthHandler{client: client}
}

func (h *HealthHandler) Name() string { return "health" }

func (h *HealthHandler) Scope() Scope { return Global() }

func (h *HealthHandler) Run(ctx context.Context, args []string, team config.Team) (*slack.Blocks, error) {
	checks := make([]string, 0, 3)

	if h.client != nil {
		// Check 1: Bot info (connectivity + token validity)
		result, err := h.client.HealthCheck(ctx)
		if err != nil {
			checks = append(checks, fmt.Sprintf("❌ **Bot Info:** %v", err))
		} else {
			checks = append(checks, fmt.Sprintf("✅ **Bot Info:** connected as *%s* (user: `%s`)", result.UserName, result.UserID))
			checks = append(checks, fmt.Sprintf("⏱️ **Latency:** %.1fms", result.LatencyMs))
		}

		// Check 2: Channel access
		channels, err := h.client.GetAccessibleChannels(ctx)
		if err != nil {
			checks = append(checks, fmt.Sprintf("❌ **Channel Access:** %v", err))
		} else {
			checks = append(checks, fmt.Sprintf("✅ **Channels:** %d accessible", len(channels)))
			for _, ch := range channels[:min(5, len(channels))] {
				checks = append(checks, fmt.Sprintf("   • `%s`", ch.Name))
			}
			if len(channels) > 5 {
				checks = append(checks, fmt.Sprintf("   • ... and %d more", len(channels)-5))
			}
		}
	} else {
		checks = append(checks, "⚠️ **No client configured** — skipping API checks")
	}

	// Build response blocks using slack-go/slack directly
	blocks := slackpkg.NewBlocks()
	blocks.BlockSet = append(blocks.BlockSet,
		slack.NewSectionBlock(&slack.TextBlockObject{Type: "mrkdwn", Text: fmt.Sprintf("🩺 Health Check · Team: **%s**", team.Name)}, nil, nil),
	)

	if len(checks) > 0 {
		blocks.BlockSet = append(blocks.BlockSet,
			slack.NewSectionBlock(&slack.TextBlockObject{Type: "mrkdwn", Text: checks[0]}, nil, nil),
		)

		for _, line := range checks[1:] {
			blocks.BlockSet = append(blocks.BlockSet,
				slack.NewSectionBlock(&slack.TextBlockObject{Type: "mrkdwn", Text: line}, nil, nil),
			)
		}
	} else {
		blocks.BlockSet = append(blocks.BlockSet,
			slack.NewSectionBlock(&slack.TextBlockObject{Type: "mrkdwn", Text: "🟢 All systems healthy"}, nil, nil),
		)
	}

	return blocks, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
