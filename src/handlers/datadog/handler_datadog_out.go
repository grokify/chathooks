package datadog

import (
	"encoding/json"

	cc "github.com/commonchat/commonchat-go"
	"github.com/grokify/chathooks/src/config"
	"github.com/grokify/chathooks/src/handlers"
	"github.com/grokify/chathooks/src/models"
)

const (
	DisplayName      = "Datadog"
	HandlerKey       = "datadog"
	MessageDirection = "out"
	DocumentationURL = "http://docs.datadoghq.com/integrations/webhooks/"
	MessageBodyType  = models.JSON
)

func NewHandler() handlers.Handler {
	return handlers.Handler{MessageBodyType: MessageBodyType, Normalize: Normalize}
}

func Normalize(cfg config.Configuration, bytes []byte) (cc.Message, error) {
	ccMsg, err := CcMessageFromBytes(bytes)
	if err != nil {
		return ccMsg, err
	}
	iconURL, err := cfg.GetAppIconURL(HandlerKey)
	if err == nil {
		ccMsg.IconURL = iconURL.String()
	}
	return ccMsg, nil
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
