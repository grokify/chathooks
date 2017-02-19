package confluence

import (
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"

	cc "github.com/grokify/commonchat"
	"github.com/grokify/webhook-proxy-go/src/adapters"
	"github.com/grokify/webhook-proxy-go/src/config"
	"github.com/grokify/webhook-proxy-go/src/util"
	"github.com/valyala/fasthttp"
)

const (
	DisplayName = "Confluence"
	HandlerKey  = "confluence"
	IconURL     = "https://wiki.ucop.edu/images/logo/default-space-logo-256.png"
)

// FastHttp request handler for Confluence outbound webhook
// https://developer.atlassian.com/static/connect/docs/beta/modules/common/webhook.html
type ConfluenceOutToGlipHandler struct {
	Config  config.Configuration
	Adapter adapters.Adapter
}

// FastHttp request handler constructor for Confluence outbound webhook
func NewConfluenceOutToGlipHandler(cfg config.Configuration, adapter adapters.Adapter) ConfluenceOutToGlipHandler {
	return ConfluenceOutToGlipHandler{Config: cfg, Adapter: adapter}
}

// HandleFastHTTP is the method to respond to a fasthttp request.
func (h *ConfluenceOutToGlipHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	ccMsg, err := Normalize(ctx.FormValue("payload"))

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
	message := cc.NewMessage()
	message.IconURL = IconURL

	src, err := ConfluenceOutMessageFromBytes(bytes)
	if err != nil {
		return message, err
	}

	if !src.IsComment() {
		if src.Page.IsCreated() {
			message.Activity = fmt.Sprintf("%v created page", src.Page.CreatorName)
		} else {
			message.Activity = fmt.Sprintf("%v updated page", src.Page.CreatorName)
		}
	} else {
		if src.Comment.IsCreated() {
			message.Activity = fmt.Sprintf("%v commented on page", src.Comment.CreatorName)
		} else {
			message.Activity = fmt.Sprintf("%v updated comment on page", src.Comment.CreatorName)
		}
	}

	attachment := cc.NewAttachment()

	if len(src.Page.Title) > 0 && len(src.Page.Self) > 0 {
		attachment.AddField(cc.Field{
			Title: "Page",
			Value: fmt.Sprintf("[%v](%v)", src.Page.Title, src.Page.Self),
			Short: true})
	}
	if len(src.Page.SpaceKey) > 0 {
		field := cc.Field{Title: "Space", Short: true}
		if src.IsComment() {
			field.Value = src.Comment.Parent.SpaceKey
		} else {
			field.Value = src.Page.SpaceKey
		}
		attachment.AddField(field)
	}

	message.AddAttachment(attachment)
	return message, nil
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
	}).Debug(fmt.Sprintf("%v message.", DisplayName))
	msg := ConfluenceOutMessage{}
	err := json.Unmarshal(bytes, &msg)
	if err != nil {
		log.WithFields(log.Fields{
			"type":  "message.json.unmarshal",
			"error": fmt.Sprintf("%v\n", err),
		}).Warn(fmt.Sprintf("%v request unmarshal failure.", DisplayName))
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
