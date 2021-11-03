package service

import (
	"context"
	"fmt"
	clog "log"
	"net/http"
	"strconv"
	"strings"

	"github.com/apex/gateway"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/buaazp/fasthttprouter"
	ccglip "github.com/grokify/commonchat/glip"
	ccslack "github.com/grokify/commonchat/slack"
	"github.com/grokify/simplego/net/anyhttp"
	hum "github.com/grokify/simplego/net/httputilmore"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"

	"github.com/grokify/chathooks/pkg/adapters"
	"github.com/grokify/chathooks/pkg/config"
	"github.com/grokify/chathooks/pkg/models"
	"github.com/grokify/chathooks/pkg/templates"

	"github.com/grokify/chathooks/pkg/handlers"
	"github.com/grokify/chathooks/pkg/handlers/aha"
	"github.com/grokify/chathooks/pkg/handlers/appsignal"
	"github.com/grokify/chathooks/pkg/handlers/apteligent"
	"github.com/grokify/chathooks/pkg/handlers/bugsnag"
	"github.com/grokify/chathooks/pkg/handlers/circleci"
	"github.com/grokify/chathooks/pkg/handlers/codeship"
	"github.com/grokify/chathooks/pkg/handlers/confluence"
	"github.com/grokify/chathooks/pkg/handlers/datadog"
	"github.com/grokify/chathooks/pkg/handlers/deskdotcom"
	"github.com/grokify/chathooks/pkg/handlers/enchant"
	"github.com/grokify/chathooks/pkg/handlers/gosquared"
	"github.com/grokify/chathooks/pkg/handlers/gosquared2"
	"github.com/grokify/chathooks/pkg/handlers/heroku"
	"github.com/grokify/chathooks/pkg/handlers/librato"
	"github.com/grokify/chathooks/pkg/handlers/magnumci"
	"github.com/grokify/chathooks/pkg/handlers/marketo"
	"github.com/grokify/chathooks/pkg/handlers/opsgenie"
	"github.com/grokify/chathooks/pkg/handlers/papertrail"
	"github.com/grokify/chathooks/pkg/handlers/pingdom"
	"github.com/grokify/chathooks/pkg/handlers/raygun"
	"github.com/grokify/chathooks/pkg/handlers/runscope"
	"github.com/grokify/chathooks/pkg/handlers/semaphore"
	"github.com/grokify/chathooks/pkg/handlers/slack"
	"github.com/grokify/chathooks/pkg/handlers/statuspage"
	"github.com/grokify/chathooks/pkg/handlers/travisci"
	"github.com/grokify/chathooks/pkg/handlers/userlike"
	"github.com/grokify/chathooks/pkg/handlers/victorops"
	"github.com/grokify/chathooks/pkg/handlers/wootric"
)

/*

Use the `CHATHOOKS_TOKENS` environment variable to load secret
tokens as a comma delimited string.

*/

// CHATHOOKS_URL=http://localhost:8080/hook CHATHOOKS_HOME_URL=http://localhost:8080 go run main.go

const (
	ParamNameInputType       = "inputType"
	ParamNameOutputType      = "outputType"
	ParamNameURL             = "url"
	ParamNameToken           = "token"
	EnvPath                  = "ENV_PATH"
	EnvEngine                = "CHATHOOKS_ENGINE" // awslambda, nethttp, fasthttp
	EnvTokens                = "CHATHOOKS_TOKENS"
	EnvWebhookUrl            = "CHATHOOKS_URL"
	EnvHomeUrl               = "CHATHOOKS_HOME_URL"
	ErrRequiredTokenNotFound = "401.01 Required Token Not Found"
	ErrRequiredTokenNotValid = "401.02 Required Token Not Valid"
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
	cfgData, err := config.NewConfigurationEnv()
	if err != nil {
		log.Fatal().Err(err)
	}

	adapterSet := adapters.NewAdapterSet()
	glipAdapter, err := ccglip.NewGlipAdapter("", adapters.GlipConfig())
	if err != nil {
		log.Fatal().Err(err)
	}
	adapterSet.Adapters["glip"] = glipAdapter
	slackAdapter, err := ccslack.NewSlackAdapter("")
	if err != nil {
		log.Fatal().Err(err)
	}
	adapterSet.Adapters["slack"] = slackAdapter

	hf := HandlerFactory{Config: cfgData, AdapterSet: adapterSet}

	handlerSet := HandlerSet{Handlers: map[string]Handler{
		"aha":        hf.InflateHandler(aha.NewHandler()),
		"appsignal":  hf.InflateHandler(appsignal.NewHandler()),
		"apteligent": hf.InflateHandler(apteligent.NewHandler()),
		"bugsnag":    hf.InflateHandler(bugsnag.NewHandler()),
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
		"wootric":    hf.InflateHandler(wootric.NewHandler())}}

	svcInfo := Service{
		Config:       cfgData,
		AdapterSet:   adapterSet,
		HandlerSet:   handlerSet,
		RequireToken: false,
		Tokens:       map[string]int{}}

	for _, token := range cfgData.Tokens {
		token = strings.TrimSpace(token)
		if len(token) > 0 {
			svcInfo.Tokens[token] = 1
		}
	}

	return svcInfo
}

