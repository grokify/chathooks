package pingdom

import (
	"encoding/json"
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"

	cc "github.com/commonchat/commonchat-go"
	"github.com/grokify/webhook-proxy-go/src/adapters"
	"github.com/grokify/webhook-proxy-go/src/config"
	"github.com/grokify/webhook-proxy-go/src/util"
	"github.com/valyala/fasthttp"
)

const (
	DisplayName      = "Pingdom"
	HandlerKey       = "pingdom"
	IconURL          = "https://a.slack-edge.com/95b9/plugins/pingdom/assets/service_512.png"
	DocumentationURL = "https://www.pingdom.com/resources/webhooks"
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
		}).Info(fmt.Sprintf("%v request is not acceptable: %v", DisplayName, err))
		return
	}

	util.SendWebhook(ctx, h.Adapter, ccMsg)
}

func Normalize(bytes []byte) (cc.Message, error) {
	message := cc.NewMessage()
	message.IconURL = IconURL

	src, err := PingdomOutMessageFromBytes(bytes)
	if err != nil {
		return message, err
	}

	descMap := map[string]string{
		"HTTP_CUSTOM": "HTTP Custom",
		"PING":        "Ping",
		"PORT_TCP":    "TCP",
		"TRANSACTION": "Transaction"}

	activity := src.CheckType
	if display, ok := descMap[activity]; ok {
		message.Activity = fmt.Sprintf("%v check", display)
	} else {
		message.Activity = fmt.Sprintf("%v check", activity)
	}

	state := strings.ToLower(src.CurrentState)
	if state == "success" {
		state = "successful"
	}
	message.Title = fmt.Sprintf("[%v](%v) is %v", src.CheckName, src.CheckURL(), state)

	attachment := cc.NewAttachment()

	if len(strings.TrimSpace(src.Description)) > 0 {
		attachment.AddField(cc.Field{Title: "Description", Value: src.Description})
	}

	if len(strings.TrimSpace(src.CheckParams.FullURL)) > 0 {
		attachment.AddField(cc.Field{Title: "URL", Value: src.CheckParams.FullURL})
	} else if len(strings.TrimSpace(src.CheckParams.URL)) > 0 {
		attachment.AddField(cc.Field{Title: "URL", Value: src.CheckParams.URL})
	} else {
		if len(strings.TrimSpace(src.CheckParams.Hostname)) > 0 {
			attachment.AddField(cc.Field{Title: "Hostname", Value: src.CheckParams.Hostname})
		}
		if src.CheckParams.Port > 0 {
			attachment.AddField(cc.Field{Title: "Port", Value: fmt.Sprintf("%v", src.CheckParams.Port)})
		}
	}

	message.AddAttachment(attachment)
	return message, nil
}

type PingdomOutMessage struct {
	CheckId               int64                 `json:"check_id,omitempty"`
	CheckName             string                `json:"check_name,omitempty"`
	CheckType             string                `json:"check_type,omitempty"`
	CheckParams           PingdomOutCheckParams `json:"check_params,omitempty"`
	Tags                  []string              `json:"tags,omitempty"`
	PreviousState         string                `json:"previous_state,omitempty"`
	CurrentState          string                `json:"current_state,omitempty"`
	StateChangedTimestamp int64                 `json:"state_changed_timestamp,omitempty"`
	StateChangedUTCTime   string                `json:"state_changed_utc_time,omitempty"`
	LongDescription       string                `json:"long_description,omitempty"`
	Description           string                `json:"description,omitempty"`
	FirstProbe            PingdomOutProbe       `json:"first_probe,omitempty"`
	SecondProbe           PingdomOutProbe       `json:"second_probe,omitempty"`
}

func PingdomOutMessageFromBytes(bytes []byte) (PingdomOutMessage, error) {
	msg := PingdomOutMessage{}
	err := json.Unmarshal(bytes, &msg)
	return msg, err
}

func (msg *PingdomOutMessage) CheckURL() string {
	return fmt.Sprintf("https://my.pingdom.com/newchecks/checks#check=%v", msg.CheckId)
}

type PingdomOutCheckParams struct {
	BasicAuth  bool   `json:"basic_auth,omitempty"`
	Encryption bool   `json:"encryption,omitempty"`
	FullURL    string `json:"full_url,omitempty"`
	Header     string `json:"header,omitempty"`
	Hostname   string `json:"hostname,omitempty"`
	IPV6       bool   `json:"ipv6,omitempty"`
	Port       int    `json:"port,omitempty"`
	URL        string `json:"url,omitempty"`
}

type PingdomOutProbe struct {
	IP       string
	IPV6     string
	Location string
}
