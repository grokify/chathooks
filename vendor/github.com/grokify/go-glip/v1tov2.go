package glipwebhook

import (
	"regexp"
	"strings"
	"time"

	v2 "github.com/grokify/go-glip/v2"
)

const (
	webhookV2Path          string = "/webhook/v2/"
	rxGlipWebhookV2Pattern string = `^https?://[^/]+/webhook/v2/[^/]+/?$`
	rxGlipWebhookV1Pattern string = `^(?i)(https?://[^/]+)/webhook/([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})/?$`
)

var rxGlipWebhookV1 = regexp.MustCompile(rxGlipWebhookV1Pattern)
var rxGlipWebhookV2 = regexp.MustCompile(rxGlipWebhookV2Pattern)

func V1ToV2WewbhookUri(input string) string {
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return input
	}
	if strings.Index(input, "/") == -1 {
		return GlipWebhookBaseURLProductionV2 + input
	}
	if rxGlipWebhookV2.MatchString(input) {
		return input
	}
	m := rxGlipWebhookV1.FindStringSubmatch(input)
	if len(m) == 3 {
		return m[1] + webhookV2Path + m[2]
	}
	return input
}

func V1ToV2WebhookBody(v1msg GlipWebhookMessage) v2.GlipWebhookMessage {
	v2msg := v2.GlipWebhookMessage{
		Activity:    v1msg.Activity,
		IconUri:     v1msg.Icon,
		Text:        v1msg.Body,
		Title:       v1msg.Title,
		Attachments: []v2.Attachment{}}
	for _, v1att := range v1msg.Attachments {
		v2msg.Attachments = append(v2msg.Attachments, V1ToV2WebhookAttachment(v1att))
	}

	return v2msg
}

func V1ToV2WebhookAttachment(v1att Attachment) v2.Attachment {
	v2att := v2.Attachment{
		Color:        v1att.Color,
		Fields:       []v2.Field{},
		ImageUri:     v1att.ImageURL,
		Intro:        v1att.Pretext,
		Text:         v1att.Text,
		ThumbnailUri: v1att.ThumbnailURL,
		Title:        v1att.Title,
		Type:         v1att.Type}
	if len(strings.TrimSpace(v2att.Type)) == 0 {
		v2att.Type = AttachmentTypeCard
	}
	if len(v1att.AuthorName) > 0 {
		v2att.Author = &v2.Author{
			Name:    v1att.AuthorName,
			IconUri: v1att.AuthorIcon,
			Uri:     v1att.AuthorLink}
	}
	if len(strings.TrimSpace(v1att.FooterIcon)) > 0 || len(strings.TrimSpace(v1att.Footer)) > 0 {
		v2att.Footnote = &v2.Footnote{
			IconUri: v1att.FooterIcon,
			Text:    v1att.Footer,
		}
		if v1att.TS > 0 {
			v2att.Footnote.Time = time.Unix(v1att.TS, 0)
		}
	}

	for _, v1field := range v1att.Fields {
		v2field := v2.Field{
			Title: v1field.Title,
			Value: v1field.Value}
		if v1field.Short {
			v2field.Style = v2.FieldStyleShort
		} else {
			v2field.Style = v2.FieldStyleLong
		}
		v2att.Fields = append(v2att.Fields, v2field)
	}

	return v2att
}
