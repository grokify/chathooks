package circleci

import (
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"

	cc "github.com/commonchat/commonchat-go"
	"github.com/grokify/chatmore/src/adapters"
	"github.com/grokify/chatmore/src/config"
	"github.com/grokify/chatmore/src/util"
	"github.com/valyala/fasthttp"
)

const (
	DisplayName = "Circlecl"
	HandlerKey  = "circleci"
	IconURL     = "https://d2rbro28ib85bu.cloudfront.net/images/integrations/128/circleci.png"
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

// HandleFastHTTP is the method to respond to a fasthttp request.
func (h *Handler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
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

func Normalize(bytes []byte) (cc.Message, error) {
	message := cc.NewMessage()
	message.IconURL = IconURL

	src, err := CircleciOutMessageFromBytes(bytes)
	if err != nil {
		return message, err
	}

	message.Activity = fmt.Sprintf("Build %v", src.Status)

	message.Title = fmt.Sprintf("[Build #%v](%s) for [**%s/%s**](%s) %s",
		src.BuildNum, src.BuildURL, src.Reponame, src.Branch, src.VCSURL, src.Status)

	attachment := cc.NewAttachment()

	if len(src.Subject) > 0 {
		if len(src.BuildURL) > 0 {
			attachment.AddField(cc.Field{
				Title: "Subject",
				Value: fmt.Sprintf("[%v](%v)", src.Subject, src.BuildURL)})
		} else {
			attachment.AddField(cc.Field{
				Title: "Subject",
				Value: src.Subject})
		}
	}

	if len(src.Branch) > 0 {
		attachment.AddField(cc.Field{
			Title: "Branch",
			Value: src.Branch})
	}
	if len(src.Username) > 0 {
		attachment.AddField(cc.Field{
			Title: "Username",
			Value: src.Username})
	}
	if len(src.CommitterName) > 0 {
		attachment.AddField(cc.Field{
			Title: "Committer",
			Value: src.CommitterName})
	}

	message.AddAttachment(attachment)
	return message, nil
}

type CircleciOutPayload struct {
	Payload CircleciOutMessage `json:"payload,omitempty"`
}

type CircleciOutMessage struct {
	VCSURL          string        `json:"vcs_url,omitempty"`
	BuildURL        string        `json:"build_url,omitempty"`
	BuildNum        int64         `json:"build_num,omitempty"`
	Branch          string        `json:"branch,omitempty"`
	VCSRevision     string        `json:"vcs_revision,omitempty"`
	CommitterName   string        `json:"committer_name,omitempty"`
	CommitterEmail  string        `json:"committer_email,omitempty"`
	Subject         string        `json:"subject,omitempty"`
	Body            string        `json:"body,omitempty"`
	Why             string        `json:"why,omitempty"`
	DontBuild       interface{}   `json:"dont_build,omitempty"`
	QueuedAt        string        `json:"queued_at,omitempty"`
	StartTime       string        `json:"start_time,omitempty"`
	StopTime        string        `json:"stop_time,omitempty"`
	BuildTimeMillis int64         `json:"build_time_millis,omitempty"`
	Username        string        `json:"username,omitempty"`
	Reponame        string        `json:"reponame,omitempty"`
	Lifecycle       string        `json:"lifecycle,omitempty"`
	Outcome         string        `json:"outcome,omitempty"`
	Status          string        `json:"status,omitempty"`
	RetryOf         interface{}   `json:"retry_of,omitempty"`
	Steps           []interface{} `json:"steps,omitempty"`
}

type CircleciOutStep struct {
	Name    string              `json:"name,omitempty"`
	Actions []CircleciOutAction `json:"actions,omitempty"`
}

type CircleciOutAction struct {
	BashCommand   interface{} `json:"bash_command,omitempty"`
	RunTimeMillis int64       `json:"run_time_millis,omitempty"`
	StartTime     string      `json:"start_time,omitempty"`
	EndTime       string      `json:"end_time,omitempty"`
	Name          string      `json:"name,omitempty"`
	ExitCode      interface{} `json:"exit_cide,omitempty"`
	Type          string      `json:"type,omitempty"`
	Index         int64       `json:"index,omitempty"`
	Status        string      `json:"status,omitempty"`
}

func CircleciOutMessageFromBytes(bytes []byte) (CircleciOutMessage, error) {
	resp := CircleciOutPayload{}
	err := json.Unmarshal(bytes, &resp)
	if err != nil {
		return CircleciOutMessage{}, err
	}
	return resp.Payload, nil
}
