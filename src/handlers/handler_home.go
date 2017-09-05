package handlers

import (
	"fmt"
	"strings"

	"github.com/eawsy/aws-lambda-go-event/service/lambda/runtime/event/apigatewayproxyevt"
	"github.com/grokify/gotilla/strings/stringsutil"
	"github.com/grokify/webhookproxy/src/adapters"
	"github.com/grokify/webhookproxy/src/config"
	"github.com/valyala/fasthttp"
)

const (
	QueryParamNamedOutputs = "adapters"
	QueryParamInputType    = "inputType"
	QueryParamOutputType   = "outputType"
	QueryParamToken        = "token"
	QueryParamOutputURL    = "url"
)

// HomeHandler is a fasthttp handler for handling the webhoo proxy homepage.
func HomeHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "%s", []byte("Webhook Proxy\nSource: https://github.com/grokify/webhookproxy"))
}

type Configuration struct {
	ConfigData config.Configuration
	AdapterSet adapters.AdapterSet
}

type HookRequestData struct {
	InputType    string
	InputBody    []byte
	OutputType   string
	OutputURL    string
	NamedOutputs []string
	Token        string
}

func HookRequestDataFromEawsyLambdaEvent(event *apigatewayproxyevt.Event) HookRequestData {
	data := HookRequestData{InputBody: []byte(event.Body)}

	if input, ok := event.QueryStringParameters[QueryParamInputType]; ok {
		data.InputType = strings.TrimSpace(input)
	}
	if output, ok := event.QueryStringParameters[QueryParamOutputType]; ok {
		data.OutputType = strings.TrimSpace(output)
	}
	if url, ok := event.QueryStringParameters[QueryParamOutputURL]; ok {
		data.OutputURL = strings.TrimSpace(url)
	}
	if token, ok := event.QueryStringParameters[QueryParamToken]; ok {
		data.Token = strings.TrimSpace(token)
	}
	if namedOutputs, ok := event.QueryStringParameters[QueryParamNamedOutputs]; ok {
		data.NamedOutputs = stringsutil.SliceTrimSpace(strings.Split(namedOutputs, ","))
	}
	return data
}

func HookRequestDataFromFastHTTPReqCtx(ctx *fasthttp.RequestCtx) HookRequestData {
	return HookRequestData{
		InputType:  strings.TrimSpace(string(ctx.QueryArgs().Peek(QueryParamInputType))),
		InputBody:  ctx.PostBody(),
		OutputType: strings.TrimSpace(string(ctx.QueryArgs().Peek(QueryParamOutputType))),
		OutputURL:  strings.TrimSpace(string(ctx.QueryArgs().Peek(QueryParamOutputURL))),
		Token:      strings.TrimSpace(string(ctx.QueryArgs().Peek(QueryParamToken))),
		NamedOutputs: stringsutil.SliceTrimSpace(strings.Split(
			string(ctx.QueryArgs().Peek(QueryParamNamedOutputs)), ",")),
	}
}
