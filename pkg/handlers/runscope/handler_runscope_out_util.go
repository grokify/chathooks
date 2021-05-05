package runscope

import (
	cc "github.com/grokify/commonchat"

	"github.com/grokify/chathooks/pkg/config"
	"github.com/grokify/chathooks/pkg/handlers"
	"github.com/grokify/chathooks/pkg/util"
)

func ExampleMessage(cfg config.Configuration, data util.ExampleData) (cc.Message, error) {
	bytes, err := data.ExampleMessageBytes(HandlerKey, "notification")
	if err != nil {
		return cc.Message{}, err
	}
	return Normalize(cfg, handlers.HandlerRequest{Body: bytes})
}
