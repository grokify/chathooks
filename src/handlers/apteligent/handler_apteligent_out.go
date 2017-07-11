package apteligent

import (
	"encoding/json"
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"

	cc "github.com/commonchat/commonchat-go"
	"github.com/grokify/webhookproxy/src/adapters"
	"github.com/grokify/webhookproxy/src/config"
	"github.com/grokify/webhookproxy/src/util"
	"github.com/valyala/fasthttp"
)

const (
	DisplayName      = "Apteligent"
	HandlerKey       = "apteligent"
	MessageDirection = "out"
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

func (h Handler) HandlerKey() string {
	return HandlerKey
}

func (h Handler) MessageDirection() string {
	return MessageDirection
}

// HandleFastHTTP is the method to respond to a fasthttp request.
func (h *Handler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	ccMsg, err := Normalize(h.Config, ctx.PostBody())

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

func Normalize(cfg config.Configuration, bytes []byte) (cc.Message, error) {
	ccMsg := cc.NewMessage()
	iconURL, err := cfg.GetAppIconURL(HandlerKey)
	if err == nil {
		ccMsg.IconURL = iconURL.String()
	}

	src, err := ApteligentOutMessageFromBytes(bytes)
	if err != nil {
		return ccMsg, err
	}

	ccMsg.Activity = fmt.Sprintf("Alert %s", strings.ToLower(src.State))

	if len(src.Description) > 0 {
		if len(src.AlertURL) > 0 {
			ccMsg.Title = fmt.Sprintf("[%s](%s)", src.Description, src.AlertURL)
		} else {
			ccMsg.Title = src.Description
		}
	} else if len(src.AlertURL) > 0 {
		ccMsg.Title = fmt.Sprintf("[%s](%s)", src.AlertURL, src.AlertURL)
	}

	if 1 == 0 {
		attachment := cc.NewAttachment()
		if len(src.ApplicationName) > 0 {
			attachment.AddField(cc.Field{
				Title: "Application",
				Value: src.ApplicationName})
		}
		ccMsg.AddAttachment(attachment)
	}

	return ccMsg, nil
}

type ApteligentOutMessage struct {
	ThresholdValue   string `json:"threshold_value,omitempty"`
	TriggeringValue  string `json:"triggering_value,omitempty"`
	IncidentTime     string `json:"incident_time,omitempty"`
	Description      string `json:"description,omitempty"`
	Metric           string `json:"metric,omitempty"`
	CrittercismAppId string `json:"crittercism_app_id,omitempty"`
	TriggerId        string `json:"trigger_id,omitempty"`
	State            string `json:"state,omitempty"`
	AlertURL         string `json:"alert_url,omitempty"`
	Filters          string `json:"filters,omitempty"`
	ApplicationName  string `json:"application_name,omitempty"`
}

/*
{
    "threshold_value":"1",
    "description":"The Crashes alert on Crittercism was resolved at 06:40 PM UTC.",
    "metric":"Crashes",
    "crittercism_app_id":"54aab27451de5e9f042ec7ee",
    "trigger_id":"54aabecc1787845ae400000f",
    "state":"RESOLVED",
    "alert_url":"https://app.crittercism.com/developers/alerts/54aab27451de5e9f042ec7ee?alertId=54aabecc1787845ae400000f",
    "filters":"{}",
    "application_name":"Crittercism"
}
{
    "threshold_value":"1",
    "triggering_value":"4",
    "incident_time":"2015-01-05T18:15:56.976000",
    "description":"Alert on Crittercism at 06:15 PM UTC. Crashes threshold 4 exceeds 1.",
    "metric":"Crashes",
    "crittercism_app_id":"54aab27451de5e9f042ec7ee",
    "trigger_id":"54aabecc1787845ae400000f",
    "state":"TRIGGERED",
    "alert_url":"https://app.crittercism.com/developers/alerts/54aab27451de5e9f042ec7ee?incidentId=54aad4dcf39917103e0041b6",
    "filters":"{}",
    "application_name":"Crittercism"
}
*/
func ApteligentOutMessageFromBytes(bytes []byte) (ApteligentOutMessage, error) {
	msg := ApteligentOutMessage{}
	err := json.Unmarshal(bytes, &msg)
	return msg, err
}
