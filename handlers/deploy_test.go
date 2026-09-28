package handlers

import (
	"context"
	"strings"
	"testing"

	"github.com/pilardi/slack-backbone/config"
	"github.com/slack-go/slack"
)

func TestDeployHandler_Name(t *testing.T) {
	h := &DeployHandler{}
	if got := h.Name(); got != "deploy" {
		t.Errorf("Name() = %q, want \"deploy\"", got)
	}
}

func TestDeployHandler_Scope(t *testing.T) {
	h := &DeployHandler{}
	s := h.Scope()
	ss, ok := s.(*ScopedScope)
	if !ok {
		t.Fatal("expected ScopedScope")
	}
	if ss.TeamName != "deploy" {
		t.Errorf("TeamName = %q, want \"deploy\"", ss.TeamName)
	}
}

func TestDeployHandler_Run_DefaultEnv(t *testing.T) {
	h := &DeployHandler{}
	blocks, err := h.Run(context.Background(), []string{}, config.Team{Name: "test-workspace"})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	text := blocksToText(blocks)
	if !strings.Contains(text, "**Deploy initiated**") {
		t.Errorf("expected deploy message, got: %s", text)
	}
	if !strings.Contains(text, "test-workspace") {
		t.Errorf("expected team name in message, got: %s", text)
	}
	if !strings.Contains(text, "production") {
		t.Errorf("expected default env 'production', got: %s", text)
	}
}

func TestDeployHandler_Run_CustomEnv(t *testing.T) {
	h := &DeployHandler{}
	blocks, err := h.Run(context.Background(), []string{"--env", "staging"}, config.Team{Name: "staging-workspace"})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	text := blocksToText(blocks)
	if !strings.Contains(text, "staging") {
		t.Errorf("expected env 'staging' in message, got: %s", text)
	}
}

func TestDeployHandler_Run_InvalidEnv(t *testing.T) {
	h := &DeployHandler{}
	_, err := h.Run(context.Background(), []string{"--env", "invalid-env"}, config.Team{Name: "test"})
	if err == nil {
		t.Error("expected error for invalid env, got nil")
	}
}

// blocksToText is a helper to extract text from Slack blocks (mirrors the mcp converter logic).
func blocksToText(b *slack.Blocks) string {
	var parts []string
	for _, block := range b.BlockSet {
		if sb, ok := block.(*slack.SectionBlock); ok && sb.Text != nil {
			parts = append(parts, sb.Text.Text)
		}
	}
	return strings.Join(parts, "\n")
}
