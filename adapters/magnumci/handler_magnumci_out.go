package magnumci

import (
	"encoding/json"
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/grokify/glip-go-webhook"
	"github.com/grokify/glip-webhook-proxy-go/config"
	"github.com/grokify/glip-webhook-proxy-go/util"
	"github.com/valyala/fasthttp"
)

const (
	DISPLAY_NAME = "Magnum CI"
	ICON_URL     = "https://a.slack-edge.com/ae7f/img/services/magnum-ci_512.png"
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
	gmsg.Activity = fmt.Sprintf("%v's %v #%v %v (%v)", src.Author, "build", src.Number, src.Status, DISPLAY_NAME)
	gmsg.Activity = fmt.Sprintf("%v (%v)", src.Title, DISPLAY_NAME)

	lines := []string{}
	//lines = append(lines, src.Title)
	lines = append(lines, fmt.Sprintf("> Commit: [%v](%v)", src.Message, src.CommitURL))
	if len(src.Author) > 0 {
		lines = append(lines, fmt.Sprintf("> Author: %v", src.Author))
	}
	if len(src.DurationString) > 0 {
		lines = append(lines, fmt.Sprintf("> Duration: %v", src.DurationString))
	}
	if len(src.BuildURL) > 0 {
		lines = append(lines, fmt.Sprintf("> [View Build](%v)", src.BuildURL))
	}
	//lines = append(lines, fmt.Sprintf("| **Build** | **Status** |\n| %v • [view](%v) | %v |", src.Number, src.Status, src.BuildURL))
	//lines = append(lines, fmt.Sprintf("| **Branch** | **Author** |\n| %v | %v |", src.Branch, src.Author))
	//lines = append(lines, fmt.Sprintf("| **Commit** |\n| %v • [view](%v) |", src.Message, src.CommitURL))
	if len(lines) > 0 {
		gmsg.Body = strings.Join(lines, "\n")
	}

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

/*

{
  "id": 1603,
  "project_id": 43,
  "title": "[PASS] project-name #130 (master - e91e132) by Dan Sosedoff",
  "number": 130,
  "commit": "e91e132612d263d95211aae6de2df9e503f22704",
  "author": "Dan Sosedoff",
  "committer": "Dan Sosedoff",
  "message": "Commit Message",
  "branch": "master",
  "state": "finished",
  "status": "pass",
  "result": 0,
  "duration": 158,
  "duration_string": "2m 38s",
  "commit_url": "http://domain.com/commit/e91e132612d263...",
  "compare_url": null,
  "build_url": "http://magnum-ci.com/projects/43/builds/1603",
  "started_at": "2013-02-14T00:09:01-06:00",
  "finished_at": "2013-02-14T00:11:39-06:00"
}

*/
