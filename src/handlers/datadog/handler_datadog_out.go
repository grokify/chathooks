package datadog

import (
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"

	cc "github.com/commonchat/commonchat-go"
	"github.com/grokify/webhook-proxy-go/src/adapters"
	"github.com/grokify/webhook-proxy-go/src/config"
	"github.com/grokify/webhook-proxy-go/src/util"
	"github.com/valyala/fasthttp"
)

const (
	DisplayName = "Datadog"
	HandlerKey  = "datadog"
	IconURL     = "https://dka575ofm4ao0.cloudfront.net/pages-favicon_logos/original/428/open-uri20140327-21944-1w47zpx"
)

// FastHttp request handler for Travis CI outbound webhook
type Handler struct {
	Config  config.Configuration
	Adapter adapters.Adapter
}

// FastHttp request handler constructor for Travis CI outbound webhook
func NewHandler(cfg config.Configuration, adapter adapters.Adapter) Handler {
	return Handler{Config: cfg, Adapter: adapter}
}

// HandleFastHTTP is the method to respond to a fasthttp request.
func (h *Handler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	ccMsg, err := Normalize(ctx.PostBody())

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotAcceptable)
		log.WithFields(log.Fields{
			"type":   "http.response",
			"status": fasthttp.StatusNotAcceptable,
		}).Info(fmt.Sprintf("%v request is not acceptable.", DisplayName))
		return
	}

	util.SendWebhook(ctx, h.Adapter, ccMsg)
}

func Normalize(bytes []byte) (cc.Message, error) {
	message, err := CcMessageFromBytes(bytes)
	if err != nil {
		return message, err
	}
	message.IconURL = IconURL
	return message, nil
}

/*

$ID	ID of the event (ex: 1234567)
$EVENT_TITLE	Title of the event (ex: [Triggered] [Memory Alert])
$EVENT_MSG	Text of the event (ex: @webhook-url Sending to the webhook)
$EVENT_TYPE	Type of the event (values: metric_alert_monitor, event_alert, or service_check)
$LAST_UPDATED	Date when the event was last updated .
$DATE	Date (epoch) where the event happened (ex: 1406662672000)
$AGGREG_KEY	ID to aggregate events belonging together (ex: 9bd4ac313a4d1e8fae2482df7b77628)
$ORG_ID	ID of your organization (ex: 11023)
$ORG_NAME	Name of your organization (ex: Datadog)
$USER	User posting the event that triggered the webhook (ex: rudy)
$SNAPSHOT	Url of the image if the event contains a snapshot (ex: https://url.to.snpashot.com/)
$LINK	Url of the event (ex: https://app.datadoghq.com/event/jump_to?event_id=123456)
$PRIORITY	Priority of the event (values: normal or low)
$TAGS	Comma-separated list of the event tags (ex: monitor, name:myService, role:computing-node)
$ALERT_ID	ID of alert (ex: 1234)
$ALERT_TITLE	Title of the alert
$ALERT_METRIC	Name of the metric if it’s an alert (ex: system.load.1)
$ALERT_SCOPE	Comma-separated list of tags triggering the alert (ex: availability-zone:us-east-1a, role:computing-node)
$ALERT_QUERY	Query of the monitor that triggered the webhook
$ALERT_STATUS	Summary of the alert status (ex: system.load.1 over host:my-host was > 0 at least once during the last 1m)
$ALERT_TRANSITION	Type of alert notification (values: Triggered or Recovered)
If you want to post your webhooks to a service requiring authentication, you can Basic HTTP authentication my modifing your URL from https://my.service.com to https://username:password@my.service.com.

Example template:

{
    "activity":"Event triggered",
    "body":"[Event $ID]($LINK): $EVENT_TITLE\n**Priority**\n$PRIORITY\n**Alert**\n$ALERT_STATUS"
}

{
    "activity":"Event triggered",
    "body":"[Event 1234567](https://app.datadoghq.com/event/jump_to?event_id=123456): [Triggered] [Memory Alert]\n**Priority**\n#PRIORITY\n**Alert**\nystem.load.1 over host:my-host was > 0 at least once during the last 1m"
}

*/

func CcMessageFromBytes(bytes []byte) (cc.Message, error) {
	msg := cc.Message{}
	err := json.Unmarshal(bytes, &msg)
	return msg, err
}
