package adapters

import (
	cc "github.com/commonchat/commonchat-go"
	"github.com/valyala/fasthttp"
)

var (
	ShowDisplayName = false
)

type Adapter interface {
	SendWebhook(url string, message cc.Message) (*fasthttp.Request, *fasthttp.Response, error)
	SendMessage(message cc.Message) (*fasthttp.Request, *fasthttp.Response, error)
	WebhookUID(ctx *fasthttp.RequestCtx) (string, error)
}

func IntegrationActivitySuffix(displayName string) string {
	if !ShowDisplayName || len(displayName) < 1 {
		return ""
	}
	return ""
}
