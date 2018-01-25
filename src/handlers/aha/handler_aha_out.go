package circleci

import (
	"encoding/json"
	"fmt"

	cc "github.com/commonchat/commonchat-go"
	"github.com/grokify/chathooks/src/config"
	"github.com/grokify/chathooks/src/handlers"
	"github.com/grokify/chathooks/src/models"
)

/*
{
    "event":"audit",
    "audit":{
        "id":"6514758955727812214",
        "audit_action":"update",
        "created_at":"2018-01-24T23:12:00.903Z",
        "interesting":true,
        "user":{
            "id":"6355516420883588191",
            "name":"John Wang",
            "email":"john.wang@ringcentral.com",
            "created_at":"2016-11-21T20:09:39.022Z",
            "updated_at":"2018-01-24T21:50:03.013Z"
        },
        "auditable_type":"release",
        "auditable_id":"6489206952503346836",
        "description":"updated release API-R-8 10.0",
        "auditable_url":"https://ringcentral.aha.io/releases/API-R-8",
        "changes":[
            {
                "field_name":"Workflow status",
                "value":"Under consideration \u0026rarr; Shipped"
            }
        ]
    }
}

Aha! release update






*/

const (
	DisplayName      = "Aha!"
	HandlerKey       = "aha"
	MessageDirection = "out"
	DocumentationURL = "https://support.aha.io/hc/en-us/articles/202000997-Integrate-with-Webhooks"
	MessageBodyType  = models.JSON
)

func NewHandler() handlers.Handler {
	return handlers.Handler{MessageBodyType: MessageBodyType, Normalize: Normalize}
}

func Normalize(cfg config.Configuration, bytes []byte) (cc.Message, error) {
	ccMsg := cc.NewMessage()
	iconURL, err := cfg.GetAppIconURL(HandlerKey)
	if err == nil {
		ccMsg.IconURL = iconURL.String()
	}

	src, err := AhaOutMessageFromBytes(bytes)
	if err != nil {
		return ccMsg, err
	}

	ccMsg.Activity = src.Activity()
	ccMsg.Title = src.Title()

	if src.Audit != nil {
		for _, change := range src.Audit.Changes {
			attachment := cc.NewAttachment()
			field := cc.Field{}
			key := strings.TrimSpace(change.FieldName)
			val := strings.TrimSpace(change.Value)
			addField := false
			if len(key) > 0 {
				field.Title = key
				addField = true
			}
			if len(val) > 0 {
				field.Value = val
				addField = true
			}
			if addField {
				attachment.AddField(field)
				ccMsg.AddAttachment(attachment)
			}
		}
	}

	return ccMsg, nil
}

func AhaOutMessageFromBytes(bytes []byte) (AhaOutMessage, error) {
	resp := AhaOutMessage{}
	err := json.Unmarshal(bytes, &resp)
	return resp, err
}

type AhaOutMessage struct {
	Event string       `json:"event,omitempty"`
	Audit *AhaOutAudit `json:"audit,omitempty"`
}

func (aom *AhaOutMessage) Activity() string {
	if aom.Audit == nil {
		return ""
	}
	return aom.Audit.Activity()
}

type AhaOutAudit struct {
	Id            string         `json:"id,omitempty"`
	AuditAction   string         `json:"audit_action,omitempty"`
	CreatedAt     time.Time      `json:"created_at,omitempty"`
	Interesting   bool           `json:"interesting,omitempty"`
	AuditableType string         `json:"auditable_type,omitempty"`
	AuditableId   string         `json:"auditable_id,omitempty"`
	User          *AhaOutUser    `json:"user,omitmpty"`
	Description   string         `json:"description,omitempty"`
	AuditableURL  string         `json:"auditable_url,omitempty"`
	Changes       []AhaOutChange `json:"changes,omitempty"`
}

func (aoa *AhaOutAudit) Activity() string {
	return fmt.Sprintf("%v %v %v", DisplayName, aoa.AuditAction, aoa.AuditableType)
}

func (aoa *AhaOutAudit) Title() string {
	return fmt.Sprintf("**%v** [%v](%v)", aoa.User.Name, aoa.Description, aoa.AuditableURL)
}

type AhaOutUser struct {
	Id        string    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.TIme `json:"updated_at,omitempty"`
}

type AhaOutChange struct {
	FieldName string `json:"field_name,omitempty"`
	Value     string `json:"value,omitempty"`
}
