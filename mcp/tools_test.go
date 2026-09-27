package mcp

import (
	"context"
	"log/slog"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/pilardi/slack-backbone/config"
)

func TestServerCreation(t *testing.T) {
	cfg := &config.Config{
		Teams: []config.Team{
			{Name: "test-workspace", BotToken: "xoxb-test", AppToken: "xapp-test"},
		},
	}
	logger := slog.Default()

	srv, err := NewServer(cfg, logger)
	if err != nil {
		t.Fatalf("NewServer failed: %v", err)
	}
	if srv == nil {
		t.Fatal("server is nil")
	}
	if srv.Teams() == nil || len(srv.Teams()) != 1 {
		t.Errorf("expected 1 team, got %d", len(srv.Teams()))
	}
}

func TestServerStartHTTP(t *testing.T) {
	cfg := &config.Config{
		Teams: []config.Team{{Name: "test", BotToken: "xoxb-test", AppToken: "xapp-test"}},
	}
	srv, err := NewServer(cfg, slog.Default())
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2000)
	defer cancel()

	err = srv.StartHTTP(ctx, 0) // port 0 = pick a random free port
	if err != nil {
		t.Fatalf("StartHTTP failed: %v", err)
	}
}

func TestRegisterTools(t *testing.T) {
	cfg := &config.Config{
		Teams: []config.Team{{Name: "test", BotToken: "xoxb-test", AppToken: "xapp-test"}},
	}
	srv, _ := NewServer(cfg, slog.Default())

	RegisterTools(srv)
}

func TestRegisterResources(t *testing.T) {
	cfg := &config.Config{
		Teams: []config.Team{{Name: "test", BotToken: "xoxb-test", AppToken: "xapp-test"}},
	}
	srv, _ := NewServer(cfg, slog.Default())

	RegisterResources(srv)
}

func TestRegisterPrompts(t *testing.T) {
	cfg := &config.Config{
		Teams: []config.Team{{Name: "test", BotToken: "xoxb-test", AppToken: "xapp-test"}},
	}
	srv, _ := NewServer(cfg, slog.Default())

	RegisterPrompts(srv)
}

func TestHandleListTeams(t *testing.T) {
	cfg := &config.Config{
		Teams: []config.Team{
			{Name: "prod-workspace"},
			{Name: "staging-workspace"},
		},
	}
	srv, _ := NewServer(cfg, slog.Default())

	result, _, err := handleListTeams(srv, context.Background())
	if err != nil {
		t.Fatalf("handleListTeams failed: %v", err)
	}
	if result == nil || len(result.Content) == 0 {
		t.Fatal("expected non-empty result")
	}
	text := result.Content[0].(*mcp.TextContent).Text
	if text == "" {
		t.Error("result text is empty")
	}
}

func TestHandleStatus(t *testing.T) {
	cfg := &config.Config{
		Teams: []config.Team{{Name: "test", BotToken: "xoxb-test", AppToken: "xapp-test"}},
	}
	srv, _ := NewServer(cfg, slog.Default())
	srv.InitClients()

	result, _, err := handleStatus(srv, context.Background(), "test")
	if err != nil {
		t.Fatalf("handleStatus failed: %v", err)
	}
	if result == nil || len(result.Content) == 0 {
		t.Fatal("expected non-empty result")
	}
}

func TestHandleDeploy(t *testing.T) {
	cfg := &config.Config{
		Teams: []config.Team{{Name: "test", BotToken: "xoxb-test", AppToken: "xapp-test"}},
	}
	srv, _ := NewServer(cfg, slog.Default())

	result, _, err := handleDeploy(srv, context.Background(), "test", "")
	if err != nil {
		t.Fatalf("handleDeploy failed: %v", err)
	}
	if result == nil || len(result.Content) == 0 {
		t.Fatal("expected non-empty result")
	}
}

func TestHandleConfirm(t *testing.T) {
	cfg := &config.Config{
		Teams: []config.Team{{Name: "test", BotToken: "xoxb-test", AppToken: "xapp-test"}},
	}
	srv, _ := NewServer(cfg, slog.Default())

	result, _, err := handleConfirm(srv, context.Background(), "test", "Deploy to production?")
	if err != nil {
		t.Fatalf("handleConfirm failed: %v", err)
	}
	if result == nil || len(result.Content) == 0 {
		t.Fatal("expected non-empty result")
	}
}

func TestHandleSummarizeChannel(t *testing.T) {
	result, err := makeSummarizeHandler()(context.Background(), &mcp.GetPromptRequest{
		Params: &mcp.GetPromptParams{
			Arguments: map[string]string{
				"team":    "prod-workspace",
				"channel": "#general",
				"count":   "5",
			},
		},
	})
	if err != nil {
		t.Fatalf("handleSummarizeChannel failed: %v", err)
	}
	if result == nil || len(result.Messages) == 0 {
		t.Fatal("expected non-empty messages")
	}
}

func TestHandleAlertOnKeyword(t *testing.T) {
	result, err := makeAlertHandler()(context.Background(), &mcp.GetPromptRequest{
		Params: &mcp.GetPromptParams{
			Arguments: map[string]string{
				"team":    "staging-workspace",
				"channel": "#deployments",
				"keyword": "rollback",
			},
		},
	})
	if err != nil {
		t.Fatalf("handleAlertOnKeyword failed: %v", err)
	}
	if result == nil || len(result.Messages) == 0 {
		t.Fatal("expected non-empty messages")
	}
}

func TestHandleDailyStatus(t *testing.T) {
	result, err := makeStatusHandler()(context.Background(), &mcp.GetPromptRequest{
		Params: &mcp.GetPromptParams{
			Arguments: map[string]string{
				"team":    "prod-workspace",
				"channel": "#ops",
			},
		},
	})
	if err != nil {
		t.Fatalf("handleDailyStatus failed: %v", err)
	}
	if result == nil || len(result.Messages) == 0 {
		t.Fatal("expected non-empty messages")
	}
}

func TestHandleListTeamsEmpty(t *testing.T) {
	cfg := &config.Config{Teams: []config.Team{}}
	srv, _ := NewServer(cfg, slog.Default())

	result, _, err := handleListTeams(srv, context.Background())
	if err != nil {
		t.Fatalf("handleListTeams failed: %v", err)
	}
	if result == nil || len(result.Content) == 0 {
		t.Fatal("expected non-empty result")
	}
}
