package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"

	log "github.com/sirupsen/logrus"
)

const (
	DocsHandlersSrcDir = "github.com/grokify/chathooks/docs/handlers"
)

func DocsHandlersDir() string {
	return path.Join(os.Getenv("GOPATH"), "src", DocsHandlersSrcDir)
}

// Configuration is the webhook proxy configuration struct.
type Configuration struct {
	Port           int
	EmojiURLFormat string
	LogrusLogLevel log.Level
	IconBaseURL    string
}

func ReadConfigurationFile(filepath string) (Configuration, error) {
	var configuration Configuration
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return configuration, err
	}
	err = json.Unmarshal(bytes, &configuration)
	return configuration, err
}

// Address returns the port address as a string with a `:` prefix
func (c *Configuration) Address() string {
	return fmt.Sprintf(":%d", c.Port)
}

func (c *Configuration) GetAppIconURL(appSlug string) (*url.URL, error) {
	return buildIconURL(c.IconBaseURL, appSlug)
}
