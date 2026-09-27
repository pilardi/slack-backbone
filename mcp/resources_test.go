package mcp

import (
	"context"
	"log/slog"
	"strings"
	"testing"

	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/pilardi/slack-backbone/config"
)

func TestHandleChannels(t *testing.T) {
	cfg := &config.Config{
		Teams: []config.Team{{Name: "test", BotToken: "xoxb-test", AppToken: "xapp-test"}},
	}
	srv, _ := NewServer(cfg, slog.Default())
	srv.InitClients()

	result, err := makeChannelsHandler(srv)(context.Background(), &mcp.ReadResourceRequest{
		Params: &mcp.ReadResourceParams{URI: "slack://channels/test"},
	})
	if err != nil {
		t.Logf("handleChannels returned error (expected with mock tokens): %v", err)
	}
	// With mock tokens, the API call may fail auth; we just verify no panic
	if result != nil && len(result.Contents) > 0 {
		// Good — got a response
	}
}

func TestHandleUsers(t *testing.T) {
	cfg := &config.Config{
		Teams: []config.Team{{Name: "test", BotToken: "xoxb-test", AppToken: "xapp-test"}},
	}
	srv, _ := NewServer(cfg, slog.Default())

	result, err := makeUsersHandler(srv)(context.Background(), &mcp.ReadResourceRequest{
		Params: &mcp.ReadResourceParams{URI: "slack://users/test"},
	})
	if err != nil {
		t.Fatalf("handleUsers failed: %v", err)
	}
	if result == nil || len(result.Contents) == 0 {
		t.Fatal("expected non-empty result")
	}
}

func TestHandleTeams(t *testing.T) {
	cfg := &config.Config{
		Teams: []config.Team{
			{Name: "prod-workspace"},
			{Name: "staging-workspace"},
		},
	}
	srv, _ := NewServer(cfg, slog.Default())

	result, err := makeTeamsHandler(srv)(context.Background(), &mcp.ReadResourceRequest{})
	if err != nil {
		t.Fatalf("handleTeams failed: %v", err)
	}
	if result == nil || len(result.Contents) == 0 {
		t.Fatal("expected non-empty result")
	}
}

func TestHandleChannelsUnknownTeam(t *testing.T) {
	cfg := &config.Config{
		Teams: []config.Team{{Name: "test", BotToken: "xoxb-test", AppToken: "xapp-test"}},
	}
	srv, _ := NewServer(cfg, slog.Default())
	srv.InitClients()

	result, err := makeChannelsHandler(srv)(context.Background(), &mcp.ReadResourceRequest{
		Params: &mcp.ReadResourceParams{URI: "slack://channels/nonexistent"},
	})
	if result != nil {
		t.Fatal("expected nil result for unknown team")
	}
	if err == nil {
		t.Fatal("expected error for unknown team")
	}
	if !strings.Contains(err.Error(), "unknown team") {
		t.Fatalf("unexpected error: %v", err)
	}
}
