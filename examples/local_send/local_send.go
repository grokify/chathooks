package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/jessevdk/go-flags"

	"github.com/grokify/chathooks/pkg/config"
	"github.com/grokify/chathooks/pkg/util"
	cc "github.com/grokify/commonchat"
	ccglip "github.com/grokify/commonchat/glip"
	ccslack "github.com/grokify/commonchat/slack"
	"github.com/valyala/fasthttp"

	"github.com/grokify/chathooks/examples"

	"github.com/grokify/chathooks/pkg/adapters"
	"github.com/grokify/chathooks/pkg/handlers/aha"
	"github.com/grokify/chathooks/pkg/handlers/appsignal"
	"github.com/grokify/chathooks/pkg/handlers/apteligent"
	"github.com/grokify/chathooks/pkg/handlers/bugsnag"
	"github.com/grokify/chathooks/pkg/handlers/circleci"
	"github.com/grokify/chathooks/pkg/handlers/codeship"
	"github.com/grokify/chathooks/pkg/handlers/confluence"
	"github.com/grokify/chathooks/pkg/handlers/datadog"
	"github.com/grokify/chathooks/pkg/handlers/deskdotcom"
	"github.com/grokify/chathooks/pkg/handlers/enchant"
	"github.com/grokify/chathooks/pkg/handlers/gosquared"
	"github.com/grokify/chathooks/pkg/handlers/gosquared2"
	"github.com/grokify/chathooks/pkg/handlers/heroku"
	"github.com/grokify/chathooks/pkg/handlers/librato"
	"github.com/grokify/chathooks/pkg/handlers/magnumci"
	"github.com/grokify/chathooks/pkg/handlers/marketo"
	"github.com/grokify/chathooks/pkg/handlers/opsgenie"
	"github.com/grokify/chathooks/pkg/handlers/papertrail"
	"github.com/grokify/chathooks/pkg/handlers/pingdom"
	"github.com/grokify/chathooks/pkg/handlers/raygun"
	"github.com/grokify/chathooks/pkg/handlers/runscope"
	"github.com/grokify/chathooks/pkg/handlers/semaphore"
	"github.com/grokify/chathooks/pkg/handlers/slack"
	"github.com/grokify/chathooks/pkg/handlers/statuspage"
	"github.com/grokify/chathooks/pkg/handlers/travisci"
	"github.com/grokify/chathooks/pkg/handlers/userlike"
	"github.com/grokify/chathooks/pkg/handlers/victorops"
	"github.com/grokify/chathooks/pkg/handlers/wootric"
)

type cliOptions struct {
	GuidOrWebhook string `short:"u" long:"url" description:"Webhook or GUID" required:"true"`
	Adapter       string `short:"a" long:"adapter" description:"Adapter" required:"true"`
	Service       string `short:"s" long:"service" description:"Service" required:"true"`
}

const (
	GLIP_WEBHOOK_ENV  = "GLIP_WEBHOOK"
	SLACK_WEBHOOK_ENV = "SLACK_WEBHOOK"
)

type Sender struct {
	Adapter cc.Adapter
}

func (sender *Sender) SendCcMessage(ccMsg cc.Message, err error) {
	if err != nil {
		panic(fmt.Sprintf("Bad Test Message: %v\n", err))
	}
	var resMsg any
	req, resp, err := sender.Adapter.SendMessage(ccMsg, &resMsg, map[string]any{})
	fmt.Printf("RESPONSE_STATUS_CODE [%v]\n", resp.StatusCode())
	if err != nil {
		fmt.Printf("ERROR [%v]\n", err)
	}

	// fmt.Println(string(resp.Body()))

	fasthttp.ReleaseRequest(req)
	fasthttp.ReleaseResponse(resp)
}

func main() {
	opts := cliOptions{}
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	dirs, _, err := examples.DocsHandlersDirInfo()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(strings.Join(dirs, ","))
	if len(opts.Service) > 0 {
		run(opts)
	}
	fmt.Println("DONE")
}

