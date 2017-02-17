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
	shortFields := []util.Field{}
	for _, att := range attachments {
		if len(att.Title) > 0 {
			lines = append(lines, fmt.Sprintf("%s**%s**", prefix, att.Title))
		}
		if len(att.Text) > 0 {
			lines = append(lines, fmt.Sprintf("%s%s", prefix, att.Text))
		}
		for _, field := range att.Fields {
			if field.Short {
				shortFields = append(shortFields, field)
				if len(shortFields) == 2 {
					fieldLines := BuildShortFieldLines(shortFields)
					if len(fieldLines) > 0 {
						lines = append(lines, fieldLines...)
					}
					shortFields = []util.Field{}
				}
				continue
			} else {
				if len(shortFields) > 0 {
					fieldLines := BuildShortFieldLines(shortFields)
					if len(fieldLines) > 0 {
						lines = append(lines, fieldLines...)
					}
				}
				shortFields = []util.Field{}
			}
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

func BuildShortFieldLines(shortFields []util.Field) []string {
	lines := []string{}
	prefix := GetGlipMarkdownBodyPrefix()
	for len(shortFields) > 0 {
		if len(shortFields) >= 2 {
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
