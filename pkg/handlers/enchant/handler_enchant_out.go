package enchant

import (
	"encoding/json"

	cc "github.com/grokify/commonchat"
	"github.com/grokify/mogo/type/stringsutil"
	"github.com/rs/zerolog/log"

	"github.com/grokify/chathooks/pkg/config"
	"github.com/grokify/chathooks/pkg/handlers"
	"github.com/grokify/chathooks/pkg/models"
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

func Normalize(cfg config.Configuration, hReq handlers.HandlerRequest) (cc.Message, error) {
	ccMsg := cc.NewMessage()
	iconURL, err := cfg.GetAppIconURL(HandlerKey)
	if err == nil {
		ccMsg.IconURL = iconURL.String()
	}

	src, err := EnchantOutMessageFromBytes(hReq.Body)
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
			Value: stringsutil.ToUpperFirst(src.Model.State, false)})
	}
	ccMsg.AddAttachment(attachment)
	return ccMsg, nil
}

type EnchantOutMessage struct {
	ID         string       `json:"id,omitempty"`
	AccountID  string       `json:"account_id,omitempty"`
	AccountURL string       `json:"account_url,omitempty"`
	CreatedAt  string       `json:"created_at,omitempty"`
	Type       string       `json:"type,omitempty"`
	Data       EnchantData  `json:"data,omitempty"`
	ActorType  string       `json:"actor_type,omitempty"`
	ActorID    string       `json:"actor_id,omitempty"`
	ActorName  string       `json:"actor_name,omitempty"`
	ModelType  string       `json:"model_type,omitempty"`
	ModelID    string       `json:"model_id,omitempty"`
	Model      EnchantModel `json:"model,omitempty"`
}

type EnchantData struct {
	LabelID    string `json:"label_id,omitempty"`
	LabelName  string `json:"label_name,omitempty"`
	LabelColor string `json:"label_color,omitempty"`
}

type EnchantModel struct {
	ID         string   `json:"id,omitempty"`
	Number     int64    `json:"number,omitempty"`
	UserID     string   `json:"user_id,omitempty"`
	State      string   `json:"state,omitempty"`
	Subject    string   `json:"subject,omitempty"`
	LabelIDs   []string `json:"label_ids,omitempty"`
	CustomerID string   `json:"customer_id,omitempty"`
	Type       string   `json:"type,omitempty"`
	ReplyTo    string   `json:"reply_to,omitempty"`
	CreatedAt  string   `json:"created_at,omitempty"`
}

func EnchantOutMessageFromBytes(bytes []byte) (EnchantOutMessage, error) {
	log.Debug().
		Str("type", "message.raw").
		Str("request_body", string(bytes)).
		Msg("handler_enchant_parse_message")

	var msg EnchantOutMessage
	err := json.Unmarshal(bytes, &msg)
	if err != nil {
		log.Warn().
			Err(err).
			Str("type", "message.json.unmarshal").
			Str("handler", DisplayName).
			Msg("FAIL - request_json_unmarshal_failure")
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
