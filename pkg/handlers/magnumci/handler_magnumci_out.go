package magnumci

import (
	"encoding/json"
	"errors"
	"fmt"

	cc "github.com/grokify/commonchat"

	"github.com/grokify/chathooks/pkg/config"
	"github.com/grokify/chathooks/pkg/handlers"
	"github.com/grokify/chathooks/pkg/models"
)

const (
	DisplayName      = "Magnum CI"
	HandlerKey       = "magnumci"
	MessageDirection = "out"
	MessageBodyType  = models.JSON
)

func NewHandler() handlers.Handler {
	return handlers.Handler{MessageBodyType: MessageBodyType, Normalize: Normalize}
}

func Normalize(cfg config.Configuration, hReq handlers.HandlerRequest) (cc.Message, error) {
	ccMsg := cc.NewMessage()
	iconURL, err := cfg.GetAppIconURL(HandlerKey)
	if err == nil {
		ccMsg.IconURL = iconURL.String()
	}

	src, err := MagnumciOutMessageFromBytes(hReq.Body)
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
		return ccMsg, errors.New("content not found")
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
