package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/buaazp/fasthttprouter"
	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
	"github.com/eawsy/aws-lambda-go-event/service/lambda/runtime/event/apigatewayproxyevt"
	//"github.com/grokify/commonchat"
	ccglip "github.com/grokify/commonchat/glip"
	ccslack "github.com/grokify/commonchat/slack"
	"github.com/grokify/gotilla/fmt/fmtutil"
	fhu "github.com/grokify/gotilla/net/fasthttputil"
	nhu "github.com/grokify/gotilla/net/nethttputil"
	"github.com/grokify/gotilla/strings/stringsutil"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"

	"github.com/grokify/chathooks/src/adapters"
	"github.com/grokify/chathooks/src/config"
	"github.com/grokify/chathooks/src/models"

	"github.com/grokify/chathooks/src/handlers"
	"github.com/grokify/chathooks/src/handlers/aha"
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

/*

Use the `CHATHOOKS_TOKENS` environment variable to load secret
tokens as a comma delimited string.

*/

const (
	ParamNameInput  = "inputType"
	ParamNameOutput = "outputType"
	ParamNameURL    = "url"
	ParamNameToken  = "token"
	EnvEnvPath      = "ENV_PATH"
	EnvEngine       = "CHATHOOKS_ENGINE" // aws, nethttp, fasthttp
	EnvTokens       = "CHATHOOKS_TOKENS"
)

type HandlerSet struct {
	Handlers map[string]Handler
}

type Handler interface {
	HandleCanonical(hookData models.HookData) []models.ErrorInfo
	HandleAwsLambda(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
	HandleEawsyLambda(event *apigatewayproxyevt.Event, ctx *runtime.Context) (events.APIGatewayProxyResponse, error)
	HandleFastHTTP(ctx *fasthttp.RequestCtx)
	HandleNetHTTP(res http.ResponseWriter, req *http.Request)
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
	"IconBaseURL":    "http://grokify.github.io/chathooks/icons/",
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
	cfgData := config.Configuration{
		Port:           8080,
		EmojiURLFormat: "https://grokify.github.io/emoji/assets/images/%s.png",
		LogrusLogLevel: 5,
		IconBaseURL:    "http://grokify.github.io/chathooks/icons/",
	}

	fmtutil.PrintJSON(cfgData)
	adapterSet := adapters.NewAdapterSet()
	glipAdapter, err := ccglip.NewGlipAdapter("")
	if err != nil {
		log.Fatal(err)
	}
	adapterSet.Adapters["glip"] = glipAdapter
	slackAdapter, err := ccslack.NewSlackAdapter("")
	if err != nil {
		log.Fatal(err)
	}
	adapterSet.Adapters["slack"] = slackAdapter

	hf := HandlerFactory{Config: cfgData, AdapterSet: adapterSet}

	handlerSet := HandlerSet{Handlers: map[string]Handler{
		"aha":        hf.InflateHandler(aha.NewHandler()),
		"appsignal":  hf.InflateHandler(appsignal.NewHandler()),
		"apteligent": hf.InflateHandler(apteligent.NewHandler()),
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

	svcInfo := ServiceInfo{
		Config:       cfgData,
		AdapterSet:   adapterSet,
		HandlerSet:   handlerSet,
		RequireToken: false,
		Tokens:       map[string]int{},
	}
	tokens := stringsutil.SplitCondenseSpace(os.Getenv(EnvTokens), ",")
	for _, token := range tokens {
		svcInfo.Tokens[token] = 1
	}

	return svcInfo
}

var serviceInfo = getConfig()

func HandleAwsLambda(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if len(serviceInfo.Tokens) > 0 {
		token, ok := req.QueryStringParameters[ParamNameToken]
		if !ok {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusUnauthorized,
				Body:       "Required Token not found"}, nil
		}
		if _, ok := serviceInfo.Tokens[token]; !ok {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusUnauthorized,
				Body:       "Required Token not valid"}, nil
		}
	}
	inputType, ok := req.QueryStringParameters[models.QueryParamInputType]
	if !ok || len(strings.TrimSpace(inputType)) == 0 {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "InputType not found"}, nil
	}

	handler, ok := serviceInfo.HandlerSet.Handlers[inputType]
	if !ok {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf("Input Handler Not found for: %v\n", inputType)}, nil
	}

	return handler.HandleAwsLambda(ctx, req)
}

