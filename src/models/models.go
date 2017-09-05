package models

import (
	"encoding/json"
	"net/url"
	"strings"

	cc "github.com/commonchat/commonchat-go"
	"github.com/eawsy/aws-lambda-go-event/service/lambda/runtime/event/apigatewayproxyevt"
	"github.com/grokify/gotilla/strings/stringsutil"
	"github.com/valyala/fasthttp"
)

const (
	QueryParamOutputAdapters = "adapters"
	QueryParamInputType      = "inputType"
	QueryParamOutputType     = "outputType"
	QueryParamToken          = "token"
	QueryParamOutputURL      = "url"
)

type MessageBodyType int

const (
	JSON MessageBodyType = iota
	URLEncoded
	URLEncodedJSONPayload
	URLEncodedJSONPayloadOrJSON
)

var intervals = [...]string{
	"json",
	"url_encoded",
	"url_encoded_json_payload",
	"url_encoded_or_json",
}

type HookData struct {
	InputType     string
	InputBody     []byte
	OutputType    string
	OutputURL     string
	OutputNames   []string
	Token         string
	InputMessage  []byte
	OutputMessage cc.Message
}

func HookDataFromEawsyLambdaEvent(bodyType MessageBodyType, event *apigatewayproxyevt.Event) HookData {
	data := HookData{
		InputBody: BodyToMessageBytesEawsyLambda(bodyType, event)}

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
	if namedOutputs, ok := event.QueryStringParameters[QueryParamOutputAdapters]; ok {
		data.OutputNames = stringsutil.SliceTrimSpace(strings.Split(namedOutputs, ","))
	}
	return data
}

func HookDataFromFastHTTPReqCtx(bodyType MessageBodyType, ctx *fasthttp.RequestCtx) HookData {
	return HookData{
		InputType:  strings.TrimSpace(string(ctx.QueryArgs().Peek(QueryParamInputType))),
		InputBody:  BodyToMessageBytesFastHTTP(bodyType, ctx),
		OutputType: strings.TrimSpace(string(ctx.QueryArgs().Peek(QueryParamOutputType))),
		OutputURL:  strings.TrimSpace(string(ctx.QueryArgs().Peek(QueryParamOutputURL))),
		Token:      strings.TrimSpace(string(ctx.QueryArgs().Peek(QueryParamToken))),
		OutputNames: stringsutil.SliceTrimSpace(strings.Split(
			string(ctx.QueryArgs().Peek(QueryParamOutputAdapters)), ",")),
	}
}

func BodyToMessageBytesEawsyLambda(bodyType MessageBodyType, event *apigatewayproxyevt.Event) []byte {
	switch bodyType {
	case URLEncodedJSONPayload:
		v, err := url.ParseQuery(event.Body)
		if err != nil {
			return []byte("")
		}
		return []byte(v.Get("payload"))
	case URLEncodedJSONPayloadOrJSON:
		if ct, ok := event.Headers["content-type"]; ok {
			ct = strings.TrimSpace(strings.ToLower(ct))
			if strings.Index(ct, `application/json`) > -1 {
				return []byte(event.Body)
			}
		}
		v, err := url.ParseQuery(event.Body)
		if err != nil {
			return []byte("")
		}
		return []byte(v.Get("payload"))
	default:
		return []byte(event.Body)
	}
}

func BodyToMessageBytesFastHTTP(bodyType MessageBodyType, ctx *fasthttp.RequestCtx) []byte {
	switch bodyType {
	case URLEncodedJSONPayload:
		return ctx.FormValue("payload")
	case URLEncodedJSONPayloadOrJSON:
		ct := strings.TrimSpace(
			strings.ToLower(
				string(ctx.Request.Header.Peek("Content-Type"))))
		if strings.Index(ct, `application/json`) > -1 {
			return ctx.PostBody()
		}
		return ctx.FormValue("payload")
	default:
		return ctx.PostBody()
	}
}

type AwsAPIGatewayProxyOutput struct {
	IsBase64Encoded bool              `json:"isBase64Encoded"`
	StatusCode      int               `json:"statusCode"`
	Body            string            `json:"body"`
	Headers         map[string]string `json:"headers"`
}

type ErrorInfo struct {
	StatusCode int
	Body       []byte
}

func ErrorInfosToAlexaResponse(errs ...ErrorInfo) AwsAPIGatewayProxyOutput {
	proxyOutput := AwsAPIGatewayProxyOutput{}
	if len(errs) == 0 {
		proxyOutput.StatusCode = 200
	} else {
		bodyBytes, err := json.Marshal(errs)
		if err != nil {
			proxyOutput.Body = err.Error()
		} else {
			proxyOutput.Body = string(bodyBytes)
		}
		if len(errs) == 1 {
			proxyOutput.StatusCode = errs[0].StatusCode
		} else {
			maxStatus := 0
			for _, errInfo := range errs {
				if errInfo.StatusCode > maxStatus {
					maxStatus = errInfo.StatusCode
				}
			}
			proxyOutput.StatusCode = maxStatus
		}
	}
	return proxyOutput
}
