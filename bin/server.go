package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/grokify/webhook-proxy-go"
	"github.com/grokify/webhook-proxy-go/src/config"
)

func main() {
	cfg := config.Configuration{
		Port:           8080,
		EmojiURLFormat: "https://grokify.github.io/emoji/assets/images/%s.png",
		LogLevel:       log.DebugLevel}

	webhookproxy.StartServer(cfg)
}