func HandleEawsyLambda(event *apigatewayproxyevt.Event, ctx *runtime.Context) (events.APIGatewayProxyResponse, error) {
	if len(serviceInfo.Tokens) > 0 {
		token, ok := event.QueryStringParameters[ParamNameToken]
		if !ok {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusUnauthorized,
				Body:       "Required Token not found"}, nil
		}
		if _, ok := serviceInfo.Tokens[token]; !ok {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusUnauthorized,
				Body:       "Required Token not valid"}, nil
		}
	}

	inputType, ok := event.QueryStringParameters[models.QueryParamInputType]
	if !ok || len(strings.TrimSpace(inputType)) == 0 {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "InputType not found"}, nil
	}

	handler, ok := serviceInfo.HandlerSet.Handlers[inputType]
	if !ok {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf("Input Handler Not found for: %v\n", inputType)}, nil
	}

	return handler.HandleEawsyLambda(event, ctx)
}

type AnyHTTPHandler struct {
	Config     config.Configuration
	AdapterSet adapters.AdapterSet
	HandlerSet HandlerSet
}

var anyHTTPHandler = AnyHTTPHandler{
	Config:     serviceInfo.Config,
	AdapterSet: serviceInfo.AdapterSet,
	HandlerSet: serviceInfo.HandlerSet,
}

func (h *AnyHTTPHandler) HandleNetHTTP(res http.ResponseWriter, req *http.Request) {
	fmt.Println("HANDLE_NetHTTP")
	if len(serviceInfo.Tokens) > 0 {
		token := nhu.GetReqHeader(req, ParamNameToken)
		if len(token) == 0 {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		if _, ok := serviceInfo.Tokens[token]; !ok {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
	}
	inputType := nhu.GetReqHeader(req, ParamNameInput)

	if handler, ok := h.HandlerSet.Handlers[inputType]; ok {
		fmt.Printf("Input_Handler_Found_Processing [%v]\n", inputType)
		handler.HandleNetHTTP(res, req)
	} else {
		fmt.Printf("Input_Handler_Not_Found [%v]\n", inputType)
	}
}

func (h *AnyHTTPHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	fmt.Println("HANDLE_FastHTTP")
	if len(serviceInfo.Tokens) > 0 {
		token := fhu.GetReqHeader(ctx, ParamNameToken)
		if len(token) == 0 {
			ctx.SetStatusCode(http.StatusUnauthorized)
			log.Warn("E_NO_TOKEN")
			return
		}
		if _, ok := serviceInfo.Tokens[token]; !ok {
			ctx.SetStatusCode(http.StatusUnauthorized)
			log.Warn("E_INCORRECT_TOKEN")
			return
		}
	}

	inputType := fhu.GetReqHeader(ctx, ParamNameInput)

	if handler, ok := h.HandlerSet.Handlers[inputType]; ok {
		fmt.Printf("Input_Handler_Found_Processing [%v]\n", inputType)
		handler.HandleFastHTTP(ctx)
	} else {
		fmt.Printf("Input_Handler_Not_Found [%v]\n", inputType)
	}
}

func serveNetHttp() {
	http.Handle("/hook", http.HandlerFunc(anyHTTPHandler.HandleNetHTTP))
	http.Handle("/hook/", http.HandlerFunc(anyHTTPHandler.HandleNetHTTP))

	log.Fatal(fasthttp.ListenAndServe(":8080", nil))
}

func serveFastHttp() {
	router := fasthttprouter.New()
	router.GET("/hook", anyHTTPHandler.HandleFastHTTP)
	router.GET("/hook/", anyHTTPHandler.HandleFastHTTP)
	router.POST("/hook", anyHTTPHandler.HandleFastHTTP)
	router.POST("/hook/", anyHTTPHandler.HandleFastHTTP)

	log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))
}

func serveAwsLambda() {
	lambda.Start(HandleAwsLambda)
}

func main() {
	if len(strings.TrimSpace(os.Getenv(EnvEnvPath))) > 0 {
		err := godotenv.Load(os.Getenv(EnvEnvPath))
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	engine := strings.ToLower(strings.TrimSpace(os.Getenv(EnvEngine)))
	switch engine {
	case "aws":
		serveAwsLambda()
	case "nethttp":
		serveNetHttp()
	case "fasthttp":
		serveFastHttp()
	}
}
