package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"

	"github.com/caarlos0/env"
	"github.com/rs/zerolog"
)

const (
	DocsHandlersSrcDir = "github.com/grokify/chathooks/docs/handlers"
	IconBaseURL        = "http://grokify.github.io/chathooks/icons/"
	EmojiURLFormat     = "https://grokify.github.io/emoji/assets/images/%s.png"

	InfoInputMessageParseBegin   = "INFO - Input Message Parse Begin"
	ErrorInputMessageParseFailed = "FAIL - Input Message Parse Failed"
)

func DocsHandlersDir() string {
	return path.Join(os.Getenv("GOPATH"), "src", DocsHandlersSrcDir)
}

// Configuration is the webhook proxy configuration struct.
type Configuration struct {
	Port           int      `env:"PORT" envDefault:"3000"`
	Engine         string   `env:"CHATHOOKS_ENGINE" envDefault:"fasthttp"`
	HomeUrl        string   `env:"CHATHOOKS_HOME_URL"`
	WebhookUrl     string   `env:"CHATHOOKS_WEBHOOK_URL"`
	Tokens         []string `env:"CHATHOOKS_TOKENS" envSeparator:","`
	LogFormat      string   `env:"CHATHOOKS_LOG_FORMAT"`
	EmojiURLFormat string
	IconBaseURL    string
	LogLevel       zerolog.Level
}

func NewConfigurationEnv() (Configuration, error) {
	cfg := Configuration{}
	if err := env.Parse(&cfg); err != nil {
		return cfg, err
	}
	cfg.EmojiURLFormat = EmojiURLFormat
	cfg.IconBaseURL = IconBaseURL
	cfg.LogLevel = 1
	return cfg, nil
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
