package main

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/buaazp/fasthttprouter"
	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
	"github.com/eawsy/aws-lambda-go-event/service/lambda/runtime/event/apigatewayproxyevt"
	"github.com/grokify/gotilla/fmt/fmtutil"
	"github.com/valyala/fasthttp"

	"github.com/grokify/chathooks/src/adapters"
	"github.com/grokify/chathooks/src/config"
	"github.com/grokify/chathooks/src/handlers"
	"github.com/grokify/chathooks/src/models"

	"github.com/grokify/chathooks/src/handlers/appsignal"
	"github.com/grokify/chathooks/src/handlers/apteligent"
	"github.com/grokify/chathooks/src/handlers/circleci"
	"github.com/grokify/chathooks/src/handlers/codeship"
	"github.com/grokify/chathooks/src/handlers/confluence"
	"github.com/grokify/chathooks/src/handlers/datadog"
	"github.com/grokify/chathooks/src/handlers/deskdotcom"
	"github.com/grokify/chathooks/src/handlers/enchant"
	"github.com/grokify/chathooks/src/handlers/gosquared"
	"github.com/grokify/chathooks/src/handlers/gosquared2"
	"github.com/grokify/chathooks/src/handlers/heroku"
	"github.com/grokify/chathooks/src/handlers/librato"
	"github.com/grokify/chathooks/src/handlers/magnumci"
	"github.com/grokify/chathooks/src/handlers/marketo"
	"github.com/grokify/chathooks/src/handlers/opsgenie"
	"github.com/grokify/chathooks/src/handlers/papertrail"
	"github.com/grokify/chathooks/src/handlers/pingdom"
	"github.com/grokify/chathooks/src/handlers/raygun"
	"github.com/grokify/chathooks/src/handlers/runscope"
	"github.com/grokify/chathooks/src/handlers/semaphore"
	"github.com/grokify/chathooks/src/handlers/slack"
	"github.com/grokify/chathooks/src/handlers/statuspage"
	"github.com/grokify/chathooks/src/handlers/travisci"
	"github.com/grokify/chathooks/src/handlers/userlike"
	"github.com/grokify/chathooks/src/handlers/victorops"
)

const (
	ParamNameInput  = "inputType"
	ParamNameOutput = "outputType"
	ParamNameURL    = "url"
	ParamNameToken  = "token"
)

type HandlerSet struct {
	Handlers map[string]Handler
}

type Handler interface {
	HandleEawsyLambda(event *apigatewayproxyevt.Event, ctx *runtime.Context) (models.AwsAPIGatewayProxyOutput, error)
	HandleFastHTTP(ctx *fasthttp.RequestCtx)
	HandleCanonical(hookData models.HookData) []models.ErrorInfo
}

type ServiceInfo struct {
	Config       config.Configuration
	AdapterSet   adapters.AdapterSet
	HandlerSet   HandlerSet
	RequireToken bool
	Tokens       map[string]int
}

/*
{
	"Port": 8080,
	"EmojiURLFormat": "https://grokify.github.io/emoji/assets/images/%s.png",
	"IconBaseURL":    "http://grokify.github.io/webhookproxy/icons/",
	"LogrusLogLevel": 5,
}
*/

type HandlerFactory struct {
	Config     config.Configuration
	AdapterSet adapters.AdapterSet
}

func (hf *HandlerFactory) NewHandler(normalize handlers.Normalize) handlers.Handler {
	return handlers.Handler{
		Config:     hf.Config,
		AdapterSet: hf.AdapterSet,
		Normalize:  normalize}
}

func (hf *HandlerFactory) InflateHandler(handler handlers.Handler) handlers.Handler {
	handler.Config = hf.Config
	handler.AdapterSet = hf.AdapterSet
	return handler
}

