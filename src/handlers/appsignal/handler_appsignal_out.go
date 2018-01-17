package appsignal

import (
	"encoding/json"
	"fmt"

	//log "github.com/Sirupsen/logrus"

	cc "github.com/commonchat/commonchat-go"
	//"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
	//"github.com/eawsy/aws-lambda-go-event/service/lambda/runtime/event/apigatewayproxyevt"
	//"github.com/grokify/chathooks/src/adapters"
	"github.com/grokify/chathooks/src/config"
	"github.com/grokify/chathooks/src/handlers"
	"github.com/grokify/chathooks/src/models"
	"github.com/grokify/gotilla/time/timeutil"
	//"github.com/valyala/fasthttp"
)

const (
	DisplayName      = "AppSignal"
	HandlerKey       = "appsignal"
	MessageDirection = "out"
	MessageBodyType  = models.JSON
)

/*
// FastHttp request handler for outbound webhook
type Handler struct {
	Config          config.Configuration
	AdapterSet      adapters.AdapterSet
	MessageBodyType models.MessageBodyType
}

// FastHttp request handler constructor for outbound webhook
func NewHandler(cfg config.Configuration, adapterSet adapters.AdapterSet) Handler {
	return Handler{Config: cfg, AdapterSet: adapterSet}
}
*/
func NewHandler() handlers.Handler {
	return handlers.Handler{MessageBodyType: MessageBodyType, Normalize: Normalize}
}

/*

func (h Handler) HandlerKey() string {
	return HandlerKey
}

func (h Handler) MessageDirection() string {
	return MessageDirection
}

// HandleEawsyLambda is the method to respond to a fasthttp request.
func (h Handler) HandleEawsyLambda(event *apigatewayproxyevt.Event, ctx *runtime.Context) (models.AwsAPIGatewayProxyOutput, error) {
	hookData := models.HookDataFromEawsyLambdaEvent(models.JSON, event)
	errs := h.HandleCanonical(hookData)
	return models.ErrorInfosToAwsAPIGatewayProxyOutput(errs...), nil
}

// HandleFastHTTP is the method to respond to a fasthttp request.
func (h Handler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	hookData := models.HookDataFromFastHTTPReqCtx(models.JSON, ctx)
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

	ccMsg, err := Normalize(h.Config, hookData.InputBody)

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
*/

func Normalize(cfg config.Configuration, bytes []byte) (cc.Message, error) {
	ccMsg := cc.NewMessage()
	iconURL, err := cfg.GetAppIconURL(HandlerKey)
	if err == nil {
		ccMsg.IconURL = iconURL.String()
	}

	src, err := AppsignalOutMessageFromBytes(bytes)
	if err != nil {
		return ccMsg, err
	}

	if len(src.Marker.URL) > 0 {
		ccMsg.Activity = "App deployed"
		ccMsg.Title = fmt.Sprintf("%v deployed ([%v](%v))", src.Marker.Site, src.Marker.Revision[:7], src.Marker.URL)

		attachment := cc.NewAttachment()
		if 1 == 0 {
			if len(src.Marker.Revision) > 0 {
				field := cc.Field{Title: "Revision", Short: true}
				if len(src.Marker.URL) > 0 {
					field.Value = fmt.Sprintf("[%v](%v)", src.Marker.Revision[:7], src.Marker.URL)
				} else {
					field.Value = src.Marker.Revision[:7]
				}
				attachment.AddField(field)
			} else if len(src.Marker.URL) > 0 {
				attachment.AddField(cc.Field{
					Title: "Build",
					Value: src.Marker.URL})
			}
		}
		if len(src.Marker.Environment) > 0 {
			attachment.AddField(cc.Field{
				Title: "Environment",
				Value: src.Marker.Environment})
		}
		if len(src.Marker.User) > 0 {
			attachment.AddField(cc.Field{
				Title: "User",
				Value: src.Marker.User,
				Short: true})
		}
		ccMsg.AddAttachment(attachment)
	} else if len(src.Exception.URL) > 0 {
		ccMsg.Activity = fmt.Sprintf("Exception incident")

		exceptionString := ""
		if len(src.Exception.URL) > 0 {
			if len(src.Exception.Exception) > 0 {
				exceptionString = fmt.Sprintf("[%v](%v)", src.Exception.Exception, src.Exception.URL)
			} else {
				exceptionString = fmt.Sprintf("[%v](%v)", src.Exception.URL, src.Exception.URL)
			}
		} else if len(src.Exception.Exception) > 0 {
			exceptionString = src.Exception.Exception
		}
		if len(exceptionString) > 0 {
			exceptionString = fmt.Sprintf(": %s", exceptionString)
		}

		ccMsg.Title = fmt.Sprintf("%v exception incident has occurred%s", src.Exception.Site, exceptionString)

		attachment := cc.NewAttachment()

		if 1 == 0 {
			if len(src.Exception.URL) > 0 {
				field := cc.Field{Title: "Exception", Short: true}
				if len(src.Exception.Exception) > 0 {
					field.Value = fmt.Sprintf("[%v](%v)", src.Exception.Exception, src.Exception.URL)
				} else {
					field.Value = fmt.Sprintf("[%v](%v)", src.Exception.URL, src.Exception.URL)
				}
				attachment.AddField(field)
			} else if len(src.Exception.Exception) > 0 {
				attachment.AddField(cc.Field{
					Title: "Exception",
					Value: src.Exception.Exception,
					Short: true})
			}
		}
		if len(src.Exception.Message) > 0 {
			attachment.AddField(cc.Field{
				Title: "Message",
				Value: src.Exception.Message,
				Short: true})
		} else if len(src.Exception.URL) > 0 {
			attachment.AddField(cc.Field{
				Title: "Exception",
				Value: src.Exception.URL,
				Short: true})
		}
		if len(src.Exception.Environment) > 0 {
			attachment.AddField(cc.Field{
				Title: "Environment",
				Value: src.Exception.Environment,
				Short: true})
		}
		if len(src.Exception.User) > 0 {
			attachment.AddField(cc.Field{
				Title: "User",
				Value: src.Exception.User,
				Short: true})
		}

		ccMsg.AddAttachment(attachment)
	} else if len(src.Performance.URL) > 0 {
		ccMsg.Activity = "Performance incident"

		if src.Performance.Duration > 0.0 {
			durationString, err := timeutil.DurationStringMinutesSeconds(int64(src.Performance.Duration))
			if err == nil {
				ccMsg.Title = fmt.Sprintf("%v performance incident has occurred for %v", src.Performance.Site, durationString)
			} else {
				ccMsg.Title = fmt.Sprintf("%v performance incident has occurred for %v", src.Performance.Site, src.Performance.Duration)
			}
		}

		attachment := cc.NewAttachment()

		if len(src.Performance.URL) > 0 {
			attachment.AddField(cc.Field{
				Title: "Action",
				Value: fmt.Sprintf("[%v](%v)", src.Performance.Action, src.Performance.URL),
				Short: true})
		} else if len(src.Performance.Action) > 0 {
			attachment.AddField(cc.Field{
				Title: "Action",
				Value: src.Performance.Action,
				Short: true})
		}
		if len(src.Performance.Hostname) > 0 {
			attachment.AddField(cc.Field{
				Title: "Hostname",
				Value: src.Performance.Hostname,
				Short: true})
		}
		if len(src.Performance.Environment) > 0 {
			attachment.AddField(cc.Field{
				Title: "Environment",
				Value: src.Performance.Environment,
				Short: true})
		}
		if len(src.Performance.User) > 0 {
			attachment.AddField(cc.Field{
				Title: "User",
				Value: src.Performance.User,
				Short: true})
		}

		ccMsg.AddAttachment(attachment)
	}

	return ccMsg, nil
}

