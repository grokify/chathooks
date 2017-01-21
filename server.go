package glipwebhookproxy

import (
	"fmt"
	"log"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func StartServer(config Configuration) {
	router := fasthttprouter.New()

	router.GET("/", HomeHandler)

	s2gHandler := NewSlackToGlipHandler(config)
	router.POST("/webhook/slack/glip/:glipguid", s2gHandler.HandleFastHTTP)
	router.POST("/webhook/slack/glip/:glipguid/", s2gHandler.HandleFastHTTP)

	log.Fatal(fasthttp.ListenAndServe(config.Address(), router.Handler))
}

type Configuration struct {
	Port           int
	EmojiURLPrefix string
	EmojiURLSuffix string
}

func (config *Configuration) Address() string {
	return fmt.Sprintf(":%v", config.Port)
}
