package semaphoreci

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/grokify/glip-go-webhook"
	"github.com/grokify/glip-webhook-proxy-go/src/adapters"
	"github.com/grokify/glip-webhook-proxy-go/src/config"
	"github.com/grokify/glip-webhook-proxy-go/src/util"
	"github.com/valyala/fasthttp"
)

const (
	DISPLAY_NAME = "Semaphore"
	HANDLER_KEY  = "semaphoreci"
	ICON_URL     = "https://a.slack-edge.com/ae7f/plugins/semaphore/assets/service_512.png"
	ICON_URL_2   = "https://s3.amazonaws.com/semaphore-media/logos/png/gear/semaphore-gear-large.png"
)

// FastHttp request handler for Semaphore CI outbound webhook
type SemaphoreciOutToGlipHandler struct {
	Config     config.Configuration
	GlipClient glipwebhook.GlipWebhookClient
}

// FastHttp request handler constructor for Semaphore CI outbound webhook
func NewSemaphoreciOutToGlipHandler(cfg config.Configuration, glip glipwebhook.GlipWebhookClient) SemaphoreciOutToGlipHandler {
	return SemaphoreciOutToGlipHandler{Config: cfg, GlipClient: glip}
}

// HandleFastHTTP is the method to respond to a fasthttp request.
func (h *SemaphoreciOutToGlipHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	glipMsg, err := NormalizeBytes(ctx.PostBody())
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

func NormalizeBytes(bytes []byte) (glipwebhook.GlipWebhookMessage, error) {
	glipMsg := glipwebhook.GlipWebhookMessage{Icon: ICON_URL}
	baseMsg, err := SemaphoreciBaseOutMessageFromBytes(bytes)
	if err != nil {
		return glipMsg, err
	}
	switch baseMsg.Event {
	case "build":
		srcMsg, err := SemaphoreciBuildOutMessageFromBytes(bytes)
		if err != nil {
			return glipMsg, err
		}
		return NormalizeSemaphoreciBuildOutMessage(srcMsg), nil
	case "deploy":
		srcMsg, err := SemaphoreciDeployOutMessageFromBytes(bytes)
		if err != nil {
			return glipMsg, err
		}
		return NormalizeSemaphoreciDeployOutMessage(srcMsg), nil
	}
	return glipwebhook.GlipWebhookMessage{Icon: ICON_URL}, errors.New("EventNotFound")
}

func NormalizeSemaphoreciBuildOutMessage(src SemaphoreciBuildOutMessage) glipwebhook.GlipWebhookMessage {
	gmsg := glipwebhook.GlipWebhookMessage{Icon: ICON_URL}

	if strings.ToLower(strings.TrimSpace(src.Event)) == "build" {
		// Joe Cool build #15 passed
		gmsg.Activity = fmt.Sprintf("%v's %v #%v %v%v", src.Commit.AuthorName, src.Event, src.BuildNumber, src.Result, glipadapter.IntegrationActivitySuffix(DISPLAY_NAME))
	} else {
		gmsg.Activity = fmt.Sprintf("%v's %v %v%v", src.Commit.AuthorName, src.Event, src.Result, glipadapter.IntegrationActivitySuffix(DISPLAY_NAME))
	}

	message := util.NewMessage()

	if len(src.Commit.Message) > 0 {
		message.AddAttachment(util.Attachment{
			Text: fmt.Sprintf("[%v/%v]: %v", src.ProjectName, src.BranchName, src.Commit.Message)})
	}
	if len(src.BuildURL) > 0 {
		message.AddAttachment(util.Attachment{
			Text: fmt.Sprintf("[View details](%v)", src.BuildURL)})
	}

	gmsg.Body = glipadapter.RenderMessage(message)
	return gmsg
}

func NormalizeSemaphoreciDeployOutMessage(src SemaphoreciDeployOutMessage) glipwebhook.GlipWebhookMessage {
	gmsg := glipwebhook.GlipWebhookMessage{Icon: ICON_URL}

	if strings.ToLower(strings.TrimSpace(src.Event)) == "build" {
		gmsg.Activity = fmt.Sprintf("%v's %v #%v %v%v",
			src.Commit.AuthorName, src.Event, src.BuildNumber, src.Result, glipadapter.IntegrationActivitySuffix(DISPLAY_NAME))
	} else {
		gmsg.Activity = fmt.Sprintf("%v's %v %v%v",
			src.Commit.AuthorName, src.Event, src.Result, glipadapter.IntegrationActivitySuffix(DISPLAY_NAME))
	}

	message := util.NewMessage()

	if len(src.Commit.Message) > 0 {
		message.AddAttachment(util.Attachment{
			Text: fmt.Sprintf("[%v/%v]: %v", src.ProjectName, src.BranchName, src.Commit.Message)})
	}
	if len(src.HtmlURL) > 0 {
		message.AddAttachment(util.Attachment{
			Text: fmt.Sprintf("[View details](%v)", src.HtmlURL)})
	}

	gmsg.Body = glipadapter.RenderMessage(message)
	return gmsg
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
