package config

import (
	"fmt"
)

type Configuration struct {
	Port           int
	EmojiURLFormat string
}

func (c *Configuration) Address() string {
	return fmt.Sprintf(":%d", c.Port)
}
