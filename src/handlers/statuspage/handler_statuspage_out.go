package statuspage

import (
	"encoding/json"
	"errors"
	"fmt"

	log "github.com/Sirupsen/logrus"

	cc "github.com/commonchat/commonchat-go"
	"github.com/grokify/gotilla/strings/stringsutil"
	"github.com/grokify/webhook-proxy-go/src/adapters"
	"github.com/grokify/webhook-proxy-go/src/config"
	"github.com/grokify/webhook-proxy-go/src/util"
	"github.com/valyala/fasthttp"
)

const (
	DisplayName        = "StatusPage"
	HandlerKey         = "statuspage"
	IconURL            = "https://pbs.twimg.com/profile_images/643541964051357696/TFppCn6j.png"
	ComponentURLFormat = "http://manage.statuspage.io/pages/%s/components"
)

// FastHttp request handler for Semaphore CI outbound webhook
type StatuspageOutToGlipHandler struct {
	Config  config.Configuration
	Adapter adapters.Adapter
}

// FastHttp request handler constructor for Semaphore CI outbound webhook
func NewStatuspageOutToGlipHandler(cfg config.Configuration, adapter adapters.Adapter) StatuspageOutToGlipHandler {
	return StatuspageOutToGlipHandler{Config: cfg, Adapter: adapter}
}

