package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Resource definitions for read-only Slack data agents can query.

// resourceChannels lists accessible channels on a team.
func resourceChannels(team string) *mcp.Resource {
	return &mcp.Resource{
		Name:        fmt.Sprintf("channels_%s", team),
		Description: fmt.Sprintf("List of accessible public and private channels for team %q.", team),
		URI:         "slack://channels/" + team,
		MIMEType:    "application/json",
	}
}

// resourceUsers lists users on a team.
func resourceUsers(team string) *mcp.Resource {
	return &mcp.Resource{
		Name:        fmt.Sprintf("users_%s", team),
		Description: fmt.Sprintf("Directory of Slack users for team %q.", team),
		URI:         "slack://users/" + team,
		MIMEType:    "application/json",
	}
}

// resourceTeams lists all configured teams.
func resourceTeams() *mcp.Resource {
	return &mcp.Resource{
		Name:        "teams",
		Description: "Metadata about all configured Slack workspaces.",
		URI:         "slack://teams",
		MIMEType:    "application/json",
	}
}

// RegisterResources adds resource definitions to the MCP server.
func RegisterResources(s *Server) {
	for _, team := range s.cfg.Teams {
		s.server.AddResource(resourceChannels(team.Name), makeChannelsHandler(s))
		s.server.AddResource(resourceUsers(team.Name), makeUsersHandler(s))
	}
	s.server.AddResource(resourceTeams(), makeTeamsHandler(s))
}

// ---- Resource Handlers (closures capturing *Server) ----

func makeChannelsHandler(s *Server) mcp.ResourceHandler {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		teamName := strings.TrimPrefix(req.Params.URI, "slack://channels/")
		client, err := s.GetClient(teamName)
		if err != nil {
			return nil, err
		}

		channels, err := client.GetAccessibleChannels(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list channels: %w", err)
		}

		var parts []string
		for _, ch := range channels {
			parts = append(parts, fmt.Sprintf(`{"id":"%s","name":"%s"}`, ch.ID, ch.Name))
		}
		return &mcp.ReadResourceResult{Contents: []*mcp.ResourceContents{{MIMEType: "application/json", Text: "[" + strings.Join(parts, ",") + "]"}}}, nil
	}
}

func makeUsersHandler(s *Server) mcp.ResourceHandler {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		teamName := strings.TrimPrefix(req.Params.URI, "slack://users/")
		// Placeholder: user directory resolved via slack_api.users_list
		return &mcp.ReadResourceResult{Contents: []*mcp.ResourceContents{{MIMEType: "application/json", Text: fmt.Sprintf(`{"team":"%s","note":"user directory resolved via slack_api.users_list"}`, teamName)}}}, nil
	}
}

func makeTeamsHandler(s *Server) mcp.ResourceHandler {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		var parts []string
		for _, t := range s.cfg.Teams {
			parts = append(parts, fmt.Sprintf(`{"name":"%s","commands":["health","status","deploy","confirm"]}`, t.Name))
		}
		return &mcp.ReadResourceResult{Contents: []*mcp.ResourceContents{{MIMEType: "application/json", Text: "[" + strings.Join(parts, ",") + "]"}}}, nil
	}
}
