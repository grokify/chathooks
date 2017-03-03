package runscope

import (
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"

	cc "github.com/commonchat/commonchat-go"
	"github.com/grokify/webhook-proxy-go/src/adapters"
	"github.com/grokify/webhook-proxy-go/src/config"
	"github.com/grokify/webhook-proxy-go/src/util"
	"github.com/valyala/fasthttp"
)

const (
	DisplayName      = "Runscope"
	HandlerKey       = "runscope"
	IconURL          = "https://pbs.twimg.com/profile_images/500425058955689986/zlcbgqTt.png"
	DocumentationURL = "https://www.runscope.com/docs/api-testing/notifications#webhook"
)

// FastHttp request handler for Travis CI outbound webhook
type RunscopeOutToGlipHandler struct {
	Config  config.Configuration
	Adapter adapters.Adapter
}

// FastHttp request handler constructor for Travis CI outbound webhook
func NewRunscopeOutToGlipHandler(cfg config.Configuration, adapter adapters.Adapter) RunscopeOutToGlipHandler {
	return RunscopeOutToGlipHandler{Config: cfg, Adapter: adapter}
}

// HandleFastHTTP is the method to respond to a fasthttp request.
func (h *RunscopeOutToGlipHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	fmt.Printf(string(ctx.PostBody()))
	ccMsg, err := Normalize(ctx.PostBody())

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotAcceptable)
		log.WithFields(log.Fields{
			"type":   "http.response",
			"status": fasthttp.StatusNotAcceptable,
		}).Info(fmt.Sprintf("%v request is not acceptable:  %v", DisplayName, err))
		return
	}

	util.SendWebhook(ctx, h.Adapter, ccMsg)
}

func Normalize(bytes []byte) (cc.Message, error) {
	message := cc.NewMessage()
	message.IconURL = IconURL

	src, err := RunscopeOutMessageFromBytes(bytes)
	if err != nil {
		return message, err
	}

	message.Activity = fmt.Sprintf("Test run %s", src.Result)

	message.Title = fmt.Sprintf("[%v](%v) test run %v ([%v](%v))",
		src.TestName,
		src.TestURL,
		src.Result,
		src.TestId[:8],
		src.TestRunURL)

	attachment := cc.NewAttachment()

	if len(src.BucketName) > 0 {
		attachment.AddField(cc.Field{
			Title: "Bucket",
			Value: fmt.Sprintf("[%s - %s](%s)", src.BucketName, src.BucketKey, src.BucketURL())})
	}
	if len(src.EnvironmentName) > 0 {
		attachment.AddField(cc.Field{
			Title: "Environment",
			Value: src.EnvironmentName})
	}
	if len(src.RegionName) > 0 {
		attachment.AddField(cc.Field{
			Title: "Region",
			Value: fmt.Sprintf("%s (%s)", src.RegionName, src.Region)})
	}
	if len(src.TeamName) > 0 {
		attachment.AddField(cc.Field{
			Title: "Team",
			Value: fmt.Sprintf("%v (%v)", src.TeamName, src.TeamId[:8])})
	}

	message.AddAttachment(attachment)
	return message, nil
}

type RunscopeOutMessage struct {
	Variables       interface{}          `json:"variables,omitempty"`
	TestId          string               `json:"test_id,omitempty"`
	TestName        string               `json:"test_name,omitempty"`
	TestRunId       string               `json:"test_run_id,omitempty"`
	TeamId          string               `json:"team_id,omitempty"`
	TeamName        string               `json:"team_name,omitempty"`
	EnvironmentUUID string               `json:"environment_uuid,omitempty"`
	EnvironmentName string               `json:"environment_name,omitempty"`
	BucketName      string               `json:"bucket_name,omitempty"`
	BucketKey       string               `json:"bucket_key,omitempty"`
	TestURL         string               `json:"test_url,omitempty"`
	TestRunURL      string               `json:"test_run_url,omitempty"`
	TriggerURL      string               `json:"trigger_url,omitempty"`
	Result          string               `json:"result,omitempty"`
	StartedAt       float64              `json:"started_at,omitempty"`
	FinishedAt      float64              `json:"finished_at,omitempty"`
	Agent           interface{}          `json:"agent,omitempty"`
	Region          string               `json:"region,omitempty"`
	RegionName      string               `json:"region_name,omitempty"`
	Requests        []RunscopeOutRequest `json:"requests,omitempty"`
}

func (msg *RunscopeOutMessage) BucketURL() string {
	return fmt.Sprintf("https://www.runscope.com/radar/%s", msg.BucketKey)
}

func (msg *RunscopeOutMessage) EnvironmentsURL() string {
	return fmt.Sprintf("https://www.runscope.com/radar/%s/environments", msg.BucketKey)
}

func RunscopeOutMessageFromBytes(bytes []byte) (RunscopeOutMessage, error) {
	msg := RunscopeOutMessage{}
	err := json.Unmarshal(bytes, &msg)
	return msg, err
}

type RunscopeOutRequest struct {
	URL                string            `json:"url,omitempty"`
	Variables          RunscopeOutStatus `json:"variables,omitempty"`
	Assertions         RunscopeOutStatus `json:"assertions,omitempty"`
	Scripts            RunscopeOutStatus `json:"scripts,omitempty"`
	Result             string            `json:"result,omitempty"`
	Method             string            `json:"method,omitempty"`
	ResponseTimeMs     int64             `json:"response_time_ms,omitempty"`
	ResponseSizeBytes  int64             `json:"response_size_bytes,omitempty"`
	ResponseStatusCode string            `json:"response_status_code,omitempty"`
	Note               string            `json:"note,omitempty"`
}

type RunscopeOutStatus struct {
	Fail  int64 `json:"fail,omitempty"`
	Total int64 `json:"total,omitempty"`
	Pass  int64 `json:"pass,omitempty"`
}
