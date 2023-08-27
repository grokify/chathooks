package gosquared

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	cc "github.com/grokify/commonchat"

	"github.com/grokify/chathooks/pkg/config"
	"github.com/grokify/chathooks/pkg/handlers"
	"github.com/grokify/chathooks/pkg/models"
)

const (
	DisplayName      = "GoSquared"
	HandlerKey       = "gosquared"
	DocumentationURL = "https://www.gosquared.com/customer/en/portal/articles/1996494-webhooks"
	MessageBodyType  = models.JSON
)

func NewHandler() handlers.Handler {
	return handlers.Handler{MessageBodyType: MessageBodyType, Normalize: Normalize}
}

func Normalize(cfg config.Configuration, hReq handlers.HandlerRequest) (cc.Message, error) {
	bytes := hReq.Body
	src, err := GosquaredOutBaseMessageFromBytes(bytes)
	if err != nil {
		return cc.NewMessage(), err
	}
	if len(src.Message.ID) > 0 {
		return NormalizeLiveMessage(cfg, bytes)
	} else if len(src.Person.ID) > 0 {
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

	if src.TriggeredAlert.Boundary == "upper" {
		ccMsg.Activity = "Site traffic spike"
	} else { // if src.TriggeredAlert.Boundary == "lower" {
		ccMsg.Activity = "Site traffic dip"
	}

	pluralSuffix := "s"
	if src.Concurrents == int64(1) {
		pluralSuffix = ""
	}

	ccMsg.Title = fmt.Sprintf("[%s](%s) has [%v visitor%s online](%s)",
		src.SiteDetails.SiteName,
		src.SiteDetails.URL,
		src.Concurrents,
		pluralSuffix,
		src.SiteDetails.DashboardURL())

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

	verb := "exited"
	if src.Boundary == "enter" {
		verb = "entered"
	}

	ccMsg.Activity = fmt.Sprintf("User has %s Smart Group", verb)
	ccMsg.Title = fmt.Sprintf("%s has %s [%s](%s)",
		src.Person.Name,
		verb,
		src.Group.Name,
		src.GroupURL())
	return ccMsg, nil
}

func NormalizeLiveMessage(cfg config.Configuration, bytes []byte) (cc.Message, error) {
	ccMsg := cc.NewMessage()
	iconURL, err := cfg.GetAppIconURL(HandlerKey)
	if err == nil {
		ccMsg.IconURL = iconURL.String()
	}

	src, err := GosquaredOutLiveMessageFromBytes(bytes)
	if err != nil {
		return ccMsg, err
	}

	ccMsg.Activity = "Live chat message"

	personInboxURL, errURL := src.PersonInboxURL()
	person := src.Person.DisplayName(false, true)
	if errURL == nil {
		person = fmt.Sprintf("[%v](%v)", person, personInboxURL)
	}

	ccMsg.Title = fmt.Sprintf("%v sent a message", person)

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
		ccMsg.AddAttachment(attachment)
	}

	return ccMsg, nil
}

type GosquaredOutBaseMessage struct {
	TriggeredAlert GosquaredOutTriggeredAlert     `json:"triggeredAlert,omitempty"`
	Concurrents    int64                          `json:"concurrents,omitempty"`
	Person         GosquaredOutPerson             `json:"person,omitempty"`
	Message        GosquaredOutLiveMessageMessage `json:"message,omitempty"`
}

func GosquaredOutBaseMessageFromBytes(bytes []byte) (GosquaredOutBaseMessage, error) {
	var msg GosquaredOutBaseMessage
	return msg, json.Unmarshal(bytes, &msg)
}

type GosquaredOutMessageSiteTraffic struct {
	TriggeredAlert GosquaredOutTriggeredAlert `json:"triggeredAlert,omitempty"`
	SiteDetails    GosquaredOutSiteDetails    `json:"siteDetails,omitempty"`
	Concurrents    int64                      `json:"concurrents,omitempty"`
}

func GosquaredOutMessageSiteTrafficFromBytes(bytes []byte) (GosquaredOutMessageSiteTraffic, error) {
	var msg GosquaredOutMessageSiteTraffic
	return msg, json.Unmarshal(bytes, &msg)
}

type GosquaredOutLiveMessage struct {
	Version   string                         `json:"version,omitempty"`
	SiteToken string                         `json:"site_token,omitempty"`
	Timestamp string                         `json:"timestamp,omitempty"`
	Person    GosquaredOutPerson             `json:"person,omitempty"`
	Message   GosquaredOutLiveMessageMessage `json:"message,omitempty"`
}

func (msg *GosquaredOutLiveMessage) PersonInboxURL() (string, error) {
	if len(strings.TrimSpace(msg.SiteToken)) == 0 || len(strings.TrimSpace(msg.Person.ID)) == 0 {
		return "", errors.New("missing information for PersonInboxURL")
	}
	personIDEsc := url.QueryEscape(strings.TrimSpace(msg.Person.ID))
	personInboxURL := fmt.Sprintf(
		"https://www.gosquared.com/inbox/%v/inbox/%v",
		strings.TrimSpace(msg.SiteToken),
		personIDEsc)
	return personInboxURL, nil
}

/*
https://www.gosquared.com/inbox/GSN-345166-V/inbox/Anon%20Chat:%20bba4b6264b073a17c74f1b0da7720114

The rule to build the url is
https://www.gosquared.com/inbox/<site_token>/inbox/<person id encoded>.
*/

func GosquaredOutLiveMessageFromBytes(bytes []byte) (GosquaredOutLiveMessage, error) {
	var msg GosquaredOutLiveMessage
	return msg, json.Unmarshal(bytes, &msg)
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
	ID        string                         `json:"id,omitempty"`
	Type      string                         `json:"message,omitempty"`
	Content   string                         `json:"content,omitempty"`
	Timestamp int64                          `json:"timestamp,omitempty"`
	From      string                         `json:"from,omitempty"`
	Private   bool                           `json:"private,omitempty"`
	Session   GosquaredOutLiveMessageSession `json:"session,omitempty"`
	Data      any                            `json:"data,omitempty"`
	Entities  []any                          `json:"entities,omitempty"`
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
	var msg GosquaredOutMessageSmartGroup
	return msg, json.Unmarshal(bytes, &msg)
}

type GosquaredOutGroup struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func (msg *GosquaredOutMessageSmartGroup) GroupURL() string {
	// https://www.gosquared.com/people/GSN-466237-B/last-seen-1-day
	return fmt.Sprintf("https://www.gosquared.com/people/%s/%s",
		msg.SiteToken, msg.Group.ID)
}

type GosquaredOutPerson struct {
	ID          string `json:"id,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	Phone       string `json:"phone,omitempty"`
	Avatar      string `json:"avatar,omitempty"`
	Description string `json:"description,omitempty"`
	Username    string `json:"username,omitempty"`
	Email       string `json:"email,omitempty"`
	Name        string `json:"name,omitempty"`
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
	ID       int64  `json:"id,omitempty"`
	Boundary string `json:"boundary,omitempty"`
	Value    string `json:"value,omitempty"`
	Type     string `json:"type,omitempty"`
}

type GosquaredOutSiteDetails struct {
	UserID    int64  `json:"user_id,omitempty"`
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
