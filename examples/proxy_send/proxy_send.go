package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/grokify/gotilla/io/ioutilmore"
	"github.com/grokify/gotilla/net/httputil"
	"github.com/grokify/gotilla/net/urlutil"
	"github.com/grokify/gotilla/strings/stringsutil"
	"github.com/grokify/webhookproxy/src/config"
	"github.com/valyala/fasthttp"
)

const (
	WebhookUrlEnvGlip  = "GLIP_WEBHOOK"
	WebhookUrlEnvSlack = "SLACK_WEBHOOK"
)

type ExampleWebhookSender struct {
	DocHandlersDir string
	BaseUrl        string
	OutputType     string
	Token          string
	Url            string
}

func (s *ExampleWebhookSender) SendExamplesForInputType(inputType string) error {
	rx := regexp.MustCompile(`^event-example_.+\.(json|txt)$`)
	inputTypeDir := path.Join(s.DocHandlersDir, inputType)
	files, err := ioutilmore.DirEntriesReSizeGt0(inputTypeDir, rx)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return errors.New(fmt.Sprintf("no ^event-example_ files found for %v", inputTypeDir))
	}
	for _, file := range files {
		filepath := path.Join(inputTypeDir, file.Name())
		err := s.SendExampleForFilepath(filepath, inputType)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *ExampleWebhookSender) SendExampleForFilepath(filepath string, inputType string) error {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	qry := url.Values{}
	qry.Add("inputType", inputType)
	qry.Add("outputType", s.OutputType)
	qry.Add("token", s.Token)
	qry.Add("url", s.Url)

	fullUrl := urlutil.BuildURL(s.BaseUrl, qry)
	fmt.Println(fullUrl)

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	req.SetBody(bytes)
	req.Header.SetRequestURI(fullUrl)
	req.Header.SetMethod("POST")
	req.Header.Set(httputil.ContentTypeHeader, httputil.ContentTypeValueJSONUTF8)

	fastClient := fasthttp.Client{}

	err = fastClient.Do(req, resp)
	fmt.Printf("RES_STATUS: %v\n", resp.StatusCode())
	if resp.StatusCode() > 299 {
		fmt.Printf("RES_BODY: %v\n", string(resp.Body()))
	}
	fasthttp.ReleaseRequest(req)
	fasthttp.ReleaseResponse(resp)
	return err
}

func main() {
	inputTypeP := flag.String("inputType", "travisci", "Example message type")
	urlP := flag.String("url", "https://hooks.glip.com/webhook/11112222-3333-4444-5555-666677778888", "Your Webhook URL")
	outputTypeP := flag.String("outputType", "glip", "Adapter name")

	flag.Parse()
	inputTypes := strings.ToLower(strings.TrimSpace(*inputTypeP))

	sender := ExampleWebhookSender{
		DocHandlersDir: config.DocsHandlersDir(),
		BaseUrl:        "http://localhost:8080/hooks",
		OutputType:     strings.ToLower(strings.TrimSpace(*outputTypeP)),
		Token:          "hello-world",
		Url:            strings.TrimSpace(*urlP),
	}
	if len(sender.Url) == 0 {
		sender.Url = os.Getenv(WebhookUrlEnvGlip)
	}

	examples := stringsutil.SliceTrimSpace(strings.Split(inputTypes, ","))

	for _, ex := range examples {
		err := sender.SendExamplesForInputType(ex)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("DONE")
}
