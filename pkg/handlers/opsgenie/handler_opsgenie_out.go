package opsgenie

import (
	"encoding/json"
	"fmt"
	"strings"

	cc "github.com/grokify/commonchat"

	"github.com/grokify/chathooks/pkg/config"
	"github.com/grokify/chathooks/pkg/handlers"
	"github.com/grokify/chathooks/pkg/models"
)

const (
	DisplayName          = "OpsGenie"
	HandlerKey           = "opsgenie"
	MessageDirection     = "out"
	AlertURLFormat       = "https://app.opsgenie.com/alert/V2#/show/%s"
	UserProfileURLFormat = "https://app.opsgenie.com/user/profile#/user/%s"
	MessageBodyType      = models.JSON
)

func NewHandler() handlers.Handler {
	return handlers.Handler{MessageBodyType: MessageBodyType, Normalize: Normalize}
}

func Normalize(cfg config.Configuration, hReq handlers.HandlerRequest) (cc.Message, error) {
	ccMsg := cc.NewMessage()
	iconURL, err := cfg.GetAppIconURL(HandlerKey)
	if err == nil {
		ccMsg.IconURL = iconURL.String()
	}

	src, err := OpsgenieOutMessageFromBytes(hReq.Body)
	if err != nil {
		return ccMsg, err
	}

	actions := map[string]string{
		"Acknowledge":     "acknowledged",
		"AddNote":         "note added",
		"AddRecipient":    "recipient added",
		"AddTags":         "tags added",
		"AddTeam":         "team added",
		"AssignOwnership": "owner assigned",
		"Close":           "closed",
		"Create":          "created",
		"Delete":          "deleted",
		"Escalate":        "escalated",
		"RemoveTags":      "tags removed",
		"TakeOwnership":   "ownership taken",
		"UnAcknowledge":   "unacknowledged"}

	if verb, ok := actions[src.Action]; ok {
		ccMsg.Activity = fmt.Sprintf("Alert %s", verb)
	} else {
		ccMsg.Activity = fmt.Sprintf("Alert %s", src.Action)
	}

	alertType := src.Source.Name
	if len(alertType) == 0 {
		alertType = src.Source.Type
	}
	if len(alertType) > 0 {
		alertType = fmt.Sprintf("%s alert", alertType)
	} else {
		alertType = "alert"
	}

	ccMsg.Title = fmt.Sprintf("%s %s ([%s](%s))",
		src.IntegrationName,
		alertType,
		src.Alert.AlertID[:8],
		src.Alert.AlertURL())

	attachment := cc.NewAttachment()

	if len(src.Alert.Message) > 0 {
		attachment.AddField(cc.Field{
			Title: "Message",
			Value: src.Alert.Message})
	}
	if len(src.EscalationNotify.ID) > 0 {
		attachment.AddField(cc.Field{
			Title: "Esclated To",
			Value: fmt.Sprintf("[%s](%s)",
				src.EscalationNotify.Name,
				src.EscalationNotify.URL())})
	}
	if len(src.Alert.Note) > 0 {
		attachment.AddField(cc.Field{
			Title: "Note",
			Value: src.Alert.Note})
	}
	if len(src.Alert.Team) > 0 {
		attachment.AddField(cc.Field{
			Title: "Team",
			Value: src.Alert.Team})
	}
	if len(src.Alert.Recipient) > 0 {
		attachment.AddField(cc.Field{
			Title: "Recipient",
			Value: src.Alert.Recipient})
	}
	if len(src.Alert.Tags) > 0 {
		attachment.AddField(cc.Field{
			Title: "Tags",
			Value: src.Alert.TagsFormatted()})
	}
	if len(src.Alert.AddedTags) > 0 {
		attachment.AddField(cc.Field{
			Title: "Added Tags",
			Value: SplitTrimSpaceJoin(src.Alert.AddedTags, ",", ", ")})
	}
	if len(src.Alert.RemovedTags) > 0 {
		attachment.AddField(cc.Field{
			Title: "Removed Tags",
			Value: SplitTrimSpaceJoin(src.Alert.RemovedTags, ",", ", ")})
	}
	if len(src.Alert.Owner) > 0 {
		attachment.AddField(cc.Field{
			Title: "Owner",
			Value: src.Alert.Owner})
	}
	if 1 == 0 {
		if len(src.Alert.Source) > 0 {
			attachment.AddField(cc.Field{
				Title: "Source",
				Value: src.Alert.Source})
		}
	}
	if len(src.Alert.Username) > 0 {
		attachment.AddField(cc.Field{
			Title: "Username / Profile",
			Value: fmt.Sprintf("[%s](%s)", src.Alert.Username, src.Alert.UserURL())})
	}

	ccMsg.AddAttachment(attachment)
	fmt.Printf("MESSAGE_BUILT %v\n", src.Action)
	return ccMsg, nil
}

