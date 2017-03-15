package papertrail

import (
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"

	cc "github.com/commonchat/commonchat-go"
	"github.com/grokify/webhook-proxy-go/src/adapters"
	"github.com/grokify/webhook-proxy-go/src/config"
	"github.com/grokify/webhook-proxy-go/src/util"
	"github.com/valyala/fasthttp"
)

const (
	DisplayName      = "Papertrail"
	HandlerKey       = "papertrail"
	IconURL          = "https://d2rbro28ib85bu.cloudfront.net/images/integrations/128/papertrail.png"
	DocumentationURL = "http://help.papertrailapp.com/kb/how-it-works/web-hooks/"
)

// FastHttp request handler for Travis CI outbound webhook
type Handler struct {
	Config  config.Configuration
	Adapter adapters.Adapter
}

// FastHttp request handler constructor for Travis CI outbound webhook
func NewHandler(cfg config.Configuration, adapter adapters.Adapter) Handler {
	return Handler{Config: cfg, Adapter: adapter}
}

// HandleFastHTTP is the method to respond to a fasthttp request.
func (h *Handler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
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

	src, err := PapertrailOutMessageFromBytes(bytes)
	if err != nil {
		return message, err
	}

	if len(src.Events) > 1 {
		message.Activity = "Events triggered"
	} else {
		message.Activity = "Event triggered"
	}

	eventCount := len(src.Events)
	searchName := src.SavedSearch.Name
	if len(src.SavedSearch.HtmlSearchURL) > 0 {
		searchName = fmt.Sprintf("[%s](%s)", src.SavedSearch.Name, src.SavedSearch.HtmlSearchURL)
	}

	if eventCount == 1 {
		message.Title = fmt.Sprintf("%s event triggered!", searchName)
	} else {
		message.Title = fmt.Sprintf("%v %s events triggered!", eventCount, searchName)
	}

	for i, event := range src.Events {
		eventNumber := i + 1
		eventNumberDisplay := ""
		if eventCount > 1 {
			eventNumberDisplay = fmt.Sprintf(" %v", eventNumber)
		}
		attachment := cc.NewAttachment()

		if len(event.Message) > 0 {
			if len(event.Severity) > 0 {
				attachment.AddField(cc.Field{
					Title: fmt.Sprintf("Event%v", eventNumberDisplay),
					Value: fmt.Sprintf("[%s] %s", event.Severity, event.Message)})
			} else {
				attachment.AddField(cc.Field{
					Title: fmt.Sprintf("Event%v", eventNumberDisplay),
					Value: fmt.Sprintf("%s", event.Message)})
			}
		}
		if len(event.SourceName) > 0 {
			source := event.SourceName
			if len(event.SourceIP) > 0 {
				source = fmt.Sprintf("%s (%s)", event.SourceName, event.SourceIP)
			}
			attachment.AddField(cc.Field{
				Title: "Source",
				Value: source})
		}
		if len(event.Program) > 0 {
			attachment.AddField(cc.Field{
				Title: "Program",
				Value: event.Program})
		}
		if len(event.Facility) > 0 {
			attachment.AddField(cc.Field{
				Title: "Facility",
				Value: event.Facility})
		}
		if len(event.ReceivedAt) > 0 {
			attachment.AddField(cc.Field{
				Title: "Received At",
				Value: event.ReceivedAt})
		}

		message.AddAttachment(attachment)
	}

	return message, nil
}

type PapertrailOutMessage struct {
	Events      []PapertrailOutEvent     `json:"events,omitempty"`
	SavedSearch PapertrailOutSavedSearch `json:"saved_search,omitempty"`
	MaxId       int64                    `json:"max_id,omitempty"`
	MinId       int64                    `json:"min_id,omitempty"`
}

func PapertrailOutMessageFromBytes(bytes []byte) (PapertrailOutMessage, error) {
	msg := PapertrailOutMessage{}
	err := json.Unmarshal(bytes, &msg)
	return msg, err
}

type PapertrailOutEvent struct {
	Id                int64
	ReceivedAt        string
	DisplayReceivedAt string
	SourceIP          string
	SourceName        string
	SourceId          int64
	Hostname          string
	Program           string
	Severity          string
	Facility          string
	Message           string
}

type PapertrailOutSavedSearch struct {
	Id            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	Query         string `json:"query,omitempty"`
	HtmlEditURL   string `json:"html_edit_url,omitempty"`
	HtmlSearchURL string `json:"html_search_url,omitempty"`
}
