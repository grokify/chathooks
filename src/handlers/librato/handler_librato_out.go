package librato

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"

	cc "github.com/commonchat/commonchat-go"
	"github.com/grokify/webhook-proxy-go/src/adapters"
	"github.com/grokify/webhook-proxy-go/src/config"
	"github.com/grokify/webhook-proxy-go/src/util"
	"github.com/valyala/fasthttp"
)

const (
	DisplayName      = "Librato"
	HandlerKey       = "librato"
	IconURL          = "https://raw.githubusercontent.com/grokify/webhook-proxy-go/master/images/icons/librato_128x128.png"
	IconURLX         = "https://a.slack-edge.com/ae7f/plugins/librato/assets/service_512.png"
	DocumentationURL = "https://www.runscope.com/docs/api-testing/notifications#webhook"
)

// FastHttp request handler for Travis CI outbound webhook
type LibratoOutToGlipHandler struct {
	Config             config.Configuration
	Adapter            adapters.Adapter
	FilterFailuresOnly bool
}

// FastHttp request handler constructor for Travis CI outbound webhook
func NewLibratoOutToGlipHandler(cfg config.Configuration, adapter adapters.Adapter) LibratoOutToGlipHandler {
	return LibratoOutToGlipHandler{Config: cfg, Adapter: adapter}
}

// HandleFastHTTP is the method to respond to a fasthttp request.
func (h *LibratoOutToGlipHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	fmt.Printf(string(ctx.PostBody()))
	ccMsg, err := Normalize(ctx.PostBody())

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotAcceptable)
		log.WithFields(log.Fields{
			"type":   "http.response",
			"status": fasthttp.StatusNotAcceptable,
		}).Info(fmt.Sprintf("%v request is not acceptable:  %v", DisplayName, err))
		return
	}

	util.SendWebhook(ctx, h.Adapter, ccMsg)
}

func Normalize(bytes []byte) (cc.Message, error) {
	src, err := LibratoOutMessageFromBytes(bytes)
	if err != nil {
		return cc.NewMessage(), err
	}

	if src.Clear == "normal" {
		return NormalizeSourceCleared(src), nil
	}
	return NormalizeSourceTriggered(src), nil
}

func NormalizeSourceTriggered(src LibratoOutMessage) cc.Message {
	src.Inflate()

	message := cc.NewMessage()
	message.IconURL = IconURL

	message.Activity = "Alert triggered"

	if len(src.Alert.Name) > 0 {
		message.Title = fmt.Sprintf("[%v](%s) alert triggered!",
			src.Alert.Name, src.Alert.RunbookURL)
	}

	for violationName, violationSet := range src.Violations {
		for _, violation := range violationSet {
			violation.Name = violationName
			message.AddAttachment(BuildViolationAttachment(src, violation))
		}
	}

	return message
}

func BuildViolationAttachment(src LibratoOutMessage, violation LibratoOutViolation) cc.Attachment {

	attachment := cc.NewAttachment()

	if len(violation.Name) > 0 {
		attachment.AddField(cc.Field{
			Title: "Violation Name",
			Value: violation.Name,
			Short: true})
	}

	if len(violation.Metric) > 0 {
		attachment.AddField(cc.Field{
			Title: "Metric",
			Value: violation.Metric,
			Short: true})
	}

	condition, err := src.GetCondition(violation.ConditionViolated)
	if err == nil {
		attachment.AddField(cc.Field{
			Title: "Threshold",
			Value: fmt.Sprintf("%v", condition.Threshold),
			Short: true})
	}
	if violation.Value > 0.0 {
		attachment.AddField(cc.Field{
			Title: "Value",
			Value: fmt.Sprintf("%v", violation.Value),
			Short: true})
	}

	if 1 == 0 {
		field := cc.Field{}

		if len(violation.Name) > 0 {
			field.Title = violation.Name
		}

		condition, err := src.GetCondition(violation.ConditionViolated)
		if err == nil {
			field.Value = fmt.Sprintf("Metric %s was above threshold %v with value %v",
				violation.Metric,
				condition.Threshold,
				violation.Value)
		}
		attachment.AddField(field)
	}

	if violation.RecordedAt > 0 {
		dt := time.Unix(violation.RecordedAt, 0).UTC()
		attachment.AddField(cc.Field{
			Title: "Recorded At",
			Value: dt.Format(time.RFC1123)})
	}

	return attachment
}

func NormalizeSourceCleared(src LibratoOutMessage) cc.Message {
	message := cc.NewMessage()
	message.IconURL = IconURL

	message.Activity = "Alert cleared"

	alertName := src.Alert.Name
	if len(alertName) < 1 {
		alertName = "Alert"
	}

	triggerTime := ""
	if src.TriggerTime > 0 {
		dt := time.Unix(src.TriggerTime, 0).UTC()
		triggerTime = fmt.Sprintf(" at %v", dt.Format(time.RFC1123))
	}

	if len(src.Alert.RunbookURL) > 0 {
		message.Title = fmt.Sprintf("[%s](%s) cleared%s",
			alertName, src.Alert.RunbookURL, triggerTime)
	} else {
		message.Title = fmt.Sprintf("%s cleared%s", alertName, triggerTime)
	}

	return message
}

type LibratoOutMessage struct {
	Alert         LibratoOutAlert                  `json:"alert,omitempty"`
	Account       string                           `json:"account,omitempty"`
	TriggerTime   int64                            `json:"trigger_time,omitempty"`
	Conditions    []LibratoOutCondition            `json:"conditions,omitempty"`
	ConditionsMap map[int64]LibratoOutCondition    `json:"-,omitempty"`
	Violations    map[string][]LibratoOutViolation `json:"violations,omitempty"`
	Clear         string                           `json:"clear,omitempty"`
}

func (msg *LibratoOutMessage) Inflate() {
	msg.ConditionsMap = map[int64]LibratoOutCondition{}
	for _, condition := range msg.Conditions {
		msg.ConditionsMap[condition.Id] = condition
	}
}

func (msg *LibratoOutMessage) GetCondition(conditionId int64) (LibratoOutCondition, error) {
	if condition, ok := msg.ConditionsMap[conditionId]; ok {
		return condition, nil
	}
	return LibratoOutCondition{},
		errors.New(fmt.Sprintf("Condition %v not found", conditionId))
}

type LibratoOutAlert struct {
	Id         int64  `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	RunbookURL string `json:"runbook_url,omitempty"`
	Version    int64  `json:"version,omitempty"`
}

type LibratoOutCondition struct {
	Id        int64   `json:"id,omitempty"`
	Type      string  `json:"type,omitempty"`
	Threshold float64 `json:"threshold,omitempty"`
	Duration  int64   `json:"duration,omitempty"`
}

type LibratoOutViolation struct {
	Name              string
	Metric            string  `json:"metric,omitempty"`
	Value             float64 `json:"value,omitempty"`
	RecordedAt        int64   `json:"recorded_at,omitempty"`
	ConditionViolated int64   `json:"condition_violated,omitempty"`
	Count             int64   `json:"count,omitempty"`
	Begin             int64   `json:"begin,omitempty"`
	End               int64   `json:"end,omitempty"`
}

func LibratoOutMessageFromBytes(bytes []byte) (LibratoOutMessage, error) {
	msg := LibratoOutMessage{}
	err := json.Unmarshal(bytes, &msg)
	return msg, err
}
