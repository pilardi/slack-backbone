package slack

import "github.com/slack-go/slack"

// Client wraps the Slack API client with multi-team awareness.
type Client struct {
	client *slack.Client
}

// NewClient creates a new API client.
func NewClient() *Client {
	return &Client{
		client: slack.New(),
	}
}

// PostMessage posts a message to a channel.
func (c *Client) PostMessage(ctx context.Context, channel string, opts ...slack.MsgOption) (*slack.PostMessageResponse, error) {
	return c.client.PostMessageContext(ctx, channel, opts...)
}

// PostEphemeral sends an ephemeral message to a user.
func (c *Client) PostEphemeral(ctx context.Context, userID string, text string, opts ...slack.EphemeralOption) (*slack.MessageUpdatedResponse, error) {
	return c.client.PostEphemeralContext(ctx, userID, text, opts...)
}
