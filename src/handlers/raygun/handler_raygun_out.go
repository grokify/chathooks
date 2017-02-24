package raygun

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
	DisplayName = "Raygun"
	HandlerKey  = "raygun"
	IconURL     = "https://raygun.com/upload/raygun-icon.svg"
	ICON_URL2   = "https://raygun.com/images/logo/raygun-logo-og.jpg"
	ICON_URL3   = "https://a.slack-edge.com/ae7f/img/services/raygun_512.png"
)

// FastHttp request handler for Travis CI outbound webhook
type RaygunOutToGlipHandler struct {
	Config  config.Configuration
	Adapter adapters.Adapter
}

// FastHttp request handler constructor for Travis CI outbound webhook
func NewRaygunOutToGlipHandler(cfg config.Configuration, adapter adapters.Adapter) RaygunOutToGlipHandler {
	return RaygunOutToGlipHandler{Config: cfg, Adapter: adapter}
}

// HandleFastHTTP is the method to respond to a fasthttp request.
func (h *RaygunOutToGlipHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
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
	message.IconURL = ICON_URL3

	src, err := RaygunOutMessageFromBytes(bytes)
	if err != nil {
		return message, err
	}

	if src.EventType == "NewErrorOccurred" {
		if len(src.Application.Name) > 0 {
			message.Activity = fmt.Sprintf("%v encountered a new error", src.Application.Name)
		} else {
			message.Activity = "A new error has occurred"
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
			message.Activity = fmt.Sprintf("%v encountered an error%v", src.Application.Name)
		} else {
			message.Activity = fmt.Sprintf("An error occurred%v", timeString)
		}
	}

	attachment := cc.NewAttachment()

	if len(src.Application.URL) > 0 {
		if len(src.Application.Name) > 0 {
			attachment.AddField(cc.Field{
				Title: "Application",
				Value: fmt.Sprintf("[%v](%v)", src.Application.Name, src.Application.URL)})
		} else {
			attachment.AddField(cc.Field{
				Value: fmt.Sprintf("[Application Details](%v)", src.Application.URL)})
		}
	}
	if len(src.Error.URL) > 0 {
		if len(src.Error.Message) > 0 {
			attachment.AddField(cc.Field{
				Title: "Error",
				Value: fmt.Sprintf("[%v](%v)", src.Error.Message, src.Error.URL)})
		} else {
			attachment.AddField(cc.Field{
				Value: fmt.Sprintf("[Error Details](%v)", src.Error.URL)})
		}
	}

	message.AddAttachment(attachment)
	return message, nil
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
