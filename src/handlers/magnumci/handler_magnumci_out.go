package magnumci

import (
	"encoding/json"
	"errors"
	"fmt"

	log "github.com/Sirupsen/logrus"

	"github.com/grokify/glip-go-webhook"
	"github.com/grokify/glip-webhook-proxy-go/src/adapters"
	"github.com/grokify/glip-webhook-proxy-go/src/config"
	"github.com/grokify/glip-webhook-proxy-go/src/util"
	"github.com/valyala/fasthttp"
)

const (
	DisplayName = "Magnum CI"
	HandlerKey  = "magnumci"
	IconURL     = "https://pbs.twimg.com/profile_images/433440931543388160/nZ3y7AB__400x400.png"
)

// FastHttp request handler for Semaphore CI outbound webhook
type MagnumciOutToGlipHandler struct {
	Config     config.Configuration
	GlipClient glipwebhook.GlipWebhookClient
}

// FastHttp request handler constructor for Semaphore CI outbound webhook
func NewMagnumciOutToGlipHandler(cfg config.Configuration, glip glipwebhook.GlipWebhookClient) MagnumciOutToGlipHandler {
	return MagnumciOutToGlipHandler{Config: cfg, GlipClient: glip}
}

// HandleFastHTTP is the method to respond to a fasthttp request.
func (h *MagnumciOutToGlipHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	bytes := ctx.PostBody()

	glipMsg, err := Normalize(bytes)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotAcceptable)
		log.WithFields(log.Fields{
			"type":   "http.response",
			"status": fasthttp.StatusNotAcceptable,
		}).Info(fmt.Sprintf("%v request is not acceptable.", DisplayName))
		return
	}

	util.SendGlipWebhookCtx(ctx, h.GlipClient, glipMsg)
}

func Normalize(bytes []byte) (glipwebhook.GlipWebhookMessage, error) {
	gmsg := glipwebhook.GlipWebhookMessage{Icon: IconURL}

	src, err := MagnumciOutMessageFromBytes(bytes)
	if err != nil {
		return gmsg, err
	}

	if len(src.Title) > 0 {
		if config.GLIP_ACTIVITY_INCLUDE_INTEGRATION_NAME {
			gmsg.Activity = fmt.Sprintf("%v (%v)", src.Title, DisplayName)
		} else {
			gmsg.Activity = fmt.Sprintf("%v", src.Title)
		}
	} else {
		gmsg.Activity = fmt.Sprintf("%s Notification", DisplayName)
	}

	attachment := util.NewAttachment()

	if len(src.Message) > 0 {
		if len(src.CommitURL) > 0 {
			attachment.AddField(util.Field{
				Title: "Commit",
				Value: fmt.Sprintf("[%v](%v)", src.Message, src.CommitURL)})
		} else {
			attachment.AddField(util.Field{
				Title: "Commit",
				Value: fmt.Sprintf("%v", src.Message)})
		}
	} else if len(src.CommitURL) > 0 {
		attachment.AddField(util.Field{
			Title: "Commit",
			Value: fmt.Sprintf("[View Commit](%v)", src.Message)})
	}

	if len(src.Author) > 0 {
		attachment.AddField(util.Field{
			Title: "Author",
			Value: src.Author,
			Short: true})
	}
	if len(src.DurationString) > 0 {
		attachment.AddField(util.Field{
			Title: "Duration",
			Value: src.DurationString,
			Short: true})
	}
	if len(src.BuildURL) > 0 {
		attachment.AddField(util.Field{
			Value: fmt.Sprintf("[View Build](%v)", src.BuildURL)})
	}

	if len(src.Title) < 1 && len(attachment.Fields) == 0 {
		return gmsg, errors.New("Content not found")
	}

	message := util.NewMessage()
	message.AddAttachment(attachment)

	gmsg.Body = glipadapter.RenderMessage(message)

	return gmsg, nil
}

type MagnumciOutMessage struct {
	Id             int64  `json:"id,omitempty"`
	ProjectId      int64  `json:"project_id,omitempty"`
	Title          string `json:"title,omitempty"`
	Number         int64  `json:"number,omitempty"`
	Commit         string `json:"commit,omitempty"`
	Author         string `json:"author,omitempty"`
	Committer      string `json:"committer,omitempty"`
	Message        string `json:"message,omitempty"`
	Branch         string `json:"branch,omitempty"`
	State          string `json:"state,omitempty"`
	Status         string `json:"status,omitempty"`
	Result         int64  `json:"result,omitempty"`
	Duration       int64  `json:"duration,omitempty"`
	DurationString string `json:"duration_string,omitempty"`
	CommitURL      string `json:"commit_url,omitempty"`
	CompareURL     string `json:"compare_url,omitempty"`
	BuildURL       string `json:"build_url,omitempty"`
	StartedAt      string `json:"started_at,omitempty"`
	FinishedAt     string `json:"finished_at,omitempty"`
}

func MagnumciOutMessageFromBytes(bytes []byte) (MagnumciOutMessage, error) {
	msg := MagnumciOutMessage{}
	err := json.Unmarshal(bytes, &msg)
	return msg, err
}
