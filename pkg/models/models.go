package models

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/grokify/commonchat"
	"github.com/grokify/gohttp/anyhttp"
	fhu "github.com/grokify/gohttp/fasthttputil"
	hum "github.com/grokify/mogo/net/http/httputilmore"
	"github.com/grokify/mogo/type/stringsutil"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"

	"github.com/grokify/chathooks/pkg/config"
)

const (
	QueryParamOutputAdapters = config.ParamNameAdapters
	QueryParamInputType      = config.ParamNameInputType
	QueryParamOutputFormat   = config.ParamNameOutputFormat
	QueryParamOutputType     = config.ParamNameOutputType
	QueryParamOutputURL      = config.ParamNameOutputURL
	QueryParamToken          = config.ParamNameToken

	ParamPayload = "payload"
)

var fixedParams = map[string]int{
	QueryParamOutputAdapters: 1,
	QueryParamInputType:      2,
	QueryParamOutputFormat:   3,
	QueryParamOutputType:     4,
	QueryParamOutputURL:      5,
	QueryParamToken:          6,
}

type RequestParams struct {
	InputType  string `url:"inputType"`
	OutputType string `url:"outputType"`
	Token      string `url:"token"`
	URL        string `url:"outputURL"`
}

type MessageBodyType int

const (
	JSON MessageBodyType = iota
	URLEncoded
	URLEncodedJSONPayload
	URLEncodedJSONPayloadOrJSON
	URLEncodedRails
)

/*
var intervals = [...]string{
	"json",
	"url_encoded",
	"url_encoded_json_payload",
	"url_encoded_or_json",
}
*/

type HookData struct {
	InputType         string             `json:"inputType,omitempty"`
	InputBody         []byte             `json:"inputBody,omitempty"`
	OutputFormat      string             `json:"outputFormat,omitempty"`
	OutputType        string             `json:"outputType,omitempty"`
	OutputURL         string             `json:"outputURL,omitempty"`
	OutputNames       []string           `json:"outputNames,omitempty"`
	Token             string             `json:"token,omitempty"`
	InputMessage      []byte             `json:"inputMessage,omitempty"`
	CustomQueryParams url.Values         `json:"customParams,omitempty"`
	CanonicalMessage  commonchat.Message `json:"canonicalMessage,omitempty"`
}

type hookDataRequest struct {
	BodyType              MessageBodyType
	Headers               map[string]string
	QueryStringParameters map[string]string
	Body                  string
	IsBase64Encoded       bool
}

// HookDataFromAwsLambdaEvent converts a Lambda event to
// generic HookData.
func HookDataFromAwsLambdaEvent(bodyType MessageBodyType, awsReq events.APIGatewayProxyRequest, messageBodyType MessageBodyType) HookData {
	hookData := newHookDataGeneric(hookDataRequest{
		BodyType:              bodyType,
		Headers:               awsReq.Headers,
		Body:                  awsReq.Body,
		IsBase64Encoded:       awsReq.IsBase64Encoded,
		QueryStringParameters: awsReq.QueryStringParameters})
	// `application/x-www-form-urlencoded` is currently not supported
	// with AWS Lambda because Lambda cannot support URL Query String
	// parameterss with this Content-Type.
	if messageBodyType == URLEncoded ||
		messageBodyType == URLEncodedJSONPayload ||
		messageBodyType == URLEncodedRails {
		var jsonData awsJSONWrapper
		err := json.Unmarshal(hookData.InputMessage, &jsonData)
		if err == nil {
			hookData.InputMessage = []byte(jsonData.Body)
		}
	}
	return hookData
}

type awsJSONWrapper struct {
	Body string `json:"body,omitempty"`
}

/*
func HookDataFromEawsyLambdaEvent(bodyType MessageBodyType, eawsyReq *apigatewayproxyevt.Event) HookData {
	return newHookDataGeneric(hookDataRequest{
		BodyType:              bodyType,
		Headers:               eawsyReq.Headers,
		Body:                  eawsyReq.Body,
		IsBase64Encoded:       eawsyReq.IsBase64Encoded,
		QueryStringParameters: eawsyReq.QueryStringParameters})
}*/

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
		CustomQueryParams: url.Values{}}
	if input, ok := queryStringParameters[QueryParamInputType]; ok {
		data.InputType = strings.TrimSpace(input)
	}
	if format, ok := queryStringParameters[QueryParamOutputFormat]; ok {
		data.OutputFormat = config.MustParseOutputFormat(format)
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
		if _, ok := fixedParams[key]; !ok {
			data.CustomQueryParams.Add(strings.ToLower(strings.TrimSpace(key)), val)
			//data.CustomParams[strings.ToLower(strings.TrimSpace(key))] = val
		}
	}
	return data
}

func HookDataFromAnyHTTPReq(bodyType MessageBodyType, aReq anyhttp.Request) HookData {
	return HookData{
		InputType:         aReq.QueryArgs().GetString(QueryParamInputType),
		InputBody:         BodyToMessageBytesAnyHTTP(bodyType, aReq),
		OutputFormat:      config.MustParseOutputFormat(aReq.QueryArgs().GetString(QueryParamOutputFormat)),
		OutputType:        aReq.QueryArgs().GetString(QueryParamOutputType),
		OutputURL:         aReq.QueryArgs().GetString(QueryParamOutputURL),
		Token:             aReq.QueryArgs().GetString(QueryParamToken),
		CustomQueryParams: aReq.QueryArgs().GetURLValues(),
		OutputNames:       strings.Split(aReq.QueryArgs().GetString(QueryParamOutputAdapters), ",")}
}

