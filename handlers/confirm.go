package handlers

import (
	"context"
	"fmt"

	"github.com/pablo/slack-backbone/config"
	"github.com/slack-go/slack"
)

// ConfirmHandler handles interactive confirmation buttons.
type ConfirmHandler struct{}

func (h *ConfirmHandler) Name() string { return "confirm" }

func (h *ConfirmHandler) Scope() Scope { return Scoped("confirm") }

func (h *ConfirmHandler) Run(ctx context.Context, args []string, team config.Team) (*slack.MessageBlock, error) {
	return &slack.MessageBlock{
		Type:   "section",
		Text:   &slack.TextBlockObject{Type: "mrkdwn", Text: "⏳ Waiting for confirmation..."},
		Blocks: []any{confirmButton(team.Name)},
	}, nil
}

func confirmButton(teamName string) any {
	return &slack.ActionBlock{
		BlockID: "confirm_" + teamName,
		Elements: &slack.ElementsBlock{
			ActionBlocks: []*slack.ActionBlock{
				{
					BlockID: "confirm_yes",
					Elements: &slack.ElementBlock{
						Type: "button",
						Text: &slack.TextBlockObject{Type: "plain_text", Text: "✅ Confirm"},
						Value: "confirmed",
					},
				},
				{
					BlockID: "confirm_no",
					Elements: &slack.ElementBlock{
						Type: "button",
						Text: &slack.TextBlockObject{Type: "plain_text", Text: "❌ Cancel"},
						Value: "cancelled",
					},
				},
			},
		},
	}
}
