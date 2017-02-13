package semaphoreci

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/grokify/glip-go-webhook"
	"github.com/grokify/glip-webhook-proxy-go/src/config"
	"github.com/grokify/glip-webhook-proxy-go/src/util"
	"github.com/valyala/fasthttp"
)

const (
	DISPLAY_NAME = "Semaphore"
	ICON_URL     = "https://a.slack-edge.com/ae7f/plugins/semaphore/assets/service_512.png"
	ICON_URL_2   = "https://s3.amazonaws.com/semaphore-media/logos/png/gear/semaphore-gear-large.png"
)

// FastHttp request handler for Semaphore CI outbound webhook
type SemaphoreciOutToGlipHandler struct {
	Config     config.Configuration
	GlipClient glipwebhook.GlipWebhookClient
}

// FastHttp request handler constructor for Semaphore CI outbound webhook
func NewSemaphoreciOutToGlipHandler(cfg config.Configuration, glip glipwebhook.GlipWebhookClient) SemaphoreciOutToGlipHandler {
	return SemaphoreciOutToGlipHandler{Config: cfg, GlipClient: glip}
}

// HandleFastHTTP is the method to respond to a fasthttp request.
func (h *SemaphoreciOutToGlipHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	glipMsg, err := NormalizeBytes(ctx.PostBody())
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotAcceptable)
		log.WithFields(log.Fields{
			"type":   "http.response",
			"status": fasthttp.StatusNotAcceptable,
		}).Info(fmt.Sprintf("%v request is not acceptable.", DISPLAY_NAME))
		return
	}

	util.SendGlipWebhookCtx(ctx, h.GlipClient, glipMsg)
}

func NormalizeBytes(bytes []byte) (glipwebhook.GlipWebhookMessage, error) {
	glipMsg := glipwebhook.GlipWebhookMessage{Icon: ICON_URL}
	baseMsg, err := SemaphoreciBaseOutMessageFromBytes(bytes)
	if err != nil {
		return glipMsg, err
	}
	switch baseMsg.Event {
	case "build":
		srcMsg, err := SemaphoreciBuildOutMessageFromBytes(bytes)
		if err != nil {
			return glipMsg, err
		}
		return NormalizeSemaphoreciBuildOutMessage(srcMsg), nil
	case "deploy":
		srcMsg, err := SemaphoreciDeployOutMessageFromBytes(bytes)
		if err != nil {
			return glipMsg, err
		}
		return NormalizeSemaphoreciDeployOutMessage(srcMsg), nil
	}
	return glipwebhook.GlipWebhookMessage{Icon: ICON_URL}, errors.New("EventNotFound")
}

func NormalizeSemaphoreciBuildOutMessage(src SemaphoreciBuildOutMessage) glipwebhook.GlipWebhookMessage {
	gmsg := glipwebhook.GlipWebhookMessage{Icon: ICON_URL}

	if strings.ToLower(strings.TrimSpace(src.Event)) == "build" {
		// Joe Cool build #15 passed
		gmsg.Activity = fmt.Sprintf("%v's %v #%v %v (%v)", src.Commit.AuthorName, src.Event, src.BuildNumber, src.Result, DISPLAY_NAME)
	} else {
		gmsg.Activity = fmt.Sprintf("%v's %v %v (%v)", src.Commit.AuthorName, src.Event, src.Result, DISPLAY_NAME)
	}

	lines := []string{}
	if len(src.Commit.Message) > 0 {
		lines = append(lines, fmt.Sprintf("> [%v/%v]: %v", src.ProjectName, src.BranchName, src.Commit.Message))
	}
	if len(src.BuildURL) > 0 {
		lines = append(lines, fmt.Sprintf("> [View details](%v)", src.BuildURL))
	}
	if len(lines) > 0 {
		gmsg.Body = strings.Join(lines, "\n")
	}
	return gmsg
}

func NormalizeSemaphoreciDeployOutMessage(src SemaphoreciDeployOutMessage) glipwebhook.GlipWebhookMessage {
	gmsg := glipwebhook.GlipWebhookMessage{Icon: ICON_URL}

	if strings.ToLower(strings.TrimSpace(src.Event)) == "build" {
		// Joe Cool build #15 passed
		gmsg.Activity = fmt.Sprintf("%v's %v #%v %v (%v)", src.Commit.AuthorName, src.Event, src.BuildNumber, src.Result, DISPLAY_NAME)
	} else {
		gmsg.Activity = fmt.Sprintf("%v's %v %v (%v)", src.Commit.AuthorName, src.Event, src.Result, DISPLAY_NAME)
	}

	lines := []string{}
	if len(src.Commit.Message) > 0 {
		lines = append(lines, fmt.Sprintf("> [%v/%v]: %v", src.ProjectName, src.BranchName, src.Commit.Message))
	}
	if len(src.HtmlURL) > 0 {
		lines = append(lines, fmt.Sprintf("> [View details](%v)", src.HtmlURL))
	}
	if len(lines) > 0 {
		gmsg.Body = strings.Join(lines, "\n")
	}
	return gmsg
}

type SemaphoreciBaseOutMessage struct {
	Event string `json:"event,omitempty"`
}

func SemaphoreciBaseOutMessageFromBytes(bytes []byte) (SemaphoreciBaseOutMessage, error) {
	msg := SemaphoreciBaseOutMessage{}
	err := json.Unmarshal(bytes, &msg)
	return msg, err
}

