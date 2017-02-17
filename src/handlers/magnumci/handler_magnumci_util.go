package magnumci

import (
	"io/ioutil"
	"path"

	"github.com/grokify/glip-go-webhook"
	"github.com/grokify/glip-webhook-proxy-go/src/config"
)

var (
	ExamplePayloadBuildFile = "example__build.json"
)

func ExampleMessageGlip() (glipwebhook.GlipWebhookMessage, error) {
	bytes, err := ExampleMessageBytes()
	if err != nil {
		return glipwebhook.GlipWebhookMessage{}, err
	}

	return Normalize(bytes)
}

func ExampleMessageBytes() ([]byte, error) {
	filepath := path.Join(
		config.DOC_HANDLERS_REL_DIR,
		HandlerKey,
		ExamplePayloadBuildFile)
	return ioutil.ReadFile(filepath)
}
