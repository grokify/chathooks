package userlike

import (
	"encoding/json"
	"errors"
	"fmt"
	//"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/grokify/glip-go-webhook"
	"github.com/grokify/glip-webhook-proxy-go/config"
	"github.com/grokify/glip-webhook-proxy-go/util"
	"github.com/valyala/fasthttp"
)

const (
	DISPLAY_NAME = "Userlike"
	ICON_URL     = "https://a.slack-edge.com/ae7f/img/services/userlike_512.png"
)

// FastHttp request handler for Userlike outbound webhook
type UserlikeOutToGlipHandler struct {
	Config     config.Configuration
	GlipClient glipwebhook.GlipWebhookClient
}

// FastHttp request handler constructor for Userlike outbound webhook
func NewUserlikeOutToGlipHandler(cfg config.Configuration, glip glipwebhook.GlipWebhookClient) UserlikeOutToGlipHandler {
	return UserlikeOutToGlipHandler{Config: cfg, GlipClient: glip}
}

// HandleFastHTTP is the method to respond to a fasthttp request.
func (h *UserlikeOutToGlipHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	glipMsg, err := BuildGlipMessageFromBytes(ctx.PostBody())

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotAcceptable)
		log.WithFields(log.Fields{
			"type":   "http.response",
			"status": fasthttp.StatusNotAcceptable,
		}).Info(fmt.Sprintf("%v request is not acceptable.", DISPLAY_NAME))
		return
	}

	util.SendGlipWebhookCtx(ctx, h.GlipClient, glipMsg)
}

func BuildGlipMessageFromBytes(bodyBytes []byte) (glipwebhook.GlipWebhookMessage, error) {
	srcMsgBase, err := UserlikeBaseOutMessageFromBytes(bodyBytes)
	if err != nil {
		return glipwebhook.GlipWebhookMessage{}, err
	}
	if srcMsgBase.Type == "offline_message" && srcMsgBase.Event == "receive" {
		srcMsg, err := UserlikeOfflineMessageReceiveOutMessageFromBytes(bodyBytes)
		if err != nil {
			return glipwebhook.GlipWebhookMessage{}, err
		}
		return NormalizeOfflineMessageReceive(srcMsg), nil
	} else if srcMsgBase.Type == "chat_meta" && srcMsgBase.Event == "start" {
		srcMsg, err := UserlikeChatMetaStartOutMessageFromBytes(bodyBytes)
		if err != nil {
			return glipwebhook.GlipWebhookMessage{}, err
		}
		return NormalizeChatMetaStart(srcMsg), nil
	}
	return glipwebhook.GlipWebhookMessage{}, errors.New("Type Not Supported")
}

/*
func Normalize(src UserlikeBaseOutMessage) (glipwebhook.GlipWebhookMessage, error) {
	if src.Type == "chat_meta" && src.Event == "start" {
		return NormalizeChatMetaStart(src)
	}
	if src.Type == "offline_message" && src.Event == "receive" {
		return NormalizeOfflineMessageReceive(src)
	}

	gmsg := glipwebhook.GlipWebhookMessage{
		Activity: DISPLAY_NAME,
		Icon:     ICON_URL}

	if strings.ToLower(strings.TrimSpace(src.Event)) == "build" {
		// Joe Cool build #15 passed
		gmsg.Activity = fmt.Sprintf("%v's %v #%v %v", src.Commit.AuthorName, src.Event, src.BuildNumber, src.Result)
	} else {
		gmsg.Activity = fmt.Sprintf("%v's %v %v", src.Commit.AuthorName, src.Event, src.Result)
	}

	lines := []string{}
	if len(src.Commit.Message) > 0 {
		lines = append(lines, fmt.Sprintf("> %v", src.Commit.Message))
	}
	if len(src.BuildURL) > 0 {
		lines = append(lines, fmt.Sprintf("> [View details](%v)", src.BuildURL))
	}
	if len(lines) > 0 {
		gmsg.Body = strings.Join(lines, "\n")
	}
	return gmsg
}
*/
func NormalizeOfflineMessageReceive(src UserlikeOfflineMessageReceiveOutMessage) glipwebhook.GlipWebhookMessage {
	activitySuffix := " from Website Visitor"
	if len(src.ClientName) > 0 {
		activitySuffix = fmt.Sprintf(" from %v", src.ClientName)
	}
	gmsg := glipwebhook.GlipWebhookMessage{
		Activity: fmt.Sprintf("Received a new Offline Message%v (%v)", activitySuffix, DISPLAY_NAME),
		Icon:     ICON_URL}
	if len(src.URL) > 0 {
		gmsg.Body = fmt.Sprintf("> [View conversation](%v)", src.URL)
	}
	return gmsg
}

