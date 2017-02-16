package glipadapter

import (
	"fmt"
	"strings"

	"github.com/grokify/glip-webhook-proxy-go/src/util"
)

var (
	AdaptersGlipActivityIncludeIntegrationName = false
	AdaptersGlipMarkdownQuote                  = false
)

func GetGlipMarkdownBodyPrefix() string {
	if AdaptersGlipMarkdownQuote {
		return "> "
	}
	return ""
}

func RenderAttachments(attachments []util.Attachment) string {
	lines := []string{}
	prefix := GetGlipMarkdownBodyPrefix()
	for _, att := range attachments {
		if len(att.Title) > 0 {
			lines = append(lines, fmt.Sprintf("%s**%s**", prefix, att.Title))
		}
		if len(att.Text) > 0 {
			lines = append(lines, fmt.Sprintf("%s%s", prefix, att.Text))
		}
		for _, field := range att.Fields {
			if len(field.Title) > 0 {
				lines = append(lines, fmt.Sprintf("%s**%s**", prefix, field.Title))
			}
			if len(field.Value) > 0 {
				lines = append(lines, fmt.Sprintf("%s%s", prefix, field.Value))
			}
		}
	}
	return strings.Join(lines, "\n")
}

func RenderMessage(message util.Message) string {
	lines := []string{}
	attachments := RenderAttachments(message.Attachments)
	if len(attachments) > 0 {
		lines = append(lines, attachments)
	}
	return strings.Join(lines, "\n")
}

func IntegrationActivitySuffix(displayName string) string {
	if AdaptersGlipActivityIncludeIntegrationName {
		return fmt.Sprintf(" (%v)", displayName)
	}
	return ""
}
