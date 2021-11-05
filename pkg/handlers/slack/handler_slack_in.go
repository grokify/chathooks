package slack

import (
	"strings"

	"github.com/grokify/chathooks/pkg/config"
	"github.com/grokify/chathooks/pkg/handlers"
	"github.com/grokify/chathooks/pkg/models"
	cc "github.com/grokify/commonchat"
	"github.com/valyala/fasthttp"

	ccslack "github.com/grokify/commonchat/slack"
)

const (
	DisplayName      = "Slack"
	HandlerKey       = "slack"
	MessageDirection = "in"
	MessageBodyType  = models.URLEncodedJSONPayloadOrJSON
)

func NewHandler() handlers.Handler {
	return handlers.Handler{MessageBodyType: MessageBodyType, Normalize: Normalize}
}

func BuildInboundMessageBytes(ctx *fasthttp.RequestCtx) []byte {
	ct := string(ctx.Request.Header.Peek("Content-Type"))
	ct = strings.TrimSpace(strings.ToLower(ct))
	if ct == "application/json" {
		return ctx.PostBody()
	}
	return ctx.FormValue("payload")
}

func Normalize(config config.Configuration, hReq handlers.HandlerRequest) (cc.Message, error) {
	slMsg, err := ccslack.ParseMessageHttpBody(hReq.Body)
	if err != nil {
		return cc.Message{}, err
	}

	ccMsg := ccslack.WebhookInBodySlackToCc(slMsg)

	return ccMsg, nil
}