func NormalizeChatMetaStart(src UserlikeChatMetaStartOutMessage) glipwebhook.GlipWebhookMessage {
	activitySuffix := " from Website Visitor"
	if len(src.ClientName) > 0 {
		activitySuffix = fmt.Sprintf(" from %v", src.ClientName)
	}
	gmsg := glipwebhook.GlipWebhookMessage{
		Activity: fmt.Sprintf("Received a new Chat%v (%v)", activitySuffix, DISPLAY_NAME),
		Icon:     ICON_URL}
	if len(src.URL) > 0 {
		gmsg.Body = fmt.Sprintf("> [View conversation](%v)", src.URL)
	}
	return gmsg
}

type UserlikeBaseOutMessage struct {
	Event string `json:"_event,omitempty"`
	Type  string `json:"_type,omitempty"`
}

func UserlikeBaseOutMessageFromBytes(bytes []byte) (UserlikeBaseOutMessage, error) {
	msg := UserlikeBaseOutMessage{}
	err := json.Unmarshal(bytes, &msg)
	return msg, err
}

type UserlikeOfflineMessageReceiveOutMessage struct {
	UserlikeBaseOutMessage
	BrowserName     string             `json:"browser_name,omitempty"`
	BrowserOS       string             `json:"browser_os,omitempty"`
	BrowserVersion  string             `json:"browser_version,omitempty"`
	ChatWidget      UserlikeChatWidget `json:"chat_widget,omitempty"`
	ClientEmail     string             `json:"client_email,omitempty"`
	ClientName      string             `json:"client_name,omitempty"`
	CreatedAt       string             `json:"created_at,omitempty"`
	Custom          interface{}        `json:"custom,omitempty"`
	DataPrivacy     interface{}        `json:"data_privacy,omitempty"`
	Id              int64              `json:"id,omitempty"`
	LocCity         string             `json:"loc_city,omitempty"`
	LocCountry      string             `json:"loc_country,omitempty"`
	LocLat          float64            `json:"loc_lat,omitempty"`
	LocLon          float64            `json:"loc_lon,omitempty"`
	MarkedRead      bool               `json:"marked_read,omitempty"`
	Message         string             `json:"message,omitempty"`
	PageImpressions int64              `json:"page_impresions,omitempty"`
	ScreenshotOID   string             `json:"screenshot_oid,omitempty"`
	ScreenshotURL   string             `json:"screenshot_url,omitempty"`
	Status          string             `json:"status,omitempty"`
	Topic           string             `json:"topic,omitempty"`
	URL             string             `json:"url,omitempty"`
	Visits          int64              `json:"visits,omitempty"`
}

func UserlikeOfflineMessageReceiveOutMessageFromBytes(bytes []byte) (UserlikeOfflineMessageReceiveOutMessage, error) {
	msg := UserlikeOfflineMessageReceiveOutMessage{}
	err := json.Unmarshal(bytes, &msg)
	return msg, err
}