func getConfig() ServiceInfo {
	/*
		cfgData, err := config.ReadConfigurationFile(configFilepath)
		if err != nil {
			log.Fatal("Configuration File [%v] not found failed with error [%v].", configFilepath, err)
		}

	*/
	cfgData := config.Configuration{
		Port:           8080,
		EmojiURLFormat: "https://grokify.github.io/emoji/assets/images/%s.png",
		LogrusLogLevel: 5,
		IconBaseURL:    "http://grokify.github.io/webhookproxy/icons/",
	}

	fmtutil.PrintJSON(cfgData)
	adapterSet := adapters.NewAdapterSet()
	glipAdapter, err := adapters.NewGlipAdapter("")
	if err != nil {
		log.Fatal(err)
	}
	adapterSet.Adapters["glip"] = glipAdapter
	slackAdapter, err := adapters.NewSlackAdapter("")
	if err != nil {
		log.Fatal(err)
	}
	adapterSet.Adapters["slack"] = slackAdapter

	hf := HandlerFactory{Config: cfgData, AdapterSet: adapterSet}

	handlerSet := HandlerSet{Handlers: map[string]Handler{
		"appsignal":  appsignal.NewHandler(cfgData, adapterSet),
		"apteligent": apteligent.NewHandler(cfgData, adapterSet),
		"circleci":   hf.InflateHandler(circleci.NewHandler()),
		"codeship":   hf.InflateHandler(codeship.NewHandler()),
		"confluence": hf.InflateHandler(confluence.NewHandler()),
		"datadog":    hf.InflateHandler(datadog.NewHandler()),
		"deskdotcom": hf.InflateHandler(deskdotcom.NewHandler()),
		"enchant":    hf.InflateHandler(enchant.NewHandler()),
		"gosquared":  hf.InflateHandler(gosquared.NewHandler()),
		"gosquared2": hf.InflateHandler(gosquared2.NewHandler()),
		"heroku":     hf.InflateHandler(heroku.NewHandler()),
		"librato":    hf.InflateHandler(librato.NewHandler()),
		"magnumci":   hf.InflateHandler(magnumci.NewHandler()),
		"marketo":    hf.InflateHandler(marketo.NewHandler()),
		"opsgenie":   hf.InflateHandler(opsgenie.NewHandler()),
		"papertrail": hf.InflateHandler(papertrail.NewHandler()),
		"pingdom":    hf.InflateHandler(pingdom.NewHandler()),
		"raygun":     hf.InflateHandler(raygun.NewHandler()),
		"runscope":   hf.InflateHandler(runscope.NewHandler()),
		"semaphore":  hf.InflateHandler(semaphore.NewHandler()),
		"slack":      hf.InflateHandler(slack.NewHandler()),
		"statuspage": hf.InflateHandler(statuspage.NewHandler()),
		"travisci":   hf.InflateHandler(travisci.NewHandler()),
		"userlike":   hf.InflateHandler(userlike.NewHandler()),
		"victorops":  hf.InflateHandler(victorops.NewHandler()),
	}}

	return ServiceInfo{
		Config:       cfgData,
		AdapterSet:   adapterSet,
		HandlerSet:   handlerSet,
		RequireToken: false,
		Tokens:       map[string]int{},
	}
}

var serviceInfo = getConfig()

func HandleEawsyLambda(event *apigatewayproxyevt.Event, ctx *runtime.Context) (models.AwsAPIGatewayProxyOutput, error) {
	if len(serviceInfo.Tokens) > 0 {
		token, ok := event.QueryStringParameters[ParamNameToken]
		if !ok {
			return models.AwsAPIGatewayProxyOutput{
				StatusCode: 401,
				Body:       "Required Token not found"}, nil
		}
		if _, ok := serviceInfo.Tokens[token]; !ok {
			return models.AwsAPIGatewayProxyOutput{
				StatusCode: 401,
				Body:       "Required Token not valid"}, nil
		}
	}

	inputType, ok := event.QueryStringParameters[models.QueryParamInputType]
	if !ok || len(strings.TrimSpace(inputType)) == 0 {
		return models.AwsAPIGatewayProxyOutput{
			StatusCode: 400,
			Body:       "InputType not found"}, nil
	}

	handler, ok := serviceInfo.HandlerSet.Handlers[inputType]
	if !ok {
		return models.AwsAPIGatewayProxyOutput{
			StatusCode: 400,
			Body:       fmt.Sprintf("Input Handler Not found for: %v\n", inputType)}, nil
	}

	return handler.HandleEawsyLambda(event, ctx)
}

type FastHTTPHandler struct {
	Config     config.Configuration
	AdapterSet adapters.AdapterSet
	HandlerSet HandlerSet
}

func (h *FastHTTPHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	fmt.Println("HANDLE_FastHTTP")
	if len(serviceInfo.Tokens) > 0 {
		token := strings.TrimSpace(string(ctx.QueryArgs().Peek(ParamNameToken)))
		if len(token) == 0 {
			ctx.SetStatusCode(401)
			return
		}
		if _, ok := serviceInfo.Tokens[token]; !ok {
			ctx.SetStatusCode(401)
			return
		}
	}

	inputType := strings.TrimSpace(string(ctx.QueryArgs().Peek(ParamNameInput)))

	if handler, ok := h.HandlerSet.Handlers[inputType]; ok {
		fmt.Printf("Input_Handler_Found_Processing [%v]\n", inputType)
		handler.HandleFastHTTP(ctx)
	} else {
		fmt.Printf("Input_Handler_Not_Found [%v]\n", inputType)
	}
}

func main() {
	fh := FastHTTPHandler{
		Config:     serviceInfo.Config,
		AdapterSet: serviceInfo.AdapterSet,
		HandlerSet: serviceInfo.HandlerSet,
	}

	router := fasthttprouter.New()
	router.GET("/", handlers.HomeHandler)
	router.GET("/hook", fh.HandleFastHTTP)
	router.GET("/hooks", fh.HandleFastHTTP)
	router.POST("/hook", fh.HandleFastHTTP)
	router.POST("/hooks", fh.HandleFastHTTP)

	log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))
}
