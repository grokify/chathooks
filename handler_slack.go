package glipwebhookproxy

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/grokify/glip-go-webhook"
	"github.com/valyala/fasthttp"
)

type SlackToGlipHandler struct {
	Config Configuration
}

func (h *SlackToGlipHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	slackMsg, err := h.BuildSlackMessage(ctx)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotAcceptable)
		return
	}

	glipMsg := h.SlackToGlip(slackMsg)
	glipWebhookGuid := fmt.Sprintf("%s", ctx.UserValue("glipguid"))

	client, err := glipwebhook.NewGlipWebhookClient(glipWebhookGuid)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotAcceptable)
		return
	}

	req, resp, err := client.PostMessageFast(glipMsg)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
		return
	}
	fmt.Fprintf(ctx, "%s", string(resp.Body()))
	fasthttp.ReleaseRequest(req)
	fasthttp.ReleaseResponse(resp)
}

func (h *SlackToGlipHandler) BuildSlackMessage(ctx *fasthttp.RequestCtx) (SlackWebhookMessage, error) {
	ct := string(ctx.Request.Header.Peek("Content-Type"))
	ct = strings.TrimSpace(strings.ToLower(ct))
	if ct == "application/json" {
		return SlackWebhookMessageFromBytes(ctx.PostBody())
	}
	return SlackWebhookMessageFromBytes(ctx.FormValue("payload"))
}

func (h *SlackToGlipHandler) SlackToGlip(slack SlackWebhookMessage) glipwebhook.GlipWebhookMessage {
	gmsg := glipwebhook.GlipWebhookMessage{
		Body:     slack.Text,
		Activity: slack.Username}
	if len(slack.IconURL) > 0 {
		gmsg.Icon = slack.IconURL
	} else {
		iconURL, err := EmojiToURL(h.Config.EmojiURLFormat, slack.IconEmoji)
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
