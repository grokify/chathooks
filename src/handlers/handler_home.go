package handlers

import (
	"fmt"
	//"strings"

	//"github.com/eawsy/aws-lambda-go-event/service/lambda/runtime/event/apigatewayproxyevt"
	//"github.com/grokify/commonchat"
	"github.com/grokify/chathooks/src/adapters"
	"github.com/grokify/chathooks/src/config"
	//cc "github.com/grokify/commonchat"
	//"github.com/grokify/gotilla/type/stringsutil"
	"github.com/valyala/fasthttp"
)

const (
	QueryParamNamedOutputs = "adapters"
	QueryParamInputType    = "inputType"
	QueryParamOutputType   = "outputType"
	QueryParamToken        = "token"
	QueryParamOutputURL    = "url"
)

var (
	ShowDisplayName = false
)

// HomeHandler is a fasthttp handler for handling the webhoo proxy homepage.
func HomeHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "%s", []byte("Chathooks\nSource: https://github.com/grokify/chathooks"))
}

type Configuration struct {
	ConfigData config.Configuration
	AdapterSet adapters.AdapterSet
}

func IntegrationActivitySuffix(displayName string) string {
	if !ShowDisplayName || len(displayName) < 1 {
		return ""
	}
	return ""
}

/*
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
*/