func run(opts cliOptions) {
	fmt.Printf("LENGUID[%v]\n", len(opts.GuidOrWebhook))
	fmt.Printf("GUID [%v]\n", opts.GuidOrWebhook)
	fmt.Printf("EXAMPLE [%v]\n", opts.Service)

	if len(opts.Service) < 1 {
		panic("Usage: send_example.go -g=<GUID> -a=glip -s=raygun")
	}

	cfg := config.Configuration{
		IconBaseURL: config.IconBaseURL}

	err := SendMessageAdapterHandler(cfg, opts)
	if err != nil {
		log.Fatal(err)
	}
}

func SendMessageAdapterHandler(cfg config.Configuration, opts cliOptions) error {
	webhookURLOrUID := opts.GuidOrWebhook
	adapterType := opts.Adapter
	service := opts.Service

	sender := Sender{}
	switch adapterType {
	case "glip":
		if len(webhookURLOrUID) < 1 {
			webhookURLOrUID = os.Getenv(GLIP_WEBHOOK_ENV)
			fmt.Printf("GLIP_GUID_ENV [%v]\n", webhookURLOrUID)
		}
		sender.Adapter = ccglip.NewGlipAdapter(webhookURLOrUID, adapters.GlipConfig())
	case "slack":
		if len(webhookURLOrUID) < 1 {
			webhookURLOrUID = os.Getenv(SLACK_WEBHOOK_ENV)
			fmt.Printf("SLACK_GUID_ENV [%v]\n", webhookURLOrUID)
		}
		adapter, err := ccslack.NewSlackAdapter(webhookURLOrUID)
		if err != nil {
			return errorsutil.Wrap(err, "incorrect webhook GUID or URL")
		}
		sender.Adapter = adapter
	default:
		return errors.New("invalid adapter")
	}

	exampleData, err := util.NewExampleData()
	if err != nil {
		return errorsutil.Wrap(err, fmt.Sprintf("invalid example data [%v]", err))
	}
	fmtutil.MustPrintJSON(exampleData)

	switch service {
	case "aha":
		source := exampleData.Data[aha.HandlerKey]
		for _, eventSlug := range source.EventSlugs {
			sender.SendCcMessage(aha.ExampleMessage(cfg, exampleData, eventSlug))
		}
	case "appsignal":
		source := exampleData.Data[appsignal.HandlerKey]
		for _, eventSlug := range source.EventSlugs {
			sender.SendCcMessage(appsignal.ExampleMessage(cfg, exampleData, eventSlug))
		}
	case "apteligent":
		source := exampleData.Data[apteligent.HandlerKey]
		for _, eventSlug := range source.EventSlugs {
			sender.SendCcMessage(apteligent.ExampleMessage(cfg, exampleData, eventSlug))
		}
	case "bugsnag":
		//sender.SendCcMessage(bugsnag.ExampleMessage(cfg, exampleData))
		source := exampleData.Data[bugsnag.HandlerKey]
		for _, eventSlug := range source.EventSlugs {
			sender.SendCcMessage(bugsnag.ExampleMessage(cfg, exampleData, eventSlug))
		}
	case "circleci":
		sender.SendCcMessage(circleci.ExampleMessage(cfg, exampleData))
	case "codeship":
		sender.SendCcMessage(codeship.ExampleMessage(cfg, exampleData))
	case "confluence":
		source := exampleData.Data[confluence.HandlerKey]
		for _, eventSlug := range source.EventSlugs {
			sender.SendCcMessage(confluence.ExampleMessage(cfg, exampleData, eventSlug))
		}
	case "datadog":
		sender.SendCcMessage(datadog.ExampleMessage(cfg, exampleData))
	case "deskdotcom":
		source := exampleData.Data[deskdotcom.HandlerKey]
		for _, eventSlug := range source.EventSlugs {
			sender.SendCcMessage(deskdotcom.ExampleMessage(cfg, exampleData, eventSlug))
		}
	case "enchant":
		sender.SendCcMessage(enchant.ExampleMessage(cfg, exampleData))
	case "gosquared":
		source := exampleData.Data[gosquared.HandlerKey]
		for _, eventSlug := range source.EventSlugs {
			sender.SendCcMessage(gosquared.ExampleMessage(cfg, exampleData, eventSlug))
		}
	case "gosquared2":
		source := exampleData.Data[gosquared.HandlerKey]
		for _, eventSlug := range source.EventSlugs {
			sender.SendCcMessage(gosquared2.ExampleMessage(cfg, exampleData, eventSlug))
		}
	case "heroku":
		sender.SendCcMessage(heroku.ExampleMessage(cfg, exampleData))
	case "librato":
		source := exampleData.Data[librato.HandlerKey]
		for _, eventSlug := range source.EventSlugs {
			sender.SendCcMessage(librato.ExampleMessage(cfg, exampleData, eventSlug))
		}
	case "magnumci":
		sender.SendCcMessage(magnumci.ExampleMessage(cfg, exampleData))
	case "marketo":
		source := exampleData.Data[marketo.HandlerKey]
		for _, eventSlug := range source.EventSlugs {
			sender.SendCcMessage(marketo.ExampleMessage(cfg, exampleData, eventSlug))
		}
	case "opsgenie":
		source := exampleData.Data[opsgenie.HandlerKey]
		for i, eventSlug := range source.EventSlugs {
			sender.SendCcMessage(opsgenie.ExampleMessage(cfg, exampleData, eventSlug))
			if i == 8 {
				time.Sleep(2000 * time.Millisecond)
			}
		}
	case "papertrail":
		source := exampleData.Data[papertrail.HandlerKey]
		for _, eventSlug := range source.EventSlugs {
			sender.SendCcMessage(papertrail.ExampleMessage(cfg, exampleData, eventSlug))
		}
	case "pingdom":
		source := exampleData.Data[pingdom.HandlerKey]
		for _, eventSlug := range source.EventSlugs {
			sender.SendCcMessage(pingdom.ExampleMessage(cfg, exampleData, eventSlug))
		}
	case "raygun":
		sender.SendCcMessage(raygun.ExampleMessage(cfg, exampleData))
	case "runscope":
		sender.SendCcMessage(runscope.ExampleMessage(cfg, exampleData))
	case "semaphore":
		source := exampleData.Data[semaphore.HandlerKey]
		for _, eventSlug := range source.EventSlugs {
			sender.SendCcMessage(semaphore.ExampleMessage(cfg, exampleData, eventSlug))
		}
	case "slack":
		source := exampleData.Data[slack.HandlerKey]
		for _, eventSlug := range source.EventSlugs {
			sender.SendCcMessage(slack.ExampleMessage(cfg, exampleData, eventSlug))
		}
	case "statuspage":
		source := exampleData.Data[statuspage.HandlerKey]
		for _, eventSlug := range source.EventSlugs {
			sender.SendCcMessage(statuspage.ExampleMessage(cfg, exampleData, eventSlug))
		}
	case "travisci":
		sender.SendCcMessage(travisci.ExampleMessage(cfg, exampleData))
	case "userlike":
		source := exampleData.Data[userlike.HandlerKey]
		for _, eventSlug := range source.EventSlugs {
			sender.SendCcMessage(userlike.ExampleMessage(cfg, exampleData, eventSlug))
		}
	case "victorops":
		sender.SendCcMessage(victorops.ExampleMessage(cfg, exampleData))
	case "wootric":
		source := exampleData.Data[wootric.HandlerKey]
		for _, eventSlug := range source.EventSlugs {
			sender.SendCcMessage(wootric.ExampleMessage(cfg, exampleData, eventSlug))
		}
	default:
		return fmt.Errorf("unknown webhook source [%s]", service)
	}
	return nil
}
