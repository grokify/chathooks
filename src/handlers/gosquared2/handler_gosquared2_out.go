package gosquared2

import (
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"

	cc "github.com/commonchat/commonchat-go"
	"github.com/grokify/webhookproxy/src/adapters"
	"github.com/grokify/webhookproxy/src/config"
	"github.com/grokify/webhookproxy/src/util"
	"github.com/valyala/fasthttp"
)

const (
	DisplayName      = "GoSquared"
	HandlerKey       = "gosquared"
	DocumentationURL = "https://www.gosquared.com/customer/en/portal/articles/1996494-webhooks"
)

// FastHttp request handler for Travis CI outbound webhook
type Handler struct {
	Config  config.Configuration
	Adapter adapters.Adapter
}

// FastHttp request handler constructor for Travis CI outbound webhook
func NewHandler(cfg config.Configuration, adapter adapters.Adapter) Handler {
	return Handler{Config: cfg, Adapter: adapter}
}

// HandleFastHTTP is the method to respond to a fasthttp request.
func (h *Handler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	ccMsg, err := Normalize(h.Config, ctx.PostBody())

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

func Normalize(cfg config.Configuration, bytes []byte) (cc.Message, error) {
	src, err := GosquaredOutBaseMessageFromBytes(bytes)
	if err != nil {
		return cc.NewMessage(), err
	}

	if len(src.Person.Id) > 0 {
		return NormalizeSmartGroup(cfg, bytes)
	}
	return NormalizeSiteTraffic(cfg, bytes)
}

func NormalizeSiteTraffic(cfg config.Configuration, bytes []byte) (cc.Message, error) {
	ccMsg := cc.NewMessage()
	iconURL, err := cfg.GetAppIconURL(HandlerKey)
	if err == nil {
		ccMsg.IconURL = iconURL.String()
	}

	src, err := GosquaredOutMessageSiteTrafficFromBytes(bytes)
	if err != nil {
		return ccMsg, err
	}

	ccMsg.Activity = "Site traffic change"

	attachment := cc.NewAttachment()

	attachment.AddField(cc.Field{
		Title: "Site",
		Value: fmt.Sprintf("[%s](%s)", src.SiteDetails.SiteName, src.SiteDetails.URL)})
	attachment.AddField(cc.Field{
		Title: "Concurrent Users Online",
		Value: fmt.Sprintf("%v", src.Concurrents)})
	attachment.AddField(cc.Field{
		Title: "Trigger Boundary",
		Value: fmt.Sprintf("%v", src.TriggeredAlert.Boundary)})

	ccMsg.AddAttachment(attachment)
	return ccMsg, nil
}

func NormalizeSmartGroup(cfg config.Configuration, bytes []byte) (cc.Message, error) {
	ccMsg := cc.NewMessage()
	iconURL, err := cfg.GetAppIconURL(HandlerKey)
	if err == nil {
		ccMsg.IconURL = iconURL.String()
	}

	src, err := GosquaredOutMessageSmartGroupFromBytes(bytes)
	if err != nil {
		return ccMsg, err
	}

	ccMsg.Activity = "Smart Group update"

	attachment := cc.NewAttachment()

	attachment.AddField(cc.Field{
		Title: "Smart Group",
		Value: fmt.Sprintf("[%v](%v)", src.Group.Name, src.GroupURL())})
	attachment.AddField(cc.Field{
		Title: "Person",
		Value: src.Person.Name})
	attachment.AddField(cc.Field{
		Title: "Action",
		Value: src.Boundary})

	ccMsg.AddAttachment(attachment)
	return ccMsg, nil
}

type GosquaredOutBaseMessage struct {
	TriggeredAlert GosquaredOutTriggeredAlert `json:"triggeredAlert,omitempty"`
	Concurrents    int64                      `json:"concurrents,omitempty"`
	Person         GosquaredOutPerson         `json:"person,omitempty"`
}

func GosquaredOutBaseMessageFromBytes(bytes []byte) (GosquaredOutBaseMessage, error) {
	msg := GosquaredOutBaseMessage{}
	err := json.Unmarshal(bytes, &msg)
	return msg, err
}

type GosquaredOutMessageSiteTraffic struct {
	TriggeredAlert GosquaredOutTriggeredAlert `json:"triggeredAlert,omitempty"`
	SiteDetails    GosquaredOutSiteDetails    `json:"siteDetails,omitempty"`
	Concurrents    int64                      `json:"concurrents,omitempty"`
}

func GosquaredOutMessageSiteTrafficFromBytes(bytes []byte) (GosquaredOutMessageSiteTraffic, error) {
	msg := GosquaredOutMessageSiteTraffic{}
	err := json.Unmarshal(bytes, &msg)
	return msg, err
}

type GosquaredOutMessageSmartGroup struct {
	Version   string             `json:"version,omitempty"`
	SiteToken string             `json:"site_token,omitempty"`
	Group     GosquaredOutGroup  `json:"group,omitempty"`
	Boundary  string             `json:"boundary,omitempty"`
	Person    GosquaredOutPerson `json:"person,omitempty"`
}

func GosquaredOutMessageSmartGroupFromBytes(bytes []byte) (GosquaredOutMessageSmartGroup, error) {
	msg := GosquaredOutMessageSmartGroup{}
	err := json.Unmarshal(bytes, &msg)
	return msg, err
}

type GosquaredOutGroup struct {
	Name string `json:"name,omitempty"`
	Id   string `json:"id,omitempty"`
}

func (msg *GosquaredOutMessageSmartGroup) GroupURL() string {
	// https://www.gosquared.com/people/GSN-466237-B/last-seen-1-day
	return fmt.Sprintf("https://www.gosquared.com/people/%s/%s",
		msg.SiteToken, msg.Group.Id)
}

type GosquaredOutPerson struct {
	CreatedAt   string `json:"created_at,omitempty"`
	Phone       string `json:"phone,omitempty"`
	Avatar      string `json:"avatar,omitempty"`
	Description string `json:"description,omitempty"`
	Username    string `json:"username,omitempty"`
	Email       string `json:"email,omitempty"`
	Name        string `json:"name,omitempty"`
	Id          string `json:"id,omitempty"`
}

type GosquaredOutMessageConcurrent struct {
	TriggeredAlert GosquaredOutTriggeredAlert `json:"triggeredAlert,omitempty"`
	Concurrents    int64                      `json:"concurrents,omitempty"`
	SiteDetails    GosquaredOutSiteDetails    `json:"siteDetails,omitempty"`
}

type GosquaredOutTriggeredAlert struct {
	Id       int64  `json:"id,omitempty"`
	Boundary string `json:"boundary,omitempty"`
	Value    string `json:"value,omitempty"`
	Type     string `json:"type,omitempty"`
}

type GosquaredOutSiteDetails struct {
	UserId    int64  `json:"user_id,omitempty"`
	Acct      string `json:"acct,omitempty"`
	Email     string `json:"email,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	SiteName  string `json:"site_name,omitempty"`
	Domain    string `json:"domain,omitempty"`
	URL       string `json:"url,omitempty"`
	Timezone  string `json:"timezone,omitempty"`
}

func (site *GosquaredOutSiteDetails) DashboardURL() string {
	return fmt.Sprintf("https://www.gosquared.com/now/%s", site.Acct)
}

func DashboardURL(siteToken string) string {
	return fmt.Sprintf("https://www.gosquared.com/now/%s", siteToken)
}

func PeopleEveryoneURL(siteToken string) string {
	// https://www.gosquared.com/people/GSN-466237-B/everyone
	return fmt.Sprintf("https://www.gosquared.com/people/%s/everyone", siteToken)
}
