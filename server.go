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
	router.POST("/slack/glip/:glipguid", s2gHandler.HandleFastHTTP)
	router.POST("/slack/glip/:glipguid/", s2gHandler.HandleFastHTTP)

	log.Fatal(fasthttp.ListenAndServe(config.FastHTTPPort(), router.Handler))
}

type Configuration struct {
	Port           int
	EmojiURLPrefix string
	EmojiURLSuffix string
}

func (config *Configuration) FastHTTPPort() string {
	return fmt.Sprintf(":%v", config.Port)
}
