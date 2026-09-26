package slack

import (
	"github.com/slack-go/slack"
)

// NewBlocks creates an empty Blocks container ready to be populated.
func NewBlocks() *slack.Blocks {
	return &slack.Blocks{}
}

// SectionBlock creates a text section block.
func SectionBlock(text string) any {
	return slack.NewSectionBlock(&slack.TextBlockObject{Type: "mrkdwn", Text: text}, nil, nil)
}

// DividerBlock creates a divider between sections.
func DividerBlock() any {
	return slack.NewDividerBlock()
}

// ActionBlockWithButton creates an action block with a single button.
func ActionBlockWithButton(blockID, text, value string) *slack.ActionBlock {
	return slack.NewActionBlock(blockID, &slack.ButtonBlockElement{
		Type:  "button",
		Text:  &slack.TextBlockObject{Type: "plain_text", Text: text},
		Value: value,
	})
}
