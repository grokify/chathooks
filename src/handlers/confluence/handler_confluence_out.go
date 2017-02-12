package confluence

import (
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"

	"github.com/grokify/glip-go-webhook"
	"github.com/grokify/glip-webhook-proxy-go/src/config"
	"github.com/grokify/glip-webhook-proxy-go/src/util"
	"github.com/valyala/fasthttp"
)

const (
	DISPLAY_NAME = "Confluence"
	ICON_URL     = "https://wiki.ucop.edu/images/logo/default-space-logo-256.png"
)

// FastHttp request handler for Confluence outbound webhook
// https://developer.atlassian.com/static/connect/docs/beta/modules/common/webhook.html
type ConfluenceOutToGlipHandler struct {
	Config     config.Configuration
	GlipClient glipwebhook.GlipWebhookClient
}

// FastHttp request handler constructor for Confluence outbound webhook
func NewConfluenceOutToGlipHandler(cfg config.Configuration, glip glipwebhook.GlipWebhookClient) ConfluenceOutToGlipHandler {
	return ConfluenceOutToGlipHandler{Config: cfg, GlipClient: glip}
}

// HandleFastHTTP is the method to respond to a fasthttp request.
func (h *ConfluenceOutToGlipHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	srcMsg, err := BuildInboundMessage(ctx)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotAcceptable)
		log.WithFields(log.Fields{
			"type":   "http.response",
			"status": fasthttp.StatusNotAcceptable,
		}).Info("Confluence request is not acceptable.")
		return
	}
	glipMsg := Normalize(srcMsg)

	util.SendGlipWebhookCtx(ctx, h.GlipClient, glipMsg)
}

func BuildInboundMessage(ctx *fasthttp.RequestCtx) (ConfluenceOutMessage, error) {
	return ConfluenceOutMessageFromBytes(ctx.FormValue("payload"))
}

func Normalize(src ConfluenceOutMessage) glipwebhook.GlipWebhookMessage {
	glip := glipwebhook.GlipWebhookMessage{Icon: ICON_URL}
	if !src.IsComment() {
		if src.Page.IsCreated() {
			glip.Activity = fmt.Sprintf("%v created page in space %v (%v)", src.Page.CreatorName, src.Page.SpaceKey, DISPLAY_NAME)
			glip.Body = fmt.Sprintf("> [%v](%v)", src.Page.Title, src.Page.Self)
		} else {
			glip.Activity = fmt.Sprintf("%v updated page in space %v (%v)", src.Page.CreatorName, src.Page.SpaceKey, DISPLAY_NAME)
			glip.Body = fmt.Sprintf("> [%v](%v)", src.Page.Title, src.Page.Self)
		}
	} else {
		if src.Comment.IsCreated() {
			glip.Activity = fmt.Sprintf("%v commented on page in space %v (%v)", src.Comment.CreatorName, src.Comment.Parent.SpaceKey, DISPLAY_NAME)
			glip.Body = fmt.Sprintf("> [%v](%v)", src.Page.Title, src.Page.Self)
		} else {
			glip.Activity = fmt.Sprintf("%v updated comment on page in space %v (%v)", src.Page.CreatorName, src.Comment.Parent.SpaceKey, DISPLAY_NAME)
			glip.Body = fmt.Sprintf("> [%v](%v)", src.Page.Title, src.Page.Self)
		}
	}
	return glip
}

type ConfluenceOutMessage struct {
	User      string            `json:"user,omitempty"`
	UserKey   string            `json:"userKey,omitempty"`
	Timestamp int64             `json:"timestamp,omitempty"`
	Username  string            `json:"username,omitempty"`
	Page      ConfluencePage    `json:"page,omitempty"`
	Comment   ConfluenceComment `json:"comment,omitempty"`
}

func ConfluenceOutMessageFromBytes(bytes []byte) (ConfluenceOutMessage, error) {
	log.WithFields(log.Fields{
		"type":    "message.raw",
		"message": string(bytes),
	}).Debug(fmt.Sprintf("%v message.", DISPLAY_NAME))
	msg := ConfluenceOutMessage{}
	err := json.Unmarshal(bytes, &msg)
	if err != nil {
		log.WithFields(log.Fields{
			"type":  "message.json.unmarshal",
			"error": fmt.Sprintf("%v\n", err),
		}).Warn(fmt.Sprintf("%v request unmarshal failure.", DISPLAY_NAME))
	}
	if msg.IsComment() {
		msg.Page = msg.Comment.Parent
	}
	return msg, err
}

