package mcp

import (
	"context"
	"log/slog"
	"strings"
	"testing"

	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/pilardi/slack-backbone/config"
)

func TestPromptSummarizeChannel(t *testing.T) {
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
		t.Fatalf("makeSummarizeHandler failed: %v", err)
	}
	if result == nil || len(result.Messages) == 0 {
		t.Fatal("expected non-empty messages")
	}
	if !strings.Contains(result.Messages[0].Content.(*mcp.TextContent).Text, "#general") {
		t.Error("expected message to reference the channel name")
	}
}

func TestPromptAlertOnKeyword(t *testing.T) {
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
		t.Fatalf("makeAlertHandler failed: %v", err)
	}
	if result == nil || len(result.Messages) == 0 {
		t.Fatal("expected non-empty messages")
	}
	text := result.Messages[0].Content.(*mcp.TextContent).Text
	if !strings.Contains(text, "rollback") {
		t.Error("expected message to reference the keyword 'rollback'")
	}
}

func TestPromptDailyStatus(t *testing.T) {
	result, err := makeStatusHandler()(context.Background(), &mcp.GetPromptRequest{
		Params: &mcp.GetPromptParams{
			Arguments: map[string]string{
				"team":    "prod-workspace",
				"channel": "#ops",
			},
		},
	})
	if err != nil {
		t.Fatalf("makeStatusHandler failed: %v", err)
	}
	if result == nil || len(result.Messages) == 0 {
		t.Fatal("expected non-empty messages")
	}
	text := result.Messages[0].Content.(*mcp.TextContent).Text
	if !strings.Contains(text, "#ops") {
		t.Error("expected message to reference the channel name")
	}
}

func TestPromptSummarizeChannelMissingTeam(t *testing.T) {
	result, err := makeSummarizeHandler()(context.Background(), &mcp.GetPromptRequest{
		Params: &mcp.GetPromptParams{
			Arguments: map[string]string{
				"channel": "#general",
				"count":   "5",
			},
		},
	})
	if err == nil {
		t.Fatal("expected error when team is missing")
	}
	if result != nil {
		t.Error("expected nil result on error")
	}
}

func TestPromptAlertOnKeywordMissingChannel(t *testing.T) {
	result, err := makeAlertHandler()(context.Background(), &mcp.GetPromptRequest{
		Params: &mcp.GetPromptParams{
			Arguments: map[string]string{
				"team":    "prod-workspace",
				"keyword": "error",
			},
		},
	})
	if err == nil {
		t.Fatal("expected error when channel is missing")
	}
	if result != nil {
		t.Error("expected nil result on error")
	}
}

func TestPromptDailyStatusMissingTeam(t *testing.T) {
	result, err := makeStatusHandler()(context.Background(), &mcp.GetPromptRequest{
		Params: &mcp.GetPromptParams{
			Arguments: map[string]string{
				"channel": "#ops",
			},
		},
	})
	if err == nil {
		t.Fatal("expected error when team is missing")
	}
	if result != nil {
		t.Error("expected nil result on error")
	}
}

func TestPromptDailyStatusEmptyArgs(t *testing.T) {
	result, err := makeStatusHandler()(context.Background(), &mcp.GetPromptRequest{
		Params: &mcp.GetPromptParams{
			Arguments: map[string]string{},
		},
	})
	if err == nil {
		t.Fatal("expected error when team and channel are missing")
	}
	if result != nil {
		t.Error("expected nil result on error")
	}
}

func TestPromptSummarizeChannelDefaultCount(t *testing.T) {
	result, err := makeSummarizeHandler()(context.Background(), &mcp.GetPromptRequest{
		Params: &mcp.GetPromptParams{
			Arguments: map[string]string{
				"team":    "prod-workspace",
				"channel": "#general",
			},
		},
	})
	if err != nil {
		t.Fatalf("makeSummarizeHandler failed with default count: %v", err)
	}
	text := result.Messages[0].Content.(*mcp.TextContent).Text
	if !strings.Contains(text, "10 messages") {
		t.Error("expected default count of 10 in message text")
	}
}

func TestPromptSummarizeChannelZeroCount(t *testing.T) {
	result, err := makeSummarizeHandler()(context.Background(), &mcp.GetPromptRequest{
		Params: &mcp.GetPromptParams{
			Arguments: map[string]string{
				"team":    "prod-workspace",
				"channel": "#general",
				"count":   "0",
			},
		},
	})
	if err != nil {
		t.Fatalf("makeSummarizeHandler failed with zero count: %v", err)
	}
	text := result.Messages[0].Content.(*mcp.TextContent).Text
	if !strings.Contains(text, "10 messages") {
		t.Error("expected fallback to default count of 10 when count is 0")
	}
}

func TestPromptSummarizeChannelNegativeCount(t *testing.T) {
	result, err := makeSummarizeHandler()(context.Background(), &mcp.GetPromptRequest{
		Params: &mcp.GetPromptParams{
			Arguments: map[string]string{
				"team":    "prod-workspace",
				"channel": "#general",
				"count":   "-5",
			},
		},
	})
	if err != nil {
		t.Fatalf("makeSummarizeHandler failed with negative count: %v", err)
	}
	text := result.Messages[0].Content.(*mcp.TextContent).Text
	if !strings.Contains(text, "10 messages") {
		t.Error("expected fallback to default count of 10 when count is negative")
	}
}

func TestPromptAlertOnKeywordEmpty(t *testing.T) {
	result, err := makeAlertHandler()(context.Background(), &mcp.GetPromptRequest{
		Params: &mcp.GetPromptParams{
			Arguments: map[string]string{},
		},
	})
	if err == nil {
		t.Fatal("expected error when all args are missing")
	}
	if result != nil {
		t.Error("expected nil result on error")
	}
}

func TestPromptRegistration(t *testing.T) {
	cfg := &config.Config{
		Teams: []config.Team{{Name: "test", BotToken: "xoxb-test", AppToken: "xapp-test"}},
	}
	srv, _ := NewServer(cfg, slog.Default())

	RegisterPrompts(srv)
	// Verify prompts were registered without panicking
	if srv == nil {
		t.Fatal("server should not be nil")
	}
}
