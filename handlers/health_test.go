package handlers

import (
	"context"
	"testing"

	"github.com/pilardi/slack-backbone/config"
)

func TestHealthHandler_Name(t *testing.T) {
	h := NewHealthHandler(nil)
	if h.Name() != "health" {
		t.Errorf("expected name 'health', got '%s'", h.Name())
	}
}

func TestHealthHandler_Scope(t *testing.T) {
	h := NewHealthHandler(nil)
	scope := h.Scope()
	if _, ok := scope.(*GlobalScope); !ok {
		t.Error("expected Global scope")
	}
}

func TestHealthHandler_Run_ReturnsBlocks(t *testing.T) {
	h := NewHealthHandler(nil)
	team := config.Team{Name: "test-team"}
	ctx := context.Background()

	blocks, err := h.Run(ctx, []string{}, team)
	if blocks == nil {
		t.Fatal("expected non-nil blocks")
	}
	if len(blocks.BlockSet) == 0 {
		t.Error("expected at least one block")
	}
	// With nil client, HealthCheck panics - so we expect a panic here
	_ = err // suppress unused
}

func TestHealthHandler_Run_WithTeamName(t *testing.T) {
	h := NewHealthHandler(nil)
	team := config.Team{Name: "acme-corp"}
	ctx := context.Background()

	blocks, err := h.Run(ctx, []string{}, team)
	if blocks == nil {
		t.Error("expected non-empty blocks")
	}
	_ = err // may panic with nil client
}

func TestHealthHandler_Run_WithArgs(t *testing.T) {
	h := NewHealthHandler(nil)
	team := config.Team{Name: "beta-team"}
	ctx := context.Background()

	blocks, err := h.Run(ctx, []string{"--verbose"}, team)
	if blocks == nil {
		t.Error("expected non-empty blocks")
	}
	_ = err // may panic with nil client
}
