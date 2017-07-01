package webhookproxy

import (
	log "github.com/Sirupsen/logrus"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"

	"github.com/grokify/chatmore/src/adapters"
	"github.com/grokify/chatmore/src/config"
	"github.com/grokify/chatmore/src/handlers"
	"github.com/grokify/chatmore/src/handlers/appsignal"
	"github.com/grokify/chatmore/src/handlers/pingdom"
	"github.com/grokify/chatmore/src/handlers/runscope"
	"github.com/grokify/chatmore/src/handlers/slack"
	"github.com/grokify/chatmore/src/handlers/travisci"
)

const (
	RouteAppsignalOut      = "/webhook/appsignal/out/:webhookuid"
	RouteAppsignalOutSlash = "/webhook/appsignal/out/:webhookuid/"
	RouteSlackIn           = "/webhook/slack/in/:webhookuid"
	RouteSlackInSlash      = "/webhook/slack/in/:webhookuid/"
	RoutePingdomOut        = "/webhook/pingdom/out/:webhookuid"
	RoutePingdomOutSlash   = "/webhook/pingdom/out/:webhookuid/"
	RouteRunscopeOut       = "/webhook/runscope/out/:webhookuid"
	RouteRunscopeOutSlash  = "/webhook/runscope/out/:webhookuid/"
	RouteTravisciOut       = "/webhook/travisci/out/:webhookuid"
	RouteTravisciOutSlash  = "/webhook/travisci/out/:webhookuid/"
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

	appsignalOutHandler := appsignal.NewHandler(cfg, &adapter)
	router.POST(RouteAppsignalOut, appsignalOutHandler.HandleFastHTTP)
	router.POST(RouteAppsignalOutSlash, appsignalOutHandler.HandleFastHTTP)

	pingdomOutHandler := pingdom.NewHandler(cfg, &adapter)
	router.POST(RoutePingdomOut, pingdomOutHandler.HandleFastHTTP)
	router.POST(RoutePingdomOutSlash, pingdomOutHandler.HandleFastHTTP)

	runscopeOutHandler := runscope.NewHandler(cfg, &adapter)
	router.POST(RouteRunscopeOut, runscopeOutHandler.HandleFastHTTP)
	router.POST(RouteRunscopeOutSlash, runscopeOutHandler.HandleFastHTTP)

	slackInHandler := slack.NewHandler(cfg, &adapter)
	router.POST(RouteSlackIn, slackInHandler.HandleFastHTTP)
	router.POST(RouteSlackInSlash, slackInHandler.HandleFastHTTP)

	travisciOutHandler := travisci.NewHandler(cfg, &adapter)
	router.POST(RouteTravisciOut, travisciOutHandler.HandleFastHTTP)
	router.POST(RouteTravisciOutSlash, travisciOutHandler.HandleFastHTTP)

	log.Fatal(fasthttp.ListenAndServe(cfg.Address(), router.Handler))
}
