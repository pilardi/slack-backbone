package slack

import "github.com/slack-go/slack"

// NewMessageBlock creates an empty message block ready to be populated.
func NewMessageBlock() *slack.MessageBlock {
	return &slack.MessageBlock{}
}

// SectionBlock creates a text section block.
func SectionBlock(text string) any {
	return &slack.SectionBlock{
		Type: "section",
		Text: &slack.TextBlockObject{Type: "mrkdwn", Text: text},
	}
}

// DividerBlock creates a divider between sections.
func DividerBlock() any {
	return &slack.DividerBlock{Type: "divider"}
}

// ActionBlockWithButton creates an action block with a single button.
func ActionBlockWithButton(blockID, text, value string) *slack.ActionBlock {
	return &slack.ActionBlock{
		BlockID: blockID,
		Elements: &slack.ElementBlock{
			Type:  "button",
			Text:  &slack.TextBlockObject{Type: "plain_text", Text: text},
			Value: value,
		},
	}
}
