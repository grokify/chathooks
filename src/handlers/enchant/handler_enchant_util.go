package enchant

import (
	cc "github.com/commonchat/commonchat-go"
	"github.com/grokify/webhook-proxy-go/src/util"
)

func ExampleMessage(data util.ExampleData) (cc.Message, error) {
	bytes, err := data.ExampleMessageBytes(HandlerKey, "notification")
	if err != nil {
		return cc.Message{}, err
	}
	return Normalize(bytes)
}

/*
import (
	cc "github.com/commonchat/commonchat-go"
)

func ExampleMessage() (cc.Message, error) {
	bytes, err := ExampleMessageBytes()
	if err != nil {
		return cc.Message{}, err
	}
	return Normalize(bytes)
}

func ExampleMessageBytes() ([]byte, error) {
	return []byte(`{
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
    "created_at": "2016-10-14T20:15:46Z"
  }
}`), nil
}

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