func HookDataFromNetHTTPReq(bodyType MessageBodyType, req *http.Request) HookData {
	return HookData{
		InputType:    hum.GetReqQueryParam(req, QueryParamInputType),
		InputBody:    BodyToMessageBytesNetHTTP(bodyType, req),
		OutputFormat: config.MustParseOutputFormat(hum.GetReqQueryParam(req, QueryParamOutputFormat)),
		OutputType:   hum.GetReqQueryParam(req, QueryParamOutputType),
		OutputURL:    hum.GetReqQueryParam(req, QueryParamOutputURL),
		Token:        hum.GetReqQueryParam(req, QueryParamToken),
		OutputNames:  hum.GetReqQueryParamSplit(req, QueryParamOutputAdapters, ",")}
}

func HookDataFromFastHTTPReqCtx(bodyType MessageBodyType, ctx *fasthttp.RequestCtx) HookData {
	return HookData{
		InputType:    fhu.GetReqQueryParam(ctx, QueryParamInputType),
		InputBody:    BodyToMessageBytesFastHTTP(bodyType, ctx),
		OutputFormat: config.MustParseOutputFormat(fhu.GetReqQueryParam(ctx, QueryParamOutputFormat)),
		OutputType:   fhu.GetReqQueryParam(ctx, QueryParamOutputType),
		OutputURL:    fhu.GetReqQueryParam(ctx, QueryParamOutputURL),
		Token:        fhu.GetReqQueryParam(ctx, QueryParamToken),
		OutputNames:  fhu.GetSplitReqQueryParam(ctx, QueryParamOutputAdapters, ",'")}
}

func bodyToMessageBytesGeneric(bodyType MessageBodyType, headers map[string]string, body string, isBase64Encoded bool) []byte {
	var bodyConverted []byte
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
		bodyConverted = []byte(v.Get(ParamPayload))
	case URLEncodedJSONPayloadOrJSON:
		if ct, ok := headers["content-type"]; ok {
			ct = strings.TrimSpace(strings.ToLower(ct))
			if strings.Contains(ct, hum.ContentTypeAppJSON) {
				return []byte(body)
			}
		}
		v, err := url.ParseQuery(body)
		if err != nil {
			return []byte("")
		}
		bodyConverted = []byte(v.Get(ParamPayload))
	default:
		bodyConverted = []byte(body)
	}
	log.Debug().
		Str("body", string(bodyConverted)).
		Msg("REQUEST_BODY")
	return bodyConverted
}

func BodyToMessageBytesAnyHTTP(bodyType MessageBodyType, aReq anyhttp.Request) []byte {
	switch bodyType {
	case URLEncodedJSONPayload:
		if err := aReq.ParseForm(); err != nil {
			return []byte{}
		}
		return aReq.PostArgs().GetBytes(ParamPayload)
	case URLEncodedJSONPayloadOrJSON:
		ct := strings.TrimSpace(strings.ToLower(aReq.HeaderString(hum.HeaderContentType)))
		if strings.Contains(ct, hum.ContentTypeAppJSON) {
			bytes, err := aReq.PostBody()
			if err != nil {
				return []byte{}
			}
			return bytes
		}
		if err := aReq.ParseForm(); err != nil {
			return []byte{}
		}
		return aReq.PostArgs().GetBytes(ParamPayload)
		//return []byte(req.Form.Get(ParamPayload))
	default:
		bytes, err := aReq.PostBody()
		if err != nil {
			return []byte{}
		}
		return bytes
	}
}

func BodyToMessageBytesNetHTTP(bodyType MessageBodyType, req *http.Request) []byte {
	switch bodyType {
	case URLEncodedJSONPayload:
		return []byte(req.Form.Get(ParamPayload))
	case URLEncodedJSONPayloadOrJSON:
		ct := strings.TrimSpace(strings.ToLower(req.Header.Get(hum.HeaderContentType)))
		if strings.Contains(ct, hum.ContentTypeAppJSON) {
			bytes, err := io.ReadAll(req.Body)
			if err != nil {
				return []byte{}
			}
			return bytes
		}
		return []byte(req.Form.Get(ParamPayload))
	default:
		bytes, err := io.ReadAll(req.Body)
		if err != nil {
			return []byte{}
		}
		return bytes
	}
}

func BodyToMessageBytesFastHTTP(bodyType MessageBodyType, ctx *fasthttp.RequestCtx) []byte {
	switch bodyType {
	case URLEncodedJSONPayload:
		return ctx.FormValue(ParamPayload)
	case URLEncodedJSONPayloadOrJSON:
		ct := strings.TrimSpace(
			strings.ToLower(
				string(ctx.Request.Header.Peek(hum.HeaderContentType))))
		if strings.Contains(ct, hum.ContentTypeAppJSON) {
			return ctx.PostBody()
		}
		return ctx.FormValue(ParamPayload)
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
	//Body       any         `json:"body,omitempty"`
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
