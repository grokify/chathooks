package travisci

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	cc "github.com/grokify/commonchat"
	"github.com/grokify/webhook-proxy-go/src/adapters"
	"github.com/grokify/webhook-proxy-go/src/config"
	"github.com/grokify/webhook-proxy-go/src/util"
	"github.com/valyala/fasthttp"
)

const (
	DisplayName = "Travis CI"
	HandlerKey  = "travisci"
	IconURL     = "https://blog.travis-ci.com/images/travis-mascot-200px.png"
)

// FastHttp request handler for Travis CI outbound webhook
type TravisciOutToGlipHandler struct {
	Config  config.Configuration
	Adapter adapters.Adapter
}

// FastHttp request handler constructor for Travis CI outbound webhook
func NewTravisciOutToGlipHandler(cfg config.Configuration, adapter adapters.Adapter) TravisciOutToGlipHandler {
	return TravisciOutToGlipHandler{Config: cfg, Adapter: adapter}
}

// HandleFastHTTP is the method to respond to a fasthttp request.
func (h *TravisciOutToGlipHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
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

func StatusMessageSuffix(statusMessage string) string {
	suffixes := map[string]string{
		"pending":       "is pending",
		"passed":        "passed",
		"broken":        "is broken",
		"fixed":         "was fixed",
		"still failing": "is still failing"}
	statusMessage = strings.ToLower(statusMessage)
	if suffix, ok := suffixes[statusMessage]; ok {
		return suffix
	}
	return statusMessage
}

func Normalize(bytes []byte) (cc.Message, error) {
	message := cc.NewMessage()
	message.IconURL = IconURL

	src, err := TravisciOutMessageFromBytes(bytes)
	if err != nil {
		return message, err
	}

	statusMessageSuffix := StatusMessageSuffix(src.StatusMessage)

	message.Activity = fmt.Sprintf("Build %s", statusMessageSuffix)

	attachment := cc.NewAttachment()
	attachment.Color = "#00ff00"

	attachment.Text = fmt.Sprintf(
		"[Build #%v](%s) for **%s/%s** %s",
		src.Number,
		src.BuildUrl,
		src.Repository.Name,
		src.Branch,
		statusMessageSuffix)

	if len(src.Message) > 0 {
		field := cc.Field{Title: "Message"}
		if len(src.CompareUrl) > 0 {
			field.Value = fmt.Sprintf("[%s](%s)", src.Message, src.CompareUrl)
		} else {
			field.Value = src.Message
		}
		attachment.AddField(field)
	}
	if len(strings.TrimSpace(src.Branch)) > 0 {
		attachment.AddField(cc.Field{Title: "Branch", Value: strings.TrimSpace(src.Branch), Short: true})
	}
	if len(src.Type) > 0 {
		attachment.AddField(cc.Field{Title: "Type", Value: src.Type, Short: true})
	}
	if len(src.AuthorName) > 0 {
		attachment.AddField(cc.Field{Title: "Author", Value: src.AuthorName, Short: true})
	}
	if len(src.CommitterName) > 0 {
		attachment.AddField(cc.Field{Title: "Committer", Value: src.CommitterName, Short: true})
	}

	message.AddAttachment(attachment)
	return message, nil
}

type TravisciOutMessage struct {
	Id                int                   `json:"id,omitempty"`
	AuthorEmail       string                `json:"author_email,omitempty"`
	AuthorName        string                `json:"author_name,omitempty"`
	Branch            string                `json:"branch,omitempty"`
	BuildUrl          string                `json:"build_url,omitempty"`
	Commit            string                `json:"commit,omitempty"`
	CommitedAt        string                `json:"committed_at,omitempty"`
	CommitterName     string                `json:"committer_name,omitempty"`
	CommitterEmail    string                `json:"committer_email,omitempty"`
	CompareUrl        string                `json:"compare_url,omitempty"`
	Config            TravisciOutConfig     `json:"config,omitempty"`
	Duration          int                   `json:"duration,omitempty"`
	FinishedAt        string                `json:"finished_at,omitempty"`
	Matrix            []TravisciOutBuild    `json:"matrix,omitempty"`
	Message           string                `json:"message,omitempty"`
	Number            string                `json:"number,omitempty"`
	PullRequest       bool                  `json:"pull_request,omitempty"`
	PullRequestNumber int                   `json:"pull_request_number,omitempty"`
	PullRequestTitle  string                `json:"pull_request_title,omitempty"`
	Repository        TravisciOutRepository `json:"repository,omitempty"`
	StartedAt         string                `json:"started_at,omitempty"`
	Status            int                   `json:"status"`
	StatusMessage     string                `json:"status_message,omitempty"`
	Type              string                `json:"type,omitempty"`
}

func TravisciOutMessageFromBytes(bytes []byte) (TravisciOutMessage, error) {
	log.WithFields(log.Fields{
		"type":    "message.raw",
		"message": string(bytes),
	}).Debug(fmt.Sprintf("%v message.", DisplayName))
	msg := TravisciOutMessage{}
	err := json.Unmarshal(bytes, &msg)
	if err != nil {
		log.WithFields(log.Fields{
			"type":  "message.json.unmarshal",
			"error": fmt.Sprintf("%v\n", err),
		}).Warn(fmt.Sprintf("%v request unmarshal failure.", DisplayName))
	}
	return msg, err
}

type TravisciOutRepository struct {
	Id        int    `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	OwnerName string `json:"owner_name,omitempty"`
	Url       string `json:"url,omitempty"`
}

type TravisciOutConfig struct {
	Language      string                   `json:"language,omitempty"`
	Notifications TravisciOutNotifications `json:"notifications,omitempty"`
}

// can Webhooks can be a string (simple) or a dictionary (secure)
type TravisciOutNotifications struct {
	// Webhooks string `json:"webhooks,omitempty"`
}

type TravisciOutBuild struct {
	Id     int `json:"id,omitempty"`
	Result int `json:"result,omitempty"`
	Status int `json:"status,omitempty"`
}

// Default template for Push Builds: "Build <%{build_url}|#%{build_number}> (<%{compare_url}|%{commit}>) of %{repository}@%{branch} by %{author} %{result} in %{duration}"

func (msg *TravisciOutMessage) PushBuildsAsMarkdown() string {
	return fmt.Sprintf("Build [#%v](%v) ([%v](%v)) of %v@%v by %v %v in %v", msg.Number, msg.BuildUrl, msg.ShortCommit(), msg.CompareUrl, msg.Repository.Name, msg.Branch, msg.AuthorName, strings.ToLower(msg.StatusMessage), msg.DurationDisplay())
}

func (msg *TravisciOutMessage) PullRequestBuildsAsMarkdown() string {
	return fmt.Sprintf("Build [#%v](%v) ([%v](%v)) of %v@%v in PR [#%v](%v) by %v %v in %v", msg.Number, msg.BuildUrl, msg.ShortCommit(), msg.CompareUrl, msg.Repository.Name, msg.Branch, msg.PullRequestNumber, msg.PullRequestURL(), msg.AuthorName, strings.ToLower(msg.StatusMessage), msg.DurationDisplay())
}

func (msg *TravisciOutMessage) AsMarkdown() string {
	if msg.Type == "pull_request" {
		return msg.PullRequestBuildsAsMarkdown()
	}
	return msg.PushBuildsAsMarkdown()
}

func (msg *TravisciOutMessage) ShortCommit() string {
	if len(msg.Commit) < 8 {
		return msg.Commit
	}
	return msg.Commit[0:7]
}

func (msg *TravisciOutMessage) DurationDisplay() string {
	if msg.Duration == 0 {
		return "0 sec"
	}
	dur, err := time.ParseDuration(fmt.Sprintf("%vs", msg.Duration))
	if err != nil {
		return "unknown"
	}
	modSeconds := math.Mod(float64(msg.Duration), float64(60))
	return fmt.Sprintf("%v min %v sec", int(dur.Minutes()), modSeconds)
}

func (msg *TravisciOutMessage) PullRequestURL() string {
	return fmt.Sprintf("%v/pull/%v", msg.Repository.Url, msg.PullRequestNumber)
}
