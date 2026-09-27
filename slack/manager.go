package slack

import (
	"context"
	"fmt"

	"github.com/Asafrose/bolt-go"
	"github.com/pilardi/slack-backbone/config"
	slackapi "github.com/slack-go/slack"
)

// Manager handles N concurrent Bolt apps, one per team.
type Manager struct {
	apps map[string]*bolt.App
}

// NewManager creates a new multi-team manager.
func NewManager() *Manager {
	return &Manager{
		apps: make(map[string]*bolt.App),
	}
}

// Register adds a team's Bolt app to the manager.
func (m *Manager) Register(ctx context.Context, team config.Team) error {
	app, err := bolt.New(bolt.AppOptions{
		Token:      team.BotToken,
		AppToken:   team.AppToken,
		SocketMode: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create Bolt app for team %q: %w", team.Name, err)
	}

	m.apps[team.Name] = app
	return nil
}

// StartAll starts listening on all registered teams.
func (m *Manager) StartAll(ctx context.Context) error {
	for name, app := range m.apps {
		go func(name string, app *bolt.App) {
			if err := app.Start(ctx); err != nil {
				fmt.Printf("⚠️  %s: %v\n", name, err)
			}
		}(name, app)
	}
	return nil
}

// PostToChannel sends a message to a specific team's channel.
func (m *Manager) PostToChannel(ctx context.Context, team config.Team, channel string, blocks []slackapi.Block) error {
	msg := &slackapi.Blocks{}
	msg.BlockSet = append(msg.BlockSet, blocks...)

	client := slackapi.New(team.BotToken)
	_, _, err := client.PostMessageContext(ctx, channel, slackapi.MsgOptionBlocks(msg.BlockSet...))
	return err
}

// PostEphemeral sends an ephemeral message to a user.
func (m *Manager) PostEphemeral(ctx context.Context, team config.Team, userID string, blocks []slackapi.Block) error {
	msg := &slackapi.Blocks{}
	msg.BlockSet = append(msg.BlockSet, blocks...)

	client := slackapi.New(team.BotToken)
	_, err := client.PostEphemeralContext(ctx, "#general", userID, slackapi.MsgOptionBlocks(msg.BlockSet...))
	return err
}
