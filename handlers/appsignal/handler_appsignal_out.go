package appsignal

import (
	"encoding/json"
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/grokify/glip-go-webhook"
	"github.com/grokify/glip-webhook-proxy-go/config"
	"github.com/grokify/glip-webhook-proxy-go/util"
	"github.com/grokify/gotilla/time/timeutil"
	"github.com/valyala/fasthttp"
)

const (
	DISPLAY_NAME = "AppSignal"
	ICON_URL     = "https://pbs.twimg.com/profile_images/3558871752/5a8d304cb458baf99a7325a9c60b8a6b_400x400.png"
)

// FastHttp request handler for Semaphore CI outbound webhook
type AppsignalOutToGlipHandler struct {
	Config     config.Configuration
	GlipClient glipwebhook.GlipWebhookClient
}

// FastHttp request handler constructor for Semaphore CI outbound webhook
func NewAppsignalOutToGlipHandler(cfg config.Configuration, glip glipwebhook.GlipWebhookClient) AppsignalOutToGlipHandler {
	return AppsignalOutToGlipHandler{Config: cfg, GlipClient: glip}
}

// HandleFastHTTP is the method to respond to a fasthttp request.
func (h *AppsignalOutToGlipHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
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

func BuildInboundMessage(ctx *fasthttp.RequestCtx) (AppsignalOutMessage, error) {
	return AppsignalOutMessageFromBytes(ctx.PostBody())
}

func Normalize(src AppsignalOutMessage) glipwebhook.GlipWebhookMessage {
	gmsg := glipwebhook.GlipWebhookMessage{Icon: ICON_URL}

	lines := []string{}

	if len(src.Marker.URL) > 0 {
		gmsg.Activity = fmt.Sprintf("%v Deployed by %v (%v)", src.Marker.Site, src.Marker.User, DISPLAY_NAME)
		if len(src.Marker.Environment) > 0 {
			lines = append(lines, fmt.Sprintf("> **Environment**\n> %v", src.Marker.Environment))
		}
		if len(src.Marker.Site) > 0 {
			lines = append(lines, fmt.Sprintf("> [View Details](%v)", src.Marker.URL))
		}
	} else if len(src.Exception.URL) > 0 {
		if len(src.Exception.Site) > 0 {
			gmsg.Activity = fmt.Sprintf("%v Exception Incident (%v)", src.Exception.Site, DISPLAY_NAME)
		} else {
			gmsg.Activity = fmt.Sprintf("Exception Incident (%v)", DISPLAY_NAME)
		}
		if len(src.Exception.Message) > 0 {
			lines = append(lines, fmt.Sprintf("> %v", src.Exception.Message))
		}
		if len(src.Exception.Environment) > 0 {
			lines = append(lines, fmt.Sprintf("> **Environment**\n> %v", src.Exception.Environment))
		}
		if len(src.Exception.Exception) > 0 {
			lines = append(lines, fmt.Sprintf("> **Exception**\n> %v", src.Exception.Exception))
		}
		if len(src.Exception.User) > 0 {
			lines = append(lines, fmt.Sprintf("> **User**\n> %v", src.Exception.User))
		}
		if len(src.Exception.URL) > 0 {
			lines = append(lines, fmt.Sprintf("> [View Details](%v)", src.Exception.URL))
		}
	} else if len(src.Performance.URL) > 0 {
		gmsg.Activity = fmt.Sprintf("%v Performance Incident (%v)", src.Performance.Site, DISPLAY_NAME)
		if len(src.Performance.Environment) > 0 {
			lines = append(lines, fmt.Sprintf("> **Environment**\n> %v", src.Performance.Environment))
		}
		if len(src.Performance.Hostname) > 0 {
			lines = append(lines, fmt.Sprintf("> **Hostname**\n> %v", src.Performance.Hostname))
		}
		if src.Performance.Duration > 0.0 {
			durationString, err := timeutil.DurationStringMinutesSeconds(int64(src.Performance.Duration))
			if err == nil {
				lines = append(lines, fmt.Sprintf("> **Duration**\n> %v", durationString))
			} else {
				lines = append(lines, fmt.Sprintf("> **Duration**\n> %v", err))
			}
		}
		if len(src.Performance.URL) > 0 {
			lines = append(lines, fmt.Sprintf("> [View Details](%v)", src.Performance.URL))
		}
	}

	if len(lines) > 0 {
		gmsg.Body = strings.Join(lines, "\n")
	}

	return gmsg
}

type AppsignalOutMessage struct {
	Marker      AppsignalMarker      `json:"marker,omitempty"`
	Exception   AppsignalException   `json:"exception,omitempty"`
	Performance AppsignalPerformance `json:"performance,omitempty"`
}

func AppsignalOutMessageFromBytes(bytes []byte) (AppsignalOutMessage, error) {
	msg := AppsignalOutMessage{}
	err := json.Unmarshal(bytes, &msg)
	return msg, err
}

type AppsignalMarker struct {
	User        string `json:"user,omitempty"`
	Site        string `json:"site,omitempty"`
	Environment string `json:"environment,omitempty"`
	Revision    string `json:"revision,omitempty"`
	Repository  string `json:"repository,omitempty"`
	URL         string `json:"url,omitempty"`
}

type AppsignalException struct {
	Exception          string `json:"exception,omitempty"`
	Site               string `json:"site,omitempty"`
	Message            string `json:"message,omitempty"`
	Action             string `json:"action,omitempty"`
	Path               string `json:"path,omitempty"`
	Revision           string `json:"revision,omitempty"`
	User               string `json:"user,omitempty"`
	FirstBacktraceLine string `json:"first_backtrace_line,omitempty"`
	URL                string `json:"url,omitempty"`
	Environment        string `json:"environment,omitempty"`
}

type AppsignalPerformance struct {
	Site        string  `json:"site,omitempty"`
	Action      string  `json:"action,omitempty"`
	Path        string  `json:"path,omitempty"`
	Duration    float64 `json:"duration,omitempty"`
	Status      int64   `json:"status,omitempty"`
	Hostname    string  `json:"hostname,omitempty"`
	Revision    string  `json:"revision,omitempty"`
	User        string  `json:"user,omitempty"`
	URL         string  `json:"url,omitempty"`
	Environment string  `json:"environment,omitempty"`
}

/*

{
  "marker":{
    "user": "thijs",
    "site": "AppSignal",
    "environment": "test",
    "revision": "3107ddc4bb053d570083b4e3e425b8d62532ddc9",
    "repository": "git@github.com:appsignal/appsignal.git",
    "url": "https://appsignal.com/test/sites/1385f7e38c5ce90000000000/web/exceptions"
  }
}

{
  "exception":{
    "exception": "ActionView::Template::Error",
    "site": "AppSignal",
    "message": "undefined method `encoding' for nil:NilClass",
    "action": "App::ErrorController#show",
    "path": "/errors",
    "revision": "3107ddc4bb053d570083b4e3e425b8d62532ddc9",
    "user": "thijs",
    "first_backtrace_line": "/usr/local/rbenv/versions/2.0.0-p353/lib/ruby/2.0.0/cgi/util.rb:7:in `escape'",
    "url": "https://appsignal.com/test/sites/1385f7e38c5ce90000000000/web/exceptions/App::SnapshotsController-show/ActionView::Template::Error",
    "environment": "test"
  }
}

{
  "performance":{
    "site": "AppSignal",
    "action": "App::ExceptionsController#index",
    "path": "/slow",
    "duration": 552.7897429999999,
    "status": 200,
    "hostname": "frontend.appsignal.com",
    "revision": "3107ddc4bb053d570083b4e3e425b8d62532ddc9",
    "user": "thijs",
    "url": "https://appsignal.com/test/sites/1385f7e38c5ce90000000000/web/performance/App::ExceptionsController-index",
    "environment": "test"
  }
}

*/
