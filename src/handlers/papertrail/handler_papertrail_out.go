package papertrail

import (
	"encoding/json"
	"fmt"
	"strings"

	cc "github.com/commonchat/commonchat-go"
	"github.com/grokify/webhookproxy/src/config"
	"github.com/grokify/webhookproxy/src/handlers"
	"github.com/grokify/webhookproxy/src/models"
)

const (
	DisplayName      = "Papertrail"
	HandlerKey       = "papertrail"
	MessageDirection = "out"
	DocumentationURL = "http://help.papertrailapp.com/kb/how-it-works/web-hooks/"
	MessageBodyType  = models.JSON
)

func NewHandler() handlers.Handler {
	return handlers.Handler{MessageBodyType: MessageBodyType, Normalize: Normalize}
}

func Normalize(cfg config.Configuration, bytes []byte) (cc.Message, error) {
	ccMsg := cc.NewMessage()
	iconURL, err := cfg.GetAppIconURL(HandlerKey)
	if err == nil {
		ccMsg.IconURL = iconURL.String()
	}

	src, err := PapertrailOutMessageFromBytes(bytes)
	if err != nil {
		return ccMsg, err
	}

	if len(src.Events) > 1 {
		ccMsg.Activity = "Events triggered"
	} else {
		ccMsg.Activity = "Event triggered"
	}

	eventCount := len(src.Events)
	searchName := src.SavedSearch.Name
	if len(src.SavedSearch.HtmlSearchURL) > 0 {
		searchName = fmt.Sprintf("[%s](%s)", src.SavedSearch.Name, src.SavedSearch.HtmlSearchURL)
	}

	if eventCount == 1 {
		ccMsg.Title = fmt.Sprintf("%s event triggered!", searchName)
	} else {
		ccMsg.Title = fmt.Sprintf("%v %s events triggered!", eventCount, searchName)
	}

	for i, event := range src.Events {
		eventNumber := i + 1
		eventNumberDisplay := ""
		if eventCount > 1 {
			eventNumberDisplay = fmt.Sprintf(" %v", eventNumber)
		}
		attachment := cc.NewAttachment()

		if len(event.Message) > 0 {
			hostString := ""
			hostParts := []string{}
			if len(event.Hostname) > 0 {
				hostParts = append(hostParts, event.Hostname)
			}
			if len(event.Facility) > 0 {
				hostParts = append(hostParts, event.Facility)
			}
			if len(hostParts) > 0 {
				hostPartsString := strings.Join(hostParts, "/")
				hostString = fmt.Sprintf(" (%v)", hostPartsString)
			}
			if len(event.Severity) > 0 {
				attachment.AddField(cc.Field{
					Title: fmt.Sprintf("Event%v", eventNumberDisplay),
					Value: fmt.Sprintf("[%s] %s%s", event.Severity, event.Message, hostString)})
			} else {
				attachment.AddField(cc.Field{
					Title: fmt.Sprintf("Event%v", eventNumberDisplay),
					Value: fmt.Sprintf("%s%s", event.Message, hostString)})
			}
		}

		ccMsg.AddAttachment(attachment)
		if 1 == 0 {
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

			ccMsg.AddAttachment(attachment)
		}
	}

	return ccMsg, nil
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