func (msg *ConfluenceOutMessage) IsComment() bool {
	if msg.Comment.ModificationDate > 0 {
		return true
	}
	return false
}

type ConfluencePage struct {
	SpaceKey         string `json:"spaceKey,omitempty"`
	ModificationDate int64  `json:"modificationDate,omitempty"`
	CreatorKey       string `json:"creatorKey,omitempty"`
	CreatorName      string `json:"creatorName,omitempty"`
	LastModifierKey  string `json:"lastModifierKey,omitempty"`
	Self             string `json:"self,omitempty"`
	LastModifierName string `json:"lastModifierName,omitempty"`
	Id               int64  `json:"id,omitempty"`
	Title            string `json:"title,omitempty"`
	CreationDate     int64  `json:"creationDate,omitempty"`
	Version          int64  `json:"version,omitempty"`
}

func (page *ConfluencePage) IsCreated() bool {
	if page.ModificationDate > 0 && page.ModificationDate == page.CreationDate {
		return true
	}
	return false
}

func (page *ConfluencePage) IsUpdated() bool {
	if page.IsCreated() {
		return false
	}
	return true
}

type ConfluenceComment struct {
	SpaceKey         string         `json:"spaceKey,omitempty"`
	Parent           ConfluencePage `json:"parent,omitempty"`
	ModificationDate int64          `json:"modificationDate,omitempty"`
	CreatorKey       string         `json:"creatorKey,omitempty"`
	CreatorName      string         `json:"creatorName,omitempty"`
	LastModifierKey  string         `json:"lastModifierKey,omitempty"`
	Self             string         `json:"self,omitempty"`
	LastModifierName string         `json:"lastModifierName,omitempty"`
	Id               int64          `json:"id,omitempty"`
	CreationDate     int64          `json:"creationDate,omitempty"`
	Version          int64          `json:"version,omitempty"`
}

func (comment *ConfluenceComment) IsCreated() bool {
	if comment.ModificationDate > 0 && comment.ModificationDate == comment.CreationDate {
		return true
	}
	return false
}

func (comment *ConfluenceComment) IsUpdated() bool {
	if comment.IsCreated() {
		return false
	}
	return true
}

/*

Confluence page_created

Activity: msg.Page.CreatorName created page in space [msg.pPge.SpaceKey]()
Body []()

{
   "page": {
     "spaceKey": "~admin",
     "modificationDate": 1471926079631,
     "creatorKey": "ff80808154510724015451074c160001",
     "creatorName": "admin",
     "lastModifierKey": "ff80808154510724015451074c160001",
     "self": "https://cloud-development-environment.atlassian.net/wiki/display/~admin/Some+random+test+page",
     "lastModifierName": "admin",
     "id": 16777227,
     "title": "Some random test page",
     "creationDate": 1471926079631,
     "version": 1
   },
   "user": "admin",
   "userKey": "ff80808154510724015451074c160001",
   "timestamp": 1471926079645,
   "username": "admin"
 }

Confluence comment_created

{
   "comment": {
     "spaceKey": "~admin",
     "parent": {
       "spaceKey": "~admin",
       "modificationDate": 1471926079631,
       "creatorKey": "ff80808154510724015451074c160001",
       "creatorName": "admin",
       "lastModifierKey": "ff80808154510724015451074c160001",
       "self": "https://cloud-development-environment.atlassian.net/wiki/display/~admin/Some+random+test+page",
       "lastModifierName": "admin",
       "id": 16777227,
       "title": "Some random test page",
       "creationDate": 1471926079631,
       "version": 1
     },
     "modificationDate": 1471926091465,
     "creatorKey": "ff80808154510724015451074c160001",
     "creatorName": "admin",
     "lastModifierKey": "ff80808154510724015451074c160001",
     "self": "https://cloud-development-environment/wiki/display/~admin/Some+random+test+page?focusedCommentId=16777228#comment-16777228",
     "lastModifierName": "admin",
     "id": 16777228,
     "creationDate": 1471926091465,
     "version": 1
   },
   "user": "admin",
   "userKey": "ff80808154510724015451074c160001",
   "timestamp": 1471926091468,
   "username": "admin"
 }

*/
