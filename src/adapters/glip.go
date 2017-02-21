package adapters

import (
	"fmt"

	cc "github.com/commonchat/commonchat-go"
	ccglip "github.com/commonchat/commonchat-go/glip"
	"github.com/grokify/glip-go-webhook"
	"github.com/valyala/fasthttp"
)

var (
	AdaptersGlipActivityIncludeIntegrationName = false
	AdaptersGlipMarkdownQuote                  = true
	AdaptersGlipUseShortFields                 = false
	EmojiURLFormat                             = ""
	WebhookURLOrUID                            = ""
)

type GlipAdapter struct {
	GlipClient      glipwebhook.GlipWebhookClient
	CommonConverter ccglip.GlipMessageConverter
	EmojiURLFormat  string
	WebhookURLOrUID string
}

func NewGlipAdapter(webhookURLOrUID string) (GlipAdapter, error) {
	glip, err := glipwebhook.NewGlipWebhookClient(webhookURLOrUID)
	return GlipAdapter{
		GlipClient:      glip,
		WebhookURLOrUID: webhookURLOrUID,
		CommonConverter: ccglip.NewGlipMessageConverter()}, err
}

func (adapter *GlipAdapter) SendWebhook(urlOrUid string, message cc.Message) (*fasthttp.Request, *fasthttp.Response, error) {
	return adapter.GlipClient.PostWebhookGUIDFast(urlOrUid, adapter.CommonConverter.ConvertCommonMessage(message))
}

func (adapter *GlipAdapter) SendMessage(message cc.Message) (*fasthttp.Request, *fasthttp.Response, error) {
	return adapter.SendWebhook(adapter.WebhookURLOrUID, message)
}

func (adapter *GlipAdapter) WebhookUID(ctx *fasthttp.RequestCtx) (string, error) {
	webhookUID := fmt.Sprintf("%s", ctx.UserValue("webhookuid"))
	return webhookUID, nil
}
