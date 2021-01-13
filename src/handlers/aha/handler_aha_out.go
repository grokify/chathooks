package aha

import (
	"encoding/json"
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/grokify/chathooks/src/config"
	"github.com/grokify/chathooks/src/handlers"
	"github.com/grokify/chathooks/src/models"
	cc "github.com/grokify/commonchat"
	"github.com/microcosm-cc/bluemonday"
	"github.com/rs/zerolog/log"
)

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

func Normalize(cfg config.Configuration, hReq handlers.HandlerRequest) (cc.Message, error) {
	ccMsg := cc.NewMessage()
	iconURL, err := cfg.GetAppIconURL(HandlerKey)
	if err == nil {
		ccMsg.IconURL = iconURL.String()
	}

	src, err := AhaOutMessageFromBytes(hReq.Body)
	if err != nil {
		return ccMsg, err
	}

	ccMsg.Title = src.Title()

	p := bluemonday.StrictPolicy()

	if src.Audit != nil && len(src.Audit.Changes) > 0 {
		attachment := cc.NewAttachment()
		for _, change := range src.Audit.Changes {
			field := cc.Field{}
			key := strings.TrimSpace(change.FieldName)
			val := strings.TrimSpace(change.Value)
			val = p.Sanitize(val)
			val = html.UnescapeString(val)
			addField := false
			if len(key) > 0 {
				field.Title = key
				addField = true
			}
			if len(val) > 0 {
				field.Value = val
				addField = true
			}
			if key != "Description" {
				field.Short = true
			}
			if addField {
				attachment.AddField(field)
			}
		}
		ccMsg.AddAttachment(attachment)
	}

	return ccMsg, nil
}

func AhaOutMessageFromBytes(bytes []byte) (AhaOutMessage, error) {
	log.Debug().
		Str("type", "message.raw").
		Str("inbound_body", string(bytes)).
		Str("handler", HandlerKey).
		Msg(config.InfoInputMessageParseBegin)

	resp := AhaOutMessage{}
	err := json.Unmarshal(bytes, &resp)
	if err != nil {
		log.Warn().
			Err(err).
			Str("handler", HandlerKey).
			Str("request_body", string(bytes)).
			Msg(config.ErrorInputMessageParseFailed)
	}
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

func (aom *AhaOutMessage) Title() string {
	if aom.Audit == nil {
		return ""
	}
	return aom.Audit.Title()
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
	username := strings.TrimSpace(aoa.User.Name)
	description := strings.TrimSpace(aoa.Description)
	itemUrl := strings.TrimSpace(aoa.AuditableURL)
	title := ""
	if len(description) > 0 && len(itemUrl) > 0 {
		title = fmt.Sprintf("[%v](%v)", description, itemUrl)
	} else if len(itemUrl) > 0 {
		title = fmt.Sprintf("[%v](%v)", itemUrl, itemUrl)
	} else if len(description) > 0 {
		title = description
	}
	if len(username) > 0 {
		if len(title) > 0 {
			title = fmt.Sprintf("**%v** %v", username, title)
		} else {
			title = fmt.Sprintf("**%v**", username)
		}
	}
	return title
}

type AhaOutUser struct {
	Id        string    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type AhaOutChange struct {
	FieldName string `json:"field_name,omitempty"`
	Value     string `json:"value,omitempty"`
}
