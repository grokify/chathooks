package gosquared

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"

	cc "github.com/commonchat/commonchat-go"
	"github.com/grokify/chatmore/src/adapters"
	"github.com/grokify/chatmore/src/config"
	"github.com/grokify/chatmore/src/util"
	"github.com/valyala/fasthttp"
)

const (
	DisplayName      = "GoSquared"
	HandlerKey       = "gosquared"
	IconURL          = "https://d2rbro28ib85bu.cloudfront.net/images/integrations/128/gosquared.png"
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

func Normalize(bytes []byte) (cc.Message, error) {
	src, err := GosquaredOutBaseMessageFromBytes(bytes)
	if err != nil {
		return cc.NewMessage(), err
	}
	if len(src.Message.Id) > 0 {
		return NormalizeLiveMessage(bytes)
	} else if len(src.Person.Id) > 0 {
		return NormalizeSmartGroup(bytes)
	}
	return NormalizeSiteTraffic(bytes)
}

func NormalizeSiteTraffic(bytes []byte) (cc.Message, error) {
	message := cc.NewMessage()
	message.IconURL = IconURL

	src, err := GosquaredOutMessageSiteTrafficFromBytes(bytes)
	if err != nil {
		return message, err
	}

	if src.TriggeredAlert.Boundary == "upper" {
		message.Activity = "Site traffic spike"
	} else { // if src.TriggeredAlert.Boundary == "lower" {
		message.Activity = "Site traffic dip"
	}

	pluralSuffix := "s"
	if src.Concurrents == int64(1) {
		pluralSuffix = ""
	}

	message.Title = fmt.Sprintf("[%s](%s) has [%v visitor%s online](%s)",
		src.SiteDetails.SiteName,
		src.SiteDetails.URL,
		src.Concurrents,
		pluralSuffix,
		src.SiteDetails.DashboardURL())

	return message, nil
}

func NormalizeSmartGroup(bytes []byte) (cc.Message, error) {
	message := cc.NewMessage()
	message.IconURL = IconURL

	src, err := GosquaredOutMessageSmartGroupFromBytes(bytes)
	if err != nil {
		return message, err
	}

	verb := "exited"
	if src.Boundary == "enter" {
		verb = "entered"
	}

	message.Activity = fmt.Sprintf("User has %s Smart Group", verb)
	message.Title = fmt.Sprintf("%s has %s [%s](%s)",
		src.Person.Name,
		verb,
		src.Group.Name,
		src.GroupURL())
	return message, nil
}

func NormalizeLiveMessage(bytes []byte) (cc.Message, error) {
	message := cc.NewMessage()
	message.IconURL = IconURL

	src, err := GosquaredOutLiveMessageFromBytes(bytes)
	if err != nil {
		return message, err
	}

	message.Activity = "Live chat message"

	personInboxURL, errURL := src.PersonInboxURL()
	person := src.Person.DisplayName(false, true)
	if errURL == nil {
		person = fmt.Sprintf("[%v](%v)", person, personInboxURL)
	}

	message.Title = fmt.Sprintf("%v sent a message", person)

	attachment := cc.NewAttachment()

	if len(strings.TrimSpace(src.Message.Content)) > 0 {
		attachment.AddField(cc.Field{
			Title: "Message",
			Value: strings.TrimSpace(src.Message.Content)})
	}

	if src.Message.Timestamp > 0 {
		epochMilli := src.Message.Timestamp
		epochNano := epochMilli * 1000000
		dt := time.Unix(0, epochNano)
		attachment.AddField(cc.Field{
			Title: "Time",
			Value: dt.Format(time.RFC822)})
	}

	if len(attachment.Fields) > 0 {
		message.AddAttachment(attachment)
	}

	return message, nil
}

type GosquaredOutBaseMessage struct {
	TriggeredAlert GosquaredOutTriggeredAlert     `json:"triggeredAlert,omitempty"`
	Concurrents    int64                          `json:"concurrents,omitempty"`
	Person         GosquaredOutPerson             `json:"person,omitempty"`
	Message        GosquaredOutLiveMessageMessage `json:"message,omitempty"`
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

type GosquaredOutLiveMessage struct {
	Version   string                         `json:"version,omitempty"`
	SiteToken string                         `json:"site_token,omitempty"`
	Timestamp string                         `json:"timestamp,omitempty"`
	Person    GosquaredOutPerson             `json:"person,omitempty"`
	Message   GosquaredOutLiveMessageMessage `json:"message,omitempty"`
}

func (msg *GosquaredOutLiveMessage) PersonInboxURL() (string, error) {
	if len(strings.TrimSpace(msg.SiteToken)) == 0 || len(strings.TrimSpace(msg.Person.Id)) == 0 {
		return "", errors.New("Information missing for PersonInboxURL")
	}
	personIdEsc := url.QueryEscape(strings.TrimSpace(msg.Person.Id))
	personInboxURL := fmt.Sprintf(
		"https://www.gosquared.com/inbox/%v/inbox/%v",
		strings.TrimSpace(msg.SiteToken),
		personIdEsc)
	return personInboxURL, nil
}

/*
https://www.gosquared.com/inbox/GSN-345166-V/inbox/Anon%20Chat:%20bba4b6264b073a17c74f1b0da7720114

The rule to build the url is
https://www.gosquared.com/inbox/<site_token>/inbox/<person id encoded>.
*/

func GosquaredOutLiveMessageFromBytes(bytes []byte) (GosquaredOutLiveMessage, error) {
	msg := GosquaredOutLiveMessage{}
	err := json.Unmarshal(bytes, &msg)
	return msg, err
}

/*
{
    "version":"0.0.1",
    "site_token":"GSN-67890-A",
    "timestamp":"2017-05-02T05:19:34.252Z",
    "person":{
        "id":"Anon Chat: bba4b6264b073a17c74f1b0da7720114",
        "email":"",
        "name":"",
        "avatar":"",
        "chat":{
            "archived":null,
            "latest":{
                "message":{
                    "type":"message",
                    "id":"167c06cf-ff7c-4d09-880c-5533795f9673",
                    "content":"hello world",
                    "timestamp":1493702374252,
                    "from":"client",
                    "private":false,
                    "session":{
                        "title":"",
                        "href":"http://127.0.0.1:2015/"
                    },
*/

type GosquaredOutLiveMessageMessage struct {
	Type      string                         `json:"message,omitempty"`
	Id        string                         `json:"id,omitempty"`
	Content   string                         `json:"content,omitempty"`
	Timestamp int64                          `json:"timestamp,omitempty"`
	From      string                         `json:"from,omitempty"`
	Private   bool                           `json:"private,omitempty"`
	Session   GosquaredOutLiveMessageSession `json:"session,omitempty"`
	Data      interface{}                    `json:"data,omitempty"`
	Entities  []interface{}                  `json:"entities,omitempty"`
	PersonID  string                         `json:"person_id,omitempty"`
}

type GosquaredOutLiveMessageSession struct {
	Title string `json:"title,omitempty"`
	Href  string `json:"href,omitempty"`
}

/*
   "message":{
       "type":"message",
       "id":"167c06cf-ff7c-4d09-880c-5533795f9673",
       "content":"hello world",
       "timestamp":1493702374252,
       "from":"client",
       "private":false,
       "session":{
           "title":"",
           "href":"http://127.0.0.1:2015/"
       },
       "data":{

       },
       "entities":[

       ],
       "person_id":"Anon Chat: bba4b6264b073a17c74f1b0da7720114"
   }
*/

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
	CreatedAt   string `json:"person,omitempty"`
	Phone       string `json:"person,omitempty"`
	Avatar      string `json:"person,omitempty"`
	Description string `json:"description,omitempty"`
	Username    string `json:"username,omitempty"`
	Email       string `json:"email,omitempty"`
	Name        string `json:"name,omitempty"`
	Id          string `json:"id,omitempty"`
}

func (person *GosquaredOutPerson) DisplayName(extraEmail bool, anonymous bool) string {
	displayName := ""
	if len(strings.TrimSpace(person.Name)) > 0 {
		displayName = strings.TrimSpace(person.Name)
	} else if len(strings.TrimSpace(person.Username)) > 0 {
		displayName = strings.TrimSpace(person.Username)
	}
	if len(strings.TrimSpace(person.Email)) > 0 {
		if len(displayName) > 0 && extraEmail {
			displayName = fmt.Sprintf("%v (%v)", displayName, strings.TrimSpace(person.Email))
		} else {
			displayName = strings.TrimSpace(person.Email)
		}
	}
	if len(displayName) == 0 && anonymous {
		displayName = "Anonymous user"
	}
	return displayName
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