type UserlikeChatWidget struct {
	Id   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type UserlikeChatMetaStartOutMessage struct {
	UserlikeBaseOutMessage
	BrowserName       string             `json:"browser_name,omitempty"`
	BrowserOS         string             `json:"browser_os,omitempty"`
	BrowserVersion    string             `json:"browser_version,omitempty"`
	ChatWidget        UserlikeChatWidget `json:"chat_widget,omitempty"`
	ClientEmail       string             `json:"client_email,omitempty"`
	ClientName        string             `json:"client_name,omitempty"`
	ClientUUID        string             `json:"client_uuid,omitempty"`
	CreatedAt         string             `json:"created_at,omitempty"`
	DataPrivacy       bool               `json:"data_privacy,omitempty"`
	Duration          string             `json:"duration,omitempty"`
	EndedAt           string             `json:"ended_at,omitempty"`
	FeedbackMessage   string             `json:"feedback_message,omitempty"`
	Id                int64              `json:"id,omitempty"`
	InitialURL        string             `json:"initial_url,omitempty"`
	LocCity           string             `json:"browser_version,omitempty"`
	LocCountry        string             `json:"browser_version,omitempty"`
	LocLat            float64            `json:"loc_lat,omitempty"`
	LocLon            float64            `json:"loc_lon,omitempty"`
	MarkedRead        bool               `json:"marked_read,omitempty"`
	OperatorCurrentId int64              `json:"operator_current_id,omitempty"`
	PageImpressions   int64              `json:"page_impressions,omitempty"`
	PostSurveyOption  string             `json:"post_survey_option,omitempty"`
	Rate              string             `json:"rate,omitempty"`
	Referrer          string             `json:"referrer,omitempty"`
	Status            string             `json:"status,omitempty"`
	Topic             string             `json:"topic,omitempty"`
	URL               string             `json:"url,omitempty"`
	Visits            int64              `json:"visits,omitempty"`
	WasProactive      bool               `json:"was_proactive,omitempty"`
}

func UserlikeChatMetaStartOutMessageFromBytes(bytes []byte) (UserlikeChatMetaStartOutMessage, error) {
	msg := UserlikeChatMetaStartOutMessage{}
	err := json.Unmarshal(bytes, &msg)
	return msg, err
}

/*

{
 "_event": "start",
 "_type": "chat_meta",
 "browser_name": "Safari",
 "browser_os": "Mac OS X",
 "browser_version": "8",
 "chat_widget": {
 "id": 9,
 "name": "Testing David"
 },
 "chat_widget_goal": {
 "id": null,
 "name": null
 },
 "client_additional01_name": null,
 "client_additional01_value": null,
 "client_additional02_name": null,
 "client_additional02_value": null,
 "client_additional03_name": null,
 "client_additional03_value": null,
 "client_email": "david@optixx.org",
 "client_name": "Jo",
 "client_uuid": "nEitxHDooFzsPhMoVn0QCw8E.L3mmogIcp+FGwXos5a4NUdZ9/uQbSCBx0wDIRVFWM+o",
 "created_at": "2014-12-29 11:26:24",
 "custom": {
    "basket": {
      "item01": {
        "desc": "33X Optical Zoom Camcorder Mini DV",
        "id": "2acefe58-91e5-11e1-beba-000c2979313a",
        "price": 139.99,
        "url": "http://application/en/electronics/34-camcorder.html"
      },
      "item02": {
        "desc": "Home Theater System",
        "id": "31aca2f2-91e5-11e1-beba-000c2979313a",
        "long": "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.",
        "price": 499.99,
        "url": "https://application/en/electronics/39-home-theater.html"
      }
    },
    "id": "428614f0-91e5-11e1-beba-000c2979313a",
    "ref": "3efd5462e"
  },
 "data_privacy": false,
 "duration": "00:00:13",
 "ended_at": null,
 "feedback_message": null,
 "id": 71,
 "inital_url": "https://devel.userlike.local/en/",
 "loc_city": null,
 "loc_country": null,
 "loc_lat": null,
 "loc_lon": null,
 "marked_read": false,
 "messages": [],
  "notes": [],
  "operator_created": {
    "email": "david@userlike.com",
    "first_name": "David",
    "id": 5,
    "last_name": "Voswinkel",
    "name": "David Voswinkel",
    "operator_group": {
      "id": 14,
      "name": "Testing David"
    }
  },
  "operator_created_id": 5,
  "operator_current": {
    "email": "david@userlike.com",
    "first_name": "David",
    "id": 5,
    "last_name": "Voswinkel",
    "name": "David Voswinkel",
    "operator_group": {
      "id": 14,
      "name": "Testing David"
    }
   },
  "operator_current_id": 5,
  "page_impressions": 2,
  "post_survey_option": null,
  "rate": null,
  "referrer": null,
  "status": "new",
  "topic": null,
  "url": "https://devel.userlike.local/en/debug/9",
  "visits": 10,
  "was_proactive": false
}

{
 "_event": "receive",
 "_type": "offline_message",
 "browser_name": "Chrome",
 "browser_os": "Mac OS X",
 "browser_version": "32",
 "chat_widget": {
   "id": 2,
   "name": "Website"
 },
 "client_email": "support@userlike.com",
 "client_name": "Userlike Support",
 "created_at": "2014-12-20 14:50:23",
 "custom": {},
 "data_privacy": null,
 "id": 3,
 "loc_city": "Cologne",
 "loc_country": "Germany",
 "loc_lat": 50.9333000183105,
 "loc_lon": 6.94999980926514,
 "marked_read": true,
 "message": "We are happy to welcome you as a Userlike user!",
 "page_impressions": 5,
 "screenshot_oid": null,
 "screenshot_url": null,
 "status": "new",
 "topic": "Support",
 "url": "http://www.userlike.com",
 "visits": 1
}

*/
