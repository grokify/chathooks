package slack

import (
	"fmt"

	"github.com/grokify/commonchat"
	//"github.com/grokify/commonchat/slack"
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

func (adapter *SlackAdapter) SendWebhook(urlOrUid string, ccMsg commonchat.Message, slackmsg interface{}) (*fasthttp.Request, *fasthttp.Response, error) {
	slackMessage := ConvertCommonMessage(ccMsg)
	slackmsg = &slackMessage
	return adapter.SlackClient.PostWebhookFast(urlOrUid, slackMessage)
}

func (adapter *SlackAdapter) SendMessage(message commonchat.Message, slackmsg interface{}) (*fasthttp.Request, *fasthttp.Response, error) {
	return adapter.SendWebhook(adapter.WebhookURLOrUID, message, slackmsg)
}

func (adapter *SlackAdapter) WebhookUID(ctx *fasthttp.RequestCtx) (string, error) {
	webhookUID := fmt.Sprintf("%s", ctx.UserValue("webhookuid"))
	return webhookUID, nil
}
