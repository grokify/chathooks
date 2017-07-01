package codeship

import (
	cc "github.com/commonchat/commonchat-go"
	"github.com/grokify/chatmore/src/util"
)

func ExampleMessage(data util.ExampleData) (cc.Message, error) {
	bytes, err := data.ExampleMessageBytes(HandlerKey, "build")
	if err != nil {
		return cc.Message{}, err
	}
	return Normalize(bytes)
}
