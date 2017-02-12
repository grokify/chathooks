package slack

import (
	"encoding/json"
	"strings"

	"github.com/grokify/glip-go-webhook"
	"github.com/grokify/glip-webhook-proxy-go/config"
	"github.com/grokify/glip-webhook-proxy-go/util"
	"github.com/valyala/fasthttp"
)

// FastHttp request handler constructor for Slack inbound webhook
type SlackToGlipHandler struct {
	Config     config.Configuration
	GlipClient glipwebhook.GlipWebhookClient
}

// FastHttp request handler constructor for Slack in bound webhook
func NewSlackToGlipHandler(config config.Configuration, glip glipwebhook.GlipWebhookClient) SlackToGlipHandler {
	return SlackToGlipHandler{Config: config, GlipClient: glip}
}

// HandleFastHTTP is the method to respond to a fasthttp request.
func (h *SlackToGlipHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	slackMsg, err := BuildInboundMessage(ctx)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotAcceptable)
		return
	}
	glipMsg := Normalize(slackMsg, h.Config.EmojiURLFormat)

	util.SendGlipWebhookCtx(ctx, h.GlipClient, glipMsg)
}

func BuildInboundMessage(ctx *fasthttp.RequestCtx) (SlackWebhookMessage, error) {
	ct := string(ctx.Request.Header.Peek("Content-Type"))
	ct = strings.TrimSpace(strings.ToLower(ct))
	if ct == "application/json" {
		return SlackWebhookMessageFromBytes(ctx.PostBody())
	}
	return SlackWebhookMessageFromBytes(ctx.FormValue("payload"))
}

func Normalize(slack SlackWebhookMessage, emojiURLFormat string) glipwebhook.GlipWebhookMessage {
	gmsg := glipwebhook.GlipWebhookMessage{
		Body:     slack.Text,
		Activity: slack.Username}
	if len(slack.IconURL) > 0 {
		gmsg.Icon = slack.IconURL
	} else {
		iconURL, err := util.EmojiToURL(emojiURLFormat, slack.IconEmoji)
		if err == nil {
			gmsg.Icon = iconURL
		}
	}
	return gmsg
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
