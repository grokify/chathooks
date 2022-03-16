package confluence

import (
	"encoding/json"
	"fmt"

	cc "github.com/grokify/commonchat"
	"github.com/rs/zerolog/log"

	"github.com/grokify/chathooks/pkg/config"
	"github.com/grokify/chathooks/pkg/handlers"
	"github.com/grokify/chathooks/pkg/models"
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

func Normalize(cfg config.Configuration, hReq handlers.HandlerRequest) (cc.Message, error) {
	ccMsg := cc.NewMessage()
	iconURL, err := cfg.GetAppIconURL(HandlerKey)
	if err == nil {
		ccMsg.IconURL = iconURL.String()
	}

	src, err := ConfluenceOutMessageFromBytes(hReq.Body)
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
	log.Debug().
		Str("type", "message.raw").
		Str("handler", HandlerKey).
		Str("request_body", string(bytes)).
		Msg(config.InfoInputMessageParseBegin)

	msg := ConfluenceOutMessage{}
	err := json.Unmarshal(bytes, &msg)
	if err != nil {
		log.Warn().
			Err(err).
			Str("type", "message.json.unmarshal").
			Str("handler", HandlerKey).
			Msg(config.ErrorInputMessageParseFailed)
	}
	if msg.IsComment() {
		msg.Page = msg.Comment.Parent
	}
	return msg, err
}

func (msg *ConfluenceOutMessage) IsComment() bool {
	return msg.Comment.ModificationDate > 0
}

type ConfluencePage struct {
	SpaceKey         string `json:"spaceKey,omitempty"`
	ModificationDate int64  `json:"modificationDate,omitempty"`
	CreatorKey       string `json:"creatorKey,omitempty"`
	CreatorName      string `json:"creatorName,omitempty"`
	LastModifierKey  string `json:"lastModifierKey,omitempty"`
	Self             string `json:"self,omitempty"`
	LastModifierName string `json:"lastModifierName,omitempty"`
	ID               int64  `json:"id,omitempty"`
	Title            string `json:"title,omitempty"`
	CreationDate     int64  `json:"creationDate,omitempty"`
	Version          int64  `json:"version,omitempty"`
}

func (page *ConfluencePage) IsCreated() bool {
	return page.ModificationDate > 0 && page.ModificationDate == page.CreationDate
}

func (page *ConfluencePage) IsUpdated() bool {
	return !page.IsCreated()
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
	return comment.ModificationDate > 0 && comment.ModificationDate == comment.CreationDate
}

func (comment *ConfluenceComment) IsUpdated() bool {
	return !comment.IsCreated()
}
