// Package mcp provides an MCP (Model Context Protocol) server that exposes
// Slack-native tools, resources, and prompts to agents.
package mcp

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/pilardi/slack-backbone/config"
	slackpkg "github.com/pilardi/slack-backbone/slack"
)

// Server wraps the MCP server and provides a clean API for agent integration.
type Server struct {
	server  *mcp.Server
	logger  *slog.Logger
	cfg     *config.Config
	clients map[string]*slackpkg.Client
}

// NewServer creates a new MCP server backed by the given config.
func NewServer(cfg *config.Config, logger *slog.Logger) (*Server, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is required")
	}
	if logger == nil {
		logger = slog.Default()
	}

	srv := &Server{
		logger:  logger,
		cfg:     cfg,
		clients: make(map[string]*slackpkg.Client),
	}

	impl := &mcp.Implementation{
		Name:    "slack-backbone",
		Version: "0.1.0",
	}

	srv.server = mcp.NewServer(impl, nil)

	return srv, nil
}

// AddTool registers a tool with the MCP server.
func (s *Server) AddTool(t *mcp.Tool, h mcp.ToolHandler) {
	s.server.AddTool(t, h)
}

// Start begins listening on stdio (default transport).
func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("starting MCP server", "transport", "stdio")
	return s.server.Run(ctx, &mcp.StdioTransport{})
}

// StartHTTP begins listening on the given HTTP port using streamable HTTP.
func (s *Server) StartHTTP(ctx context.Context, port int) error {
	s.logger.Info("starting MCP server", "transport", "http", "port", port)

	handler := mcp.NewStreamableHTTPHandler(
		func(req *http.Request) *mcp.Server { return s.server },
		nil,
	)

	addr := fmt.Sprintf(":%d", port)
	srv := &http.Server{
		Addr:        addr,
		Handler:     handler,
		ReadTimeout: 30 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("HTTP server failed", "error", err)
		}
	}()

	s.logger.Info("MCP HTTP listening", "addr", addr)
	return nil
}

// InitClients pre-creates API clients for each team.
func (s *Server) InitClients() {
	for _, team := range s.cfg.Teams {
		s.clients[team.Name] = slackpkg.NewClient(team.BotToken)
	}
}

// GetClient returns a client for the given team, creating one if needed.
func (s *Server) GetClient(teamName string) (*slackpkg.Client, error) {
	if client, ok := s.clients[teamName]; ok {
		return client, nil
	}
	return nil, fmt.Errorf("unknown team: %s", teamName)
}

// Teams returns the configured teams.
func (s *Server) Teams() []config.Team {
	return s.cfg.Teams
}

// Timeout is the default operation timeout for tools.
const Timeout = 30 * time.Second
