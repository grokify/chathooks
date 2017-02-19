package userlike

import (
	"fmt"
	"io/ioutil"
	"path"

	cc "github.com/grokify/commonchat"
	"github.com/grokify/glip-webhook-proxy-go/src/config"
)

var (
	ExampleFileOfflineMessageReceive = "example__offline_message__receive.json"
	ExampleFileChatWidgetConfig      = "example__chat_widget__config.json"
)

func ExampleMessageOfflineMessageReceive() (cc.Message, error) {
	bytes, err := ExampleMessageOfflineMessageReceiveBytes()
	if err != nil {
		return cc.Message{}, err
	}
	return Normalize(bytes)
}

func ExampleMessageOfflineMessageReceiveBytes() ([]byte, error) {
	filepath := path.Join(
		config.DOC_HANDLERS_REL_DIR,
		HandlerKey,
		ExampleFileOfflineMessageReceive)
	return ioutil.ReadFile(filepath)
}

func ExampleMessageChatWidgetConfig() (cc.Message, error) {
	bytes, err := ExampleMessageChatWidgetConfigBytes()
	if err != nil {
		return cc.Message{}, err
	}
	return Normalize(bytes)
}

func ExampleMessageChatWidgetConfigBytes() ([]byte, error) {
	filepath := path.Join(
		config.DOC_HANDLERS_REL_DIR,
		HandlerKey,
		ExampleFileChatWidgetConfig)
	return ioutil.ReadFile(filepath)
}

func ExampleMessageChatMeta(event string) (cc.Message, error) {
	filename := fmt.Sprintf("example__chat_meta__%s.json", event)
	filepath := path.Join(
		config.DOC_HANDLERS_REL_DIR,
		HandlerKey,
		filename)

	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return cc.Message{}, err
	}

	return Normalize(bytes)
}

func ExampleMessageOperator(event string) (cc.Message, error) {
	filename := fmt.Sprintf("example__operator__%s.json", event)
	filepath := path.Join(
		config.DOC_HANDLERS_REL_DIR,
		HandlerKey,
		filename)

	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return cc.Message{}, err
	}

	chatMsg, err := UserlikeOperatorOutMessageFromBytes(bytes)
	if err != nil {
		return cc.Message{}, err
	}
	return NormalizeOperator(chatMsg), nil
}
