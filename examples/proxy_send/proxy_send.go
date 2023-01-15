package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/google/go-querystring/query"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/net/http/httputilmore"
	"github.com/grokify/mogo/os/osutil"
	"github.com/grokify/mogo/type/stringsutil"
	"github.com/jessevdk/go-flags"
	"github.com/joho/godotenv"
	"github.com/valyala/fasthttp"

	"github.com/grokify/chathooks/pkg/config"
	"github.com/grokify/chathooks/pkg/models"
)

const (
	EnvWebhookURLGlip         = "GLIP_WEBHOOK"
	EnvWebhookURLSlack        = "SLACK_WEBHOOK"
	EnvChathooksReqInputType  = "CHATHOOKS_REQ_INPUT_TYPE"
	EnvChathooksReqOutputType = "CHATHOOKS_REQ_OUTPUT_TYPE"
	EnvChathooksReqToken      = "CHATHOOKS_REQ_TOKEN"
	EnvChathooksReqURL        = "CHATHOOKS_REQ_URL"
	EnvPath                   = "ENV_PATH"
)

type cliOptions struct {
	URLOrGUID    string `short:"u" long:"url" description:"Webhook URL or GUID" required:"true"`
	Input        string `short:"i" long:"input" description:"Input Service"`
	Output       string `short:"o" long:"output" description:"Output Adapter" required:"true"`
	Token        string `short:"t" long:"token" description:"Token"`
	ChathooksURL string `short:"c" long:"chathooks_url" description:"Chathooks URL"`
}

type ExampleWebhookSender struct {
	DocHandlersDir string
	BaseURL        string
	RequestParams  models.RequestParams
}

func (s *ExampleWebhookSender) SendExamplesForInputType(inputType string) error {
	rx := regexp.MustCompile(`^event-example_.+\.(json|txt)$`)
	inputTypeDir := path.Join(s.DocHandlersDir, inputType)
	entries, err := osutil.ReadDirMore(inputTypeDir, rx, false, true, false)
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		return fmt.Errorf("no ^event-example_ files found for [%v]", inputTypeDir)
	}
	for _, entry := range entries {
		filepath := path.Join(inputTypeDir, entry.Name())
		err := s.SendExampleForFilepath(filepath, inputType)
		if err != nil {
			return err
		}
	}
	return nil
}

func BuildURLQueryString(baseURL string, qry interface{}) string {
	v, _ := query.Values(qry)
	qryString := v.Encode()
	if len(qryString) > 0 {
		return baseURL + "?" + qryString
	}
	return baseURL
}

func (s *ExampleWebhookSender) SendExampleForFilepath(filepath string, inputType string) error {
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	qry := models.RequestParams{
		InputType:  inputType,
		OutputType: s.RequestParams.OutputType,
		Token:      s.RequestParams.Token,
		URL:        s.RequestParams.URL}

	fullURL := BuildURLQueryString(s.BaseURL, qry)
	// fmt.Printf("FULL_URL: %v\n", fullURL)

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	req.SetBody(bytes)
	req.Header.SetRequestURI(fullURL)
	req.Header.SetMethod(http.MethodPost)
	req.Header.Set(httputilmore.HeaderContentType, httputilmore.ContentTypeAppJSONUtf8)

	fastClient := fasthttp.Client{}

	err = fastClient.Do(req, resp)
	/*
		fmt.Printf("RES_STATUS: %v\n", resp.StatusCode())
		if resp.StatusCode() >= 300 || 1 == 1 {
			fmt.Printf("RES_BODY: %v\n", string(resp.Body()))
		}
	*/
	fasthttp.ReleaseRequest(req)
	fasthttp.ReleaseResponse(resp)
	return err
}

func main() {
	opts := cliOptions{}
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	qry := models.RequestParams{
		InputType:  opts.Input,
		OutputType: opts.Output,
		Token:      opts.Token,
		URL:        opts.URLOrGUID}
	fmtutil.PrintJSON(qry)

	if len(os.Getenv(EnvPath)) > 0 {
		err := godotenv.Load(os.Getenv(EnvPath))
		if err != nil {
			panic(err)
		}

		if len(os.Getenv(EnvChathooksReqInputType)) > 0 {
			qry.InputType = os.Getenv(EnvChathooksReqInputType)
		}
		if len(os.Getenv(EnvChathooksReqOutputType)) > 0 {
			qry.OutputType = os.Getenv(EnvChathooksReqOutputType)
		}
		if len(os.Getenv(EnvChathooksReqToken)) > 0 {
			qry.Token = os.Getenv(EnvChathooksReqToken)
		}
		if len(os.Getenv(EnvChathooksReqURL)) > 0 {
			qry.URL = os.Getenv(EnvChathooksReqURL)
		}
	}

	fmtutil.PrintJSON(qry)

	chathooksURL := "http://localhost:8080/hook"
	if len(strings.TrimSpace(opts.ChathooksURL)) > 0 {
		chathooksURL = opts.ChathooksURL
	}

	sender := ExampleWebhookSender{
		DocHandlersDir: config.DocsHandlersDir(),
		BaseURL:        chathooksURL,
		RequestParams:  qry}

	if len(sender.RequestParams.URL) == 0 {
		sender.RequestParams.URL = os.Getenv(EnvWebhookURLGlip)
	}

	examples := stringsutil.SplitCondenseSpace(qry.InputType, ",")

	for _, ex := range examples {
		err := sender.SendExamplesForInputType(strings.ToLower(ex))
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("DONE")
}
