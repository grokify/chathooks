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
        "confluence":{
            "event_slugs": ["page-created","comment-created"]
        },
        "heroku":{
            "file_extension": "txt",
            "event_slugs":["build"]
        },
        "librato":{
            "event_slugs":["alert-triggered","alert-cleared"]
        },
        "papertrail":{
            "event_slugs":["notifications-array-len-1","notifications-array"]
        },
        "pingdom":{
        	"event_slugs":["dns-check","http-check","http-custom-check","imap-check","ping-check","pop3-check","smtp-check","tcp-check","transaction-check","udp-check"]
        },
        "semaphore":{
        	"event_slugs":["build","deploy"]
        },
        "userlike":{
        	"event_slugs":["chat-meta_feedback","chat-meta_forward","chat-meta_rating","chat-meta_receive","chat-meta_start","chat-meta_survey","chat-widget_config","offline-message_receive","operator_away","operator_back","operator_offline","operator_online"]
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
