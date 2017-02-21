package semaphoreci

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"

	cc "github.com/commonchat/commonchat-go"
	"github.com/grokify/webhook-proxy-go/src/adapters"
	"github.com/grokify/webhook-proxy-go/src/config"
	"github.com/grokify/webhook-proxy-go/src/util"
	"github.com/valyala/fasthttp"
)

const (
	DisplayName = "Semaphore"
	HandlerKey  = "semaphoreci"
	IconURL     = "https://a.slack-edge.com/ae7f/plugins/semaphore/assets/service_512.png"
	ICON_URL_2  = "https://s3.amazonaws.com/semaphore-media/logos/png/gear/semaphore-gear-large.png"
)

// FastHttp request handler for Semaphore CI outbound webhook
type SemaphoreciOutToGlipHandler struct {
	Config  config.Configuration
	Adapter adapters.Adapter
}

// FastHttp request handler constructor for Semaphore CI outbound webhook
func NewSemaphoreciOutToGlipHandler(cfg config.Configuration, adapter adapters.Adapter) SemaphoreciOutToGlipHandler {
	return SemaphoreciOutToGlipHandler{Config: cfg, Adapter: adapter}
}

// HandleFastHTTP is the method to respond to a fasthttp request.
func (h *SemaphoreciOutToGlipHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
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

//func NormalizeBytes(bytes []byte) (glipwebhook.GlipWebhookMessage, error) {
func Normalize(bytes []byte) (cc.Message, error) {
	message := cc.NewMessage()
	message.IconURL = IconURL

	baseMsg, err := SemaphoreciBaseOutMessageFromBytes(bytes)
	if err != nil {
		return message, err
	}

	switch baseMsg.Event {
	case "build":
		srcMsg, err := SemaphoreciBuildOutMessageFromBytes(bytes)
		if err != nil {
			return message, err
		}
		return NormalizeSemaphoreciBuildOutMessage(srcMsg), nil
	case "deploy":
		srcMsg, err := SemaphoreciDeployOutMessageFromBytes(bytes)
		if err != nil {
			return message, err
		}
		return NormalizeSemaphoreciDeployOutMessage(srcMsg), nil
	}
	return cc.Message{IconURL: IconURL}, errors.New("EventNotFound")
}

func NormalizeSemaphoreciBuildOutMessage(src SemaphoreciBuildOutMessage) cc.Message {
	message := cc.NewMessage()
	message.IconURL = IconURL

	if strings.ToLower(strings.TrimSpace(src.Event)) == "build" {
		// Joe Cool build #15 passed
		message.Activity = fmt.Sprintf("%v's %v #%v %v%v", src.Commit.AuthorName, src.Event, src.BuildNumber, src.Result, adapters.IntegrationActivitySuffix(DisplayName))
	} else {
		message.Activity = fmt.Sprintf("%v's %v %v%v", src.Commit.AuthorName, src.Event, src.Result, adapters.IntegrationActivitySuffix(DisplayName))
	}

	attachment := cc.NewAttachment()

	if len(src.Commit.Message) > 0 {
		attachment.Text = src.Commit.Message
	}
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
	if len(src.BuildURL) > 0 {
		attachment.AddField(cc.Field{
			Value: fmt.Sprintf("[View details](%v)", src.BuildURL)})
	}

	message.AddAttachment(attachment)
	return message
}

func NormalizeSemaphoreciDeployOutMessage(src SemaphoreciDeployOutMessage) cc.Message {
	message := cc.NewMessage()
	message.IconURL = IconURL

	if strings.ToLower(strings.TrimSpace(src.Event)) == "build" {
		message.Activity = fmt.Sprintf("%v's %v #%v %v%v",
			src.Commit.AuthorName, src.Event, src.BuildNumber, src.Result, adapters.IntegrationActivitySuffix(DisplayName))
	} else {
		message.Activity = fmt.Sprintf("%v's %v %v%v",
			src.Commit.AuthorName, src.Event, src.Result, adapters.IntegrationActivitySuffix(DisplayName))
	}

	attachment := cc.NewAttachment()

	if len(src.Commit.Message) > 0 {
		attachment.Text = src.Commit.Message
	}
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
	if len(src.HtmlURL) > 0 {
		attachment.AddField(cc.Field{
			Value: fmt.Sprintf("[View details](%v)", src.HtmlURL)})
	}

	message.AddAttachment(attachment)
	return message
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
