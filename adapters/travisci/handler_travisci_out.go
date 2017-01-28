package travisci

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/grokify/glip-go-webhook"
	"github.com/grokify/glip-webhook-proxy/config"
	"github.com/valyala/fasthttp"
)

const (
	DISPLAY_NAME = "Travis CI"
	ICON_URL     = "https://blog.travis-ci.com/images/travis-mascot-200px.png"
)

type TravisciOutToGlipHandler struct {
	Config     config.Configuration
	GlipClient glipwebhook.GlipWebhookClient
}

func NewTravisciOutToGlipHandler(cfg config.Configuration, glip glipwebhook.GlipWebhookClient) TravisciOutToGlipHandler {
	return TravisciOutToGlipHandler{Config: cfg, GlipClient: glip}
}

func (h *TravisciOutToGlipHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	srcMsg, err := h.BuildTravisciOutMessage(ctx)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotAcceptable)
		fmt.Printf("ERROR:\n%v\n", err)
		return
	}
	glipMsg := h.TravisciOutToGlip(srcMsg)

	glipWebhookGuid := fmt.Sprintf("%s", ctx.UserValue("glipguid"))
	glipWebhookGuid = strings.TrimSpace(glipWebhookGuid)

	req, resp, err := h.GlipClient.PostWebhookGUIDFast(glipWebhookGuid, glipMsg)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
		return
	}
	fmt.Fprintf(ctx, "%s", string(resp.Body()))
	fasthttp.ReleaseRequest(req)
	fasthttp.ReleaseResponse(resp)
}

func (h *TravisciOutToGlipHandler) BuildTravisciOutMessage(ctx *fasthttp.RequestCtx) (TravisciOutMessage, error) {
	return TravisciOutMessageFromBytes(ctx.FormValue("payload"))
}

func (h *TravisciOutToGlipHandler) TravisciOutToGlip(src TravisciOutMessage) glipwebhook.GlipWebhookMessage {
	gmsg := glipwebhook.GlipWebhookMessage{
		Body:     strings.Join([]string{">", src.PushBuildsAsMarkdown()}, " "),
		Activity: DISPLAY_NAME,
		Icon:     ICON_URL}
	return gmsg
}

type TravisciOutMessage struct {
	Id             int                   `json:"id,omitempty"`
	AuthorEmail    string                `json:"author_email,omitempty"`
	AuthorName     string                `json:"author_name,omitempty"`
	Branch         string                `json:"branch,omitempty"`
	BuildUrl       string                `json:"build_url,omitempty"`
	Commit         string                `json:"commit,omitempty"`
	CommitedAt     string                `json:"committed_at,omitempty"`
	CommitterName  string                `json:"committer_name,omitempty"`
	CommitterEmail string                `json:"committer_email,omitempty"`
	CompareUrl     string                `json:"compare_url,omitempty"`
	Config         TravisciOutConfig     `json:"config,omitempty"`
	Duration       int                   `json:"duration,omitempty"`
	FinishedAt     string                `json:"finished_at,omitempty"`
	Matrix         []TravisciOutBuild    `json:"matrix,omitempty"`
	Message        string                `json:"message,omitempty"`
	Number         string                `json:"number,omitempty"`
	Repository     TravisciOutRepository `json:"repository,omitempty"`
	StartedAt      string                `json:"started_at,omitempty"`
	Status         int                   `json:"status"`
	StatusMessage  string                `json:"status_message,omitempty"`
	Type           string                `json:"type,omitempty"`
}

func TravisciOutMessageFromBytes(bytes []byte) (TravisciOutMessage, error) {
	fmt.Println("HERE")
	fmt.Println(string(bytes))
	msg := TravisciOutMessage{}
	err := json.Unmarshal(bytes, &msg)
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

// can Webhooks be array?
type TravisciOutNotifications struct {
	Webhooks string `json:"webhooks,omitempty"`
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

/*

Webhook Notification Reference

https://docs.travis-ci.com/user/notifications#Configuring-webhook-notifications

Format:

"Build <%{build_url}|#%{build_number}> (<%{compare_url}|%{commit}>) of %{repository}@%{branch} by %{author} %{result} in %{duration}"

Payload:

{
  "id": 1,
  "number": "1",
  "status": null,
  "started_at": null,
  "finished_at": null,
  "status_message": "Passed",
  "commit": "62aae5f70ceee39123ef",
  "branch": "master",
  "message": "the commit message",
  "compare_url": "https://github.com/svenfuchs/minimal/compare/master...develop",
  "committed_at": "2011-11-11T11: 11: 11Z",
  "committer_name": "Sven Fuchs",
  "committer_email": "svenfuchs@artweb-design.de",
  "author_name": "Sven Fuchs",
  "author_email": "svenfuchs@artweb-design.de",
  "type": "push",
  "build_url": "https://travis-ci.org/svenfuchs/minimal/builds/1",
  "repository": {
    "id": 1,
    "name": "minimal",
    "owner_name": "svenfuchs",
    "url": "http://github.com/svenfuchs/minimal"
   },
  "config": {
    "notifications": {
      "webhooks": ["http://evome.fr/notifications", "http://example.com/"]
    }
  },
  "matrix": [
    {
      "id": 2,
      "repository_id": 1,
      "number": "1.1",
      "state": "created",
      "started_at": null,
      "finished_at": null,
      "config": {
        "notifications": {
          "webhooks": ["http://evome.fr/notifications", "http://example.com/"]
        }
      },
      "status": null,
      "log": "",
      "result": null,
      "parent_id": 1,
      "commit": "62aae5f70ceee39123ef",
      "branch": "master",
      "message": "the commit message",
      "committed_at": "2011-11-11T11: 11: 11Z",
      "committer_name": "Sven Fuchs",
      "committer_email": "svenfuchs@artweb-design.de",
      "author_name": "Sven Fuchs",
      "author_email": "svenfuchs@artweb-design.de",
      "compare_url": "https://github.com/svenfuchs/minimal/compare/master...develop"
    }
  ]
}

*/