func (svc *Service) HandleAwsLambda(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Info().Msg("FUNC_HandleAwsLambda__BEGIN")
	if len(svc.Tokens) > 0 {
		token, ok := req.QueryStringParameters[ParamNameToken]
		if !ok {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusUnauthorized,
				Body:       ErrRequiredTokenNotFound}, nil
		}
		if _, ok := svc.Tokens[token]; !ok {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusUnauthorized,
				Body:       ErrRequiredTokenNotValid}, nil
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

func (svc *Service) HandleAnyRequest(aRes anyhttp.Response, aReq anyhttp.Request) {
	log.Info().Msg("FUNC_HandleAnyRequest__BEGIN")

	if err := aReq.ParseForm(); err != nil {
		aRes.SetStatusCode(http.StatusInternalServerError)
		log.Warn().Msg("E_CANNOT_PARSE_FORM")
		return
	}

	if len(svc.Tokens) > 0 {
		token := strings.TrimSpace(aReq.QueryArgs().GetString(ParamNameToken))

		if len(token) == 0 {
			aRes.SetStatusCode(http.StatusUnauthorized)
			log.Warn().Msg("E_NO_TOKEN")
			return
		}
		if _, ok := svc.Tokens[token]; !ok {
			aRes.SetStatusCode(http.StatusUnauthorized)
			log.Warn().Msg("E_INCORRECT_TOKEN")
			return
		}
	}

	inputType := aReq.QueryArgs().GetString(ParamNameInputType)

	if handler, ok := svc.HandlerSet.Handlers[inputType]; ok {
		log.Info().
			Str("handler_input_type", inputType).
			Msg("Input_Handler_Found_Processing")
		handler.HandleAnyHTTP(aRes, aReq)
	} else {
		fmt.Printf("Input_Handler_Not_Found [%v]\n", inputType)
	}
}

func (svc *Service) HandleHookNetHTTP(res http.ResponseWriter, req *http.Request) {
	log.Info().Msg("FUNC_HandleNetHTTP__BEGIN")
	svc.HandleAnyRequest(anyhttp.NewResReqNetHttp(res, req))
}

func (svc *Service) HandleHookFastHTTP(ctx *fasthttp.RequestCtx) {
	log.Info().Msg("HANDLE_FastHTTP")
	svc.HandleAnyRequest(anyhttp.NewResReqFastHttp(ctx))
}

func (svc *Service) HandleHomeAnyRequest(aRes anyhttp.Response, aReq anyhttp.Request) {
	log.Info().Msg("HANDLE_HOME_AnyHTTP")
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
	log.Debug().Msg("HANDLE_NetHTTP")
	svc.HandleHomeAnyRequest(anyhttp.NewResReqNetHttp(res, req))
}

func (svc *Service) HandleHomeFastHTTP(ctx *fasthttp.RequestCtx) {
	log.Debug().Msg("HANDLE_FastHTTP")
	svc.HandleHomeAnyRequest(anyhttp.NewResReqFastHttp(ctx))
}

func (svc Service) PortInt() int {
	return svc.Config.Port
}

func (svc Service) HttpEngine() string {
	return svc.Config.Engine
}

func (svc Service) Router() http.Handler {
	return getHttpServeMux(svc)
}

func (svc Service) RouterFast() *fasthttprouter.Router {
	router := fasthttprouter.New()
	router.GET("/", svc.HandleHomeFastHTTP)
	router.POST("/hook", svc.HandleHookFastHTTP)
	router.POST("/hook/", svc.HandleHookFastHTTP)
	router.POST("/webhook", svc.HandleHookFastHTTP)
	router.POST("/webhook/", svc.HandleHookFastHTTP)
	return router
}

func getHttpServeMux(svc Service) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", http.HandlerFunc(svc.HandleHomeNetHTTP))
	mux.HandleFunc("/hook", http.HandlerFunc(svc.HandleHookNetHTTP))
	mux.HandleFunc("/hook/", http.HandlerFunc(svc.HandleHookNetHTTP))
	mux.HandleFunc("/webhook", http.HandlerFunc(svc.HandleHookNetHTTP))
	mux.HandleFunc("/webhook/", http.HandlerFunc(svc.HandleHookNetHTTP))
	return mux
}

func ServeNetHttp(svc Service) {
	log.Info().
		Int("port", svc.Config.Port).
		Msg("STARTING_NET_HTTP")
	http.ListenAndServe(portAddress(svc.Config.Port), getHttpServeMux(svc))
}

func ServeFastHttp(svc Service) {
	log.Info().
		Int("port", svc.Config.Port).
		Msg("STARTING_FAST_HTTP")
	router := svc.RouterFast()
	clog.Fatal(fasthttp.ListenAndServe(portAddress(svc.Config.Port), router.Handler))
}

func ServeAwsLambda(svc Service) {
	clog.Fatal(gateway.ListenAndServe(portAddress(svc.Config.Port), getHttpServeMux(svc)))
}

func serveAwsLambdaSimple(svc Service) { lambda.Start(svc.HandleAwsLambda) }
func portAddress(port int) string      { return ":" + strconv.Itoa(port) }
