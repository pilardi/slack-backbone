package handlers

import (
	"context"
	"fmt"

	"github.com/pilardi/slack-backbone/config"
	"github.com/slack-go/slack"
)

// ConfirmHandler handles interactive confirmation buttons.
type ConfirmHandler struct{}

func (h *ConfirmHandler) Name() string { return "confirm" }

func (h *ConfirmHandler) Scope() Scope { return Scoped("confirm") }

func (h *ConfirmHandler) Run(ctx context.Context, args []string, team config.Team) (*slack.Blocks, error) {
	section := slack.NewSectionBlock(&slack.TextBlockObject{Type: "mrkdwn", Text: "⏳ Waiting for confirmation..."}, nil, nil)
	action := confirmButton(team.Name)

	return &slack.Blocks{
		BlockSet: []slack.Block{section, action},
	}, nil
}

func confirmButton(teamName string) *slack.ActionBlock {
	return slack.NewActionBlock("confirm_"+teamName,
		&slack.ButtonBlockElement{
			Type:  "button",
			Text:  &slack.TextBlockObject{Type: "plain_text", Text: "✅ Confirm"},
			Value: "confirmed",
		},
		&slack.ButtonBlockElement{
			Type:  "button",
			Text:  &slack.TextBlockObject{Type: "plain_text", Text: "❌ Cancel"},
			Value: "cancelled",
		},
	)
}

// ConfirmResponse builds a response block for the confirm action.
func ConfirmResponse(teamName string, confirmed bool) *slack.Blocks {
	status := "✅ **Confirmed!**"
	if !confirmed {
		status = "❌ **Cancelled.**"
	}
	return &slack.Blocks{
		BlockSet: []slack.Block{
			slack.NewSectionBlock(&slack.TextBlockObject{Type: "mrkdwn", Text: fmt.Sprintf("%s · Team: **%s**", status, teamName)}, nil, nil),
		},
	}
}
