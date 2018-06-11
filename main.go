package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/apex/gateway"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/buaazp/fasthttprouter"
	ccglip "github.com/grokify/commonchat/glip"
	ccslack "github.com/grokify/commonchat/slack"
	cfg "github.com/grokify/gotilla/config"
	"github.com/grokify/gotilla/net/anyhttp"
	hum "github.com/grokify/gotilla/net/httputilmore"
	"github.com/grokify/gotilla/strings/stringsutil"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"

	"github.com/grokify/chathooks/src/adapters"
	"github.com/grokify/chathooks/src/config"
	"github.com/grokify/chathooks/src/models"
	"github.com/grokify/chathooks/src/templates"

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

// CHATHOOKS_URL=http://localhost:8080/hook CHATHOOKS_HOME_URL=http://localhost:8080 go run main.go

const (
	ParamNameInput  = "inputType"
	ParamNameOutput = "outputType"
	ParamNameURL    = "url"
	ParamNameToken  = "token"
	EnvPath         = "ENV_PATH"
	EnvEngine       = "CHATHOOKS_ENGINE" // awslambda, nethttp, fasthttp
	EnvTokens       = "CHATHOOKS_TOKENS"
	EnvWebhookUrl   = "CHATHOOKS_URL"
	EnvHomeUrl      = "CHATHOOKS_HOME_URL"
)

type HandlerSet struct {
	Handlers map[string]Handler
}

type Handler interface {
	HandleCanonical(hookData models.HookData) []models.ErrorInfo
	HandleAwsLambda(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
	HandleFastHTTP(ctx *fasthttp.RequestCtx)
	HandleNetHTTP(res http.ResponseWriter, req *http.Request)
	HandleAnyHTTP(aRes anyhttp.Response, aReq anyhttp.Request)
}

type Service struct {
	Config       config.Configuration
	AdapterSet   adapters.AdapterSet
	HandlerSet   HandlerSet
	RequireToken bool
	Tokens       map[string]int
}

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

func NewService() Service {
	/*
		cfgData := config.Configuration{
			Port:           8080,
			LogrusLogLevel: 5,
			EmojiURLFormat: config.EmojiURLFormat,
			IconBaseURL:    config.IconBaseURL}
	*/

	cfgData, err := config.NewConfigurationEnv()
	if err != nil {
		log.Fatal(err)
	}

	//fmtutil.PrintJSON(cfgData)
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
		"victorops":  hf.InflateHandler(victorops.NewHandler())}}

	svcInfo := Service{
		Config:       cfgData,
		AdapterSet:   adapterSet,
		HandlerSet:   handlerSet,
		RequireToken: false,
		Tokens:       map[string]int{}}
	tokens := stringsutil.SplitCondenseSpace(os.Getenv(EnvTokens), ",")
	for _, token := range tokens {
		svcInfo.Tokens[token] = 1
	}

	return svcInfo
}

func (svc *Service) HandleAwsLambda(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if len(svc.Tokens) > 0 {
		token, ok := req.QueryStringParameters[ParamNameToken]
		if !ok {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusUnauthorized,
				Body:       "Required Token not found"}, nil
		}
		if _, ok := svc.Tokens[token]; !ok {
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

	handler, ok := svc.HandlerSet.Handlers[inputType]
	if !ok {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf("Input Handler Not found for: %v\n", inputType)}, nil
	}

	return handler.HandleAwsLambda(ctx, req)
}

/*
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
*/

func (svc *Service) HandleAnyRequest(aRes anyhttp.Response, aReq anyhttp.Request) {
	log.Info("HANDLE_AnyHTTP")
	if len(svc.Tokens) > 0 {
		token := aReq.QueryArgs().GetString(ParamNameToken)

		if len(token) == 0 {
			aRes.SetStatusCode(http.StatusUnauthorized)
			log.Warn("E_NO_TOKEN")
			return
		}
		fmt.Println("HANDLE_NetHTTP_S2b")
		if _, ok := svc.Tokens[token]; !ok {
			aRes.SetStatusCode(http.StatusUnauthorized)
			log.Warn("E_INCORRECT_TOKEN")
			return
		}
	}
	if err := aReq.ParseForm(); err != nil {
		aRes.SetStatusCode(http.StatusInternalServerError)
		log.Warn("E_CANNOT_PARSE_FORM")
		return
	}
	inputType := aReq.QueryArgs().GetString(ParamNameInput)

	if handler, ok := svc.HandlerSet.Handlers[inputType]; ok {
		fmt.Printf("Input_Handler_Found_Processing [%v]\n", inputType)
		handler.HandleAnyHTTP(aRes, aReq)
	} else {
		fmt.Printf("Input_Handler_Not_Found [%v]\n", inputType)
	}
}

