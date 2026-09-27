package handlers

import (
	"context"

	"github.com/pilardi/slack-backbone/config"
	"github.com/slack-go/slack"
)

// Handler is the interface all command handlers must implement.
type Handler interface {
	Name() string
	Scope() Scope
	Run(ctx context.Context, args []string, team config.Team) (*slack.Blocks, error)
}
