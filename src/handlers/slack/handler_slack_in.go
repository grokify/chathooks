package slack

import (
	"strings"

	"github.com/grokify/chathooks/src/config"
	"github.com/grokify/chathooks/src/handlers"
	"github.com/grokify/chathooks/src/models"
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

/*
// FastHttp request handler constructor for Slack inbound webhook
type Handler struct {
	Config  config.Configuration
	Adapter adapters.Adapter
}

// FastHttp request handler constructor for Slack in bound webhook
func NewHandler(config config.Configuration, adapter adapters.Adapter) Handler {
	return Handler{Config: config, Adapter: adapter}
}

func (h Handler) HandlerKey() string {
	return HandlerKey
}

func (h Handler) MessageDirection() string {
	return MessageDirection
}

// HandleFastHTTP is the method to respond to a fasthttp request.
func (h Handler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	ccMsg, err := Normalize(h.Config, BuildInboundMessageBytes(ctx))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotAcceptable)
		log.WithFields(log.Fields{
			"type":   "http.response",
			"status": fasthttp.StatusNotAcceptable,
		}).Info(fmt.Sprintf("%v request is not acceptable.", DisplayName))
		return
	}

	util.SendWebhook(ctx, h.Adapter, ccMsg)
}
*/

func BuildInboundMessageBytes(ctx *fasthttp.RequestCtx) []byte {
	ct := string(ctx.Request.Header.Peek("Content-Type"))
	ct = strings.TrimSpace(strings.ToLower(ct))
	if ct == "application/json" {
		return ctx.PostBody()
	}
	return ctx.FormValue("payload")
}

func Normalize(config config.Configuration, hReq handlers.HandlerRequest) (cc.Message, error) {
	slMsg, err := ccslack.NewMessageFromBytes(hReq.Body)
	if err != nil {
		return cc.Message{}, err
	}

	ccMsg := ccslack.WebhookInBodySlackToCc(slMsg)

	return ccMsg, nil
	/*
		message := cc.Message{
			Activity:  slack.Username,
			Text:      slack.Text,
			IconEmoji: slack.IconEmoji,
			IconURL:   slack.IconURL}
		return message, nil*/
}

/*
type SlackWebhookMessage struct {
	Username  string `json:"username"`
	IconEmoji string `json:"icon_emoji"`
	IconURL   string `json:"icon_url"`
	Text      string `json:"text"`
}

func SlackWebhookMessageFromBytes(bytes []byte) (SlackWebhookMessage, error) {
	msg := webhook.SlackWebhookMessage{}
	err := json.Unmarshal(bytes, &msg)
	return msg, err
}
*/
