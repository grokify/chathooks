package glipwebhookproxy

import (
	"fmt"
	"log"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

const (
	ROUTE_SLACK_GLIP       = "/webhook/slack/glip/:glipguid"
	ROUTE_SLACK_GLIP_SLASH = "/webhook/slack/glip/:glipguid/"
)

func StartServer(config Configuration) {
	router := fasthttprouter.New()

	router.GET("/", HomeHandler)

	s2gHandler := SlackToGlipHandler{Config: config}
	router.POST(ROUTE_SLACK_GLIP, s2gHandler.HandleFastHTTP)
	router.POST(ROUTE_SLACK_GLIP_SLASH, s2gHandler.HandleFastHTTP)

	log.Fatal(fasthttp.ListenAndServe(config.Address(), router.Handler))
}

type Configuration struct {
	Port           int
	EmojiURLFormat string
}

func (config *Configuration) Address() string {
	return fmt.Sprintf(":%d", config.Port)
}
