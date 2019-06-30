package glipwebhook

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	v2 "github.com/grokify/go-glip/v2"
	"github.com/grokify/gotilla/net/httputilmore"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

const (
	GlipWebhookBaseURLProduction   string = "https://hooks.glip.com/webhook/"
	GlipWebhookBaseURLProductionV2 string = "https://hooks.glip.com/webhook/v2/"
	GlipWebhookBaseURLSandbox      string = "https://hooks-glip.devtest.ringcentral.com/webhook/"
	GlipWebhookBaseURLSandboxV2    string = "https://hooks-glip.devtest.ringcentral.com/webhook/v2/"
	AttachmentTypeCard             string = "Card"
)

var (
	WebhookBaseURL string = "https://hooks.glip.com/webhook/"
)

type GlipWebhookClient struct {
	HttpClient     *http.Client
	FastClient     fasthttp.Client
	WebhookUrl     string
	webhookVersion int
}

func newGlipWebhookClientCore(urlOrGuid string) GlipWebhookClient {
	client := GlipWebhookClient{webhookVersion: 2}
	if len(urlOrGuid) > 0 {
		client.WebhookUrl = client.buildWebhookURL(urlOrGuid)
	}
	return client
}

func NewGlipWebhookClient(urlOrGuid string) (GlipWebhookClient, error) {
	client := newGlipWebhookClientCore(urlOrGuid)
	client.HttpClient = httputilmore.NewHttpClient()
	return client, nil
}

func NewGlipWebhookClientFast(urlOrGuid string) (GlipWebhookClient, error) {
	client := newGlipWebhookClientCore(urlOrGuid)
	client.FastClient = fasthttp.Client{}
	return client, nil
}

func (client *GlipWebhookClient) buildWebhookURL(urlOrUid string) string {
	rx := regexp.MustCompile(`^https?://`)
	rs := rx.FindString(urlOrUid)
	if len(rs) > 0 {
		log.WithFields(log.Fields{
			"lib":                    "go-glip",
			"request_url_http_match": urlOrUid}).Debug("Webhook URL has scheme.")
		return urlOrUid
	}
	return strings.Join([]string{WebhookBaseURL, urlOrUid}, "")
}

