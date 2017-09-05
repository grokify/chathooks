package handlers

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	cc "github.com/commonchat/commonchat-go"
	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
	"github.com/eawsy/aws-lambda-go-event/service/lambda/runtime/event/apigatewayproxyevt"
	"github.com/grokify/webhookproxy/src/adapters"
	"github.com/grokify/webhookproxy/src/config"
	"github.com/grokify/webhookproxy/src/models"
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

// HandleEawsyLambda is the method to respond to a fasthttp request.
func (h Handler) HandleEawsyLambda(event *apigatewayproxyevt.Event, ctx *runtime.Context) (models.AwsAPIGatewayProxyOutput, error) {
	hookData := models.HookDataFromEawsyLambdaEvent(h.MessageBodyType, event)
	errs := h.HandleCanonical(hookData)
	return models.ErrorInfosToAlexaResponse(errs...), nil
}

// HandleFastHTTP is the method to respond to a fasthttp request.
func (h Handler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	hookData := models.HookDataFromFastHTTPReqCtx(h.MessageBodyType, ctx)
	errs := h.HandleCanonical(hookData)

	proxyOutput := models.ErrorInfosToAlexaResponse(errs...)
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
		return []models.ErrorInfo{models.ErrorInfo{StatusCode: 500, Body: []byte(err.Error())}}
	}
	hookData.OutputMessage = ccMsg
	return h.AdapterSet.SendWebhooks(hookData)
}
