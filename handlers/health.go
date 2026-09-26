package handlers

import (
	"context"
	"fmt"

	"github.com/pablo/slack-backbone/config"
	"github.com/slack-go/slack"
)

// HealthHandler handles the /slack-backbone health command.
type HealthHandler struct{}

func (h *HealthHandler) Name() string { return "health" }

func (h *HealthHandler) Scope() Scope { return Global() }

func (h *HealthHandler) Run(ctx context.Context, args []string, team config.Team) (*slack.MessageBlock, error) {
	return &slack.MessageBlock{
		Type: "section",
		Text: &slack.TextBlockObject{
			Type: "mrkdwn",
			Text: fmt.Sprintf("🟢 Healthy · Team: **%s**", team.Name),
		},
	}, nil
}
