package codeship

import (
	"encoding/json"
	"fmt"

	cc "github.com/commonchat/commonchat-go"
	"github.com/grokify/chathooks/src/config"
	"github.com/grokify/chathooks/src/handlers"
	"github.com/grokify/chathooks/src/models"
)

const (
	DisplayName      = "Codeship"
	HandlerKey       = "codeship"
	MessageDirection = "out"
	DocumentationURL = "https://documentation.codeship.com/basic/getting-started/webhooks/"
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

	src, err := CodeshipOutMessageFromBytes(bytes)
	if err != nil {
		return ccMsg, err
	}

	build := src.Build

	status := build.Status
	if status == "infrastructure_failure" {
		status = "failed due to infrastructure error"
	}

	ccMsg.Activity = fmt.Sprintf("Build %v", status)
	ccMsg.Title = fmt.Sprintf("[Build #%v](%s) for **%s** %s ([%s](%s))",
		build.BuildId,
		build.BuildURL,
		build.ProjectName,
		status,
		build.ShortCommitId,
		build.CommitURL)

	attachment := cc.NewAttachment()

	if len(build.Message) > 0 {
		if len(build.CommitURL) > 0 {
			attachment.AddField(cc.Field{
				Title: "Message",
				Value: fmt.Sprintf("[%v](%v)", build.Message, build.CommitURL)})
		} else {
			attachment.AddField(cc.Field{
				Title: "Message",
				Value: build.Message})
		}
	}
	if len(build.Branch) > 0 {
		attachment.AddField(cc.Field{
			Title: "Branch",
			Value: build.Branch,
			Short: true})
	}
	if len(build.Committer) > 0 {
		attachment.AddField(cc.Field{
			Title: "Committer",
			Value: build.Committer,
			Short: true})
	}

	ccMsg.AddAttachment(attachment)
	return ccMsg, nil
}

type CodeshipOutMessage struct {
	Build CodeshipOutBuild `json:"build,omitempty"`
}

type CodeshipOutBuild struct {
	BuildURL        string `json:"build_url,omitempty"`
	CommitURL       string `json:"commit_url,omitempty"`
	ProjectId       int64  `json:"project_id,omitempty"`
	BuildId         int64  `json:"build_id,omitempty"`
	Status          string `json:"status,omitempty"`
	ProjectFullName string `json:"project_full_name,omitempty"`
	ProjectName     string `json:"project_name,omitempty"`
	CommitId        string `json:"commit_id,omitempty"`
	ShortCommitId   string `json:"short_commit_id,omitempty"`
	Message         string `json:"message,omitempty"`
	Committer       string `json:"committer,omitempty"`
	Branch          string `json:"branch,omitempty"`
}

func CodeshipOutMessageFromBytes(bytes []byte) (CodeshipOutMessage, error) {
	msg := CodeshipOutMessage{}
	err := json.Unmarshal(bytes, &msg)
	return msg, err
}

/*
{
  "build": {
    "build_url":"https://www.codeship.com/projects/10213/builds/973711",
    "commit_url":"https://github.com/codeship/docs/
                  commit/96943dc5269634c211b6fbb18896ecdcbd40a047",
    "project_id":10213,
    "build_id":973711,
    "status":"testing",
    # PROJECT_FULL_NAME IS DEPRECATED AND WILL BE REMOVED IN THE FUTURE
    "project_full_name":"codeship/docs",
    "project_name":"codeship/docs",
    "commit_id":"96943dc5269634c211b6fbb18896ecdcbd40a047",
    "short_commit_id":"96943",
    "message":"Merge pull request #34 from codeship/feature/shallow-clone",
    "committer":"beanieboi",
    "branch":"master"
  }
}

The status field can have one of the following values:
testing for newly started build

error for failed builds
success for passed builds
stopped for stopped builds
waiting for waiting builds
ignored for builds ignored because the account is over the monthly build limit
blocked for builds blocked because of excessive resource consumption
infrastructure_failure for builds which failed because of an internal error on the build VM
*/
