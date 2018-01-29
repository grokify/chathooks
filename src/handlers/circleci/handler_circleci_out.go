package circleci

import (
	"encoding/json"
	"fmt"

	"github.com/grokify/chathooks/src/config"
	"github.com/grokify/chathooks/src/handlers"
	"github.com/grokify/chathooks/src/models"
	cc "github.com/grokify/commonchat"
)

const (
	DisplayName      = "Circlecl"
	HandlerKey       = "circleci"
	MessageDirection = "out"
	DocumentationURL = "https://circleci.com/docs/1.0/configuration/#notify"
	MessageBodyType  = models.JSON
)

func NewHandler() handlers.Handler {
	return handlers.Handler{MessageBodyType: MessageBodyType, Normalize: Normalize}
}

func Normalize(cfg config.Configuration, bytes []byte) (cc.Message, error) {
	ccMsg := cc.NewMessage()
	iconURL, err := cfg.GetAppIconURL(HandlerKey)
	if err == nil {
		ccMsg.IconURL = iconURL.String()
	}

	src, err := CircleciOutMessageFromBytes(bytes)
	if err != nil {
		return ccMsg, err
	}

	ccMsg.Activity = fmt.Sprintf("Build %v", src.Status)

	ccMsg.Title = fmt.Sprintf("[Build #%v](%s) for [**%s/%s**](%s) %s",
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

	ccMsg.AddAttachment(attachment)
	return ccMsg, nil
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
