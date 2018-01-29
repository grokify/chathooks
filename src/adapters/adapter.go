package adapters

import (
	"github.com/grokify/chathooks/src/models"
	cc "github.com/grokify/commonchat"
	"github.com/valyala/fasthttp"
)

var (
	ShowDisplayName = false
)

type AdapterSet struct {
	Adapters map[string]Adapter
}

func NewAdapterSet() AdapterSet {
	return AdapterSet{Adapters: map[string]Adapter{}}
}

func (set *AdapterSet) SendWebhooks(hookData models.HookData) []models.ErrorInfo {
	errs := []models.ErrorInfo{}
	if len(hookData.OutputType) > 0 && len(hookData.OutputURL) > 0 {
		if adapter, ok := set.Adapters[hookData.OutputType]; ok {
			var msg interface{}
			req, res, err := adapter.SendWebhook(hookData.OutputURL, hookData.CanonicalMessage, &msg)
			errs = set.procResponse(errs, req, res, err)
		}
	}
	for _, namedAdapter := range hookData.OutputNames {
		if adapter, ok := set.Adapters[namedAdapter]; ok {
			var msg interface{}
			req, res, err := adapter.SendMessage(hookData.CanonicalMessage, &msg)
			errs = set.procResponse(errs, req, res, err)
		}
	}
	return errs
}

func (set *AdapterSet) procResponse(errs []models.ErrorInfo, req *fasthttp.Request, res *fasthttp.Response, err error) []models.ErrorInfo {
	if err != nil {
		errs = append(errs, models.ErrorInfo{StatusCode: 500, Body: []byte(err.Error())})
	} else if res.StatusCode() > 299 {
		errs = append(errs, models.ErrorInfo{
			StatusCode: res.StatusCode(),
			Body:       res.Body(),
		})
	}
	fasthttp.ReleaseRequest(req)
	fasthttp.ReleaseResponse(res)
	return errs
}

type Adapter interface {
	SendWebhook(url string, message cc.Message, finalMsg interface{}) (*fasthttp.Request, *fasthttp.Response, error)
	SendMessage(message cc.Message, finalMsg interface{}) (*fasthttp.Request, *fasthttp.Response, error)
	WebhookUID(ctx *fasthttp.RequestCtx) (string, error)
}

func IntegrationActivitySuffix(displayName string) string {
	if !ShowDisplayName || len(displayName) < 1 {
		return ""
	}
	return ""
}
