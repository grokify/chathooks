package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/grokify/webhook-proxy-go/src/config"
)

const (
	DefaultExtension = "json"
)

type ExampleData struct {
	Data map[string]ExampleSource `json:"data,omitempty"`
}

type ExampleSource struct {
	FileExtension string   `json:"file_extension,omitempty"`
	EventSlugs    []string `json:"event_slugs,omitempty"`
}

func NewExampleData() (ExampleData, error) {
	data := ExampleData{}
	err := json.Unmarshal(ExampleDataRaw(), &data)
	return data, err
}

func ExampleDataRaw() []byte {
	return []byte(`{
    "data": {
        "appsignal": {
            "file_extension": "json",
            "event_slugs": ["marker","exception","performance"]
        },
        "apteligent":{
            "event_slugs": ["alert","alert-open","alert-close"]
        },
        "confluence":{
            "event_slugs": ["page-created","comment-created"]
        },
        "gosquared":{
            "event_slugs": ["site-traffic","smart-group","live-chat"]
        },
        "heroku":{
            "file_extension": "txt",
            "event_slugs":["build"]
        },
        "librato":{
            "event_slugs":["2","alert-triggered","alert-cleared"]
        },
        "opsgenie":{
            "event_slugs_":["create","close","delete",
            "acknowledge","unacknowledge"],

            "event_slugs__": [
            "add-note","add-recipient","add-tags","add-team"],

            "event_slugs":["remove-tags","assign-ownership","take-ownership", "escalate",
            "custom-action-test-action"]
        },
        "papertrail":{
            "event_slugs":["notifications-array-len-1","notifications-array"]
        },
        "pingdom":{
            "event_slugs":["http-check"],
        	"event_slugs_":["dns-check","http-check","http-custom-check","imap-check","ping-check","pop3-check","smtp-check","tcp-check","transaction-check","udp-check"]
        },
        "semaphore":{
            "event_slugs":["build","deploy"]
        },
        "statuspage":{
            "event_slugs":["incident-updates","incident-updates-create","component-updates"]
        },
        "userlike":{
        	"event_slugs_":["chat-meta_feedback","chat-meta_forward","chat-meta_rating","chat-meta_receive","chat-meta_start","chat-meta_survey"],
            "event_slugs":["chat-widget_config","offline-message_receive","operator_away","operator_back","operator_offline","operator_online"]
        }
    }
}`)
}

func (data *ExampleData) ExampleMessageBytes(handlerKey string, eventSlug string) ([]byte, error) {
	filepath := path.Join(
		config.DocHandlersDir,
		handlerKey,
		data.BuildFilename(handlerKey, eventSlug))
	return ioutil.ReadFile(filepath)
}

func (data *ExampleData) BuildFilename(handlerKey string, eventSlug string) string {
	ext := DefaultExtension
	if src, ok := data.Data[handlerKey]; ok {
		if len(src.FileExtension) > 0 {
			ext = src.FileExtension
		}
	}
	return fmt.Sprintf("event-example_%s.%s", eventSlug, ext)
}
