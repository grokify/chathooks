package librato

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	cc "github.com/grokify/commonchat"

	"github.com/grokify/chathooks/pkg/config"
	"github.com/grokify/chathooks/pkg/handlers"
	"github.com/grokify/chathooks/pkg/models"
)

const (
	DisplayName      = "Librato"
	HandlerKey       = "librato"
	MessageDirection = "out"
	DocumentationURL = "https://www.runscope.com/docs/api-testing/notifications#webhook"
	MessageBodyType  = models.JSON
)

var (
	IncludeRecordedAt = false
)

func NewHandler() handlers.Handler {
	return handlers.Handler{MessageBodyType: MessageBodyType, Normalize: Normalize}
}

func Normalize(cfg config.Configuration, hReq handlers.HandlerRequest) (cc.Message, error) {
	src, err := LibratoOutMessageFromBytes(hReq.Body)
	if err != nil {
		return cc.NewMessage(), err
	}

	if src.Clear == "normal" {
		return NormalizeSourceCleared(cfg, src), nil
	}
	return NormalizeSourceTriggered(cfg, src), nil
}

func NormalizeSourceTriggered(cfg config.Configuration, src LibratoOutMessage) cc.Message {
	src.Inflate()

	ccMsg := cc.NewMessage()
	iconURL, err := cfg.GetAppIconURL(HandlerKey)
	if err == nil {
		ccMsg.IconURL = iconURL.String()
	}

	ccMsg.Activity = "Alert triggered"

	if len(src.Alert.Name) > 0 {
		if len(src.Alert.RunbookURL) > 0 {
			ccMsg.Title = fmt.Sprintf("Alert [%v](%s) has triggered!",
				src.Alert.Name, src.Alert.RunbookURL)
		} else {
			ccMsg.Title = fmt.Sprintf("Alert %v has triggered!",
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
			ccMsg.AddAttachment(BuildViolationAttachment(src, violation, violationSuffix))
		}
	}

	return ccMsg
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

func NormalizeSourceCleared(cfg config.Configuration, src LibratoOutMessage) cc.Message {
	ccMsg := cc.NewMessage()
	iconURL, err := cfg.GetAppIconURL(HandlerKey)
	if err == nil {
		ccMsg.IconURL = iconURL.String()
	}

	ccMsg.Activity = "Alert cleared"

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
		ccMsg.Title = fmt.Sprintf("[%s](%s) cleared%s",
			alertName, src.Alert.RunbookURL, triggerTime)
	} else {
		ccMsg.Title = fmt.Sprintf("%s cleared%s", alertName, triggerTime)
	}

	return ccMsg
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