type SemaphoreciBuildOutMessage struct {
	BranchName    string            `json:"branch_name,omitempty"`
	BranchURL     string            `json:"branch_url,omitempty"`
	ProjectName   string            `json:"project_name,omitempty"`
	ProjectHashId string            `json:"project_hash_id,omitempty"`
	BuildURL      string            `json:"build_url,omitempty"`
	BuildNumber   int64             `json:"build_number,omitempty"`
	Result        string            `json:"result,omitempty"`
	Event         string            `json:"event,omitempty"`
	StartedAt     string            `json:"started_at,omitempty"`
	FinishedAt    string            `json:"finished_at,omitempty"`
	Commit        SemaphoreciCommit `json:"commit,omitempty"`
}

type SemaphoreciCommit struct {
	Id          string `json:"id,omitempty"`
	URL         string `json:"url,omitempty"`
	AuthorName  string `json:"author_name,omitempty"`
	AuthorEmail string `json:"author_email,omitempty"`
	Message     string `json:"message,omitempty"`
	Timestamp   string `json:"timestamp,omitempty"`
}

func SemaphoreciBuildOutMessageFromBytes(bytes []byte) (SemaphoreciBuildOutMessage, error) {
	msg := SemaphoreciBuildOutMessage{}
	err := json.Unmarshal(bytes, &msg)
	if err == nil {
		msg.Commit.Message = strings.ToLower(strings.TrimSpace(msg.Commit.Message))
	}
	return msg, err
}

type SemaphoreciDeployOutMessage struct {
	ProjectName   string            `json:"project_name,omitempty"`
	ProjectHashId string            `json:"project_hash_id,omitempty"`
	Result        string            `json:"result,omitempty"`
	Event         string            `json:"event,omitempty"`
	ServerName    string            `json:"server_name,omitempty"`
	Number        int64             `json:"number,omitempty"`
	CreatedAt     string            `json:"created_at,omitempty"`
	UpdatedAt     string            `json:"updated_at,omitempty"`
	StartedAt     string            `json:"started_at,omitempty"`
	FinishedAt    string            `json:"finished_at,omitempty"`
	HtmlURL       string            `json:"html_url,omitempty"`
	BuildNumber   int64             `json:"build_number,omitempty"`
	BranchName    string            `json:"branch_name,omitempty"`
	BranchHtmlURL string            `json:"branch_html_url,omitempty"`
	BuildHtmlURL  string            `json:"bulid_html_url,omitempty"`
	Commit        SemaphoreciCommit `json:"commit,omitempty"`
}

func SemaphoreciDeployOutMessageFromBytes(bytes []byte) (SemaphoreciDeployOutMessage, error) {
	msg := SemaphoreciDeployOutMessage{}
	err := json.Unmarshal(bytes, &msg)
	if err == nil {
		msg.Commit.Message = strings.ToLower(strings.TrimSpace(msg.Commit.Message))
	}
	return msg, err
}

/*

if strings.ToLower(msg.Event)== "build"
%v %v %v
msg.Commit.AuthorName
msg.ProjectName
Build #65 passed
fmt.Sprintf("%v %v #%v %v", msg.ProjectName, msg.BuildNumber, msg. Result)
%v", msg.COmmit.Message
[View Details](%v)


Webhook Notification Reference

Format:

"Build <%{build_url}|#%{build_number}> (<%{compare_url}|%{commit}>) of %{repository}@%{branch} by %{author} %{result} in %{duration}"

Payload:

{
  "branch_name": "gem_updates",
  "branch_url": "https://semaphoreci.com/projects/44/branches/50",
  "project_name": "base-app",
  "project_hash_id": "123-aga-471-6a8",
  "build_url": "https://semaphoreci.com/projects/44/branches/50/builds/15",
  "build_number": 15,
  "result": "passed",
  "event": "build",
  "started_at": "2012-07-09T15:23:53Z",
  "finished_at": "2012-07-09T15:30:16Z",
  "commit": {
    "id": "dc395381e650f3bac18457909880829fc20e34ba",
    "url": "https://github.com/renderedtext/base-app/commit/dc395381e650f3bac18457909880829fc20e34ba",
    "author_name": "Vladimir Saric",
    "author_email": "vladimir@renderedtext.com",
    "message": "Update 'shoulda' gem.",
    "timestamp": "2012-07-04T18:14:08Z"
  }
}

Deploy Webhook

{
  "project_name": "heroku-deploy-test",
  "project_hash_id": "123-aga-471-6a8",
  "result": "passed",
  "event": "deploy",
  "server_name": "server-heroku-master-automatic-2",
  "number": 2,
  "created_at": "2013-07-30T13:52:33Z",
  "updated_at": "2013-07-30T13:53:21Z",
  "started_at": "2013-07-30T13:52:38Z",
  "finished_at": "2013-07-30T13:53:21Z",
  "html_url": "https://semaphoreci.com/projects/2420/servers/81/deploys/2",
  "build_number": 10,
  "branch_name": "master",
  "branch_html_url": "https://semaphoreci.com/projects/2420/branches/58394",
  "build_html_url": "https://semaphoreci.com/projects/2420/branches/58394/builds/7",
  "commit": {
    "author_email": "rastasheep3@gmail.com",
    "author_name": "Aleksandar Diklic",
    "id": "43ddb7516ecc743f0563abd7418f0bd3617348c4",
    "message": "One more time",
    "timestamp": "2013-07-19T12:56:25Z",
    "url": "https://github.com/rastasheep/heroku-deploy-test/commit/43ddb7516ecc743f0563abd7418f0bd3617348c4"
  }
}

*/
