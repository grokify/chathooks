package magnumci

import (
	"encoding/json"
	"errors"
	"fmt"

	log "github.com/Sirupsen/logrus"

	cc "github.com/commonchat/commonchat-go"
	"github.com/grokify/webhookproxy/src/adapters"
	"github.com/grokify/webhookproxy/src/config"
	"github.com/grokify/webhookproxy/src/util"
	"github.com/valyala/fasthttp"
)

const (
	DisplayName      = "Magnum CI"
	HandlerKey       = "magnumci"
	MessageDirection = "out"
)

// FastHttp request handler for outbound webhook
type Handler struct {
	Config  config.Configuration
	Adapter adapters.Adapter
}

// FastHttp request handler constructor for outbound webhook
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
func (h *Handler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
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

func Normalize(cfg config.Configuration, bytes []byte) (cc.Message, error) {
	ccMsg := cc.NewMessage()
	iconURL, err := cfg.GetAppIconURL(HandlerKey)
	if err == nil {
		ccMsg.IconURL = iconURL.String()
	}

	src, err := MagnumciOutMessageFromBytes(bytes)
	if err != nil {
		return ccMsg, err
	}

	ccMsg.Activity = fmt.Sprintf("Build %v", src.State)

	if len(src.Title) > 0 {
		ccMsg.Title = fmt.Sprintf("[Build #%v](%v) **%v**", src.Number, src.BuildURL, src.Title)
	} else {
		ccMsg.Title = fmt.Sprintf("Build #%v](%v)", src.Number, src.BuildURL)
	}

	attachment := cc.NewAttachment()

	if len(src.Message) > 0 {
		if len(src.CommitURL) > 0 {
			attachment.AddField(cc.Field{
				Title: "Message",
				Value: fmt.Sprintf("[%v](%v)", src.Message, src.CommitURL)})
		} else {
			attachment.AddField(cc.Field{
				Title: "Message",
				Value: fmt.Sprintf("%v", src.Message)})
		}
	} else if len(src.CommitURL) > 0 {
		attachment.AddField(cc.Field{
			Title: "Commit",
			Value: fmt.Sprintf("[View Commit](%v)", src.Message)})
	}

	if len(src.Author) > 0 {
		attachment.AddField(cc.Field{
			Title: "Author",
			Value: src.Author,
			Short: true})
	}
	if len(src.Committer) > 0 {
		attachment.AddField(cc.Field{
			Title: "Committer",
			Value: src.Committer,
			Short: true})
	}
	if len(src.DurationString) > 0 {
		attachment.AddField(cc.Field{
			Title: "Duration",
			Value: src.DurationString,
			Short: true})
	}

	if len(src.Title) < 1 && len(attachment.Fields) == 0 {
		return ccMsg, errors.New("Content not found")
	}

	ccMsg.AddAttachment(attachment)
	return ccMsg, nil
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
