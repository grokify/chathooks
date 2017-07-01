package slack

import (
	"encoding/json"
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	cc "github.com/commonchat/commonchat-go"
	"github.com/grokify/chatmore/src/adapters"
	"github.com/grokify/chatmore/src/config"
	"github.com/grokify/chatmore/src/util"
	"github.com/valyala/fasthttp"
)

const (
	DisplayName = "Slack"
	HandlerKey  = "slack"
)

// FastHttp request handler constructor for Slack inbound webhook
type Handler struct {
	Config  config.Configuration
	Adapter adapters.Adapter
}

// FastHttp request handler constructor for Slack in bound webhook
func NewHandler(config config.Configuration, adapter adapters.Adapter) Handler {
	return Handler{Config: config, Adapter: adapter}
}

// HandleFastHTTP is the method to respond to a fasthttp request.
func (h *Handler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	ccMsg, err := Normalize(BuildInboundMessageBytes(ctx))
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

func BuildInboundMessageBytes(ctx *fasthttp.RequestCtx) []byte {
	ct := string(ctx.Request.Header.Peek("Content-Type"))
	ct = strings.TrimSpace(strings.ToLower(ct))
	if ct == "application/json" {
		return ctx.PostBody()
	}
	return ctx.FormValue("payload")
}

func Normalize(bytes []byte) (cc.Message, error) {
	slack, err := SlackWebhookMessageFromBytes(bytes)
	if err != nil {
		return cc.Message{}, err
	}

	message := cc.Message{
		Activity:  slack.Username,
		Text:      slack.Text,
		IconEmoji: slack.IconEmoji,
		IconURL:   slack.IconURL}
	return message, nil
}

type SlackWebhookMessage struct {
	Username  string `json:"username"`
	IconEmoji string `json:"icon_emoji"`
	IconURL   string `json:"icon_url"`
	Text      string `json:"text"`
}

func SlackWebhookMessageFromBytes(bytes []byte) (SlackWebhookMessage, error) {
	msg := SlackWebhookMessage{}
	err := json.Unmarshal(bytes, &msg)
	return msg, err
}
