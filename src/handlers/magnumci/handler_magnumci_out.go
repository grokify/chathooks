package magnumci

import (
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"

	"github.com/grokify/glip-go-webhook"
	"github.com/grokify/glip-webhook-proxy-go/src/adapters"
	"github.com/grokify/glip-webhook-proxy-go/src/config"
	"github.com/grokify/glip-webhook-proxy-go/src/util"
	"github.com/valyala/fasthttp"
)

const (
	DISPLAY_NAME = "Magnum CI"
	HANDLER_KEY  = "magnumci"
	ICON_URL     = "https://pbs.twimg.com/profile_images/433440931543388160/nZ3y7AB__400x400.png"
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
	srcMsg, err := BuildInboundMessage(ctx)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotAcceptable)
		log.WithFields(log.Fields{
			"type":   "http.response",
			"status": fasthttp.StatusNotAcceptable,
		}).Info(fmt.Sprintf("%v request is not acceptable.", DISPLAY_NAME))
		return
	}
	glipMsg := Normalize(srcMsg)

	util.SendGlipWebhookCtx(ctx, h.GlipClient, glipMsg)
}

func BuildInboundMessage(ctx *fasthttp.RequestCtx) (MagnumciOutMessage, error) {
	return MagnumciOutMessageFromBytes(ctx.PostBody())
}

func Normalize(src MagnumciOutMessage) glipwebhook.GlipWebhookMessage {
	gmsg := glipwebhook.GlipWebhookMessage{Icon: ICON_URL}
	if config.GLIP_ACTIVITY_INCLUDE_INTEGRATION_NAME {
		gmsg.Activity = fmt.Sprintf("%v (%v)", src.Title, DISPLAY_NAME)
	} else {
		gmsg.Activity = fmt.Sprintf("%v", src.Title)
	}

	message := util.NewMessage()

	message.AddAttachment(util.Attachment{
		Title: "Commit",
		Text:  fmt.Sprintf("[%v](%v)", src.Message, src.CommitURL)})

	if len(src.Author) > 0 {
		message.AddAttachment(util.Attachment{
			Title: "Author",
			Text:  src.Author})
	}
	if len(src.DurationString) > 0 {
		message.AddAttachment(util.Attachment{
			Title: "Duration",
			Text:  src.DurationString})
	}
	if len(src.BuildURL) > 0 {
		message.AddAttachment(util.Attachment{
			Text: fmt.Sprintf("[View Build](%v)", src.BuildURL)})
	}

	gmsg.Body = glipadapter.RenderMessage(message)
	return gmsg
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
