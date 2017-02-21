package semaphoreci

import (
	"io/ioutil"
	"path"

	"github.com/commonchat/commonchat-go"
	"github.com/grokify/webhook-proxy-go/src/config"
)

var (
	ExampleFileBuild  = "example__build.json"
	ExampleFileDeploy = "example__deploy.json"
)

func ExampleMessageBuild() (commonchat.Message, error) {
	bytes, err := ExampleMessageBuildBytes()
	if err != nil {
		return commonchat.Message{}, err
	}
	return Normalize(bytes)
}

func ExampleMessageBuildBytes() ([]byte, error) {
	filepath := path.Join(
		config.DOC_HANDLERS_REL_DIR,
		HandlerKey,
		ExampleFileBuild)
	return ioutil.ReadFile(filepath)
}

func ExampleMessageDeploy() (commonchat.Message, error) {
	bytes, err := ExampleMessageDeployBytes()
	if err != nil {
		return commonchat.Message{}, err
	}
	return Normalize(bytes)
}

func ExampleMessageDeployBytes() ([]byte, error) {
	filepath := path.Join(
		config.DOC_HANDLERS_REL_DIR,
		HandlerKey,
		ExampleFileDeploy)
	return ioutil.ReadFile(filepath)
}
