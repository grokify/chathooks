package adapters

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"github.com/commonchat/commonchat-go/slack"
	"github.com/grokify/gotilla/net/httputil"
	"github.com/valyala/fasthttp"
)

const (
	ContentTypeHeader = "Content-Type"
	ContentTypeValue  = "application/json"
	HTTPMethod        = "POST"
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
	client := SlackWebhookClient{UrlPrefix: regexp.MustCompile(`^https?:`)}
	client.WebhookUrl = client.BuildWebhookURL(urlOrUid)
	if clientType == "fast" {
		client.FastClient = fasthttp.Client{}
	} else {
		client.HttpClient = httputil.NewHttpClient()
	}
	return client, nil
}

func (client *SlackWebhookClient) BuildWebhookURL(urlOrUid string) string {
	rs := client.UrlPrefix.FindString(urlOrUid)
	if len(rs) > 0 {
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
	req.Header.SetRequestURI(client.BuildWebhookURL(url))
	req.Header.Set(ContentTypeHeader, ContentTypeValue)

	err = client.FastClient.Do(req, resp)
	return req, resp, err
}

func (client *SlackWebhookClient) PostWebhookGUIDFast(guid string, message slack.Message) (*fasthttp.Request, *fasthttp.Response, error) {
	return client.PostWebhookFast(strings.Join([]string{WebhookBaseURL, guid}, ""), message)
}