type AppsignalOutMessage struct {
	Marker      AppsignalMarker      `json:"marker,omitempty"`
	Exception   AppsignalException   `json:"exception,omitempty"`
	Performance AppsignalPerformance `json:"performance,omitempty"`
	Test        string               `json:"test,omitempty"`
}

func AppsignalOutMessageFromBytes(bytes []byte) (AppsignalOutMessage, error) {
	msg := AppsignalOutMessage{}
	err := json.Unmarshal(bytes, &msg)
	return msg, err
}

type AppsignalMarker struct {
	User        string `json:"user,omitempty"`
	Site        string `json:"site,omitempty"`
	Environment string `json:"environment,omitempty"`
	Revision    string `json:"revision,omitempty"`
	Repository  string `json:"repository,omitempty"`
	URL         string `json:"url,omitempty"`
}

type AppsignalException struct {
	Exception          string `json:"exception,omitempty"`
	Site               string `json:"site,omitempty"`
	Message            string `json:"message,omitempty"`
	Action             string `json:"action,omitempty"`
	Path               string `json:"path,omitempty"`
	Revision           string `json:"revision,omitempty"`
	User               string `json:"user,omitempty"`
	Hostname           string `json:"hostname,omitempty"`
	FirstBacktraceLine string `json:"first_backtrace_line,omitempty"`
	URL                string `json:"url,omitempty"`
	Environment        string `json:"environment,omitempty"`
	Namespace          string `json:"namespace,omitempty"`
}

/*

{
  "exception":{
    "exception":"RuntimeError",
    "site":"My Glip App",
    "message":"Test Exception",
    "action":"GET /",
    "path":"/","revision":"No deploy yet",
    "user":"N/A",
    "hostname":"lmrc6152.rcoffice.ringcentral.com",
    "first_backtrace_line":"oauth2.rb:17:in `block in \u003cmain\u003e'",
    "url":"https://appsignal.com/grokbase/sites/58bdbb7c16b7e2656bfc3bed/web/exceptions/GET%20-slash-/RuntimeError",
    "environment":"development",
    "metadata":{},
    "namespace":"web"
  }
}
*/

type AppsignalPerformance struct {
	Site        string  `json:"site,omitempty"`
	Action      string  `json:"action,omitempty"`
	Path        string  `json:"path,omitempty"`
	Duration    float64 `json:"duration,omitempty"`
	Status      int64   `json:"status,omitempty"`
	Hostname    string  `json:"hostname,omitempty"`
	Revision    string  `json:"revision,omitempty"`
	User        string  `json:"user,omitempty"`
	URL         string  `json:"url,omitempty"`
	Environment string  `json:"environment,omitempty"`
}
