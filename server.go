package webhookproxy

import (
	log "github.com/Sirupsen/logrus"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"

	"github.com/grokify/webhook-proxy-go/src/adapters"
	"github.com/grokify/webhook-proxy-go/src/config"
	"github.com/grokify/webhook-proxy-go/src/handlers"
	"github.com/grokify/webhook-proxy-go/src/handlers/runscope"
	"github.com/grokify/webhook-proxy-go/src/handlers/slack"
	"github.com/grokify/webhook-proxy-go/src/handlers/travisci"
)

const (
	RouteSlackInToGlip          = "/webhook/slack/in/glip/:webhookuid"
	RouteSlackInToGlipSlash     = "/webhook/slack/in/glip/:webhookuid/"
	RouteTravisciOutToGlip      = "/webhook/travisci/out/glip/:webhookuid"
	RouteTravisciOutToGlipSlash = "/webhook/travisci/out/glip/:webhookuid/"
	RouteRunscopeOutToGlip      = "/webhook/runscope/out/glip/:webhookuid"
	RouteRunscopeOutToGlipSlash = "/webhook/runscope/out/glip/:webhookuid/"
)

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

	runscopeOutHandler := runscope.NewRunscopeOutToGlipHandler(cfg, &adapter)
	router.POST(RouteRunscopeOutToGlip, runscopeOutHandler.HandleFastHTTP)
	router.POST(RouteRunscopeOutToGlipSlash, runscopeOutHandler.HandleFastHTTP)

	slackInHandler := slack.NewSlackToGlipHandler(cfg, &adapter)
	router.POST(RouteSlackInToGlip, slackInHandler.HandleFastHTTP)
	router.POST(RouteSlackInToGlipSlash, slackInHandler.HandleFastHTTP)

	travisciOutHandler := travisci.NewTravisciOutToGlipHandler(cfg, &adapter)
	router.POST(RouteTravisciOutToGlip, travisciOutHandler.HandleFastHTTP)
	router.POST(RouteTravisciOutToGlipSlash, travisciOutHandler.HandleFastHTTP)

	log.Fatal(fasthttp.ListenAndServe(cfg.Address(), router.Handler))
}
