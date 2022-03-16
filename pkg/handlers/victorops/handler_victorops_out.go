package victorops

import (
	"encoding/json"

	cc "github.com/grokify/commonchat"

	"github.com/grokify/chathooks/pkg/config"
	"github.com/grokify/chathooks/pkg/handlers"
	"github.com/grokify/chathooks/pkg/models"
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

func Normalize(cfg config.Configuration, hReq handlers.HandlerRequest) (cc.Message, error) {
	ccMsg, err := CcMessageFromBytes(hReq.Body)
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
	var msg cc.Message
	return msg, json.Unmarshal(bytes, &msg)
}
