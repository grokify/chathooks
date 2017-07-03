package enchant

import (
	cc "github.com/commonchat/commonchat-go"
	"github.com/grokify/webhookproxy/src/config"
	"github.com/grokify/webhookproxy/src/util"
)

func ExampleMessage(cfg config.Configuration, data util.ExampleData) (cc.Message, error) {
	bytes, err := data.ExampleMessageBytes(HandlerKey, "notification")
	if err != nil {
		return cc.Message{}, err
	}
	return Normalize(cfg, bytes)
}
