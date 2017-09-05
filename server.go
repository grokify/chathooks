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

	"github.com/grokify/webhookproxy/src/adapters"
	"github.com/grokify/webhookproxy/src/config"
	"github.com/grokify/webhookproxy/src/handlers"
	"github.com/grokify/webhookproxy/src/models"

	"github.com/grokify/webhookproxy/src/handlers/appsignal"
	"github.com/grokify/webhookproxy/src/handlers/apteligent"
	"github.com/grokify/webhookproxy/src/handlers/circleci"
	"github.com/grokify/webhookproxy/src/handlers/codeship"
	"github.com/grokify/webhookproxy/src/handlers/confluence"
	"github.com/grokify/webhookproxy/src/handlers/datadog"
	"github.com/grokify/webhookproxy/src/handlers/deskdotcom"
	"github.com/grokify/webhookproxy/src/handlers/enchant"
	"github.com/grokify/webhookproxy/src/handlers/gosquared"
	"github.com/grokify/webhookproxy/src/handlers/gosquared2"
	"github.com/grokify/webhookproxy/src/handlers/heroku"
	"github.com/grokify/webhookproxy/src/handlers/librato"
	"github.com/grokify/webhookproxy/src/handlers/magnumci"
	"github.com/grokify/webhookproxy/src/handlers/marketo"
	"github.com/grokify/webhookproxy/src/handlers/opsgenie"
	"github.com/grokify/webhookproxy/src/handlers/papertrail"
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

type Base struct {
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

func getConfig() Base {
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

	return Base{
		Config:       cfgData,
		AdapterSet:   adapterSet,
		HandlerSet:   handlerSet,
		RequireToken: false,
		Tokens:       map[string]int{},
	}
}

var base = getConfig()

func HandleEawsyLambda(event *apigatewayproxyevt.Event, ctx *runtime.Context) (models.AwsAPIGatewayProxyOutput, error) {
	if len(base.Tokens) > 0 {
		token, ok := event.QueryStringParameters[ParamNameToken]
		if !ok {
			return models.AwsAPIGatewayProxyOutput{
				StatusCode: 401,
				Body:       "Required Token not found"}, nil
		}
		if _, ok := base.Tokens[token]; !ok {
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

	handler, ok := base.HandlerSet.Handlers[inputType]
	if !ok {
		return models.AwsAPIGatewayProxyOutput{
			StatusCode: 400,
			Body:       fmt.Sprintf("Input Handler Not found for: %v\n")}, nil
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
	if len(base.Tokens) > 0 {
		token := strings.TrimSpace(string(ctx.QueryArgs().Peek(ParamNameToken)))
		if len(token) == 0 {
			ctx.SetStatusCode(401)
			return
		}
		if _, ok := base.Tokens[token]; !ok {
			ctx.SetStatusCode(401)
			return
		}
	}

	inputType := strings.TrimSpace(string(ctx.QueryArgs().Peek(ParamNameInput)))

	fmt.Printf("INPUT_Type [%v]\n", inputType)

	if handler, ok := h.HandlerSet.Handlers[inputType]; ok {
		fmt.Printf("Input_Handler_Found_Processing [%v]\n", inputType)
		handler.HandleFastHTTP(ctx)
	} else {
		fmt.Printf("Input_Handler_Not_Found [%v]\n", inputType)
	}
}

func main() {
	fh := FastHTTPHandler{
		Config:     base.Config,
		AdapterSet: base.AdapterSet,
		HandlerSet: base.HandlerSet,
	}

	router := fasthttprouter.New()
	router.GET("/", handlers.HomeHandler)
	router.GET("/hook", fh.HandleFastHTTP)
	router.GET("/hooks", fh.HandleFastHTTP)
	router.POST("/hook", fh.HandleFastHTTP)
	router.POST("/hooks", fh.HandleFastHTTP)

	log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))
}
