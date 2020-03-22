package models

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/eawsy/aws-lambda-go-event/service/lambda/runtime/event/apigatewayproxyevt"
	cc "github.com/grokify/commonchat"
	"github.com/grokify/gotilla/net/anyhttp"
	fhu "github.com/grokify/gotilla/net/fasthttputil"
	nhu "github.com/grokify/gotilla/net/nethttputil"
	"github.com/grokify/gotilla/type/stringsutil"
	"github.com/valyala/fasthttp"
)

const (
	QueryParamOutputAdapters = "adapters"
	QueryParamInputType      = "inputType"
	QueryParamOutputType     = "outputType"
	QueryParamToken          = "token"
	QueryParamOutputURL      = "url"
)

var FixedParams = map[string]int{
	QueryParamOutputAdapters: 1,
	QueryParamInputType:      2,
	QueryParamOutputType:     3,
	QueryParamToken:          4,
	QueryParamOutputURL:      5}

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
	URLEncodedRails
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
	CustomParams     url.Values `json:"customParams,omitempty"`
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
		QueryStringParameters: awsReq.QueryStringParameters})
}

func HookDataFromEawsyLambdaEvent(bodyType MessageBodyType, eawsyReq *apigatewayproxyevt.Event) HookData {
	return newHookDataGeneric(hookDataRequest{
		BodyType:              bodyType,
		Headers:               eawsyReq.Headers,
		Body:                  eawsyReq.Body,
		IsBase64Encoded:       eawsyReq.IsBase64Encoded,
		QueryStringParameters: eawsyReq.QueryStringParameters})
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
	data := HookData{
		CustomParams: url.Values{}}
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
		data.OutputNames = stringsutil.SliceCondenseSpace(strings.Split(namedOutputs, ","), true, false)
	}
	// Include any other parameter as a custom param.
	for key, val := range queryStringParameters {
		if _, ok := FixedParams[key]; !ok {
			data.CustomParams.Add(strings.ToLower(strings.TrimSpace(key)), val)
			//data.CustomParams[strings.ToLower(strings.TrimSpace(key))] = val
		}
	}
	return data
}

func HookDataFromAnyHTTPReq(bodyType MessageBodyType, aReq anyhttp.Request) HookData {
	return HookData{
		InputType:    aReq.QueryArgs().GetString(QueryParamInputType),
		InputBody:    BodyToMessageBytesAnyHTTP(bodyType, aReq),
		OutputType:   aReq.QueryArgs().GetString(QueryParamOutputType),
		OutputURL:    aReq.QueryArgs().GetString(QueryParamOutputURL),
		Token:        aReq.QueryArgs().GetString(QueryParamToken),
		CustomParams: aReq.QueryArgs().GetURLValues(),
		OutputNames:  strings.Split(aReq.QueryArgs().GetString(QueryParamOutputAdapters), ",")}
}

func HookDataFromNetHTTPReq(bodyType MessageBodyType, req *http.Request) HookData {
	return HookData{
		InputType:   nhu.GetReqQueryParam(req, QueryParamInputType),
		InputBody:   BodyToMessageBytesNetHTTP(bodyType, req),
		OutputType:  nhu.GetReqQueryParam(req, QueryParamOutputType),
		OutputURL:   nhu.GetReqQueryParam(req, QueryParamOutputURL),
		Token:       nhu.GetReqQueryParam(req, QueryParamToken),
		OutputNames: nhu.GetSplitReqQueryParam(req, QueryParamOutputAdapters, ",")}
}

func HookDataFromFastHTTPReqCtx(bodyType MessageBodyType, ctx *fasthttp.RequestCtx) HookData {
	return HookData{
		InputType:   fhu.GetReqQueryParam(ctx, QueryParamInputType),
		InputBody:   BodyToMessageBytesFastHTTP(bodyType, ctx),
		OutputType:  fhu.GetReqQueryParam(ctx, QueryParamOutputType),
		OutputURL:   fhu.GetReqQueryParam(ctx, QueryParamOutputURL),
		Token:       fhu.GetReqQueryParam(ctx, QueryParamToken),
		OutputNames: fhu.GetSplitReqQueryParam(ctx, QueryParamOutputAdapters, ",'")}
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

func BodyToMessageBytesAnyHTTP(bodyType MessageBodyType, aReq anyhttp.Request) []byte {
	switch bodyType {
	case URLEncodedJSONPayload:
		if err := aReq.ParseForm(); err != nil {
			return []byte("")
		}
		return aReq.PostArgs().GetBytes("payload")
	case URLEncodedJSONPayloadOrJSON:
		ct := strings.TrimSpace(strings.ToLower(aReq.HeaderString("Content-Type")))
		if strings.Index(ct, `application/json`) > -1 {
			bytes, err := aReq.PostBody()
			if err != nil {
				return []byte("")
			}
			return bytes
		}
		if err := aReq.ParseForm(); err != nil {
			return []byte("")
		}
		return aReq.PostArgs().GetBytes("payload")
		//return []byte(req.Form.Get("payload"))
	default:
		bytes, err := aReq.PostBody()
		if err != nil {
			return []byte("")
		}
		return bytes
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

func BuildAwsAPIGatewayProxyResponse(hookData HookData, errs ...ErrorInfo) (events.APIGatewayProxyResponse, error) {
	resInfo := ResponseInfo{
		HookData:   hookData,
		Responses:  errs,
		StatusCode: GetMaxStatusCode(errs...)}
	return resInfo.ToAPIGatewayProxyResponse()
}
