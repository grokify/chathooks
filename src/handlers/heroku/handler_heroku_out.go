package heroku

import (
	"fmt"
	"net/url"
	"strings"

	log "github.com/Sirupsen/logrus"

	cc "github.com/commonchat/commonchat-go"
	"github.com/grokify/glip-go-webhook"
	"github.com/grokify/webhookproxy/src/adapters"
	"github.com/grokify/webhookproxy/src/config"
	"github.com/grokify/webhookproxy/src/util"
	"github.com/valyala/fasthttp"
)

const (
	DisplayName      = "Heroku"
	HandlerKey       = "heroku"
	MessageDirection = "out"
	IconURL          = "https://a.slack-edge.com/ae7f/plugins/heroku/assets/service_512.png"
)

// FastHttp request handler for Heroku outbound webhook
// https://devcenter.heroku.com/articles/deploy-hooks#http-post-hook
type Handler struct {
	Config     config.Configuration
	GlipClient glipwebhook.GlipWebhookClient
	Adapter    adapters.Adapter
}

// FastHttp request handler constructor for Confluence outbound webhook
func NewHandler(cfg config.Configuration, adapter adapters.Adapter) Handler {
	return Handler{Config: cfg, Adapter: adapter}
}

func (h Handler) HandlerKey() string {
	return HandlerKey
}

func (h Handler) MessageDirection() string {
	return MessageDirection
}

// HandleFastHTTP is the method to respond to a fasthttp request.
func (h *Handler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	srcMsg, err := BuildInboundMessage(ctx)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotAcceptable)
		log.WithFields(log.Fields{
			"type":   "http.response",
			"status": fasthttp.StatusNotAcceptable,
		}).Info(fmt.Sprintf("%v request is not acceptable.", DisplayName))
		return
	}

	ccMsg, err := NormalizeHerokuMessage(srcMsg)
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

func BuildInboundMessage(ctx *fasthttp.RequestCtx) (HerokuOutMessage, error) {
	return HerokuOutMessage{
		App:      string(ctx.FormValue("app")),
		User:     string(ctx.FormValue("user")),
		URL:      string(ctx.FormValue("url")),
		Head:     string(ctx.FormValue("head")),
		HeadLong: string(ctx.FormValue("head_long")),
		PrevHead: string(ctx.FormValue("prev_head")),
		GitLog:   string(ctx.FormValue("git_log")),
		Release:  string(ctx.FormValue("release"))}, nil
}

//func Normalize(src HerokuOutMessage) glipwebhook.GlipWebhookMessage {
func Normalize(bytes []byte) (cc.Message, error) {
	src, err := HerokuOutMessageFromQueryString(string(bytes))
	if err != nil {
		return cc.Message{}, err
	}
	return NormalizeHerokuMessage(src)
}

func NormalizeHerokuMessage(src HerokuOutMessage) (cc.Message, error) {
	message := cc.NewMessage()
	message.IconURL = IconURL

	if len(strings.TrimSpace(src.App)) > 0 {
		message.Activity = fmt.Sprintf("%v deployed on %v\n\n", src.App, DisplayName)
	} else {
		message.Activity = fmt.Sprintf("An app has been deployed on %v", DisplayName)
	}

	attachment := cc.NewAttachment()

	if len(strings.TrimSpace(src.App)) > 0 {
		field := cc.Field{Title: "Application"}
		if len(src.App) < 35 {
			field.Short = true
		}
		if len(strings.TrimSpace(src.URL)) > 0 {
			field.Value = fmt.Sprintf("[%s](%s)", src.App, src.URL)
		} else {
			field.Value = src.App
		}
		attachment.AddField(field)
	}
	if len(strings.TrimSpace(src.Release)) > 0 {
		attachment.AddField(cc.Field{Title: "Release", Value: src.Release, Short: true})
	}
	if len(strings.TrimSpace(src.User)) > 0 {
		attachment.AddField(cc.Field{Title: "User", Value: src.User, Short: true})
	}

	message.AddAttachment(attachment)
	return message, nil
}

type HerokuOutMessage struct {
	App      string `json:"app,omitempty"`
	User     string `json:"user,omitempty"`
	URL      string `json:"url,omitempty"`
	Head     string `json:"head,omitempty"`
	HeadLong string `json:"head_long,omitempty"`
	PrevHead string `json:"prev_head,omitempty"`
	GitLog   string `json:"git_log,omitempty"`
	Release  string `json:"release,omitempty"`
}

func HerokuOutMessageFromQueryString(queryString string) (HerokuOutMessage, error) {
	msg := HerokuOutMessage{}
	values, err := url.ParseQuery(queryString)
	if err != nil {
		return msg, err
	}
	if len(strings.TrimSpace(values.Get("app"))) > 0 {
		msg.App = strings.TrimSpace(values.Get("app"))
	}
	if len(strings.TrimSpace(values.Get("user"))) > 0 {
		msg.User = strings.TrimSpace(values.Get("user"))
	}
	if len(strings.TrimSpace(values.Get("url"))) > 0 {
		msg.URL = strings.TrimSpace(values.Get("url"))
	}
	if len(strings.TrimSpace(values.Get("head"))) > 0 {
		msg.Head = strings.TrimSpace(values.Get("head"))
	}
	if len(strings.TrimSpace(values.Get("head_long"))) > 0 {
		msg.HeadLong = strings.TrimSpace(values.Get("head_long"))
	}
	if len(strings.TrimSpace(values.Get("prev_head"))) > 0 {
		msg.PrevHead = strings.TrimSpace(values.Get("prev_head"))
	}
	if len(strings.TrimSpace(values.Get("git_log"))) > 0 {
		msg.GitLog = strings.TrimSpace(values.Get("git_log"))
	}
	if len(strings.TrimSpace(values.Get("release"))) > 0 {
		msg.Release = strings.TrimSpace(values.Get("release"))
	}
	return msg, nil
}
