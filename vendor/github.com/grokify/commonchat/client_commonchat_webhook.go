package commonchat

import (
	"github.com/valyala/fasthttp"
)

type Adapter interface {
	SendWebhook(url string, ccMsg Message, formattedMsg interface{}) (*fasthttp.Request, *fasthttp.Response, error)
	SendMessage(ccMsg Message, formattedMsg interface{}) (*fasthttp.Request, *fasthttp.Response, error)
	WebhookUID(ctx *fasthttp.RequestCtx) (string, error)
}
