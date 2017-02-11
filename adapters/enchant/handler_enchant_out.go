package enchant

import (
	"encoding/json"
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/grokify/gotilla/strings/stringsutil"

	"github.com/grokify/glip-go-webhook"
	"github.com/grokify/glip-webhook-proxy-go/config"
	"github.com/grokify/glip-webhook-proxy-go/util"
	"github.com/valyala/fasthttp"
)

const (
	DISPLAY_NAME = "Enchant"
	ICON_URL     = "https://pbs.twimg.com/profile_images/530790354966962176/2trsSpWz_400x400.png"
)

// FastHttp request handler for Enchant outbound webhook
// https://dev.enchant.com/webhooks
type EnchantOutToGlipHandler struct {
	Config     config.Configuration
	GlipClient glipwebhook.GlipWebhookClient
}

// FastHttp request handler constructor for Confluence outbound webhook
func NewEnchantOutToGlipHandler(cfg config.Configuration, glip glipwebhook.GlipWebhookClient) EnchantOutToGlipHandler {
	return EnchantOutToGlipHandler{Config: cfg, GlipClient: glip}
}

// HandleFastHTTP is the method to respond to a fasthttp request.
func (h *EnchantOutToGlipHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	srcMsg, err := BuildInboundMessage(ctx)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotAcceptable)
		log.WithFields(log.Fields{
			"type":   "http.response",
			"status": fasthttp.StatusNotAcceptable,
		}).Info("Confluence request is not acceptable.")
		return
	}
	glipMsg := Normalize(srcMsg)

	util.SendGlipWebhookCtx(ctx, h.GlipClient, glipMsg)
}

func BuildInboundMessage(ctx *fasthttp.RequestCtx) (EnchantOutMessage, error) {
	return EnchantOutMessageFromBytes(ctx.FormValue("payload"))
}

func Normalize(src EnchantOutMessage) glipwebhook.GlipWebhookMessage {
	glip := glipwebhook.GlipWebhookMessage{Icon: ICON_URL}
	glip.Activity = fmt.Sprintf("%v (%v)", src.ActorName, DISPLAY_NAME)
	lines := []string{}
	lines = append(lines, fmt.Sprintf("> %v", src.Model.Subject))
	lines = append(lines, fmt.Sprintf("| **State** |\n| %v |", stringsutil.ToUpperFirst(src.Model.State)))
	glip.Body = strings.Join(lines, "\n")
	return glip
}

type EnchantOutMessage struct {
	Id         string       `json:"id,omitempty"`
	AccountId  string       `json:"account_id,omitempty"`
	AccountURL string       `json:"account_url,omitempty"`
	CreatedAt  string       `json:"created_at,omitempty"`
	Type       string       `json:"type,omitempty"`
	Data       EnchantData  `json:"data,omitempty"`
	ActorType  string       `json:"actor_type,omitempty"`
	ActorId    string       `json:"actor_id,omitempty"`
	ActorName  string       `json:"actor_name,omitempty"`
	ModelType  string       `json:"model_type,omitempty"`
	ModelId    string       `json:"model_id,omitempty"`
	Model      EnchantModel `json:"model,omitempty"`
}

type EnchantData struct {
	LabelId    string `json:"label_id,omitempty"`
	LabelName  string `json:"label_name,omitempty"`
	LabelColor string `json:"label_color,omitempty"`
}

type EnchantModel struct {
	Id         string   `json:"id,omitempty"`
	Number     int64    `json:"number,omitempty"`
	UserId     string   `json:"user_id,omitempty"`
	State      string   `json:"state,omitempty"`
	Subject    string   `json:"subject,omitempty"`
	LabelIds   []string `json:"label_ids,omitempty"`
	CustomerId string   `json:"customer_id,omitempty"`
	Type       string   `json:"type,omitempty"`
	ReplyTo    string   `json:"reply_to,omitempty"`
	CreatedAt  string   `json:"created_at,omitempty"`
}

func EnchantOutMessageFromBytes(bytes []byte) (EnchantOutMessage, error) {
	log.WithFields(log.Fields{
		"type":    "message.raw",
		"message": string(bytes),
	}).Debug(fmt.Sprintf("%v message.", DISPLAY_NAME))
	msg := EnchantOutMessage{}
	err := json.Unmarshal(bytes, &msg)
	if err != nil {
		log.WithFields(log.Fields{
			"type":  "message.json.unmarshal",
			"error": fmt.Sprintf("%v\n", err),
		}).Warn(fmt.Sprintf("%v request unmarshal failure.", DISPLAY_NAME))
	}
	return msg, err
}

/*

{
  "id": "7f94629",
  "account_id": "a91bb74",
  "account_url": "company.enchant.com",
  "created_at": "2016-10-17T19:52:43Z",
  "type": "ticket.label_added",
  "data": {
    "label_id": "97b0a40",
    "label_name": "High Priority",
    "label_color": "red"
  },
  "actor_type": "user",
  "actor_id": "a91bb75",
  "actor_name": "Michelle Han",
  "model_type": "ticket",
  "model_id": "a52ec86",
  "model": {
    "id": "a52ec86",
    "number": 53249,
    "user_id": "a91bb75",
    "state": "open",
    "subject": "email from customer",
    "label_ids": [
      "97b0a3e",
      "97b0a40"
    ],
    "customer_id": "97b0a43",
    "type": "email",
    "reply_to": "john@smith.com",
    "created_at": "2016-10-14T20:15:46Z",
    ... # truncated
  }
}

*/
