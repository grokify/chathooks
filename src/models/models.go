package models

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/grokify/gotilla/fmt/fmtutil"

	"github.com/aws/aws-lambda-go/events"
	cc "github.com/commonchat/commonchat-go"
	"github.com/eawsy/aws-lambda-go-event/service/lambda/runtime/event/apigatewayproxyevt"
	fhu "github.com/grokify/gotilla/net/fasthttputil"
	nhu "github.com/grokify/gotilla/net/nethttputil"
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

type RequestParams struct {
	InputType  string `url:"inputType"`
	OutputType string `url:"outputType"`
	Token      string `url:"token"`
	URL        string `url:"url"`
}

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
	InputType        string     `json:"inputType,omitempty"`
	InputBody        []byte     `json:"inputBody,omitempty"`
	OutputType       string     `json:"outputType,omitempty"`
	OutputURL        string     `json:"outputUrl,omitempty"`
	OutputNames      []string   `json:"outputNames,omitempty"`
	Token            string     `json:"token,omitempty"`
	InputMessage     []byte     `json:"inputMessage,omitempty"`
	CanonicalMessage cc.Message `json:"canonicalMessage,omitempty"`
}

type hookDataRequest struct {
	BodyType              MessageBodyType
	Headers               map[string]string
	QueryStringParameters map[string]string
	Body                  string
	IsBase64Encoded       bool
}

func HookDataFromAwsLambdaEvent(bodyType MessageBodyType, awsReq events.APIGatewayProxyRequest) HookData {
	return newHookDataGeneric(hookDataRequest{
		BodyType:              bodyType,
		Headers:               awsReq.Headers,
		Body:                  awsReq.Body,
		IsBase64Encoded:       awsReq.IsBase64Encoded,
		QueryStringParameters: awsReq.QueryStringParameters,
	})
}

func HookDataFromEawsyLambdaEvent(bodyType MessageBodyType, eawsyReq *apigatewayproxyevt.Event) HookData {
	return newHookDataGeneric(hookDataRequest{
		BodyType:              bodyType,
		Headers:               eawsyReq.Headers,
		Body:                  eawsyReq.Body,
		IsBase64Encoded:       eawsyReq.IsBase64Encoded,
		QueryStringParameters: eawsyReq.QueryStringParameters,
	})
}

func newHookDataGeneric(req hookDataRequest) HookData {
	data := newHookDataForQueryString(req.QueryStringParameters)
	data.InputBody = bodyToMessageBytesGeneric(
		req.BodyType,
		req.Headers,
		req.Body,
		req.IsBase64Encoded)
	return data
}

func GetMapString2Simple(mapSS map[string]string, key string) string {
	if value, ok := mapSS[key]; ok {
		return value
	}
	return ""
}

func newHookDataForQueryString(queryStringParameters map[string]string) HookData {
	data := HookData{}
	if input, ok := queryStringParameters[QueryParamInputType]; ok {
		data.InputType = strings.TrimSpace(input)
	}
	if output, ok := queryStringParameters[QueryParamOutputType]; ok {
		data.OutputType = strings.TrimSpace(output)
	}
	if url, ok := queryStringParameters[QueryParamOutputURL]; ok {
		data.OutputURL = strings.TrimSpace(url)
	}
	if token, ok := queryStringParameters[QueryParamToken]; ok {
		data.Token = strings.TrimSpace(token)
	}
	if namedOutputs, ok := queryStringParameters[QueryParamOutputAdapters]; ok {
		data.OutputNames = stringsutil.SliceTrimSpace(strings.Split(namedOutputs, ","))
	}
	return data
}

func HookDataFromNetHTTPReq(bodyType MessageBodyType, req *http.Request) HookData {
	return HookData{
		InputType:   nhu.GetReqHeader(req, QueryParamInputType),
		InputBody:   BodyToMessageBytesNetHTTP(bodyType, req),
		OutputType:  nhu.GetReqHeader(req, QueryParamOutputType),
		OutputURL:   nhu.GetReqHeader(req, QueryParamOutputURL),
		Token:       nhu.GetReqHeader(req, QueryParamToken),
		OutputNames: nhu.GetSplitReqHeader(req, QueryParamOutputAdapters, ","),
	}
}

func HookDataFromFastHTTPReqCtx(bodyType MessageBodyType, ctx *fasthttp.RequestCtx) HookData {
	return HookData{
		InputType:   fhu.GetReqHeader(ctx, QueryParamInputType),
		InputBody:   BodyToMessageBytesFastHTTP(bodyType, ctx),
		OutputType:  fhu.GetReqHeader(ctx, QueryParamOutputType),
		OutputURL:   fhu.GetReqHeader(ctx, QueryParamOutputURL),
		Token:       fhu.GetReqHeader(ctx, QueryParamToken),
		OutputNames: fhu.GetSplitReqHeader(ctx, QueryParamOutputAdapters, ",'"),
	}
}

