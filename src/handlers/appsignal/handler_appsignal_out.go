package appsignal

import (
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"

	cc "github.com/commonchat/commonchat-go"
	"github.com/grokify/gotilla/time/timeutil"
	"github.com/grokify/webhook-proxy-go/src/adapters"
	"github.com/grokify/webhook-proxy-go/src/config"
	"github.com/grokify/webhook-proxy-go/src/util"
	"github.com/valyala/fasthttp"
)

const (
	DisplayName = "AppSignal"
	HandlerKey  = "appsignal"
	IconURL     = "https://pbs.twimg.com/profile_images/3558871752/5a8d304cb458baf99a7325a9c60b8a6b_400x400.png"
)

// FastHttp request handler for Semaphore CI outbound webhook
type AppsignalOutToGlipHandler struct {
	Config  config.Configuration
	Adapter adapters.Adapter
}

// FastHttp request handler constructor for Semaphore CI outbound webhook
func NewAppsignalOutToGlipHandler(cfg config.Configuration, adapter adapters.Adapter) AppsignalOutToGlipHandler {
	return AppsignalOutToGlipHandler{Config: cfg, Adapter: adapter}
}

// HandleFastHTTP is the method to respond to a fasthttp request.
func (h *AppsignalOutToGlipHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	ccMsg, err := Normalize(ctx.PostBody())

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

func Normalize(bytes []byte) (cc.Message, error) {
	message := cc.NewMessage()
	message.IconURL = IconURL

	src, err := AppsignalOutMessageFromBytes(bytes)
	if err != nil {
		return message, err
	}

	if len(src.Marker.URL) > 0 {
		message.Activity = fmt.Sprintf("%v Deployed by %v (%v)", src.Marker.Site, src.Marker.User, DisplayName)

		attachment := cc.NewAttachment()
		if len(src.Marker.Environment) > 0 {
			attachment.AddField(cc.Field{
				Title: "Environment",
				Value: src.Marker.Environment})
		}
		if len(src.Marker.URL) > 0 {
			attachment.AddField(cc.Field{
				Value: fmt.Sprintf("[View Details](%v)", src.Marker.URL)})
		}
		message.AddAttachment(attachment)
	} else if len(src.Exception.URL) > 0 {
		if len(src.Exception.Site) > 0 {
			message.Activity = fmt.Sprintf("%v Exception Incident (%v)", src.Exception.Site, DisplayName)
		} else {
			message.Activity = fmt.Sprintf("Exception Incident (%v)", DisplayName)
		}

		attachment := cc.NewAttachment()

		if len(src.Exception.Message) > 0 {
			attachment.AddField(cc.Field{
				Value: src.Exception.Message,
				Short: false})
		}
		if len(src.Exception.Environment) > 0 {
			attachment.AddField(cc.Field{
				Title: "Environment",
				Value: src.Exception.Environment,
				Short: true})
		}
		if len(src.Exception.Exception) > 0 {
			attachment.AddField(cc.Field{
				Title: "Exception",
				Value: src.Exception.Exception,
				Short: true})
		}
		if len(src.Exception.User) > 0 {
			attachment.AddField(cc.Field{
				Title: "User",
				Value: src.Exception.User,
				Short: true})
		}
		if len(src.Exception.URL) > 0 {
			attachment.AddField(cc.Field{
				Value: fmt.Sprintf("[View Details](%v)", src.Exception.URL),
				Short: true})
		}
		message.AddAttachment(attachment)
	} else if len(src.Performance.URL) > 0 {
		message.Activity = fmt.Sprintf("%v Performance Incident (%v)", src.Performance.Site, DisplayName)

		attachment := cc.NewAttachment()

		if len(src.Performance.Environment) > 0 {
			attachment.AddField(cc.Field{
				Title: "Environment",
				Value: src.Performance.Environment,
				Short: true})
		}
		if len(src.Performance.Hostname) > 0 {
			attachment.AddField(cc.Field{
				Title: "Hostname",
				Value: src.Performance.Hostname,
				Short: true})
		}
		if src.Performance.Duration > 0.0 {
			durationString, err := timeutil.DurationStringMinutesSeconds(int64(src.Performance.Duration))
			if err == nil {
				attachment.AddField(cc.Field{
					Title: "Duration",
					Value: durationString,
					Short: true})
			} else {
				attachment.AddField(cc.Field{
					Title: "Duration",
					Value: fmt.Sprintf("%v", src.Performance.Duration),
					Short: true})
			}
		}
		if len(src.Performance.URL) > 0 {
			attachment.AddField(cc.Field{
				Value: fmt.Sprintf("[View Details](%v)", src.Performance.URL)})
		}
		message.AddAttachment(attachment)
	}

	return message, nil
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
