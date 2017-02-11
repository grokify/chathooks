package heroku

import (
	"fmt"
	"net/url"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/grokify/glip-go-webhook"
	"github.com/grokify/glip-webhook-proxy-go/config"
	"github.com/grokify/glip-webhook-proxy-go/util"
	"github.com/valyala/fasthttp"
)

const (
	DISPLAY_NAME = "Heroku"
	ICON_URL     = "http://atlassian.wpengine.netdna-cdn.com/wp-content/uploads/heroku.jpeg"
)

// FastHttp request handler for Heroku outbound webhook
// https://devcenter.heroku.com/articles/deploy-hooks#http-post-hook
type HerokuOutToGlipHandler struct {
	Config     config.Configuration
	GlipClient glipwebhook.GlipWebhookClient
}

// FastHttp request handler constructor for Confluence outbound webhook
func NewHerokuOutToGlipHandler(cfg config.Configuration, glip glipwebhook.GlipWebhookClient) HerokuOutToGlipHandler {
	return HerokuOutToGlipHandler{Config: cfg, GlipClient: glip}
}

// HandleFastHTTP is the method to respond to a fasthttp request.
func (h *HerokuOutToGlipHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	srcMsg, err := BuildInboundMessage(ctx)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotAcceptable)
		log.WithFields(log.Fields{
			"type":   "http.response",
			"status": fasthttp.StatusNotAcceptable,
		}).Info(fmt.Sprintf("%v request is not acceptable.", DISPLAY_NAME))
		return
	}
	glipMsg := Normalize(srcMsg)

	util.SendGlipWebhookCtx(ctx, h.GlipClient, glipMsg)
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

func Normalize(src HerokuOutMessage) glipwebhook.GlipWebhookMessage {
	glip := glipwebhook.GlipWebhookMessage{Icon: ICON_URL}
	if len(strings.TrimSpace(src.User)) > 0 {
		if len(strings.TrimSpace(src.App)) > 0 {
			glip.Activity = fmt.Sprintf("%v deployed %v (%v)", src.User, src.App, DISPLAY_NAME)
		} else {
			glip.Activity = fmt.Sprintf("%v deployed an app (%v)", src.User, DISPLAY_NAME)
		}
	} else {
		if len(strings.TrimSpace(src.App)) > 0 {
			glip.Activity = fmt.Sprintf("%v deployed", src.App)
		} else {
			glip.Activity = "An app has been deployed"
		}
	}
	lines := []string{}
	if len(strings.TrimSpace(src.Release)) > 0 {
		lines = append(lines, fmt.Sprintf("> Release: %v", src.Release))
	}
	if len(strings.TrimSpace(src.URL)) > 0 {
		lines = append(lines, fmt.Sprintf("> App: [%v](%v)", src.URL, src.URL))
	}
	if len(lines) > 0 {
		glip.Body = strings.Join(lines, "\n")
	}
	return glip
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

func HerokuOutMessageFromQueryString(query string) (HerokuOutMessage, error) {
	msg := HerokuOutMessage{}
	values, err := url.ParseQuery(query)
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

/*

https://devcenter.heroku.com/articles/deploy-hooks#http-post-hook

app=secure-woodland-9775&user=example%40example.com&url=http%3A%2F%2Fsecure-woodland-9775.herokuapp.com&head=4f20bdd&head_long=4f20bdd&prev_head=&git_log=%20%20*%20Michael%20Friis%3A%20add%20bar&release=v7

*/