// HandleFastHTTP is the method to respond to a fasthttp request.
func (h *StatuspageOutToGlipHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
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

// {$component.name} status changed from {$component_update.old_status} to {$component_update.new_status}. [(Manage your Components)]({http://manage.statuspage.io/pages/{$page.id}/components})

func Normalize(bytes []byte) (cc.Message, error) {
	src, err := StatuspageOutMessageFromBytes(bytes)
	if err != nil {
		return cc.NewMessage(), err
	}
	if len(src.ComponentUpdate.CreatedAt) > 0 {
		return NormalizeComponentUpdate(src)
	}
	return NormalizeIncidentUpdate(src)
}

func NormalizeComponentUpdate(src StatuspageOutMessage) (cc.Message, error) {
	message := cc.NewMessage()
	message.IconURL = IconURL

	message.Activity = "Component status changed"

	componentURL, err := src.PageURL()
	if err == nil {
		message.Title = fmt.Sprintf("[%s](%s) component status udpated from **%s** to **%s**",
			src.Component.Name,
			componentURL,
			src.ComponentUpdate.OldStatus,
			src.ComponentUpdate.NewStatus)
	} else {
		message.Title = fmt.Sprintf("%s component status udpated from **%s** to **%s**",
			src.Component.Name,
			src.ComponentUpdate.OldStatus,
			src.ComponentUpdate.NewStatus)
	}

	return message, nil
}

func ToUpperFirstWorlds(input string, sep1 string, sep2 string) string {
	return ""
}

func NormalizeIncidentUpdate(src StatuspageOutMessage) (cc.Message, error) {
	message := cc.NewMessage()
	message.IconURL = IconURL

	if len(src.Incident.IncidentUpdates) == 0 {
		return message, errors.New("No incident update found")
	} else if len(src.Incident.IncidentUpdates) == 1 {
		message.Activity = "New incident created"
	} else {
		message.Activity = "Incident updated"
	}

	thisUpdate := src.Incident.IncidentUpdates[0]

	if len(src.Incident.IncidentUpdates) == 1 {
		message.Title = fmt.Sprintf("[%s](%s) incident created with status **%s**",
			src.Page.StatusDescription,
			src.Incident.Shortlink,
			stringsutil.ToUpperFirst(thisUpdate.Status))
	} else if len(src.Incident.IncidentUpdates) > 1 {
		prevUpdate := src.Incident.IncidentUpdates[1]
		message.Title = fmt.Sprintf("[%s](%s) incident updated from **%s** to **%s**",
			src.Page.StatusDescription,
			src.Incident.Shortlink,
			stringsutil.ToUpperFirst(prevUpdate.Status),
			stringsutil.ToUpperFirst(thisUpdate.Status))
	}

	attachment := cc.NewAttachment()
	if len(thisUpdate.Body) > 0 {
		attachment.AddField(cc.Field{
			Title: "Message",
			Value: thisUpdate.Body})
	}

	if len(attachment.Fields) > 0 {
		message.AddAttachment(attachment)
	}
	return message, nil
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
	if len(msg.Page.Id) < 1 {
		return "", errors.New("Page Id Not Found")
	}
	return fmt.Sprintf(ComponentURLFormat, msg.Page.Id), nil
}

func (msg *StatuspageOutMessage) IncidentURL() string {
	return fmt.Sprintf("")
}

type StatuspageOutMeta struct {
	Unsubscribe   string `json:"unsubscribe,omitempty"`
	Documentation string `json:"documentation,omitempty"`
}

type StatuspageOutPage struct {
	Id                string `json:"id,omitempty"`
	StatusIndicator   string `json:"status_indicator,omitempty"`
	StatusDescription string `json:"status_description,omitempty"`
}

type StatuspageOutComponentUpdate struct {
	CreatedAt   string `json:"created_at,omitempty"`
	NewStatus   string `json:"new_status,omitempty"`
	OldStatus   string `json:"old_status,omitempty"`
	Id          string `json:"id,omitempty"`
	ComponentId string `json:"component_id,omitempty"`
}

type StatuspageOutComponent struct {
	CreatedAt string `json:"created_at,omitempty"`
	Id        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Status    string `json:"status,omitempty"`
}

type StatuspageOutIncident struct {
	Backfilled                    bool                          `json:"backfilled,omitempty"`
	Impact                        string                        `json:"impact,omitempty"`
	ImpactOverride                interface{}                   `json:"impact_override,omitempty"`
	MonitoringAt                  string                        `json:"monitoring_at,omitempty"`
	PostmortemBody                interface{}                   `json:"postmortem_body,omitempty"`
	PostmortemBodyLastUpdatedAt   string                        `json:"postmortem_body_last_updated_at,omitempty"`
	PostmortemIgnored             bool                          `json:"postmortem_ignored,omitempty"`
	PostmortemNotifiedSubscribers bool                          `json:"postmortem_notified_subscribers,omitempty"`
	PostmortemNotifiedTwitter     bool                          `json:"postmortem_notified_twitter,omitempty"`
	PostmortemPublishedAt         string                        `json:"postmortem_published_at,omitempty"`
	ResovledAt                    string                        `json:"resolved_at,omitempty"`
	ScheduledAutoTransition       bool                          `json:"scheduled_auto_transition,omitempty"`
	ScheduledFor                  interface{}                   `json:"scheduled_for,omitempty"`
	ScheduledRemindPrior          bool                          `json:"scheduled_remind_prior,omitempty"`
	ScheduledRemindedAt           interface{}                   `json:"scheduled_reminded_at,omitempty"`
	ScheduledUntil                interface{}                   `json:"scheduled_until,omitempty"`
	Shortlink                     string                        `json:"shortlink,omitempty"`
	Status                        string                        `json:"status,omitempty"`
	UpdatedAt                     string                        `json:"updated_at,omitempty"`
	Id                            string                        `json:"id,omitempty"`
	OrganizationId                string                        `json:"organization_id,omitempty"`
	IncidentUpdates               []StatuspageOutIncidentUpdate `json:"incident_updates,omitempty"`
	Name                          string                        `json:"name,omitempty"`
}

type StatuspageOutIncidentUpdate struct {
	Body               string `json:"body,omitempty"`
	CreatedAt          string `json:"created_at,omitempty"`
	DisplayAt          string `json:"display_at,omitempty"`
	Status             string `json:"status,omitempty"`
	TwitterUpdatedAt   string `json:"twitter_updated_at,omitempty"`
	UpdatedAt          string `json:"updated_at,omitempty"`
	WantsTwitterUpdate bool   `json:"wants_twitter_update,omitempty"`
	Id                 string `json:"id,omitempty"`
	IncidentId         string `json:"incident_id,omitempty"`
}

func StatuspageOutMessageFromBytes(bytes []byte) (StatuspageOutMessage, error) {
	msg := StatuspageOutMessage{}
	err := json.Unmarshal(bytes, &msg)
	return msg, err
}
