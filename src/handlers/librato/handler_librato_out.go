package librato

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"

	cc "github.com/commonchat/commonchat-go"
	"github.com/grokify/webhookproxy/src/adapters"
	"github.com/grokify/webhookproxy/src/config"
	"github.com/grokify/webhookproxy/src/util"
	"github.com/valyala/fasthttp"
)

const (
	DisplayName      = "Librato"
	HandlerKey       = "librato"
	MessageDirection = "out"
	IconURL          = "https://raw.githubusercontent.com/grokify/webhookproxy/master/images/icons/librato_128x128.png"
	IconURLX         = "https://a.slack-edge.com/ae7f/plugins/librato/assets/service_512.png"
	DocumentationURL = "https://www.runscope.com/docs/api-testing/notifications#webhook"
)

var (
	IncludeRecordedAt = false
)

// FastHttp request handler for outbound webhook
type Handler struct {
	Config             config.Configuration
	Adapter            adapters.Adapter
	FilterFailuresOnly bool
}

// FastHttp request handler constructor for outbound webhook
func NewHandler(cfg config.Configuration, adapter adapters.Adapter) Handler {
	return Handler{Config: cfg, Adapter: adapter}
}

func (h Handler) HandlerKey() string {
	return HandlerKey
}

func (h Handler) MessageDirection() string {
	return MessageDirection
}

// HandleFastHTTP is the method to respond to a fasthttp request.
func (h *Handler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
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
		if len(src.Alert.RunbookURL) > 0 {
			message.Title = fmt.Sprintf("Alert [%v](%s) has triggered!",
				src.Alert.Name, src.Alert.RunbookURL)
		} else {
			message.Title = fmt.Sprintf("Alert %v has triggered!",
				src.Alert.Name)
		}
	}

	for violationName, violationSet := range src.Violations {
		n := len(violationSet)
		for i, violation := range violationSet {
			violation.Name = violationName
			violationSuffix := ""
			if n > 1 {
				violationSuffix = fmt.Sprintf(" %v", i+1)
			}
			message.AddAttachment(BuildViolationAttachment(src, violation, violationSuffix))
		}
	}

	return message
}

func BuildViolationAttachment(src LibratoOutMessage, violation LibratoOutViolation, violationSuffix string) cc.Attachment {
	attachment := cc.NewAttachment()

	condition, errNoCondition := src.GetCondition(violation.ConditionViolated)

	IncludeRecordedAt = true
	violationRecordedAtSuffix := ""
	if IncludeRecordedAt && violation.RecordedAt > 0 {
		dt := time.Unix(violation.RecordedAt, 0).UTC()
		violationRecordedAtSuffix = fmt.Sprintf(" recorded at %v", dt.Format(time.RFC1123))
	}

	if errNoCondition == nil {
		conditionComparison := "above"
		if float64(violation.Value) < condition.Threshold {
			conditionComparison = "below"
		}

		attachment.AddField(cc.Field{
			Title: fmt.Sprintf("Violation%s", violationSuffix),
			Value: fmt.Sprintf("%s metric `%v` was **%s** threshold %v with value %v%s",
				violation.Name,
				violation.Metric,
				conditionComparison,
				strconv.FormatFloat(condition.Threshold, 'f', -1, 64),
				violation.Value,
				violationRecordedAtSuffix)})
	} else {
		attachment.AddField(cc.Field{
			Title: "Violation",
			Value: fmt.Sprintf("%v: metric `%v` with value %v%s",
				violation.Name,
				violation.Metric,
				violation.Value,
				violationRecordedAtSuffix)})
	}

	if 1 == 0 {
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
