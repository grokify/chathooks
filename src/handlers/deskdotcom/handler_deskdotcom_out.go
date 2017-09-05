package deskdotcom

import (
	"encoding/json"

	cc "github.com/commonchat/commonchat-go"
	"github.com/grokify/webhookproxy/src/config"
	"github.com/grokify/webhookproxy/src/handlers"
	"github.com/grokify/webhookproxy/src/models"
)

const (
	DisplayName      = "Desk.com"
	HandlerKey       = "deskdotcom"
	MessageDirection = "out"
	DocumentationURL = "https://support.desk.com/customer/portal/articles/869334-configuring-webhooks-in-desk-com-apps"
	MessageBodyType  = models.JSON
)

func NewHandler() handlers.Handler {
	return handlers.Handler{MessageBodyType: MessageBodyType, Normalize: Normalize}
}

func Normalize(cfg config.Configuration, bytes []byte) (cc.Message, error) {
	ccMsg, err := CcMessageFromBytes(bytes)
	if err != nil {
		return ccMsg, err
	}
	iconURL, err := cfg.GetAppIconURL(HandlerKey)
	if err == nil {
		ccMsg.IconURL = iconURL.String()
	}
	return ccMsg, nil
}

/*

Example template:

{
  "activity":"Case {{case.id}} updated",
  "body":"**Case Type**\n{{case.status}} {{case.type}}\n **Subject**\n[{{case.subject}}]({{case.direct_url}})"
}

*/

func CcMessageFromBytes(bytes []byte) (cc.Message, error) {
	msg := cc.Message{}
	err := json.Unmarshal(bytes, &msg)
	return msg, err
}
