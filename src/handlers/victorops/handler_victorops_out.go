package victorops

import (
	"encoding/json"

	cc "github.com/commonchat/commonchat-go"
	"github.com/grokify/chathooks/src/config"
	"github.com/grokify/chathooks/src/handlers"
	"github.com/grokify/chathooks/src/models"
)

const (
	DisplayName      = "VictorOps"
	HandlerKey       = "victorops"
	MessageDirection = "out"
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

func CcMessageFromBytes(bytes []byte) (cc.Message, error) {
	msg := cc.Message{}
	err := json.Unmarshal(bytes, &msg)
	return msg, err
}
