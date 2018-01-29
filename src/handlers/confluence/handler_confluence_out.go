package confluence

import (
	"encoding/json"
	"fmt"

	"github.com/grokify/chathooks/src/config"
	"github.com/grokify/chathooks/src/handlers"
	"github.com/grokify/chathooks/src/models"
	cc "github.com/grokify/commonchat"
	log "github.com/sirupsen/logrus"
)

const (
	DisplayName      = "Confluence"
	HandlerKey       = "confluence"
	MessageDirection = "out"
	MessageBodyType  = models.JSON
)

func NewHandler() handlers.Handler {
	return handlers.Handler{MessageBodyType: MessageBodyType, Normalize: Normalize}
}

func Normalize(cfg config.Configuration, bytes []byte) (cc.Message, error) {
	ccMsg := cc.NewMessage()
	iconURL, err := cfg.GetAppIconURL(HandlerKey)
	if err == nil {
		ccMsg.IconURL = iconURL.String()
	}

	src, err := ConfluenceOutMessageFromBytes(bytes)
	if err != nil {
		return ccMsg, err
	}

	if !src.IsComment() {
		if src.Page.IsCreated() {
			ccMsg.Activity = fmt.Sprintf("%v created page", src.Page.CreatorName)
		} else {
			ccMsg.Activity = fmt.Sprintf("%v updated page", src.Page.CreatorName)
		}
	} else {
		if src.Comment.IsCreated() {
			ccMsg.Activity = fmt.Sprintf("%v commented on page", src.Comment.CreatorName)
		} else {
			ccMsg.Activity = fmt.Sprintf("%v updated comment on page", src.Comment.CreatorName)
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

	ccMsg.AddAttachment(attachment)
	return ccMsg, nil
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
