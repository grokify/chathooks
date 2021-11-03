package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	cc "github.com/grokify/commonchat"
	"github.com/grokify/simplego/net/anyhttp"
	"github.com/grokify/simplego/net/urlutil"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"

	"github.com/grokify/chathooks/pkg/adapters"
	"github.com/grokify/chathooks/pkg/config"
	"github.com/grokify/chathooks/pkg/models"
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
	Env         map[string]string // handler environment
	QueryParams url.Values        // query string params
	Body        []byte            // message, e.g. request body
}

func NewHandlerRequest() HandlerRequest {
	return HandlerRequest{
		Env:         map[string]string{},
		QueryParams: url.Values{},
		Body:        []byte("")}
}

type Normalize func(config.Configuration, HandlerRequest) (cc.Message, error)

// HandleAwsLambda is the method to respond to a fasthttp request.
func (h Handler) HandleAwsLambda(ctx context.Context, awsReq events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	hookData := models.HookDataFromAwsLambdaEvent(h.MessageBodyType, awsReq, h.MessageBodyType)
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
		log.Info().
			Err(err).
			Str("event", "outgoing.webhook.error").
			Msg("ERROR")
	} else {
		if bytes, err := json.Marshal(awsRes.Body); err != nil {
			aRes.SetStatusCode(http.StatusInternalServerError)
		} else {
			_, err := aRes.SetBodyBytes(bytes)
			if err != nil {
				log.Info().
					Err(err).
					Str("event", "outgoing.webhook.error").
					Msg("ERROR")
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
		log.Info().
			Err(err).
			Str("event", "outgoing.webhook.error").
			Msg("ERROR")
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
		log.Info().
			Err(err).
			Str("event", "outgoing.webhook.error").
			Msg("ERROR")
	} else {
		ctx.SetStatusCode(awsRes.StatusCode)
		fmt.Fprint(ctx, awsRes.Body)
	}
}

// HandleCanonical is the method to handle a processed request.
func (h Handler) HandleCanonical(hookData models.HookData) []models.ErrorInfo {
	log.Debug().
		Str("event", "incoming.webhook").
		Str("handler", DisplayName).
		Str("input_body", string(hookData.InputBody)).
		Msg("HANDLE_CANONICAL")

	ccMsg, err := h.Normalize(h.Config,
		HandlerRequest{
			QueryParams: hookData.CustomQueryParams,
			Body:        hookData.InputBody})
	activityUrl := strings.TrimSpace(hookData.CustomQueryParams.Get("activity"))
	if len(activityUrl) > 0 {
		ccMsg.Activity = activityUrl
	}
	iconQry := strings.TrimSpace(hookData.CustomQueryParams.Get("icon"))
	if len(iconQry) > 0 {
		if urlutil.IsHttp(iconQry, true, true) {
			ccMsg.IconURL = iconQry
		} else {
			ccMsg.IconEmoji = iconQry
		}
	}

	if err != nil {
		log.Info().
			Err(err).
			Str("type", "http.response").
			Int("http_status", fasthttp.StatusNotAcceptable).
			Str("handler", DisplayName).
			Msg("request conversion failed")

		return []models.ErrorInfo{{StatusCode: 500, Body: []byte(err.Error())}}
	}
	hookData.CanonicalMessage = ccMsg
	return h.AdapterSet.SendWebhooks(hookData)
}