func (client *GlipWebhookClient) SendMessage(message GlipWebhookMessage) ([]byte, error) {
	resp, err := client.PostMessage(message)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (client *GlipWebhookClient) PostMessage(message GlipWebhookMessage) (*http.Response, error) {
	return client.PostWebhook(client.WebhookUrl, message)
}

func (client *GlipWebhookClient) PostWebhook(url string, message GlipWebhookMessage) (*http.Response, error) {
	if client.webhookVersion == 2 {
		return client.PostWebhookV2(V1ToV2WewbhookUri(url), webhookBodyV1ToV2(message))
	}
	return httputilmore.PostJsonSimple(url, message)
}

func (client *GlipWebhookClient) PostWebhookV2(url string, message v2.GlipWebhookMessage) (*http.Response, error) {
	return httputilmore.PostJsonSimple(url, message)
}

func (client *GlipWebhookClient) PostWebhookGUID(guid string, message GlipWebhookMessage) (*http.Response, error) {
	return client.PostWebhook(strings.Join([]string{WebhookBaseURL, guid}, ""), message)
}

// Request using fasthttp
// Recycle request and response using fasthttp.ReleaseRequest(req) and
// fasthttp.ReleaseResponse(resp)
func (client *GlipWebhookClient) PostMessageFast(message GlipWebhookMessage) (*fasthttp.Request, *fasthttp.Response, error) {
	return client.PostWebhookFast(client.WebhookUrl, message)
}

func (client *GlipWebhookClient) PostWebhookFast(url string, message GlipWebhookMessage) (*fasthttp.Request, *fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	if client.webhookVersion == 2 {
		url = V1ToV2WewbhookUri(url)
		bytes, err := json.Marshal(webhookBodyV1ToV2(message))
		if err != nil {
			return req, resp, err
		}
		req.SetBody(bytes)
	} else {
		bytes, err := json.Marshal(message)
		if err != nil {
			return req, resp, err
		}
		req.SetBody(bytes)
	}

	req.Header.SetRequestURI(url)
	req.Header.SetMethod(http.MethodPost)
	req.Header.Set(httputilmore.HeaderContentType, httputilmore.ContentTypeAppJsonUtf8)

	err := client.FastClient.Do(req, resp)
	return req, resp, err
}

func (client *GlipWebhookClient) PostWebhookGUIDFast(guidOrURL string, message GlipWebhookMessage) (*fasthttp.Request, *fasthttp.Response, error) {
	return client.PostWebhookFast(client.buildWebhookURL(guidOrURL), message)
}

func webhookBodyV1ToV2(v1msg GlipWebhookMessage) v2.GlipWebhookMessage {
	v2msg := v2.GlipWebhookMessage{
		Activity:    v1msg.Activity,
		IconUri:     v1msg.Icon,
		Text:        v1msg.Body,
		Title:       v1msg.Title,
		Attachments: []v2.Attachment{}}
	for _, v1att := range v1msg.Attachments {
		v2msg.Attachments = append(v2msg.Attachments, attachmentV1ToV2(v1att))
	}

	return v2msg
}

func attachmentV1ToV2(v1att Attachment) v2.Attachment {
	v2att := v2.Attachment{
		Color:        v1att.Color,
		Fields:       []v2.Field{},
		ImageUri:     v1att.ImageURL,
		Intro:        v1att.Pretext,
		Text:         v1att.Text,
		ThumbnailUri: v1att.ThumbnailURL,
		Title:        v1att.Title,
		Type:         v1att.Type,
	}
	if len(strings.TrimSpace(v2att.Type)) == 0 {
		v2att.Type = "Card"
	}
	if len(v1att.AuthorName) > 0 {
		v2att.Author = &v2.Author{
			Name:    v1att.AuthorName,
			IconUri: v1att.AuthorIcon,
			Uri:     v1att.AuthorLink}
	}
	if len(strings.TrimSpace(v1att.FooterIcon)) > 0 || len(strings.TrimSpace(v1att.Footer)) > 0 {
		v2att.Footnote = &v2.Footnote{
			IconUri: v1att.FooterIcon,
			Text:    v1att.Footer,
		}
		if v1att.TS > 0 {
			v2att.Footnote.Time = time.Unix(v1att.TS, 0)
		}
	}

	for _, v1field := range v1att.Fields {
		v2field := v2.Field{
			Title: v1field.Title,
			Value: v1field.Value}
		if v1field.Short {
			v2field.Style = v2.FieldStyleShort
		} else {
			v2field.Style = v2.FieldStyleLong
		}
		v2att.Fields = append(v2att.Fields, v2field)
	}

	return v2att
}

type GlipWebhookMessage struct {
	Icon           string       `json:"icon,omitempty"`
	Activity       string       `json:"activity,omitempty"`
	Title          string       `json:"title,omitempty"`
	Body           string       `json:"body,omitempty"`
	AttachmentType string       `json:"attachment_type,omitempty"`
	Attachments    []Attachment `json:"attachments,omitempty"`
}

type Attachment struct {
	Type         string  `json:"card,omitempty"`
	Color        string  `json:"color,omitempty"`
	Pretext      string  `json:"pretext,omitempty"`
	AuthorName   string  `json:"author_name,omitempty"`
	AuthorLink   string  `json:"author_link,omitempty"`
	AuthorIcon   string  `json:"author_icon,omitempty"`
	Title        string  `json:"title,omitempty"`
	TitleLink    string  `json:"title_link,omitempty"`
	Fallback     string  `json:"fallback,omitempty"`
	Fields       []Field `json:"fields,omitempty"`
	Text         string  `json:"text,omitempty"`
	ImageURL     string  `json:"image_url,omitempty"`
	ThumbnailURL string  `json:"thumbnail_url,omitempty"`
	Footer       string  `json:"footer,omitempty"`
	FooterIcon   string  `json:"footer_icon,omitempty"`
	TS           int64   `json:"ts,omitempty"`
}

type Author struct {
	Name    string `json:"name,omitempty"`
	URI     string `json:"uri,omitempty"`
	IconURI string `json:"iconUri,omitempty"`
}

type Footnote struct {
	Text    string `json:"text,omitempty"`
	IconURI string `json:"iconUri,omitempty"`
}

type Field struct {
	Title string `json:"title,omitempty"`
	Value string `json:"value,omitempty"`
	Short bool   `json:"short,omitempty"`
	Style string `json:"style,omitempty"`
}

type GlipWebhookResponse struct {
	Status  string           `json:"status,omitempty"`
	Message string           `json:"message,omitempty"`
	Error   GlipWebhookError `json:"error,omitempty"`
}

type GlipWebhookError struct {
	Code           string                   `json:"code,omitempty"`
	Message        string                   `json:"message,omitempty"`
	HttpStatusCode int                      `json:"http_status_code,omitempty"`
	ResponseData   string                   `json:"response_data,omitempty"`
	Response       GlipWebhookErrorResponse `json:"response,omitempty"`
}

func (gwerr *GlipWebhookError) Inflate() {
	if len(gwerr.ResponseData) > 2 {
		res := GlipWebhookErrorResponse{}
		err := json.Unmarshal([]byte(gwerr.ResponseData), &res)
		if err == nil {
			gwerr.Response = res
		}
	}
}

type GlipWebhookErrorResponse struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	Validation bool   `json:"validation"`
}
