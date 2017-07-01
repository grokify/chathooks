package raygun

import (
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"

	cc "github.com/commonchat/commonchat-go"
	"github.com/grokify/chatmore/src/adapters"
	"github.com/grokify/chatmore/src/config"
	"github.com/grokify/chatmore/src/util"
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
		}
		if len(src.Application.Name) > 0 {
			message.Activity = fmt.Sprintf("%v encountered an error%v", src.Application.Name, timeString)
		} else {
			message.Activity = fmt.Sprintf("An error occurred%v", timeString)
		}
	}

	attachment := cc.NewAttachment()

	followups := map[string]string{
		"NewErrorOccurred":     "New Error",
		"ErrorReoccurred":      "Error Reoccurred",
		"OneMinuteFollowUp":    "One Minute Follow Up",
		"FiveMinuteFollowUp":   "5 Minute Follow Up",
		"TenMinuteFollowUp":    "10 Minute Follow Up",
		"ThirtyMinuteFollowUp": "30 Minute Follow Up",
		"HourlyFollowUp":       "Hourly Follow Up"}
	if len(src.EventType) > 0 {
		if desc, ok := followups[src.EventType]; ok {
			attachment.AddField(cc.Field{
				Title: "Message Type",
				Value: desc})
		} else {
			attachment.AddField(cc.Field{
				Title: "Message Type",
				Value: src.EventType})
		}
	}

	if len(src.Application.URL) > 0 {
		if len(src.Application.Name) > 0 {
			attachment.AddField(cc.Field{
				Title: "Application",
				Value: fmt.Sprintf("[%v](%v)", src.Application.Name, src.Application.URL)})
		} else {
			attachment.AddField(cc.Field{
				Title: "Application",
				Value: fmt.Sprintf("[%v](%v)", src.Application.URL, src.Application.URL)})
		}
	}

	if len(src.Error.URL) > 0 {
		if len(src.Error.Message) > 0 {
			attachment.AddField(cc.Field{
				Title: "Error",
				Value: fmt.Sprintf("[%v](%v)", src.Error.Message, src.Error.URL)})
		} else {
			attachment.AddField(cc.Field{
				Title: "Error",
				Value: fmt.Sprintf("[%v](%v)", src.Error.URL, src.Error.URL)})
		}
		attachment.AddField(cc.Field{Title: "Users Affected", Value: fmt.Sprintf("%v", src.Error.UsersAffected), Short: true})
		attachment.AddField(cc.Field{Title: "Total Occurrences", Value: fmt.Sprintf("%v", src.Error.TotalOccurrences), Short: true})
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
