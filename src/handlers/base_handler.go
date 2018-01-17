package handlers

import (
	"context"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-lambda-go/events"
	cc "github.com/commonchat/commonchat-go"
	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
	"github.com/eawsy/aws-lambda-go-event/service/lambda/runtime/event/apigatewayproxyevt"
	"github.com/grokify/chathooks/src/adapters"
	"github.com/grokify/chathooks/src/config"
	"github.com/grokify/chathooks/src/models"
	"github.com/valyala/fasthttp"
)

const (
	DisplayName = "base_handler"
)

type Handler struct {
	Config          config.Configuration
	AdapterSet      adapters.AdapterSet
	Normalize       Normalize
	MessageBodyType models.MessageBodyType
}

type Normalize func(config.Configuration, []byte) (cc.Message, error)

// HandleAwsLambda is the method to respond to a fasthttp request.
func (h Handler) HandleAwsLambda(ctx context.Context, awsReq events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	hookData := models.HookDataFromAwsLambdaEvent(h.MessageBodyType, awsReq)
	errs := h.HandleCanonical(hookData)
	return models.ErrorInfosToAwsAPIGatewayProxyResponse(errs...), nil
}

// HandleEawsyLambda is the method to respond to a fasthttp request.
func (h Handler) HandleEawsyLambda(event *apigatewayproxyevt.Event, ctx *runtime.Context) (events.APIGatewayProxyResponse, error) {
	hookData := models.HookDataFromEawsyLambdaEvent(h.MessageBodyType, event)
	errs := h.HandleCanonical(hookData)
	return models.ErrorInfosToAwsAPIGatewayProxyResponse(errs...), nil
}

// HandleNetHTTP is the method to respond to a fasthttp request.
func (h Handler) HandleNetHTTP(res http.ResponseWriter, req *http.Request) {
	hookData := models.HookDataFromNetHTTPReq(h.MessageBodyType, req)
	errs := h.HandleCanonical(hookData)

	resInfo := models.ErrorsInfoToResponseInfo(errs...)
	res.WriteHeader(resInfo.StatusCode)
	res.Write(resInfo.Body)
}

// HandleFastHTTP is the method to respond to a fasthttp request.
func (h Handler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	hookData := models.HookDataFromFastHTTPReqCtx(h.MessageBodyType, ctx)
	errs := h.HandleCanonical(hookData)

	proxyOutput := models.ErrorInfosToAwsAPIGatewayProxyOutput(errs...)
	ctx.SetStatusCode(proxyOutput.StatusCode)
	if proxyOutput.StatusCode > 399 {
		fmt.Fprintf(ctx, "%s", proxyOutput.Body)
	}
}

// HandleCanonical is the method to handle a processed request.
func (h Handler) HandleCanonical(hookData models.HookData) []models.ErrorInfo {
	log.WithFields(log.Fields{
		"event":   "incoming.webhook",
		"handler": DisplayName}).Info("HANDLE_FASTHTTP")
	log.WithFields(log.Fields{
		"event":   "incoming.webhook",
		"handler": DisplayName}).Info(string(hookData.InputBody))

	ccMsg, err := h.Normalize(h.Config, hookData.InputBody)

	if err != nil {
		//ctx.SetStatusCode(fasthttp.StatusNotAcceptable)
		log.WithFields(log.Fields{
			"type":         "http.response",
			"status":       fasthttp.StatusNotAcceptable,
			"errorMessage": err.Error(),
		}).Info(fmt.Sprintf("%v request conversion failed.", DisplayName))
		return []models.ErrorInfo{{StatusCode: 500, Body: []byte(err.Error())}}
	}
	hookData.OutputMessage = ccMsg
	return h.AdapterSet.SendWebhooks(hookData)
}
