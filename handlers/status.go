package handlers

import (
	"context"
	"fmt"

	"github.com/pablo/slack-backbone/config"
	"github.com/slack-go/slack"
)

// StatusHandler handles the /slack-backbone status command.
type StatusHandler struct{}

func (h *StatusHandler) Name() string { return "status" }

func (h *StatusHandler) Scope() Scope { return Global() }

func (h *StatusHandler) Run(ctx context.Context, args []string, team config.Team) (*slack.MessageBlock, error) {
	return &slack.MessageBlock{
		Type: "section",
		Text: &slack.TextBlockObject{
			Type: "mrkdwn",
			Text: fmt.Sprintf("✅ Running · Team: **%s**", team.Name),
		},
	}, nil
}
