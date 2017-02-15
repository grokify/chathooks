package appsignal

import (
	"io/ioutil"
	"path"

	"github.com/grokify/glip-go-webhook"
	"github.com/grokify/glip-webhook-proxy-go/src/config"
)

const (
	EXAMPLE_MESSAGE_EXCEPTION   = "example__exception.json"
	EXAMPLE_MESSAGE_MARKER      = "example__marker.json"
	EXAMPLE_MESSAGE_PERFORMANCE = "example__performance.json"
)

func ExampleMarkerMessageGlip() (glipwebhook.GlipWebhookMessage, error) {
	msg, err := ExampleMarkerMessageSource()
	if err != nil {
		return glipwebhook.GlipWebhookMessage{}, err
	}
	return Normalize(msg), nil
}

func ExampleMarkerMessageSource() (AppsignalOutMessage, error) {
	bytes, err := ExampleMarkerMessageBytes()
	if err != nil {
		return AppsignalOutMessage{}, err
	}
	return AppsignalOutMessageFromBytes(bytes)
}

func ExampleMarkerMessageBytes() ([]byte, error) {
	filepath := path.Join(
		config.DOC_HANDLERS_REL_DIR,
		HANDLER_KEY,
		EXAMPLE_MESSAGE_MARKER)
	return ioutil.ReadFile(filepath)
}

func ExampleExceptionMessageGlip() (glipwebhook.GlipWebhookMessage, error) {
	msg, err := ExampleExceptionMessageSource()
	if err != nil {
		return glipwebhook.GlipWebhookMessage{}, err
	}
	return Normalize(msg), nil
}

func ExampleExceptionMessageSource() (AppsignalOutMessage, error) {
	bytes, err := ExampleExceptionMessageBytes()
	if err != nil {
		return AppsignalOutMessage{}, err
	}
	return AppsignalOutMessageFromBytes(bytes)
}

func ExampleExceptionMessageBytes() ([]byte, error) {
	filepath := path.Join(
		config.DOC_HANDLERS_REL_DIR,
		HANDLER_KEY,
		EXAMPLE_MESSAGE_EXCEPTION)
	return ioutil.ReadFile(filepath)
}

func ExamplePerformanceMessageGlip() (glipwebhook.GlipWebhookMessage, error) {
	msg, err := ExamplePerformanceMessageSource()
	if err != nil {
		return glipwebhook.GlipWebhookMessage{}, err
	}
	return Normalize(msg), nil
}

func ExamplePerformanceMessageSource() (AppsignalOutMessage, error) {
	bytes, err := ExamplePerformanceMessageBytes()
	if err != nil {
		return AppsignalOutMessage{}, err
	}
	return AppsignalOutMessageFromBytes(bytes)
}

func ExamplePerformanceMessageBytes() ([]byte, error) {
	filepath := path.Join(
		config.DOC_HANDLERS_REL_DIR,
		HANDLER_KEY,
		EXAMPLE_MESSAGE_PERFORMANCE)
	return ioutil.ReadFile(filepath)
}
