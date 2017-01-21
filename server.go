package glipwebhookproxy

import (
	"log"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func StartServer(config Configuration) {
	router := fasthttprouter.New()

	router.GET("/", HomeHandler)

	s2gHandler := NewSlackToGlipHandler(config)
	router.POST("/slack/glip/:glipguid", s2gHandler.HandleFastHTTP)
	router.POST("/slack/glip/:glipguid/", s2gHandler.HandleFastHTTP)

	log.Fatal(fasthttp.ListenAndServe(config.Port, router.Handler))
}

type Configuration struct {
	Port           string
	EmojiURLPrefix string
	EmojiURLSuffix string
}
