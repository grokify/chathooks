package enchant

import (
	"encoding/json"
	"fmt"

	"github.com/grokify/chathooks/src/config"
	"github.com/grokify/chathooks/src/handlers"
	"github.com/grokify/chathooks/src/models"
	cc "github.com/grokify/commonchat"
	"github.com/grokify/gotilla/strings/stringsutil"
	log "github.com/sirupsen/logrus"
)

const (
	DisplayName      = "Enchant"
	HandlerKey       = "enchant"
	MessageDirection = "out"
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

	src, err := EnchantOutMessageFromBytes(bytes)
	if err != nil {
		return ccMsg, err
	}

	if len(src.ActorName) > 0 {
		ccMsg.Activity = src.ActorName
	}

	attachment := cc.NewAttachment()

	if len(src.Model.Subject) > 0 {
		attachment.Text = src.Model.Subject
	}
	if len(src.Model.State) > 0 {
		attachment.AddField(cc.Field{
			Title: "State",
			Value: stringsutil.ToUpperFirst(src.Model.State)})
	}
	ccMsg.AddAttachment(attachment)
	return ccMsg, nil
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
	}).Debug(fmt.Sprintf("%v message.", DisplayName))
	msg := EnchantOutMessage{}
	err := json.Unmarshal(bytes, &msg)
	if err != nil {
		log.WithFields(log.Fields{
			"type":  "message.json.unmarshal",
			"error": fmt.Sprintf("%v\n", err),
		}).Warn(fmt.Sprintf("%v request unmarshal failure.", DisplayName))
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
