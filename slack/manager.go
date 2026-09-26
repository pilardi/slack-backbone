package slack

import (
	"context"
	"fmt"

	"github.com/pablo/slack-backbone/config"
	"github.com/slack-go/bolt"
	"github.com/slack-go/slack"
)

// Manager handles N concurrent Bolt apps, one per team.
type Manager struct {
	apps  map[string]*bolt.App
	client *slack.Client
}

// NewManager creates a new multi-team manager.
func NewManager() *Manager {
	return &Manager{
		apps:   make(map[string]*bolt.App),
		client: slack.New(),
	}
}

// Register adds a team's Bolt app to the manager.
func (m *Manager) Register(ctx context.Context, team config.Team) error {
	app, err := bolt.New(
		bolt.OptionToken(team.BotToken),
		bolt.OptionAppLevelToken(team.AppToken),
	)
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
func (m *Manager) PostToChannel(ctx context.Context, team config.Team, channel string, blocks []any) error {
	msg := slack.NewMessageBlock()
	msg.Blocks = append(msg.Blocks, blocks...)

	_, _, err := m.client.PostMessageContext(ctx, channel, slack.MsgOptionBlocks(msg.Blocks...))
	return err
}

// PostEphemeral sends an ephemeral message to a user.
func (m *Manager) PostEphemeral(ctx context.Context, team config.Team, userID string, blocks []any) error {
	msg := slack.NewMessageBlock()
	msg.Blocks = append(msg.Blocks, blocks...)

	_, err := m.client.PostEphemeralContext(ctx, userID, msg.Text.Text, slack.EphemeralOptionBlocks(msg.Blocks...))
	return err
}
