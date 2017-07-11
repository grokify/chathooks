package semaphore

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"

	cc "github.com/commonchat/commonchat-go"
	"github.com/grokify/gotilla/strings/stringsutil"
	"github.com/grokify/webhookproxy/src/adapters"
	"github.com/grokify/webhookproxy/src/config"
	"github.com/grokify/webhookproxy/src/util"
	"github.com/valyala/fasthttp"
)

const (
	DisplayName      = "Semaphore"
	HandlerKey       = "semaphore"
	MessageDirection = "out"
)

// FastHttp request handler for outbound webhook
type Handler struct {
	Config  config.Configuration
	Adapter adapters.Adapter
}

// FastHttp request handler constructor for Soutbound webhook
func NewHandler(cfg config.Configuration, adapter adapters.Adapter) Handler {
	return Handler{Config: cfg, Adapter: adapter}
}

func (h Handler) HandlerKey() string {
	return HandlerKey
}

func (h Handler) MessageDirection() string {
	return MessageDirection
}

// HandleFastHTTP is the method to respond to a fasthttp request.
func (h Handler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
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

//func NormalizeBytes(bytes []byte) (glipwebhook.GlipWebhookMessage, error) {
func Normalize(cfg config.Configuration, bytes []byte) (cc.Message, error) {
	ccMsg := cc.NewMessage()
	iconURL, err := cfg.GetAppIconURL(HandlerKey)
	if err == nil {
		ccMsg.IconURL = iconURL.String()
	}

	baseMsg, err := SemaphoreciBaseOutMessageFromBytes(bytes)
	if err != nil {
		return ccMsg, err
	}

	switch baseMsg.Event {
	case "build":
		srcMsg, err := SemaphoreciBuildOutMessageFromBytes(bytes)
		if err != nil {
			return ccMsg, err
		}
		return NormalizeSemaphoreciBuildOutMessage(cfg, srcMsg), nil
	case "deploy":
		srcMsg, err := SemaphoreciDeployOutMessageFromBytes(bytes)
		if err != nil {
			return ccMsg, err
		}
		return NormalizeSemaphoreciDeployOutMessage(cfg, srcMsg), nil
	}
	return cc.Message{IconURL: ""}, errors.New("EventNotFound")
}

func NormalizeSemaphoreciBuildOutMessage(cfg config.Configuration, src SemaphoreciBuildOutMessage) cc.Message {
	ccMsg := cc.NewMessage()
	iconURL, err := cfg.GetAppIconURL(HandlerKey)
	if err == nil {
		ccMsg.IconURL = iconURL.String()
	}

	ccMsg.Activity = fmt.Sprintf("%v %v %v", src.ProjectName, src.Event, src.Result)

	ccMsg.Title = fmt.Sprintf("[%v #%v](%v) for **%v/%v** %v ([%v](%v))",
		stringsutil.ToUpperFirst(src.Event),
		src.BuildNumber,
		src.BuildURL,
		src.ProjectName,
		src.BranchName,
		src.Result,
		src.Commit.Id[:7],
		src.Commit.URL)

	attachment := cc.NewAttachment()

	if len(src.Commit.Message) > 0 {
		attachment.AddField(cc.Field{
			Title: "Message",
			Value: src.Commit.Message,
			Short: true})
	}
	if 1 == 0 {
		if len(src.ProjectName) > 0 {
			attachment.AddField(cc.Field{
				Title: "Project",
				Value: src.ProjectName,
				Short: true})
		}
		if len(src.BranchName) > 0 {
			attachment.AddField(cc.Field{
				Title: "Branch",
				Value: src.BranchName,
				Short: true})
		}
		if len(src.Event) > 0 {
			attachment.AddField(cc.Field{
				Title: "Event",
				Value: src.Event,
				Short: true})
		}
	}
	if len(src.Commit.AuthorName) > 0 {
		attachment.AddField(cc.Field{
			Title: "Committer",
			Value: src.Commit.AuthorName,
			Short: true})
	}

	ccMsg.AddAttachment(attachment)
	return ccMsg
}

func NormalizeSemaphoreciDeployOutMessage(cfg config.Configuration, src SemaphoreciDeployOutMessage) cc.Message {
	ccMsg := cc.NewMessage()
	iconURL, err := cfg.GetAppIconURL(HandlerKey)
	if err == nil {
		ccMsg.IconURL = iconURL.String()
	}

	ccMsg.Activity = fmt.Sprintf("%v %v %v", src.ProjectName, src.Event, src.Result)

	ccMsg.Title = fmt.Sprintf("[%v #%v](%v) for **%v/%v** %v ([%v](%v))",
		stringsutil.ToUpperFirst(src.Event),
		src.Number, src.HtmlURL,
		src.ProjectName,
		src.BranchName,
		src.Result,
		src.Commit.Id[:7],
		src.Commit.URL)

	attachment := cc.NewAttachment()

	if len(src.Commit.Message) > 0 {
		attachment.AddField(cc.Field{
			Title: "Message",
			Value: src.Commit.Message})
	}
	if 1 == 0 {
		if len(src.ProjectName) > 0 {
			attachment.AddField(cc.Field{
				Title: "Project",
				Value: src.ProjectName,
				Short: true})
		}
		if len(src.BranchName) > 0 {
			attachment.AddField(cc.Field{
				Title: "Branch",
				Value: src.BranchName,
				Short: true})
		}
		if len(src.Event) > 0 {
			attachment.AddField(cc.Field{
				Title: "Event",
				Value: src.Event,
				Short: true})
		}
	}
	if len(src.Commit.AuthorName) > 0 {
		attachment.AddField(cc.Field{
			Title: "Committer",
			Value: src.Commit.AuthorName,
			Short: true})
	}

	ccMsg.AddAttachment(attachment)
	return ccMsg
}

type SemaphoreciBaseOutMessage struct {
	Event string `json:"event,omitempty"`
}

func SemaphoreciBaseOutMessageFromBytes(bytes []byte) (SemaphoreciBaseOutMessage, error) {
	msg := SemaphoreciBaseOutMessage{}
	err := json.Unmarshal(bytes, &msg)
	return msg, err
}

type SemaphoreciBuildOutMessage struct {
	BranchName    string            `json:"branch_name,omitempty"`
	BranchURL     string            `json:"branch_url,omitempty"`
	ProjectName   string            `json:"project_name,omitempty"`
	ProjectHashId string            `json:"project_hash_id,omitempty"`
	BuildURL      string            `json:"build_url,omitempty"`
	BuildNumber   int64             `json:"build_number,omitempty"`
	Result        string            `json:"result,omitempty"`
	Event         string            `json:"event,omitempty"`
	StartedAt     string            `json:"started_at,omitempty"`
	FinishedAt    string            `json:"finished_at,omitempty"`
	Commit        SemaphoreciCommit `json:"commit,omitempty"`
}

type SemaphoreciCommit struct {
	Id          string `json:"id,omitempty"`
	URL         string `json:"url,omitempty"`
	AuthorName  string `json:"author_name,omitempty"`
	AuthorEmail string `json:"author_email,omitempty"`
	Message     string `json:"message,omitempty"`
	Timestamp   string `json:"timestamp,omitempty"`
}

func SemaphoreciBuildOutMessageFromBytes(bytes []byte) (SemaphoreciBuildOutMessage, error) {
	msg := SemaphoreciBuildOutMessage{}
	err := json.Unmarshal(bytes, &msg)
	if err == nil {
		msg.Commit.Message = strings.ToLower(strings.TrimSpace(msg.Commit.Message))
	}
	return msg, err
}

type SemaphoreciDeployOutMessage struct {
	ProjectName   string            `json:"project_name,omitempty"`
	ProjectHashId string            `json:"project_hash_id,omitempty"`
	Result        string            `json:"result,omitempty"`
	Event         string            `json:"event,omitempty"`
	ServerName    string            `json:"server_name,omitempty"`
	Number        int64             `json:"number,omitempty"`
	CreatedAt     string            `json:"created_at,omitempty"`
	UpdatedAt     string            `json:"updated_at,omitempty"`
	StartedAt     string            `json:"started_at,omitempty"`
	FinishedAt    string            `json:"finished_at,omitempty"`
	HtmlURL       string            `json:"html_url,omitempty"`
	BuildNumber   int64             `json:"build_number,omitempty"`
	BranchName    string            `json:"branch_name,omitempty"`
	BranchHtmlURL string            `json:"branch_html_url,omitempty"`
	BuildHtmlURL  string            `json:"bulid_html_url,omitempty"`
	Commit        SemaphoreciCommit `json:"commit,omitempty"`
}

func SemaphoreciDeployOutMessageFromBytes(bytes []byte) (SemaphoreciDeployOutMessage, error) {
	msg := SemaphoreciDeployOutMessage{}
	err := json.Unmarshal(bytes, &msg)
	if err == nil {
		msg.Commit.Message = strings.ToLower(strings.TrimSpace(msg.Commit.Message))
	}
	return msg, err
}
