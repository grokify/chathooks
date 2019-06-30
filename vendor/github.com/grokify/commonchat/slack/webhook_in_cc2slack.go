package slack

import (
	"regexp"
	"strings"

	cc "github.com/grokify/commonchat"
	"github.com/grokify/gotilla/fmt/fmtutil"
)

func ConvertCommonMessage(ccMsg cc.Message) Message {
	slackMessage := Message{
		Attachments: ConvertAttachmentsSlack(ccMsg.Attachments),
		IconEmoji:   ccMsg.IconEmoji,
		IconURL:     ccMsg.IconURL,
		Mrkdwn:      true}

	textLines := []string{}
	if len(ccMsg.Activity) > 0 {
		textLines = append(textLines, ccMsg.Activity)
	}
	if len(ccMsg.Title) > 0 {
		textLines = append(textLines, ccMsg.Title)
	}
	if len(ccMsg.Text) > 0 {
		textLines = append(textLines, ccMsg.Text)
	}
	if len(textLines) > 0 {
		text := strings.Join(textLines, "\n")
		text = ConvertMarkdownSlack(text)
		slackMessage.Text = text
	}
	fmtutil.PrintJSON(slackMessage)
	return slackMessage
}

func ConvertMarkdownSlack(markdown string) string {
	slack := markdown
	re1 := regexp.MustCompile(`\[([^\[\]]+)\]\((.*?)\)`)
	slack = re1.ReplaceAllString(slack, "<$2|$1>")
	re2 := regexp.MustCompile(`\*\*([^\*]+?)\*\*`)
	slack = re2.ReplaceAllString(slack, "*$1*")
	return slack
}

func ConvertAttachmentsSlack(commonAttachments []cc.Attachment) []Attachment {
	slackAttachments := []Attachment{}
	for _, commonAttachment := range commonAttachments {
		slackAttachments = append(slackAttachments, ConvertAttachmentSlack(commonAttachment))
	}
	return slackAttachments
}

func ConvertAttachmentSlack(commonAttachment cc.Attachment) Attachment {
	slackAttachment := Attachment{
		Title:        ConvertMarkdownSlack(commonAttachment.Title),
		Pretext:      ConvertMarkdownSlack(commonAttachment.Pretext),
		Text:         ConvertMarkdownSlack(commonAttachment.Text),
		Color:        commonAttachment.Color,
		Fields:       ConvertFieldsSlack(commonAttachment.Fields),
		ThumbnailURL: commonAttachment.ThumbnailURL,
		MarkdownIn:   []string{"text"}}
	return slackAttachment
}

func ConvertFieldsSlack(commonFields []cc.Field) []Field {
	slackFields := []Field{}
	for _, commonField := range commonFields {
		slackFields = append(slackFields, ConvertFieldSlack(commonField))
	}
	return slackFields
}

func ConvertFieldSlack(commonField cc.Field) Field {
	slackField := Field{
		Title: ConvertMarkdownSlack(commonField.Title),
		Value: ConvertMarkdownSlack(commonField.Value),
		Short: commonField.Short}
	return slackField
}