func (svc *Service) HandleNetHTTP(res http.ResponseWriter, req *http.Request) {
	log.Info("HANDLE_NetHTTP")
	svc.HandleAnyRequest(anyhttp.NewResReqNetHttp(res, req))
}

func (svc *Service) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	log.Info("HANDLE_FastHTTP")
	svc.HandleAnyRequest(anyhttp.NewResReqFastHttp(ctx))
}

func (svc *Service) HandleHomeAnyRequest(aRes anyhttp.Response, aReq anyhttp.Request) {
	log.Info("HANDLE_HOME_AnyHTTP")
	fmt.Println(svc.Config.WebhookUrl)
	data := templates.HomeData{
		HomeUrl:    svc.Config.HomeUrl,
		WebhookUrl: svc.Config.WebhookUrl}
	if _, err := aRes.SetBodyBytes([]byte(templates.HomePage(data))); err != nil {
		aRes.SetStatusCode(http.StatusInternalServerError)
	} else {
		aRes.SetStatusCode(http.StatusOK)
		aRes.SetContentType(hum.ContentTypeTextHtmlUtf8)
	}
}

func (svc *Service) HandleHomeNetHTTP(res http.ResponseWriter, req *http.Request) {
	log.Info("HANDLE_NetHTTP")
	svc.HandleHomeAnyRequest(anyhttp.NewResReqNetHttp(res, req))
}

func (svc *Service) HandleHomeFastHTTP(ctx *fasthttp.RequestCtx) {
	log.Info("HANDLE_FastHTTP")
	svc.HandleHomeAnyRequest(anyhttp.NewResReqFastHttp(ctx))
}

func getHttpServeMux(svc Service) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", http.HandlerFunc(svc.HandleHomeNetHTTP))
	mux.HandleFunc("/hook", http.HandlerFunc(svc.HandleNetHTTP))
	mux.HandleFunc("/hook/", http.HandlerFunc(svc.HandleNetHTTP))
	return mux
}

func serveNetHttp(svc Service) {
	log.Info(fmt.Sprintf("STARTING_NET_HTTP [%v]", svc.Config.Port))
	http.ListenAndServe(portAddress(svc.Config.Port), getHttpServeMux(svc))
}

func serveFastHttp(svc Service) {
	log.Info(fmt.Sprintf("STARTING_FAST_HTTP [%v]", svc.Config.Port))
	router := fasthttprouter.New()
	router.GET("/", svc.HandleHomeFastHTTP)
	router.POST("/hook", svc.HandleFastHTTP)
	router.POST("/hook/", svc.HandleFastHTTP)
	log.Fatal(fasthttp.ListenAndServe(portAddress(svc.Config.Port), router.Handler))
}

func serveAwsLambda(svc Service) {
	log.Fatal(gateway.ListenAndServe(portAddress(svc.Config.Port), getHttpServeMux(svc)))
}

func serveAwsLambdaSimple(svc Service) { lambda.Start(svc.HandleAwsLambda) }
func portAddress(port int) string      { return ":" + strconv.Itoa(port) }

func main() {
	if err := cfg.LoadDotEnvSkipEmpty(os.Getenv("ENV_PATH"), "./.env"); err != nil {
		panic(err)
	}

	svc := NewService()

	engine := svc.Config.Engine

	switch engine {
	case "awslambda":
		serveAwsLambda(svc)
	case "nethttp":
		serveNetHttp(svc)
	case "fasthttp":
		serveFastHttp(svc)
	default:
		log.Fatal(fmt.Sprintf("Engine Not Supported [%v]", engine))
	}
}
