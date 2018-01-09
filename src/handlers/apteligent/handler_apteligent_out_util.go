package apteligent

import (
	cc "github.com/commonchat/commonchat-go"
	"github.com/grokify/chathooks/src/config"
	"github.com/grokify/chathooks/src/util"
)

func ExampleMessage(cfg config.Configuration, data util.ExampleData, eventSlug string) (cc.Message, error) {
	bytes, err := data.ExampleMessageBytes(HandlerKey, eventSlug)
	if err != nil {
		return cc.Message{}, err
	}
	return Normalize(cfg, bytes)
}
