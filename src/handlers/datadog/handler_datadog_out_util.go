package datadog

import (
	cc "github.com/commonchat/commonchat-go"
	"github.com/grokify/webhook-proxy-go/src/util"
)

func ExampleMessage(data util.ExampleData) (cc.Message, error) {
	bytes, err := data.ExampleMessageBytes(HandlerKey, "formatted1")
	if err != nil {
		return cc.Message{}, err
	}
	return Normalize(bytes)
}
