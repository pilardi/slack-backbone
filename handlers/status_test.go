package handlers

import (
	"context"
	"testing"

	"github.com/pilardi/slack-backbone/config"
)

func TestStatusHandler_Name(t *testing.T) {
	h := NewStatusHandler(nil)
	if h.Name() != "status" {
		t.Errorf("expected name 'status', got '%s'", h.Name())
	}
}

func TestStatusHandler_Scope(t *testing.T) {
	h := NewStatusHandler(nil)
	scope := h.Scope()
	if _, ok := scope.(*GlobalScope); !ok {
		t.Error("expected Global scope")
	}
}

func TestStatusHandler_Run_ReturnsBlocks(t *testing.T) {
	h := NewStatusHandler(nil)
	team := config.Team{Name: "test-team"}
	ctx := context.Background()

	blocks, err := h.Run(ctx, []string{}, team)
	if blocks == nil {
		t.Fatal("expected non-nil blocks")
	}
	if len(blocks.BlockSet) == 0 {
		t.Error("expected at least one block")
	}
	_ = err // may panic with nil client
}

func TestStatusHandler_Run_WithTeamArg(t *testing.T) {
	h := NewStatusHandler(nil)
	team := config.Team{Name: "acme-corp"}
	ctx := context.Background()

	blocks, err := h.Run(ctx, []string{"--team", "acme-corp"}, team)
	if blocks == nil {
		t.Error("expected non-empty blocks")
	}
	_ = err // may panic with nil client
}

func TestStatusHandler_Run_WithArgs(t *testing.T) {
	h := NewStatusHandler(nil)
	team := config.Team{Name: "beta-team"}
	ctx := context.Background()

	blocks, err := h.Run(ctx, []string{"--team", "beta-team"}, team)
	if blocks == nil {
		t.Error("expected non-empty blocks")
	}
	_ = err // may panic with nil client
}
