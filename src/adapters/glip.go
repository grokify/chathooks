package adapters

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"
	cc "github.com/commonchat/commonchat-go"
	ccglip "github.com/commonchat/commonchat-go/glip"
	"github.com/grokify/glip-go-webhook"
	"github.com/valyala/fasthttp"
)

var (
	AdaptersGlipActivityIncludeIntegrationName = false
	AdaptersGlipMarkdownQuote                  = false
	AdaptersGlipUseAttachments                 = true
	AdaptersGlipUseShortFields                 = false
	AdatpersGlipUseFieldExtraSpacing           = true
	EmojiURLFormat                             = ""
	WebhookURLOrUID                            = ""
)

type GlipAdapter struct {
	GlipClient      glipwebhook.GlipWebhookClient
	CommonConverter ccglip.GlipMessageConverter
	EmojiURLFormat  string
	WebhookURLOrUID string
}

func NewGlipAdapter(webhookURLOrUID string) (*GlipAdapter, error) {
	glip, err := glipwebhook.NewGlipWebhookClient(webhookURLOrUID)
	converter := ccglip.NewGlipMessageConverter()
	converter.UseAttachments = AdaptersGlipUseAttachments
	converter.UseShortFields = AdaptersGlipUseShortFields
	converter.UseFieldExtraSpacing = AdatpersGlipUseFieldExtraSpacing
	return &GlipAdapter{
		GlipClient:      glip,
		WebhookURLOrUID: webhookURLOrUID,
		CommonConverter: converter}, err
}

func (adapter *GlipAdapter) SendWebhook(urlOrUid string, message cc.Message) (*fasthttp.Request, *fasthttp.Response, error) {
	glipMessage := adapter.CommonConverter.ConvertCommonMessage(message)

	glipMessageString, err := json.Marshal(glipMessage)
	if err == nil {
		log.WithFields(log.Fields{
			"event":   "outgoing.webhook.glip",
			"handler": "Glip Adapter"}).Info(string(glipMessageString))
	}
	return adapter.GlipClient.PostWebhookGUIDFast(urlOrUid, glipMessage)
}

func (adapter *GlipAdapter) SendMessage(message cc.Message) (*fasthttp.Request, *fasthttp.Response, error) {
	return adapter.SendWebhook(adapter.WebhookURLOrUID, message)
}

func (adapter *GlipAdapter) WebhookUID(ctx *fasthttp.RequestCtx) (string, error) {
	webhookUID := adapter.WebhookURLOrUID
	return webhookUID, nil
}
