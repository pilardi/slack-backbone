package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/pilardi/slack-backbone/config"
	"github.com/slack-go/slack"
)

// DeployHandler handles the /slack-backbone deploy command.
type DeployHandler struct{}

func (h *DeployHandler) Name() string { return "deploy" }

func (h *DeployHandler) Scope() Scope { return Scoped("deploy") }

func (h *DeployHandler) Run(ctx context.Context, args []string, team config.Team) (*slack.Blocks, error) {
	env := "production"
	channel := ""
	for i := 0; i < len(args)-1; i++ {
		if args[i] == "--env" && i+1 < len(args) {
			env = args[i+1]
			i++
		} else if args[i] == "--channel" && i+1 < len(args) {
			channel = args[i+1]
			i++
		}
	}

	// Validate known environments
	validEnvs := []string{"production", "staging", "development", "qa"}
	known := false
	for _, e := range validEnvs {
		if e == env {
			known = true
			break
		}
	}
	if !known {
		return nil, fmt.Errorf("unknown environment: %s (valid: %v)", env, validEnvs)
	}

	// Log the deployment intent
	slog.InfoContext(ctx, "deploy initiated",
		"team", team.Name,
		"env", env,
		"channel", channel,
		"args", args,
	)

	// Simulate a deployment ID (in production this would come from CI/CD)
	deployID := fmt.Sprintf("dep-%d", time.Now().UnixMilli())

	var msg string
	if channel != "" {
		msg = fmt.Sprintf(
			"🚀 **Deploy initiated**\n\n"+
				"• **Team:** `%s`\n"+
				"• **Environment:** `%s`\n"+
				"• **Channel:** `%s`\n"+
				"• **Deploy ID:** `%s`\n"+
				"• **Status:** ⏳ *In progress*\n\n"+
				"_A real integration would trigger a CI/CD pipeline here._",
			team.Name, env, channel, deployID,
		)
	} else {
		msg = fmt.Sprintf(
			"🚀 **Deploy initiated**\n\n"+
				"• **Team:** `%s`\n"+
				"• **Environment:** `%s`\n"+
				"• **Deploy ID:** `%s`\n"+
				"• **Status:** ⏳ *In progress*\n\n"+
				"_A real integration would trigger a CI/CD pipeline here._",
			team.Name, env, deployID,
		)
	}

	return &slack.Blocks{
		BlockSet: []slack.Block{
			slack.NewSectionBlock(&slack.TextBlockObject{
				Type: "mrkdwn",
				Text: msg,
			}, nil, nil),
		},
	}, nil
}
