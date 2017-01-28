package config

import (
	"fmt"

	"github.com/Sirupsen/logrus"
)

type Configuration struct {
	Port           int
	EmojiURLFormat string
	LogLevel       logrus.Level
}

// Address returns the port address as a string with a `:` prefix
func (c *Configuration) Address() string {
	return fmt.Sprintf(":%d", c.Port)
}
