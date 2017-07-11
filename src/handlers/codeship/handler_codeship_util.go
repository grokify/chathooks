package codeship

import (
	cc "github.com/commonchat/commonchat-go"
	"github.com/grokify/webhookproxy/src/config"
	"github.com/grokify/webhookproxy/src/util"

	"github.com/grokify/gotilla/fmt/fmtutil"
)

func ExampleMessage(cfg config.Configuration, data util.ExampleData) (cc.Message, error) {
	bytes, err := data.ExampleMessageBytes(HandlerKey, "build")
	if err != nil {
		return cc.Message{}, err
	}

	ccMsg, err := Normalize(cfg, bytes)
	if err == nil {
		fmtutil.PrintJSON(ccMsg)
		panic("A")
	}

	return Normalize(cfg, bytes)
}
