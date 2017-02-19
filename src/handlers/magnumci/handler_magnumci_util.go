package magnumci

import (
	"io/ioutil"
	"path"

	"github.com/grokify/commonchat"
	"github.com/grokify/webhook-proxy-go/src/config"
)

var (
	ExamplePayloadBuildFile = "example__build.json"
)

func ExampleMessage() (commonchat.Message, error) {
	bytes, err := ExampleMessageBytes()
	if err != nil {
		return commonchat.Message{}, err
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
