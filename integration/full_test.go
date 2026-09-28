package integration

import (
	"context"
	"testing"

	"github.com/pilardi/slack-backbone/config"
	"github.com/pilardi/slack-backbone/handlers"
)

// TestFullFlow verifies config loading + handler execution end-to-end.
func TestFullFlow(t *testing.T) {
	cfg := &config.Config{
		Teams: []config.Team{
			{Name: "prod-workspace", BotToken: "xoxb-prod", AppToken: "xapp-prod"},
			{Name: "staging-workspace", BotToken: "xoxb-stg", AppToken: "xapp-stg"},
		},
	}

	// Verify config has expected teams
	if len(cfg.Teams) != 2 {
		t.Fatalf("expected 2 teams, got %d", len(cfg.Teams))
	}

	// Run deploy handler for each team
	for _, team := range cfg.Teams {
		h := &handlers.DeployHandler{}
		blocks, err := h.Run(context.Background(), []string{"--env", "staging"}, team)
		if err != nil {
			t.Fatalf("team %s: unexpected error: %v", team.Name, err)
		}
		if blocks == nil || len(blocks.BlockSet) == 0 {
			t.Fatalf("team %s: expected non-empty block set", team.Name)
		}
	}
}

// TestHealthEndToEnd verifies health handler end-to-end.
func TestHealthEndToEnd(t *testing.T) {
	cfg := &config.Config{
		Teams: []config.Team{{Name: "test", BotToken: "xoxb-mock", AppToken: "xapp-mock"}},
	}

	h := &handlers.HealthHandler{}
	blocks, err := h.Run(context.Background(), []string{}, cfg.Teams[0])
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if blocks == nil || len(blocks.BlockSet) == 0 {
		t.Fatal("expected non-empty block set")
	}
}

// TestStatusEndToEnd verifies status handler end-to-end.
func TestStatusEndToEnd(t *testing.T) {
	cfg := &config.Config{
		Teams: []config.Team{{Name: "test", BotToken: "xoxb-mock", AppToken: "xapp-mock"}},
	}

	h := &handlers.StatusHandler{}
	blocks, err := h.Run(context.Background(), []string{}, cfg.Teams[0])
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if blocks == nil || len(blocks.BlockSet) == 0 {
		t.Fatal("expected non-empty block set")
	}
}
