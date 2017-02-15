package appsignal

import (
	"encoding/json"
	"fmt"
	//"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/grokify/glip-go-webhook"
	"github.com/grokify/glip-webhook-proxy-go/src/adapters"
	"github.com/grokify/glip-webhook-proxy-go/src/config"
	"github.com/grokify/glip-webhook-proxy-go/src/util"
	"github.com/grokify/gotilla/time/timeutil"
	"github.com/valyala/fasthttp"
)

const (
	DISPLAY_NAME = "AppSignal"
	HANDLER_KEY  = "appsignal"
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

	message := util.NewMessage()

	if len(src.Marker.URL) > 0 {
		gmsg.Activity = fmt.Sprintf("%v Deployed by %v (%v)", src.Marker.Site, src.Marker.User, DISPLAY_NAME)

		if len(src.Marker.Environment) > 0 {
			message.AddAttachment(util.Attachment{
				Title: "Environment",
				Text:  src.Marker.Environment})
		}
		if len(src.Marker.URL) > 0 {
			message.AddAttachment(util.Attachment{
				Text: fmt.Sprintf("[View Details](%v)", src.Marker.URL)})
		}
	} else if len(src.Exception.URL) > 0 {
		if len(src.Exception.Site) > 0 {
			gmsg.Activity = fmt.Sprintf("%v Exception Incident (%v)", src.Exception.Site, DISPLAY_NAME)
		} else {
			gmsg.Activity = fmt.Sprintf("Exception Incident (%v)", DISPLAY_NAME)
		}

		if len(src.Exception.Message) > 0 {
			message.AddAttachment(util.Attachment{
				Text: src.Exception.Message})
		}
		if len(src.Exception.Environment) > 0 {
			message.AddAttachment(util.Attachment{
				Title: "Environment",
				Text:  src.Exception.Environment})
		}
		if len(src.Exception.Exception) > 0 {
			message.AddAttachment(util.Attachment{
				Title: "Exception",
				Text:  src.Exception.Exception})
		}
		if len(src.Exception.User) > 0 {
			message.AddAttachment(util.Attachment{
				Title: "User",
				Text:  src.Exception.User})
		}
		if len(src.Exception.URL) > 0 {
			message.AddAttachment(util.Attachment{
				Text: fmt.Sprintf("[View Details](%v)", src.Exception.URL)})
		}
	} else if len(src.Performance.URL) > 0 {
		gmsg.Activity = fmt.Sprintf("%v Performance Incident (%v)", src.Performance.Site, DISPLAY_NAME)

		if len(src.Performance.Environment) > 0 {
			message.AddAttachment(util.Attachment{
				Title: "Environment",
				Text:  src.Performance.Environment})
		}
		if len(src.Performance.Hostname) > 0 {
			message.AddAttachment(util.Attachment{
				Title: "Hostname",
				Text:  src.Performance.Hostname})
		}
		if src.Performance.Duration > 0.0 {
			durationString, err := timeutil.DurationStringMinutesSeconds(int64(src.Performance.Duration))
			if err == nil {
				message.AddAttachment(util.Attachment{
					Title: "Duration",
					Text:  durationString})
			} else {
				message.AddAttachment(util.Attachment{
					Title: "Duration",
					Text:  fmt.Sprintf("%v", src.Performance.Duration)})
			}
		}
		if len(src.Performance.URL) > 0 {
			message.AddAttachment(util.Attachment{
				Text: fmt.Sprintf("[View Details](%v)", src.Performance.URL)})
		}
	}

	gmsg.Body = glipadapter.RenderMessage(message)

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
