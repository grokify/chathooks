package raygun

import (
	"encoding/json"
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/grokify/glip-go-webhook"
	"github.com/grokify/glip-webhook-proxy-go/src/config"
	"github.com/grokify/glip-webhook-proxy-go/src/util"
	"github.com/valyala/fasthttp"
)

const (
	DISPLAY_NAME = "Raygun"
	ICON_URL     = "https://raygun.com/images/logo/raygun-logo-og.jpg"
)

// FastHttp request handler for Travis CI outbound webhook
type RaygunOutToGlipHandler struct {
	Config     config.Configuration
	GlipClient glipwebhook.GlipWebhookClient
}

// FastHttp request handler constructor for Travis CI outbound webhook
func NewRaygunOutToGlipHandler(cfg config.Configuration, glip glipwebhook.GlipWebhookClient) RaygunOutToGlipHandler {
	return RaygunOutToGlipHandler{Config: cfg, GlipClient: glip}
}

// HandleFastHTTP is the method to respond to a fasthttp request.
func (h *RaygunOutToGlipHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
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

func BuildInboundMessage(ctx *fasthttp.RequestCtx) (RaygunOutMessage, error) {
	return RaygunOutMessageFromBytes(ctx.PostBody())
}

func Normalize(src RaygunOutMessage) glipwebhook.GlipWebhookMessage {
	glipMsg := glipwebhook.GlipWebhookMessage{
		Body:     strings.Join([]string{">", src.AsMarkdown()}, " "),
		Activity: DISPLAY_NAME,
		Icon:     ICON_URL}

	if src.EventType == "NewErrorOccurred" {
		if len(src.Application.Name) > 0 {
			glipMsg.Activity = fmt.Sprintf("%v encountered a new error (%v)", src.Application.Name, DISPLAY_NAME)
		} else {
			glipMsg.Activity = fmt.Sprintf("A new error has occurred (%v)", DISPLAY_NAME)
		}
	} else {
		timeString := ""
		if src.EventType == "ErrorReoccurred" {
			timeString = " again"
		} else if src.EventType == "OneMinuteFollowUp" {
			timeString = " 1 minute ago"
		} else if src.Event == "FiveMinuteFollowUp" {
			timeString = " 5 minutes ago"
		} else if src.Event == "TenMinuteFollowUp" {
			timeString = " 10 minutes ago"
		} else if src.Event == "ThirtyMinuteFollowUp" {
			timeString = " 30 minutes ago"
		} else if src.Event == "HourlyFollowUp" {
			timeString = " 1 hour ago"
		}
		if len(src.Application.Name) > 0 {
			glipMsg.Activity = fmt.Sprintf("%v encountered an error%v (%v)", src.Application.Name, DISPLAY_NAME)
		} else {
			glipMsg.Activity = fmt.Sprintf("An error occurred%v (%v)", timeString, DISPLAY_NAME)
		}
	}
	lines := []string{}
	if len(src.Application.URL) > 0 {
		if len(src.Application.Name) > 0 {
			lines = append(lines, fmt.Sprintf("> [Application details: %v](%v)", src.Application.Name, src.Application.URL))
		} else {
			lines = append(lines, fmt.Sprintf("> [Application details](%v)", src.Application.URL))
		}
	}
	if len(src.Error.URL) > 0 {
		if len(src.Error.Message) > 0 {
			lines = append(lines, fmt.Sprintf("> [Error details: %v](%v)", src.Error.Message, src.Error.URL))
		} else {
			lines = append(lines, fmt.Sprintf("> [Error details](%v)", src.Error.URL))
		}
	}
	if len(lines) > 0 {
		glipMsg.Body = strings.Join(lines, "\n")
	}
	return glipMsg
}

type RaygunOutMessage struct {
	Event       string            `json:"event,omitempty"`
	EventType   string            `json:"eventType,omitempty"`
	Error       RaygunError       `json:"error,omitempty"`
	Application RaygunApplication `json:"application,omitempty"`
}

func RaygunOutMessageFromBytes(bytes []byte) (RaygunOutMessage, error) {
	msg := RaygunOutMessage{}
	err := json.Unmarshal(bytes, &msg)
	return msg, err
}

type RaygunError struct {
	URL              string `json:"url,omitempty"`
	Message          string `json:"message,omitempty"`
	FirstOccurredOn  string `json:"firstOccurredOn,omitempty"`
	LastOccurredOn   string `json:"lastOccurredOn,omitempty"`
	UsersAffected    int    `json:"usersAffected,omitempty"`
	TotalOccurrences int    `json:"totalOccurrences,omitempty"`
}

type RaygunApplication struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

func (msg *RaygunOutMessage) AsMarkdown() string {
	if msg.EventType == "NewErrorOccurred" {
		return fmt.Sprintf("[]() encoutered a new error")
	}
	return ""
}

/*

A new error has been reported for test, the error was []()

{
  "event":"error_notification",
  "eventType":"NewErrorOccurred",
  "error": {
    "url":"http://app.raygun.io/error-url",
    "message":"",
    "firstOccurredOn":"1970-01-28T01:49:36Z",
    "lastOccurredOn":"1970-01-28T01:49:36Z",
    "usersAffected":1,
    "totalOccurrences":1
  },
  "application": {
    "name":"application name",
    "url":"http://app.raygun.io/application-url"
  }
}
*/
