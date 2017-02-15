package glipadapter

import (
	"fmt"
	"strings"

	"github.com/grokify/glip-webhook-proxy-go/src/util"
)

var (
	AdaptersGlipMarkdownQuote = false
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
