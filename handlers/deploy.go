package handlers

import (
	"context"
	"fmt"

	"github.com/pilardi/slack-backbone/config"
	"github.com/slack-go/slack"
)

// DeployHandler handles the /slack-backbone deploy command.
type DeployHandler struct{}

func (h *DeployHandler) Name() string { return "deploy" }

func (h *DeployHandler) Scope() Scope { return Scoped("deploy") }

func (h *DeployHandler) Run(ctx context.Context, args []string, team config.Team) (*slack.Blocks, error) {
	env := "production"
	for i := 0; i < len(args)-1; i++ {
		if args[i] == "--env" && i+1 < len(args) {
			env = args[i+1]
			i++
		}
	}

	return &slack.Blocks{
		BlockSet: []slack.Block{
			slack.NewSectionBlock(&slack.TextBlockObject{
				Type: "mrkdwn",
				Text: fmt.Sprintf("🚀 Deploying to **%s**...", env),
			}, nil, nil),
		},
	}, nil
}
