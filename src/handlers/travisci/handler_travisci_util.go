package travisci

import (
	"io/ioutil"
	"path"

	"github.com/grokify/glip-go-webhook"
	"github.com/grokify/glip-webhook-proxy-go/src/config"
)

var (
	EXAMPLE_BUILD = "example__travisci__build.json"
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
		HANDLER_KEY,
		EXAMPLE_BUILD)
	return ioutil.ReadFile(filepath)
}

/*

Webhook Notification Reference

https://docs.travis-ci.com/user/notifications#Configuring-webhook-notifications

Format:

"Build <%{build_url}|#%{build_number}> (<%{compare_url}|%{commit}>) of %{repository}@%{branch} by %{author} %{result} in %{duration}"

*/
