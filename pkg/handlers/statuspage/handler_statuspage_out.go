package statuspage

import (
	"encoding/json"
	"errors"
	"fmt"

	cc "github.com/grokify/commonchat"
	"github.com/grokify/mogo/type/stringsutil"

	"github.com/grokify/chathooks/pkg/config"
	"github.com/grokify/chathooks/pkg/handlers"
	"github.com/grokify/chathooks/pkg/models"
)

const (
	DisplayName        = "StatusPage"
	HandlerKey         = "statuspage"
	MessageDirection   = "out"
	ComponentURLFormat = "http://manage.statuspage.io/pages/%s/components"
	MessageBodyType    = models.JSON
)

func NewHandler() handlers.Handler {
	return handlers.Handler{MessageBodyType: MessageBodyType, Normalize: Normalize}
}

// {$component.name} status changed from {$component_update.old_status} to {$component_update.new_status}. [(Manage your Components)]({http://manage.statuspage.io/pages/{$page.id}/components})

func Normalize(cfg config.Configuration, hReq handlers.HandlerRequest) (cc.Message, error) {
	src, err := StatuspageOutMessageFromBytes(hReq.Body)
	if err != nil {
		return cc.NewMessage(), err
	}
	if len(src.ComponentUpdate.CreatedAt) > 0 {
		return NormalizeComponentUpdate(cfg, src)
	}
	return NormalizeIncidentUpdate(cfg, src)
}

func NormalizeComponentUpdate(cfg config.Configuration, src StatuspageOutMessage) (cc.Message, error) {
	ccMsg := cc.NewMessage()
	iconURL, err := cfg.GetAppIconURL(HandlerKey)
	if err == nil {
		ccMsg.IconURL = iconURL.String()
	}

	ccMsg.Activity = "Component status changed"

	componentURL, err := src.PageURL()
	if err == nil {
		ccMsg.Title = fmt.Sprintf("[%s](%s) component status updated from **%s** to **%s**",
			src.Component.Name,
			componentURL,
			src.ComponentUpdate.OldStatus,
			src.ComponentUpdate.NewStatus)
	} else {
		ccMsg.Title = fmt.Sprintf("%s component status updated from **%s** to **%s**",
			src.Component.Name,
			src.ComponentUpdate.OldStatus,
			src.ComponentUpdate.NewStatus)
	}

	return ccMsg, nil
}

func ToUpperFirstWorlds(input string, sep1 string, sep2 string) string {
	return ""
}

func NormalizeIncidentUpdate(cfg config.Configuration, src StatuspageOutMessage) (cc.Message, error) {
	ccMsg := cc.NewMessage()
	iconURL, err := cfg.GetAppIconURL(HandlerKey)
	if err == nil {
		ccMsg.IconURL = iconURL.String()
	}

	if len(src.Incident.IncidentUpdates) == 0 {
		return ccMsg, errors.New("no incident update found")
	} else if len(src.Incident.IncidentUpdates) == 1 {
		ccMsg.Activity = "New incident created"
	} else {
		ccMsg.Activity = "Incident updated"
	}

	thisUpdate := src.Incident.IncidentUpdates[0]

	if len(src.Incident.IncidentUpdates) == 1 {
		ccMsg.Title = fmt.Sprintf("[%s](%s) incident created with status **%s**",
			src.Page.StatusDescription,
			src.Incident.Shortlink,
			stringsutil.ToUpperFirst(thisUpdate.Status, false))
	} else if len(src.Incident.IncidentUpdates) > 1 {
		prevUpdate := src.Incident.IncidentUpdates[1]
		ccMsg.Title = fmt.Sprintf("[%s](%s) incident updated from **%s** to **%s**",
			src.Page.StatusDescription,
			src.Incident.Shortlink,
			stringsutil.ToUpperFirst(prevUpdate.Status, false),
			stringsutil.ToUpperFirst(thisUpdate.Status, false))
	}

	attachment := cc.NewAttachment()
	if len(thisUpdate.Body) > 0 {
		attachment.AddField(cc.Field{
			Title: "Message",
			Value: thisUpdate.Body})
	}

	if len(attachment.Fields) > 0 {
		ccMsg.AddAttachment(attachment)
	}
	return ccMsg, nil
}

