package handlers

import "github.com/slack-go/slack"

// Handler is the interface all command handlers must implement.
type Handler interface {
	Name() string
	Scope() Scope
	Run(ctx context.Context, args []string, team Team) (*slack.MessageBlock, error)
}
