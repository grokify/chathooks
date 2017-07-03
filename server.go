package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"

	"github.com/grokify/webhookproxy/src/adapters"
	"github.com/grokify/webhookproxy/src/config"
	"github.com/grokify/webhookproxy/src/handlers"
	"github.com/grokify/webhookproxy/src/handlers/appsignal"
	"github.com/grokify/webhookproxy/src/handlers/pingdom"
	"github.com/grokify/webhookproxy/src/handlers/raygun"
	"github.com/grokify/webhookproxy/src/handlers/runscope"
	"github.com/grokify/webhookproxy/src/handlers/semaphore"
	"github.com/grokify/webhookproxy/src/handlers/slack"
	"github.com/grokify/webhookproxy/src/handlers/statuspage"
	"github.com/grokify/webhookproxy/src/handlers/travisci"
	"github.com/grokify/webhookproxy/src/handlers/userlike"
	"github.com/grokify/webhookproxy/src/handlers/victorops"
)

const (
	RouteSlackIn      = "/hook/slack/in/:webhookuid"
	RouteSlackInSlash = "/hook/slack/in/:webhookuid/"
)

func buildHookOutRoutes(handlerKey string, msgDir string) []string {
	routes := []string{}
	routes = append(routes, fmt.Sprintf("/hook/%s/%s/:webhookuid", handlerKey, msgDir))
	routes = append(routes, fmt.Sprintf("/hook/%s/%s/:webhookuid/", handlerKey, msgDir))
	return routes
}

func addRoutes(router *fasthttprouter.Router, handler WebhookHandler) *fasthttprouter.Router {
	routes := buildHookOutRoutes(handler.HandlerKey(), handler.MessageDirection())
	for _, route := range routes {
		router.POST(route, handler.HandleFastHTTP)
	}
	return router
}

type WebhookHandler interface {
	HandlerKey() string
	MessageDirection() string
	HandleFastHTTP(*fasthttp.RequestCtx)
}

// StartServer initializes and starts the webhook proxy server
func StartServer(cfg config.Configuration) {
	//log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(cfg.LogLevel)

	adapter, err := adapters.NewGlipAdapter("")
	if err != nil {
		panic("Cannot build adapter")
	}

	router := fasthttprouter.New()

	router.GET("/", handlers.HomeHandler)

	// Add handlers for services sending outbound formatted hooks
	router = addRoutes(router, appsignal.NewHandler(cfg, &adapter))
	router = addRoutes(router, pingdom.NewHandler(cfg, &adapter))
	router = addRoutes(router, raygun.NewHandler(cfg, &adapter))
	router = addRoutes(router, semaphore.NewHandler(cfg, &adapter))
	router = addRoutes(router, runscope.NewHandler(cfg, &adapter))
	router = addRoutes(router, statuspage.NewHandler(cfg, &adapter))
	router = addRoutes(router, travisci.NewHandler(cfg, &adapter))
	router = addRoutes(router, userlike.NewHandler(cfg, &adapter))
	router = addRoutes(router, victorops.NewHandler(cfg, &adapter))

	// Add handlers for services sending inbound hooks for Slack
	router = addRoutes(router, slack.NewHandler(cfg, &adapter))

	log.WithFields(log.Fields{
		"type": "http.server.start"}).
		Info(fmt.Sprintf("Listening on port %v", cfg.Port))

	log.Fatal(fasthttp.ListenAndServe(cfg.Address(), router.Handler))

}

func main() {
	cfg := config.Configuration{
		Port:           8080,
		EmojiURLFormat: "https://grokify.github.io/emoji/assets/images/%s.png",
		IconBaseURL:    "http://grokify.github.io/webhookproxy/icons/",
		LogLevel:       log.DebugLevel}

	StartServer(cfg)
}
