package adapters

import (
	"fmt"

	"github.com/commonchat/commonchat-go"
	"github.com/commonchat/commonchat-go/slack"
	"github.com/valyala/fasthttp"
)

type SlackAdapter struct {
	SlackClient     SlackWebhookClient
	EmojiURLFormat  string
	WebhookURLOrUID string
}

func NewSlackAdapter(webhookURLOrUID string) (*SlackAdapter, error) {
	client, err := NewSlackWebhookClient(webhookURLOrUID, FastHTTP)
	return &SlackAdapter{
		SlackClient:     client,
		WebhookURLOrUID: webhookURLOrUID}, err
}

func (adapter *SlackAdapter) SendWebhook(urlOrUid string, message commonchat.Message) (*fasthttp.Request, *fasthttp.Response, error) {
	return adapter.SlackClient.PostWebhookFast(urlOrUid, slack.ConvertCommonMessage(message))
}

func (adapter *SlackAdapter) SendMessage(message commonchat.Message) (*fasthttp.Request, *fasthttp.Response, error) {
	return adapter.SendWebhook(adapter.WebhookURLOrUID, message)
}

func (adapter *SlackAdapter) WebhookUID(ctx *fasthttp.RequestCtx) (string, error) {
	webhookUID := fmt.Sprintf("%s", ctx.UserValue("webhookuid"))
	return webhookUID, nil
}
