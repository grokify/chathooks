package travisci

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	cc "github.com/grokify/commonchat"
	"github.com/rs/zerolog/log"

	"github.com/grokify/chathooks/pkg/config"
	"github.com/grokify/chathooks/pkg/handlers"
	"github.com/grokify/chathooks/pkg/models"
)

const (
	DisplayName      = "Travis CI"
	HandlerKey       = "travisci"
	MessageDirection = "out"
	MessageBodyType  = models.JSON
)

func NewHandler() handlers.Handler {
	return handlers.Handler{MessageBodyType: MessageBodyType, Normalize: Normalize}
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

func Normalize(cfg config.Configuration, hReq handlers.HandlerRequest) (cc.Message, error) {
	bytes := hReq.Body
	ccMsg := cc.NewMessage()
	iconURL, err := cfg.GetAppIconURL(HandlerKey)
	if err == nil {
		ccMsg.IconURL = iconURL.String()
	}

	src, err := TravisciOutMessageFromBytes(bytes)
	if err != nil {
		return ccMsg, err
	}

	statusMessageSuffix := StatusMessageSuffix(src.StatusMessage)

	ccMsg.Activity = fmt.Sprintf("Build %s", statusMessageSuffix)

	attachment := cc.NewAttachment()
	attachment.Color = "#00ff00"

	attachment.Text = fmt.Sprintf(
		"[Build #%v](%s) for **%s/%s** %s",
		src.Number,
		src.BuildURL,
		src.Repository.Name,
		src.Branch,
		statusMessageSuffix)

	if len(src.Message) > 0 {
		field := cc.Field{Title: "Message"}
		if len(src.CompareURL) > 0 {
			field.Value = fmt.Sprintf("[%s](%s)", src.Message, src.CompareURL)
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

	ccMsg.AddAttachment(attachment)
	return ccMsg, nil
}

type TravisciOutMessage struct {
	ID                int                   `json:"id,omitempty"`
	AuthorEmail       string                `json:"author_email,omitempty"`
	AuthorName        string                `json:"author_name,omitempty"`
	Branch            string                `json:"branch,omitempty"`
	BuildURL          string                `json:"build_url,omitempty"`
	Commit            string                `json:"commit,omitempty"`
	CommitedAt        string                `json:"committed_at,omitempty"`
	CommitterName     string                `json:"committer_name,omitempty"`
	CommitterEmail    string                `json:"committer_email,omitempty"`
	CompareURL        string                `json:"compare_url,omitempty"`
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
	log.Debug().
		Str("type", "message.raw").
		Str("handler", HandlerKey).
		Str("request_body", string(bytes)).
		Msg(config.InfoInputMessageParseBegin)

	var msg TravisciOutMessage
	err := json.Unmarshal(bytes, &msg)
	if err != nil {
		log.Warn().
			Err(err).
			Str("type", "message.json.unmarshal").
			Str("handler", HandlerKey).
			Msg(config.ErrorInputMessageParseFailed)
	}
	return msg, err
}

type TravisciOutRepository struct {
	ID        int    `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	OwnerName string `json:"owner_name,omitempty"`
	URL       string `json:"url,omitempty"`
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
	ID     int `json:"id,omitempty"`
	Result int `json:"result,omitempty"`
	Status int `json:"status,omitempty"`
}

// Default template for Push Builds: "Build <%{build_url}|#%{build_number}> (<%{compare_url}|%{commit}>) of %{repository}@%{branch} by %{author} %{result} in %{duration}"

func (msg *TravisciOutMessage) PushBuildsAsMarkdown() string {
	return fmt.Sprintf("Build [#%v](%v) ([%v](%v)) of %v@%v by %v %v in %v", msg.Number, msg.BuildURL, msg.ShortCommit(), msg.CompareURL, msg.Repository.Name, msg.Branch, msg.AuthorName, strings.ToLower(msg.StatusMessage), msg.DurationDisplay())
}

func (msg *TravisciOutMessage) PullRequestBuildsAsMarkdown() string {
	return fmt.Sprintf("Build [#%v](%v) ([%v](%v)) of %v@%v in PR [#%v](%v) by %v %v in %v", msg.Number, msg.BuildURL, msg.ShortCommit(), msg.CompareURL, msg.Repository.Name, msg.Branch, msg.PullRequestNumber, msg.PullRequestURL(), msg.AuthorName, strings.ToLower(msg.StatusMessage), msg.DurationDisplay())
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
	return fmt.Sprintf("%v/pull/%v", msg.Repository.URL, msg.PullRequestNumber)
}