func bodyToMessageBytesGeneric(bodyType MessageBodyType, headers map[string]string, body string, isBase64Encoded bool) []byte {
	if isBase64Encoded {
		decoded, err := base64.StdEncoding.DecodeString(body)
		if err != nil {
			return []byte("")
		}
		body = string(decoded)
	}
	switch bodyType {
	case URLEncodedJSONPayload:
		v, err := url.ParseQuery(body)
		if err != nil {
			return []byte("")
		}
		return []byte(v.Get("payload"))
	case URLEncodedJSONPayloadOrJSON:
		if ct, ok := headers["content-type"]; ok {
			ct = strings.TrimSpace(strings.ToLower(ct))
			if strings.Index(ct, `application/json`) > -1 {
				return []byte(body)
			}
		}
		v, err := url.ParseQuery(body)
		if err != nil {
			return []byte("")
		}
		return []byte(v.Get("payload"))
	default:
		return []byte(body)
	}
}

func BodyToMessageBytesNetHTTP(bodyType MessageBodyType, req *http.Request) []byte {
	switch bodyType {
	case URLEncodedJSONPayload:
		return []byte(req.Form.Get("payload"))
	case URLEncodedJSONPayloadOrJSON:
		ct := strings.TrimSpace(strings.ToLower(req.Header.Get("Content-Type")))
		if strings.Index(ct, `application/json`) > -1 {
			bytes, err := ioutil.ReadAll(req.Body)
			if err != nil {
				return []byte("")
			}
			return bytes
		}
		return []byte(req.Form.Get("payload"))
	default:
		bytes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return []byte("")
		}
		return bytes
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

type ResponseInfo struct {
	HookData   HookData    `json:"hookData,omitempty"`
	Responses  []ErrorInfo `json:"responses,omitempty"`
	StatusCode int         `json:"statusCode,omitempty"`
	//URL        string      `json:"url,omitempty"`
	//Body       interface{} `json:"body,omitempty"`
	//InputType  string      `json:"inputType,omitempty"`
	//OutputType string      `json:"outputType,omitempty"`
}

func (ri *ResponseInfo) ToAPIGatewayProxyResponse() (events.APIGatewayProxyResponse, error) {
	res := events.APIGatewayProxyResponse{
		StatusCode: ri.StatusCode,
	}

	bodyBytes, err := json.Marshal(ri)
	if err != nil {
		return res, nil
	}
	res.Body = string(bodyBytes)

	return res, nil
}

/*
func ErrorsInfoToResponseInfo(errs ...ErrorInfo) ErrorInfo {
	resInfo := ResponseInfo{
		Responses: errs,
	}
	return resInfo
}
*/

func GetMaxStatusCode(errs ...ErrorInfo) int {
	if len(errs) == 0 {
		return http.StatusOK
	} else if len(errs) == 1 {
		return errs[0].StatusCode
	}
	maxStatus := 200
	for _, errInfo := range errs {
		if errInfo.StatusCode > maxStatus {
			maxStatus = errInfo.StatusCode
		}
	}
	return maxStatus
}

func ErrorsInfoToResponseInfoOld(errs ...ErrorInfo) ErrorInfo {
	resInfo := ErrorInfo{}
	if len(errs) == 0 {
		resInfo.StatusCode = http.StatusOK
		return resInfo
	}
	bodyBytes, err := json.Marshal(errs)
	if err != nil {
		resInfo.Body = []byte(err.Error())
	} else {
		resInfo.Body = bodyBytes
	}
	if len(errs) == 1 {
		resInfo.StatusCode = errs[0].StatusCode
	} else {
		maxStatus := 0
		for _, errInfo := range errs {
			if errInfo.StatusCode > maxStatus {
				maxStatus = errInfo.StatusCode
			}
		}
		resInfo.StatusCode = maxStatus
	}
	return resInfo
}

/*
//func ErrorInfosToAlexaResponse(errs ...ErrorInfo) AwsAPIGatewayProxyOutput {
func ErrorInfosToAwsAPIGatewayProxyOutput(errs ...ErrorInfo) AwsAPIGatewayProxyOutput {
	resInfo := ErrorsInfoToResponseInfo()

	return AwsAPIGatewayProxyOutput{
		StatusCode: resInfo.StatusCode,
		Body:       string(resInfo.Body),
	}
}
*/

func BuildAwsAPIGatewayProxyResponse(hookData HookData, errs ...ErrorInfo) (events.APIGatewayProxyResponse, error) {
	resInfo := ResponseInfo{
		HookData:   hookData,
		Responses:  errs,
		StatusCode: GetMaxStatusCode(errs...),
	}
	fmtutil.PrintJSON(resInfo)
	return resInfo.ToAPIGatewayProxyResponse()
	//resInfo := ErrorsInfoToResponseInfo()
	/*
		return events.APIGatewayProxyResponse{
			StatusCode: resInfo.StatusCode,
			Body:       string(resInfo.Body),
		}
	*/
}

/*
func ErrorInfosToAwsAPIGatewayProxyResponse(errs ...ErrorInfo) events.APIGatewayProxyResponse {
	resInfo := ErrorsInfoToResponseInfo()

	return events.APIGatewayProxyResponse{
		StatusCode: resInfo.StatusCode,
		Body:       string(resInfo.Body),
	}
}
*/
