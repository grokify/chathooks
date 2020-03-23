package adapters

import (
	"github.com/grokify/chathooks/src/models"
	cc "github.com/grokify/commonchat"
	"github.com/valyala/fasthttp"

	log "github.com/sirupsen/logrus"
)

var (
	ShowDisplayName = false
)

type AdapterSet struct {
	Adapters map[string]cc.Adapter
}

func NewAdapterSet() AdapterSet {
	return AdapterSet{Adapters: map[string]cc.Adapter{}}
}

func (set *AdapterSet) SendWebhooks(hookData models.HookData) []models.ErrorInfo {
	errs := []models.ErrorInfo{}
	if len(hookData.OutputType) > 0 && len(hookData.OutputURL) > 0 {
		if adapter, ok := set.Adapters[hookData.OutputType]; ok {
			var msg interface{}
			req, res, err := adapter.SendWebhook(hookData.OutputURL, hookData.CanonicalMessage, &msg)
			log.Infof("ADAPTER_API_STATUS_CODE [%s][%v][%s][%s]", hookData.OutputType, res.StatusCode(), hookData.OutputURL, string(res.Body()))
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