type StatuspageOutMessage struct {
	Meta            StatuspageOutMeta            `json:"meta,omitempty"`
	Page            StatuspageOutPage            `json:"page,omitempty"`
	ComponentUpdate StatuspageOutComponentUpdate `json:"component_update,omitempty"`
	Component       StatuspageOutComponent       `json:"component,omitempty"`
	Incident        StatuspageOutIncident        `json:"incident,omitempty"`
}

func (msg *StatuspageOutMessage) PageURL() (string, error) {
	// http://manage.statuspage.io/pages/{$page.id}/components
	if len(msg.Page.ID) < 1 {
		return "", errors.New("page.ID not found")
	}
	return fmt.Sprintf(ComponentURLFormat, msg.Page.ID), nil
}

func (msg *StatuspageOutMessage) IncidentURL() string {
	return ""
}

type StatuspageOutMeta struct {
	Unsubscribe   string `json:"unsubscribe,omitempty"`
	Documentation string `json:"documentation,omitempty"`
}

type StatuspageOutPage struct {
	ID                string `json:"id,omitempty"`
	StatusIndicator   string `json:"status_indicator,omitempty"`
	StatusDescription string `json:"status_description,omitempty"`
}

type StatuspageOutComponentUpdate struct {
	ID          string `json:"id,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	NewStatus   string `json:"new_status,omitempty"`
	OldStatus   string `json:"old_status,omitempty"`
	ComponentID string `json:"component_id,omitempty"`
}

type StatuspageOutComponent struct {
	ID        string `json:"id,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	Name      string `json:"name,omitempty"`
	Status    string `json:"status,omitempty"`
}

type StatuspageOutIncident struct {
	ID                            string                        `json:"id,omitempty"`
	Backfilled                    bool                          `json:"backfilled,omitempty"`
	Impact                        string                        `json:"impact,omitempty"`
	ImpactOverride                any                           `json:"impact_override,omitempty"`
	MonitoringAt                  string                        `json:"monitoring_at,omitempty"`
	PostmortemBody                any                           `json:"postmortem_body,omitempty"`
	PostmortemBodyLastUpdatedAt   string                        `json:"postmortem_body_last_updated_at,omitempty"`
	PostmortemIgnored             bool                          `json:"postmortem_ignored,omitempty"`
	PostmortemNotifiedSubscribers bool                          `json:"postmortem_notified_subscribers,omitempty"`
	PostmortemNotifiedTwitter     bool                          `json:"postmortem_notified_twitter,omitempty"`
	PostmortemPublishedAt         string                        `json:"postmortem_published_at,omitempty"`
	ResovledAt                    string                        `json:"resolved_at,omitempty"`
	ScheduledAutoTransition       bool                          `json:"scheduled_auto_transition,omitempty"`
	ScheduledFor                  any                           `json:"scheduled_for,omitempty"`
	ScheduledRemindPrior          bool                          `json:"scheduled_remind_prior,omitempty"`
	ScheduledRemindedAt           any                           `json:"scheduled_reminded_at,omitempty"`
	ScheduledUntil                any                           `json:"scheduled_until,omitempty"`
	Shortlink                     string                        `json:"shortlink,omitempty"`
	Status                        string                        `json:"status,omitempty"`
	UpdatedAt                     string                        `json:"updated_at,omitempty"`
	OrganizationID                string                        `json:"organization_id,omitempty"`
	IncidentUpdates               []StatuspageOutIncidentUpdate `json:"incident_updates,omitempty"`
	Name                          string                        `json:"name,omitempty"`
}

type StatuspageOutIncidentUpdate struct {
	ID                 string `json:"id,omitempty"`
	Body               string `json:"body,omitempty"`
	CreatedAt          string `json:"created_at,omitempty"`
	DisplayAt          string `json:"display_at,omitempty"`
	Status             string `json:"status,omitempty"`
	TwitterUpdatedAt   string `json:"twitter_updated_at,omitempty"`
	UpdatedAt          string `json:"updated_at,omitempty"`
	WantsTwitterUpdate bool   `json:"wants_twitter_update,omitempty"`
	IncidentID         string `json:"incident_id,omitempty"`
}

func StatuspageOutMessageFromBytes(bytes []byte) (StatuspageOutMessage, error) {
	msg := StatuspageOutMessage{}
	return msg, json.Unmarshal(bytes, &msg)
}
