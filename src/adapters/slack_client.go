package adapters

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/commonchat/commonchat-go/slack"
	"github.com/grokify/gotilla/net/httputil"
	"github.com/valyala/fasthttp"
)

const (
	HTTPMethod = "POST"
)

var (
	WebhookBaseURL = "https://hooks.slack.com/services/"
)

type SlackWebhookClient struct {
	HttpClient *http.Client
	FastClient fasthttp.Client
	WebhookUrl string
	UrlPrefix  *regexp.Regexp
}

func NewSlackWebhookClient(urlOrUid string, clientType string) (SlackWebhookClient, error) {
	log.WithFields(log.Fields{
		"lib": "slack_client.go",
		"request_url_client_init": urlOrUid}).Debug("")

	client := SlackWebhookClient{UrlPrefix: regexp.MustCompile(`^https:`)}
	client.WebhookUrl = client.BuildWebhookURL(urlOrUid)
	if clientType == "fast" {
		client.FastClient = fasthttp.Client{}
	} else {
		client.HttpClient = httputil.NewHttpClient()
	}
	return client, nil
}

func (client *SlackWebhookClient) BuildWebhookURL(urlOrUid string) string {
	rx := regexp.MustCompile(`^https:`)
	rs := rx.FindString(urlOrUid)
	if len(rs) > 0 {
		log.WithFields(log.Fields{
			"lib": "slack_client.go",
			"request_url_http_match": urlOrUid}).Debug("")
		return urlOrUid
	}
	return strings.Join([]string{WebhookBaseURL, urlOrUid}, "")
}

func (client *SlackWebhookClient) PostWebhookFast(url string, message slack.Message) (*fasthttp.Request, *fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	bytes, err := json.Marshal(message)
	if err != nil {
		return req, resp, err
	}
	req.SetBody(bytes)

	req.Header.SetMethod(HTTPMethod)
	req.Header.SetRequestURI(url)

	req.Header.Set(httputil.ContentTypeHeader, httputil.ContentTypeValueJSONUTF8)

	err = client.FastClient.Do(req, resp)
	return req, resp, err
}

func (client *SlackWebhookClient) PostWebhookGUIDFast(urlOrUid string, message slack.Message) (*fasthttp.Request, *fasthttp.Response, error) {
	return client.PostWebhookFast(client.BuildWebhookURL(urlOrUid), message)
}
