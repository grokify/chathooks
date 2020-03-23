package heroku

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/grokify/chathooks/src/config"
	"github.com/grokify/chathooks/src/handlers"
	"github.com/grokify/chathooks/src/models"
	cc "github.com/grokify/commonchat"
	"github.com/valyala/fasthttp"
)

const (
	DisplayName      = "Heroku"
	HandlerKey       = "heroku"
	MessageDirection = "out"
	MessageBodyType  = models.URLEncoded

	WebhookDocsUrl = "https://devcenter.heroku.com/articles/app-webhooks-tutorial"
)

func NewHandler() handlers.Handler {
	return handlers.Handler{MessageBodyType: MessageBodyType, Normalize: Normalize}
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
func Normalize(cfg config.Configuration, hReq handlers.HandlerRequest) (cc.Message, error) {
	src, err := HerokuOutMessageFromQuery(hReq.Body)
	if err != nil {
		return cc.Message{}, err
	}
	return NormalizeHerokuMessage(cfg, src)
}

func NormalizeHerokuMessage(cfg config.Configuration, src HerokuOutMessage) (cc.Message, error) {
	ccMsg := cc.NewMessage()
	iconURL, err := cfg.GetAppIconURL(HandlerKey)
	if err == nil {
		ccMsg.IconURL = iconURL.String()
	}

	if len(strings.TrimSpace(src.App)) > 0 {
		ccMsg.Activity = fmt.Sprintf("%v deployed on %v\n\n", src.App, DisplayName)
	} else {
		ccMsg.Activity = fmt.Sprintf("An app has been deployed on %v", DisplayName)
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

	ccMsg.AddAttachment(attachment)
	return ccMsg, nil
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

func HerokuOutMessageFromQuery(queryString []byte) (HerokuOutMessage, error) {
	msg := HerokuOutMessage{}
	values, err := url.ParseQuery(string(queryString))
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
