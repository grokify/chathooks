package appsignal

import (
	"io/ioutil"
	"path"

	cc "github.com/commonchat/commonchat-go"
	"github.com/grokify/webhook-proxy-go/src/config"
)

const (
	EXAMPLE_MESSAGE_EXCEPTION   = "example__exception.json"
	EXAMPLE_MESSAGE_MARKER      = "example__marker.json"
	EXAMPLE_MESSAGE_PERFORMANCE = "example__performance.json"
)

func ExampleMessageMarker() (cc.Message, error) {
	bytes, err := ExampleMessageMarkerBytes()
	if err != nil {
		return cc.Message{}, err
	}
	return Normalize(bytes)
}

func ExampleMessageMarkerBytes() ([]byte, error) {
	filepath := path.Join(
		config.DOC_HANDLERS_REL_DIR,
		HandlerKey,
		EXAMPLE_MESSAGE_MARKER)
	return ioutil.ReadFile(filepath)
}

func ExampleMessageException() (cc.Message, error) {
	bytes, err := ExampleMessageExceptionBytes()
	if err != nil {
		return cc.Message{}, err
	}
	return Normalize(bytes)
}

func ExampleMessageExceptionBytes() ([]byte, error) {
	filepath := path.Join(
		config.DOC_HANDLERS_REL_DIR,
		HandlerKey,
		EXAMPLE_MESSAGE_EXCEPTION)
	return ioutil.ReadFile(filepath)
}

func ExampleMessagePerformance() (cc.Message, error) {
	bytes, err := ExampleMessagePerformanceBytes()
	if err != nil {
		return cc.Message{}, err
	}
	return Normalize(bytes)
}

func ExampleMessagePerformanceBytes() ([]byte, error) {
	filepath := path.Join(
		config.DOC_HANDLERS_REL_DIR,
		HandlerKey,
		EXAMPLE_MESSAGE_PERFORMANCE)
	return ioutil.ReadFile(filepath)
}
