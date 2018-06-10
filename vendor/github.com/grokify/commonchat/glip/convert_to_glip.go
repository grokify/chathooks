package glip

import (
	"fmt"
	"strings"

	cc "github.com/grokify/commonchat"
	"github.com/grokify/go-glip"
)

type GlipMessageConverter struct {
	EmojiURLFormat                 string
	ActivityIncludeIntegrationName bool
	UseAttachments                 bool // overrides other 'use' options
	UseMarkdownQuote               bool
	UseShortFields                 bool
	UseFieldExtraSpacing           bool
}

func NewGlipMessageConverter() GlipMessageConverter {
	return GlipMessageConverter{}
}

func (cv *GlipMessageConverter) ConvertCommonMessage(commonMessage cc.Message) glipwebhook.GlipWebhookMessage {
	glip := glipwebhook.GlipWebhookMessage{
		Activity: commonMessage.Activity,
		Title:    commonMessage.Title,
		Icon:     commonMessage.IconURL}

	if len(commonMessage.IconURL) > 0 {
		glip.Icon = commonMessage.IconURL
	} else if len(commonMessage.IconEmoji) > 0 {
		iconURL, err := cc.EmojiToURL(cv.EmojiURLFormat, commonMessage.IconEmoji)
		if err == nil {
			glip.Icon = iconURL
		}
	}
	bodyLines := []string{}
	if len(commonMessage.Text) > 0 {
		bodyLines = append(bodyLines, commonMessage.Text)
	}

	if len(commonMessage.Attachments) > 0 {
		if cv.UseAttachments {
			glip.Attachments = convertAttachments(commonMessage.Attachments)
		} else {
			attachmentText := cv.renderAttachmentsAsMarkdown(commonMessage.Attachments)
			if len(attachmentText) > 0 {
				bodyLines = append(bodyLines, attachmentText)
			}
		}
	}

	if len(bodyLines) > 0 {
		glip.Body = strings.Join(bodyLines, "\n")
	}
	return glip
}

func (cv *GlipMessageConverter) getMarkdownBodyPrefix() string {
	if cv.UseMarkdownQuote {
		return "> "
	}
	return ""
}

func convertAttachments(commonAttachments []cc.Attachment) []glipwebhook.Attachment {
	glipAttachments := []glipwebhook.Attachment{}
	for _, commonAttachment := range commonAttachments {
		glipAttachments = append(glipAttachments, convertAttachment(commonAttachment))
	}
	return glipAttachments
}

func convertAttachment(commonAttachment cc.Attachment) glipwebhook.Attachment {
	return glipwebhook.Attachment{
		Type:    "Card",
		Title:   commonAttachment.Title,
		Pretext: commonAttachment.Pretext,
		Text:    commonAttachment.Text,
		Color:   commonAttachment.Color,
		Fields:  convertFields(commonAttachment.Fields)}
}

func convertFields(commonFields []cc.Field) []glipwebhook.Field {
	glipFields := []glipwebhook.Field{}
	for _, commonField := range commonFields {
		glipFields = append(glipFields, convertField(commonField))
	}
	return glipFields
}

func convertField(commonField cc.Field) glipwebhook.Field {
	return glipwebhook.Field{
		Title: commonField.Title,
		Value: commonField.Value,
		Short: commonField.Short}
}

func (cv *GlipMessageConverter) renderAttachmentsAsMarkdown(attachments []cc.Attachment) string {
	lines := []string{}
	prefix := cv.getMarkdownBodyPrefix()
	shortFields := []cc.Field{}
	for _, att := range attachments {
		if len(att.Title) > 0 {
			lines = append(lines, fmt.Sprintf("%s**%s**", prefix, att.Title))
		}
		if len(att.Text) > 0 {
			lines = append(lines, fmt.Sprintf("%s%s", prefix, att.Text))
		}
		for _, field := range att.Fields {
			if !cv.UseShortFields {
				field.Short = false
			}
			if field.Short {
				shortFields = append(shortFields, field)
				if len(shortFields) == 2 {
					fieldLines := cv.buildMarkdownShortFieldLines(shortFields)
					if len(fieldLines) > 0 {
						lines = cv.appendEmptyLine(lines)
						lines = append(lines, fieldLines...)
					}
					shortFields = []cc.Field{}
				}
				continue
			} else {
				if len(shortFields) > 0 {
					fieldLines := cv.buildMarkdownShortFieldLines(shortFields)
					if len(fieldLines) > 0 {
						lines = cv.appendEmptyLine(lines)
						lines = append(lines, fieldLines...)
					}
				}
				shortFields = []cc.Field{}
			}
			if len(field.Title) > 0 || len(field.Value) > 0 {
				lines = cv.appendEmptyLine(lines)
				if len(field.Title) > 0 {
					lines = append(lines, fmt.Sprintf("%s**%s**", prefix, field.Title))
				}
				if len(field.Value) > 0 {
					lines = append(lines, fmt.Sprintf("%s%s", prefix, field.Value))
				}
			}
		}
	}
	return strings.Join(lines, "\n")
}

func (cv *GlipMessageConverter) buildMarkdownShortFieldLines(shortFields []cc.Field) []string {
	lines := []string{}
	prefix := cv.getMarkdownBodyPrefix()
	for len(shortFields) > 0 {
		if len(shortFields) >= 2 {
			lines = cv.appendEmptyLine(lines)
			field1 := shortFields[0]
			field2 := shortFields[1]
			if len(field2.Title) > 0 || len(field2.Title) > 0 {
				lines = append(lines, fmt.Sprintf("%s| **%v** | **%v** |", prefix, field1.Title, field2.Title))
			}
			if len(field2.Value) > 0 || len(field2.Value) > 0 {
				lines = append(lines, fmt.Sprintf("%s| %v | %v |", prefix, field1.Value, field2.Value))
			}
			shortFields = shortFields[2:]
		} else {
			lines = cv.appendEmptyLine(lines)
			field1 := shortFields[0]
			if len(field1.Title) > 0 {
				lines = append(lines, fmt.Sprintf("%s**%s**", prefix, field1.Title))
			}
			if len(field1.Value) > 0 {
				lines = append(lines, fmt.Sprintf("%s%s", prefix, field1.Value))
			}
			shortFields = shortFields[1:]
		}
	}
	return lines
}

func (cv *GlipMessageConverter) appendEmptyLine(lines []string) []string {
	if cv.UseFieldExtraSpacing {
		if len(lines) > 0 {
			if len(lines[len(lines)-1]) > 0 {
				lines = append(lines, "")
			}
		}
	}
	return lines
}

/*
func (cv *GlipMessageConverter) RenderMessage(message cc.Message) string {
	lines := []string{}
	attachments := cv.RenderAttachments(message.Attachments)
	if len(attachments) > 0 {
		lines = append(lines, attachments)
	}
	return strings.Join(lines, "\n")
}
*/

func (cv *GlipMessageConverter) integrationActivitySuffix(displayName string) string {
	if cv.ActivityIncludeIntegrationName {
		return fmt.Sprintf(" (%v)", displayName)
	}
	return ""
}
