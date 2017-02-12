package glipwebhookproxy

import (
	log "github.com/Sirupsen/logrus"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"

	"github.com/grokify/glip-go-webhook"
	"github.com/grokify/glip-webhook-proxy-go/src/config"
	"github.com/grokify/glip-webhook-proxy-go/src/handlers/slack"
	"github.com/grokify/glip-webhook-proxy-go/src/handlers/travisci"
)

const (
	ROUTE_SLACK_IN_GLIP           = "/webhook/slack/in/glip/:glipguid"
	ROUTE_SLACK_IN_GLIP_SLASH     = "/webhook/slack/in/glip/:glipguid/"
	ROUTE_TRAVISCI_OUT_GLIP       = "/webhook/travisci/out/glip/:glipguid"
	ROUTE_TRAVISCI_OUT_GLIP_SLASH = "/webhook/travisci/out/glip/:glipguid/"
)

// StartServer initializes and starts the webhook proxy server
func StartServer(cfg config.Configuration) {
	//log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(cfg.LogLevel)

	glip, _ := glipwebhook.NewGlipWebhookClient("")

	router := fasthttprouter.New()

	router.GET("/", HomeHandler)

	slackInHandler := slack.NewSlackToGlipHandler(cfg, glip)
	router.POST(ROUTE_SLACK_IN_GLIP, slackInHandler.HandleFastHTTP)
	router.POST(ROUTE_SLACK_IN_GLIP_SLASH, slackInHandler.HandleFastHTTP)

	travisciOutHandler := travisci.NewTravisciOutToGlipHandler(cfg, glip)
	router.POST(ROUTE_TRAVISCI_OUT_GLIP, travisciOutHandler.HandleFastHTTP)
	router.POST(ROUTE_TRAVISCI_OUT_GLIP_SLASH, travisciOutHandler.HandleFastHTTP)

	log.Fatal(fasthttp.ListenAndServe(cfg.Address(), router.Handler))
}
