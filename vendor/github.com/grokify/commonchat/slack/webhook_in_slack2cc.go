package slack

import (
	cc "github.com/grokify/commonchat"
	"github.com/grokify/gotilla/text/markdown"
)

func WebhookInBodySlackToCc(slMsg Message) cc.Message {
	ccMsg := cc.Message{
		Activity:    slMsg.Username,
		Attachments: []cc.Attachment{},
		Text:        slMsg.Text,
		IconEmoji:   slMsg.IconEmoji,
		IconURL:     slMsg.IconURL}
	for _, slAtt := range slMsg.Attachments {
		ccMsg.Attachments = append(ccMsg.Attachments, attachmentSlackToCc(slAtt))
	}
	return ccMsg
}

func attachmentSlackToCc(slAtt Attachment) cc.Attachment {
	ccAtt := cc.Attachment{
		AuthorIcon:   slAtt.AuthorIcon,
		AuthorLink:   slAtt.AuthorLink,
		AuthorName:   slAtt.AuthorName,
		Color:        slAtt.Color,
		Fallback:     slAtt.Fallback,
		Fields:       []cc.Field{},
		MarkdownIn:   slAtt.MarkdownIn,
		Pretext:      markdown.SkypeToMarkdown(slAtt.Pretext, true),
		Text:         markdown.SkypeToMarkdown(slAtt.Text, true),
		ThumbnailURL: slAtt.ThumbnailURL,
		Title:        slAtt.Title,
	}
	for _, slField := range slAtt.Fields {
		ccAtt.Fields = append(ccAtt.Fields, fieldSlackToCc(slField))
	}
	return ccAtt
}

func fieldSlackToCc(slField Field) cc.Field {
	ccField := cc.Field{
		Short: slField.Short,
		Title: slField.Title,
		Value: markdown.SkypeToMarkdown(slField.Value, true)}
	return ccField
}
