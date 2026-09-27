package slack

import (
	"context"
	"fmt"
	"time"

	"github.com/slack-go/slack"
)

// HealthResult holds the result of a health check.
type HealthResult struct {
	UserID    string
	UserName  string
	LatencyMs float64
}

// Client wraps the Slack API client with multi-team awareness.
type Client struct {
	client *slack.Client
}

// NewClient creates a new API client.
func NewClient(token string) *Client {
	return &Client{
		client: slack.New(token),
	}
}

// HealthCheck verifies connectivity by calling auth.test and fetching bot user info.
func (c *Client) HealthCheck(ctx context.Context) (*HealthResult, error) {
	start := time.Now()

	resp, err := c.client.AuthTestContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("auth test: %w", err)
	}

	userID := resp.UserID
	name := resp.User

	latencyMs := float64(time.Since(start).Milliseconds())

	return &HealthResult{
		UserID:    userID,
		UserName:  name,
		LatencyMs: latencyMs,
	}, nil
}

// GetAccessibleChannels returns a list of channels the bot has access to.
func (c *Client) GetAccessibleChannels(ctx context.Context) ([]ChannelInfo, error) {
	var channels []ChannelInfo
	cursor := ""

	for {
		chs, nextCursor, err := c.client.GetConversationsContext(ctx, &slack.GetConversationsParameters{
			Cursor:          cursor,
			ExcludeArchived: true,
			Limit:           100,
			Types:           []string{"public_channel", "private_channel"},
		})
		if err != nil {
			return nil, fmt.Errorf("conversations list: %w", err)
		}

		for _, ch := range chs {
			channels = append(channels, ChannelInfo{ID: ch.ID, Name: ch.Name})
		}

		cursor = nextCursor
		if cursor == "" || len(chs) == 0 {
			break
		}
	}

	return channels, nil
}

// PostMessage posts a message to a channel.
func (c *Client) PostMessage(ctx context.Context, channel string, opts ...slack.MsgOption) (string, string, error) {
	return c.client.PostMessageContext(ctx, channel, opts...)
}

// PostEphemeral sends an ephemeral message to a user.
func (c *Client) PostEphemeral(ctx context.Context, channelID, userID string, opts ...slack.MsgOption) (string, error) {
	return c.client.PostEphemeralContext(ctx, channelID, userID, opts...)
}

// ChannelInfo holds basic info about a Slack channel.
type ChannelInfo struct {
	ID   string
	Name string
}
