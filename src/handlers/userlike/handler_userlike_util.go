package userlike

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/grokify/glip-go-webhook"
	"github.com/grokify/glip-webhook-proxy-go/src/config"
)

var (
	EXAMPLE_FILE_OFFLINE_MESSAGE_RECEIVE = "example__offline_message__receive.json"
	EXAMPLE_FILE_CHAT_WIDGET_CONFIG      = "example__chat_widget__config.json"
)

func ExampleOfflineMessageReceiveMessageGlip() (glipwebhook.GlipWebhookMessage, error) {
	msg, err := ExampleOfflineMessageReceiveSource()
	if err != nil {
		return glipwebhook.GlipWebhookMessage{}, err
	}
	return NormalizeOfflineMessage(msg), nil
}

func ExampleOfflineMessageReceiveSource() (UserlikeOfflineMessageOutMessage, error) {
	bytes, err := ExampleOfflineMessageReceiveBytes()
	if err != nil {
		return UserlikeOfflineMessageOutMessage{}, err
	}
	return UserlikeOfflineMessageOutMessageFromBytes(bytes)
}

func ExampleOfflineMessageReceiveBytes() ([]byte, error) {
	filepath := path.Join(
		config.DOC_HANDLERS_REL_DIR,
		HANDLER_KEY,
		EXAMPLE_FILE_OFFLINE_MESSAGE_RECEIVE)
	return ioutil.ReadFile(filepath)
}

func ExampleChatWidgetConfigMessageGlip() (glipwebhook.GlipWebhookMessage, error) {
	msg, err := ExampleChatWidgetConfigSource()
	if err != nil {
		return glipwebhook.GlipWebhookMessage{}, err
	}
	return NormalizeChatWidget(msg), nil
}

func ExampleChatWidgetConfigSource() (UserlikeChatWidgetOutMessage, error) {
	bytes, err := ExampleChatWidgetConfigBytes()
	if err != nil {
		return UserlikeChatWidgetOutMessage{}, err
	}
	return UserlikeChatWidgetOutMessageFromBytes(bytes)
}

func ExampleChatWidgetConfigBytes() ([]byte, error) {
	filepath := path.Join(
		config.DOC_HANDLERS_REL_DIR,
		HANDLER_KEY,
		EXAMPLE_FILE_CHAT_WIDGET_CONFIG)
	return ioutil.ReadFile(filepath)
}

func ExampleUserlikeChatMetaStartOutMessageGlip() (glipwebhook.GlipWebhookMessage, error) {
	return ExampleUserlikeChatMetaMessageGlip("start")
}

func ExampleUserlikeChatMetaMessageGlip(event string) (glipwebhook.GlipWebhookMessage, error) {
	filename := fmt.Sprintf("example__chat_meta__%s.json", event)
	filepath := path.Join(
		config.DOC_HANDLERS_REL_DIR,
		HANDLER_KEY,
		filename)

	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return glipwebhook.GlipWebhookMessage{}, err
	}

	chatMsg, err := UserlikeChatMetaStartOutMessageFromBytes(bytes)
	if err != nil {
		return glipwebhook.GlipWebhookMessage{}, err
	}
	return NormalizeChatMeta(chatMsg), nil
}

func ExampleUserlikeOperatorMessageGlip(event string) (glipwebhook.GlipWebhookMessage, error) {
	filename := fmt.Sprintf("example__operator__%s.json", event)
	filepath := path.Join(
		config.DOC_HANDLERS_REL_DIR,
		HANDLER_KEY,
		filename)

	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return glipwebhook.GlipWebhookMessage{}, err
	}

	chatMsg, err := UserlikeOperatorOutMessageFromBytes(bytes)
	if err != nil {
		return glipwebhook.GlipWebhookMessage{}, err
	}
	return NormalizeOperator(chatMsg), nil
}
