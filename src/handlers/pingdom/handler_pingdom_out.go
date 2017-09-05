package pingdom

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
	DisplayName      = "Pingdom"
	HandlerKey       = "pingdom"
	MessageDirection = "out"
	DocumentationURL = "https://www.pingdom.com/resources/webhooks"
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

	src, err := PingdomOutMessageFromBytes(bytes)
	if err != nil {
		return ccMsg, err
	}

	descMap := map[string]string{
		"HTTP_CUSTOM": "HTTP Custom",
		"PING":        "Ping",
		"PORT_TCP":    "TCP",
		"TRANSACTION": "Transaction"}

	activity := src.CheckType
	if display, ok := descMap[activity]; ok {
		ccMsg.Activity = fmt.Sprintf("%v check", display)
	} else {
		ccMsg.Activity = fmt.Sprintf("%v check", activity)
	}

	state := strings.ToLower(src.CurrentState)
	if state == "success" {
		state = "successful"
	}
	ccMsg.Title = fmt.Sprintf("[%v](%v) is %v", src.CheckName, src.CheckURL(), state)

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

	ccMsg.AddAttachment(attachment)
	return ccMsg, nil
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