func SplitTrimSpaceJoin(input string, sep1 string, sep2 string) string {
	inputParts := strings.Split(input, sep1)
	outputParts := []string{}
	for _, part := range inputParts {
		partTrimed := strings.TrimSpace(part)
		if len(partTrimed) > 0 {
			outputParts = append(outputParts, partTrimed)
		}
	}
	return strings.Join(outputParts, sep2)
}

type OpsgenieOutMessage struct {
	Source           OpsgenieOutSource           `json:"source,omitempty"`
	Alert            OpsgenieOutAlert            `json:"alert,omitempty"`
	Action           string                      `json:"action,omitempty"`
	IntegrationID    string                      `json:"integrationId,omitempty"`
	IntegrationName  string                      `json:"integrationName,omitempty"`
	EscalationNotify OpsgenieOutEscalationNotify `json:"escalationNotify,omitempty"`
}

func OpsgenieOutMessageFromBytes(bytes []byte) (OpsgenieOutMessage, error) {
	msg := OpsgenieOutMessage{}
	err := json.Unmarshal(bytes, &msg)
	return msg, err
}

type OpsgenieOutAlert struct {
	UpdatedAt   int64    `json:"updatedAt,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Teams       []string `json:"teams,omitempty"`
	Recipients  []string `json:"recipients,omitempty"`
	Message     string   `json:"message,omitempty"`
	Username    string   `json:"username,omitempty"`
	AlertID     string   `json:"alertId,omitempty"`
	Source      string   `json:"source,omitempty"`
	Alias       string   `json:"alias,omitempty"`
	TinyID      string   `json:"tinyId,omitempty"`
	CreatedAt   int64    `json:"createdAt,omitempty"`
	UserID      string   `json:"userId,omitempty"`
	Entity      string   `json:"entity,omitempty"`
	Owner       string   `json:"owner,omitempty"`
	AddedTags   string   `json:"addedTags,omitempty"`
	RemovedTags string   `json:"removedTags,omitempty"`
	Note        string   `json:"note,omitempty"`
	Recipient   string   `json:"recipient,omitempty"`
	Team        string   `json:"team,omitempty"`
}

func (alert *OpsgenieOutAlert) UserURL() string {
	return fmt.Sprintf(UserProfileURLFormat, alert.UserID)
}

func (alert *OpsgenieOutAlert) AlertURL() string {
	return fmt.Sprintf(AlertURLFormat, alert.UserID)
}

func (alert *OpsgenieOutAlert) TagsFormatted() string {
	return strings.Join(alert.Tags, ", ")
}

type OpsgenieOutSource struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

type OpsgenieOutEscalationNotify struct {
	Name   string `json:"name,omitempty"`
	ID     string `json:"id,omitempty"`
	Type   string `json:"type,omitempty"`
	Entity string `json:"entity,omitempty"`
}

func (alert *OpsgenieOutEscalationNotify) URL() string {
	return fmt.Sprintf(UserProfileURLFormat, alert.ID)
}
