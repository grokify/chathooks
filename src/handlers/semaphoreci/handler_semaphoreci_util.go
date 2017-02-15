package semaphoreci

import (
	"io/ioutil"
	"path"

	"github.com/grokify/glip-go-webhook"
	"github.com/grokify/glip-webhook-proxy-go/src/config"
)

var (
	EXAMPLE_FILE_BUILD  = "example__build.json"
	EXAMPLE_FILE_DEPLOY = "example__deploy.json"
)

func ExampleBuildMessageGlip() (glipwebhook.GlipWebhookMessage, error) {
	msg, err := ExampleBuildMessageSource()
	if err != nil {
		return glipwebhook.GlipWebhookMessage{}, err
	}
	return NormalizeSemaphoreciBuildOutMessage(msg), nil
}

func ExampleBuildMessageSource() (SemaphoreciBuildOutMessage, error) {
	bytes, err := ExampleBuildMessageBytes()
	if err != nil {
		return SemaphoreciBuildOutMessage{}, err
	}
	return SemaphoreciBuildOutMessageFromBytes(bytes)
}

func ExampleBuildMessageBytes() ([]byte, error) {
	filepath := path.Join(
		config.DOC_HANDLERS_REL_DIR,
		HANDLER_KEY,
		EXAMPLE_FILE_BUILD)
	return ioutil.ReadFile(filepath)
}

func ExampleDeployMessageGlip() (glipwebhook.GlipWebhookMessage, error) {
	msg, err := ExampleDeployMessageSource()
	if err != nil {
		return glipwebhook.GlipWebhookMessage{}, err
	}
	return NormalizeSemaphoreciDeployOutMessage(msg), nil
}

func ExampleDeployMessageSource() (SemaphoreciDeployOutMessage, error) {
	bytes, err := ExampleDeployMessageBytes()
	if err != nil {
		return SemaphoreciDeployOutMessage{}, err
	}
	return SemaphoreciDeployOutMessageFromBytes(bytes)
}

func ExampleDeployMessageBytes() ([]byte, error) {
	filepath := path.Join(
		config.DOC_HANDLERS_REL_DIR,
		HANDLER_KEY,
		EXAMPLE_FILE_DEPLOY)
	return ioutil.ReadFile(filepath)
}
