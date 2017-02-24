package config

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
)

var (
	//GLIP_ACTIVITY_INCLUDE_INTEGRATION_NAME = true
	DocHandlersDir = "../docs/handlers"
)

// Configuration is the webhook proxy configuration struct.
type Configuration struct {
	Port           int
	EmojiURLFormat string
	LogLevel       log.Level
}

// Address returns the port address as a string with a `:` prefix
func (c *Configuration) Address() string {
	return fmt.Sprintf(":%d", c.Port)
}
