package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/pilardi/slack-backbone/config"
	slackpkg "github.com/pilardi/slack-backbone/slack"
	"github.com/slack-go/slack"
)

// StatusHandler handles the /slack-backbone status command.
type StatusHandler struct {
	client *slackpkg.Client
}

// NewStatusHandler creates a new status handler with the given Slack client.
func NewStatusHandler(client *slackpkg.Client) *StatusHandler {
	return &StatusHandler{client: client}
}

func (h *StatusHandler) Name() string { return "status" }

func (h *StatusHandler) Scope() Scope { return Global() }

func (h *StatusHandler) Run(ctx context.Context, args []string, team config.Team) (*slack.Blocks, error) {
	// Parse --team flag if present
	var targetTeam string
	for i := 0; i < len(args)-1; i++ {
		if args[i] == "--team" && i+1 < len(args) {
			targetTeam = args[i+1]
			i++
		}
	}

	blocks := slackpkg.NewBlocks()

	if targetTeam != "" {
		// Show status for a specific team
		blocks.BlockSet = append(blocks.BlockSet,
			slack.NewSectionBlock(&slack.TextBlockObject{Type: "mrkdwn", Text: fmt.Sprintf("📊 **Status · Team: %s**", targetTeam)}, nil, nil),
		)

		if h.client != nil {
			result, err := h.client.HealthCheck(ctx)
			if err != nil {
				blocks.BlockSet = append(blocks.BlockSet,
					slack.NewSectionBlock(&slack.TextBlockObject{Type: "mrkdwn", Text: fmt.Sprintf("❌ **Bot:** %v", err)}, nil, nil),
				)
			} else {
				blocks.BlockSet = append(blocks.BlockSet,
					slack.NewSectionBlock(&slack.TextBlockObject{Type: "mrkdwn", Text: fmt.Sprintf("✅ **Bot:** connected as *%s* (`%s`)", result.UserName, result.UserID)}, nil, nil),
					slack.NewSectionBlock(&slack.TextBlockObject{Type: "mrkdwn", Text: fmt.Sprintf("⏱️ **Latency:** %.1fms", result.LatencyMs)}, nil, nil),
				)

				channels, err := h.client.GetAccessibleChannels(ctx)
				if err != nil {
					blocks.BlockSet = append(blocks.BlockSet,
						slack.NewSectionBlock(&slack.TextBlockObject{Type: "mrkdwn", Text: fmt.Sprintf("❌ **Channels:** %v", err)}, nil, nil),
					)
				} else {
					blocks.BlockSet = append(blocks.BlockSet,
						slack.NewSectionBlock(&slack.TextBlockObject{Type: "mrkdwn", Text: fmt.Sprintf("✅ **Channels:** %d accessible", len(channels))}, nil, nil),
					)

					if len(channels) > 0 {
						chList := make([]string, 0, len(channels))
						for _, ch := range channels[:min(5, len(channels))] {
							chList = append(chList, fmt.Sprintf("`%s`", ch.Name))
						}
						if len(channels) > 5 {
							chList = append(chList, fmt.Sprintf("... and %d more", len(channels)-5))
						}
						blocks.BlockSet = append(blocks.BlockSet,
							slack.NewSectionBlock(&slack.TextBlockObject{Type: "mrkdwn", Text: strings.Join(chList, ", ")}, nil, nil),
						)
					}
				}

				blocks.BlockSet = append(blocks.BlockSet, slack.NewDividerBlock())
			}
		} else {
			blocks.BlockSet = append(blocks.BlockSet,
				slack.NewSectionBlock(&slack.TextBlockObject{Type: "mrkdwn", Text: "⚠️ **No client configured** — skipping API checks"}, nil, nil),
			)
		}

		return blocks, nil
	}

	// No target team: show overview of all teams
	blocks.BlockSet = append(blocks.BlockSet,
		slack.NewSectionBlock(&slack.TextBlockObject{Type: "mrkdwn", Text: "📋 **Team Overview**"}, nil, nil),
	)

	if h.client != nil {
		result, err := h.client.HealthCheck(ctx)
		if err != nil {
			blocks.BlockSet = append(blocks.BlockSet,
				slack.NewSectionBlock(&slack.TextBlockObject{Type: "mrkdwn", Text: fmt.Sprintf("❌ **Bot:** %v", err)}, nil, nil),
			)
		} else {
			blocks.BlockSet = append(blocks.BlockSet,
				slack.NewSectionBlock(&slack.TextBlockObject{Type: "mrkdwn", Text: fmt.Sprintf("✅ **Bot:** connected as *%s* (`%s`)", result.UserName, result.UserID)}, nil, nil),
				slack.NewSectionBlock(&slack.TextBlockObject{Type: "mrkdwn", Text: fmt.Sprintf("⏱️ **Latency:** %.1fms", result.LatencyMs)}, nil, nil),
			)

			channels, err := h.client.GetAccessibleChannels(ctx)
			if err != nil {
				blocks.BlockSet = append(blocks.BlockSet,
					slack.NewSectionBlock(&slack.TextBlockObject{Type: "mrkdwn", Text: fmt.Sprintf("❌ **Channels:** %v", err)}, nil, nil),
				)
			} else {
				blocks.BlockSet = append(blocks.BlockSet,
					slack.NewSectionBlock(&slack.TextBlockObject{Type: "mrkdwn", Text: fmt.Sprintf("✅ **Total Channels:** %d accessible", len(channels))}, nil, nil),
				)

				if len(channels) > 0 {
					chList := make([]string, 0, len(channels))
					for _, ch := range channels[:min(5, len(channels))] {
						chList = append(chList, fmt.Sprintf("`%s`", ch.Name))
					}
					if len(channels) > 5 {
						chList = append(chList, fmt.Sprintf("... and %d more", len(channels)-5))
					}
					blocks.BlockSet = append(blocks.BlockSet,
						slack.NewSectionBlock(&slack.TextBlockObject{Type: "mrkdwn", Text: strings.Join(chList, ", ")}, nil, nil),
					)
				}
			}

			blocks.BlockSet = append(blocks.BlockSet, slack.NewDividerBlock())
		}
	} else {
		blocks.BlockSet = append(blocks.BlockSet,
			slack.NewSectionBlock(&slack.TextBlockObject{Type: "mrkdwn", Text: "⚠️ **No client configured** — skipping API checks"}, nil, nil),
		)
	}

	return blocks, nil
}
