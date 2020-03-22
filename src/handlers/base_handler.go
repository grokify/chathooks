package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
	"github.com/eawsy/aws-lambda-go-event/service/lambda/runtime/event/apigatewayproxyevt"
	"github.com/grokify/chathooks/src/adapters"
	"github.com/grokify/chathooks/src/config"
	"github.com/grokify/chathooks/src/models"
	cc "github.com/grokify/commonchat"
	"github.com/grokify/gotilla/net/anyhttp"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

const (
	DisplayName = "base_handler"
)

type Handler struct {
	Config          config.Configuration
	AdapterSet      adapters.AdapterSet
	Key             string
	Normalize       Normalize
	MessageBodyType models.MessageBodyType
}

type HandlerRequest struct {
	Env    map[string]string // handler environment
	Params map[string]string // query string params
	Body   []byte            // message, e.g. request body
}

func NewHandlerRequest() HandlerRequest {
	return HandlerRequest{
		Env:    map[string]string{},
		Params: map[string]string{},
		Body:   []byte("")}
}

type Normalize func(config.Configuration, HandlerRequest) (cc.Message, error)

//type Normalize func(config.Configuration, []byte) (cc.Message, error)

// HandleAwsLambda is the method to respond to a fasthttp request.
func (h Handler) HandleAwsLambda(ctx context.Context, awsReq events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	hookData := models.HookDataFromAwsLambdaEvent(h.MessageBodyType, awsReq)
	errs := h.HandleCanonical(hookData)
	awsRes, err := models.BuildAwsAPIGatewayProxyResponse(hookData, errs...)
	return awsRes, err
}

// HandleEawsyLambda is the method to respond to a fasthttp request.
func (h Handler) HandleEawsyLambda(event *apigatewayproxyevt.Event, ctx *runtime.Context) (events.APIGatewayProxyResponse, error) {
	hookData := models.HookDataFromEawsyLambdaEvent(h.MessageBodyType, event)
	errs := h.HandleCanonical(hookData)
	awsRes, err := models.BuildAwsAPIGatewayProxyResponse(hookData, errs...)
	return awsRes, err
}

// HandleNetHTTP is the method to respond to a fasthttp request.
func (h Handler) HandleAnyHTTP(aRes anyhttp.Response, aReq anyhttp.Request) {
	hookData := models.HookDataFromAnyHTTPReq(h.MessageBodyType, aReq)
	errs := h.HandleCanonical(hookData)

	awsRes, err := models.BuildAwsAPIGatewayProxyResponse(hookData, errs...)

	if err != nil {
		aRes.SetStatusCode(http.StatusInternalServerError)
		log.WithFields(log.Fields{
			"event":   "outgoing.webhook.error",
			"handler": err.Error()}).Info("ERROR")
	} else {
		if bytes, err := json.Marshal(awsRes.Body); err != nil {
			aRes.SetStatusCode(http.StatusInternalServerError)
		} else {
			_, err := aRes.SetBodyBytes(bytes)
			if err != nil {
				log.WithFields(log.Fields{
					"event":   "outgoing.webhook.error",
					"handler": err.Error()}).Info("ERROR")
				aRes.SetStatusCode(http.StatusInternalServerError)
			} else {
				aRes.SetStatusCode(awsRes.StatusCode)
			}
		}
	}
}

// HandleNetHTTP is the method to respond to a fasthttp request.
func (h Handler) HandleNetHTTP(res http.ResponseWriter, req *http.Request) {
	hookData := models.HookDataFromNetHTTPReq(h.MessageBodyType, req)
	errs := h.HandleCanonical(hookData)

	awsRes, err := models.BuildAwsAPIGatewayProxyResponse(hookData, errs...)

	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.WithFields(log.Fields{
			"event":   "outgoing.webhook.error",
			"handler": err.Error()}).Info("ERROR")
	} else {
		res.WriteHeader(awsRes.StatusCode)
		fmt.Fprint(res, awsRes.Body)
	}
}

// HandleFastHTTP is the method to respond to a fasthttp request.
func (h Handler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	hookData := models.HookDataFromFastHTTPReqCtx(h.MessageBodyType, ctx)
	errs := h.HandleCanonical(hookData)

	awsRes, err := models.BuildAwsAPIGatewayProxyResponse(hookData, errs...)

	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		log.WithFields(log.Fields{
			"event":   "outgoing.webhook.error",
			"handler": err.Error()}).Info("ERROR")
	} else {
		ctx.SetStatusCode(awsRes.StatusCode)
		fmt.Fprint(ctx, awsRes.Body)
	}
}

// HandleCanonical is the method to handle a processed request.
func (h Handler) HandleCanonical(hookData models.HookData) []models.ErrorInfo {
	log.WithFields(log.Fields{
		"event":   "incoming.webhook",
		"handler": DisplayName}).Info("HANDLE_CANONICAL")
	log.WithFields(log.Fields{
		"event":   "incoming.webhook",
		"handler": DisplayName}).Info(string(hookData.InputBody))

	ccMsg, err := h.Normalize(h.Config, HandlerRequest{Body: hookData.InputBody})

	if err != nil {
		log.WithFields(log.Fields{
			"type":         "http.response",
			"status":       fasthttp.StatusNotAcceptable,
			"errorMessage": err.Error(),
		}).Info(fmt.Sprintf("%v request conversion failed.", DisplayName))
		return []models.ErrorInfo{{StatusCode: 500, Body: []byte(err.Error())}}
	}
	hookData.CanonicalMessage = ccMsg
	return h.AdapterSet.SendWebhooks(hookData)
}
