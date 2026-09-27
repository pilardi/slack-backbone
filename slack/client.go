package slack

import (
	"context"

	"github.com/slack-go/slack"
)

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

// PostMessage posts a message to a channel.
func (c *Client) PostMessage(ctx context.Context, channel string, opts ...slack.MsgOption) (string, string, error) {
	return c.client.PostMessageContext(ctx, channel, opts...)
}

// PostEphemeral sends an ephemeral message to a user.
func (c *Client) PostEphemeral(ctx context.Context, channelID, userID string, opts ...slack.MsgOption) (string, error) {
	return c.client.PostEphemeralContext(ctx, channelID, userID, opts...)
}
