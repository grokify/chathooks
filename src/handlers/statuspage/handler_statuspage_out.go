package statuspage

import (
	"encoding/json"
	"errors"
	"fmt"

	log "github.com/Sirupsen/logrus"

	cc "github.com/commonchat/commonchat-go"
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

	message.Activity = "Status Changed"
	componentURL, err := src.PageURL()
	if err == nil {
		message.Title = fmt.Sprintf("[%s](%s) status changed", src.Component.Name)
	} else {
		message.Title = fmt.Sprintf("%s status changed", src.Component.Name)
	}
	attachment := cc.Attachment{}

	if len(src.ComponentUpdate.NewStatus) > 0 {
		attachment.AddField(cc.Field{
			Title: "New Status",
			Value: src.ComponentUpdate.NewStatus})
	}

	if len(src.ComponentUpdate.OldStatus) > 0 {
		attachment.AddField(cc.Field{
			Title: "Old Status",
			Value: src.ComponentUpdate.OldStatus})
	}

	message.AddAttachment(attachment)
	return message
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
	return fmt.Sprintf(ComponentURLFormat, msg.Page.Id)
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

/*
type MagnumciOutMessage struct {
	Id             int64  `json:"id,omitempty"`
	ProjectId      int64  `json:"project_id,omitempty"`
	Title          string `json:"title,omitempty"`
	Number         int64  `json:"number,omitempty"`
	Commit         string `json:"commit,omitempty"`
	Author         string `json:"author,omitempty"`
	Committer      string `json:"committer,omitempty"`
	Message        string `json:"message,omitempty"`
	Branch         string `json:"branch,omitempty"`
	State          string `json:"state,omitempty"`
	Status         string `json:"status,omitempty"`
	Result         int64  `json:"result,omitempty"`
	Duration       int64  `json:"duration,omitempty"`
	DurationString string `json:"duration_string,omitempty"`
	CommitURL      string `json:"commit_url,omitempty"`
	CompareURL     string `json:"compare_url,omitempty"`
	BuildURL       string `json:"build_url,omitempty"`
	StartedAt      string `json:"started_at,omitempty"`
	FinishedAt     string `json:"finished_at,omitempty"`
}

func MagnumciOutMessageFromBytes(bytes []byte) (MagnumciOutMessage, error) {
	msg := MagnumciOutMessage{}
	err := json.Unmarshal(bytes, &msg)
	return msg, err
}
*/
