package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/grokify/gotilla/fmt/fmtutil"

	cc "github.com/commonchat/commonchat-go"
	"github.com/grokify/webhook-proxy-go/src/adapters"
	"github.com/grokify/webhook-proxy-go/src/util"
	"github.com/valyala/fasthttp"

	"github.com/grokify/webhook-proxy-go/src/handlers/appsignal"
	"github.com/grokify/webhook-proxy-go/src/handlers/circleci"
	"github.com/grokify/webhook-proxy-go/src/handlers/codeship"
	"github.com/grokify/webhook-proxy-go/src/handlers/confluence"
	"github.com/grokify/webhook-proxy-go/src/handlers/enchant"
	"github.com/grokify/webhook-proxy-go/src/handlers/gosquared"
	"github.com/grokify/webhook-proxy-go/src/handlers/heroku"
	"github.com/grokify/webhook-proxy-go/src/handlers/librato"
	"github.com/grokify/webhook-proxy-go/src/handlers/magnumci"
	"github.com/grokify/webhook-proxy-go/src/handlers/opsgenie"
	"github.com/grokify/webhook-proxy-go/src/handlers/papertrail"
	"github.com/grokify/webhook-proxy-go/src/handlers/pingdom"
	"github.com/grokify/webhook-proxy-go/src/handlers/raygun"
	"github.com/grokify/webhook-proxy-go/src/handlers/runscope"
	"github.com/grokify/webhook-proxy-go/src/handlers/semaphore"
	"github.com/grokify/webhook-proxy-go/src/handlers/statuspage"
	"github.com/grokify/webhook-proxy-go/src/handlers/travisci"
	"github.com/grokify/webhook-proxy-go/src/handlers/userlike"
)

const (
	GLIP_WEBHOOK_ENV  = "GLIP_WEBHOOK"
	SLACK_WEBHOOK_ENV = "SLACK_WEBHOOK"
)

type Sender struct {
	Adapter adapters.Adapter
}

func (sender *Sender) SendCcMessage(ccMsg cc.Message, err error) {
	if err != nil {
		panic(fmt.Sprintf("Bad Test Message: %v\n", err))
	}
	req, resp, err := sender.Adapter.SendMessage(ccMsg)
	fmt.Printf("RESPONSE_STATUS_CODE [%v]\n", resp.StatusCode())
	if err != nil {
		fmt.Printf("ERROR [%v]\n", err)
	}
	fasthttp.ReleaseRequest(req)
	fasthttp.ReleaseResponse(resp)
}

func main() {
	guidPointer := flag.String("guid", "", "Glip webhook GUID or URL")
	examplePointer := flag.String("example", "", "Example message type")
	adapterType := flag.String("adapter", "", "Adapter")

	flag.Parse()
	guid := strings.TrimSpace(*guidPointer)
	example := strings.ToLower(strings.TrimSpace(*examplePointer))

	fmt.Printf("LENGUID[%v]\n", len(guid))
	fmt.Printf("GUID [%v]\n", guid)
	fmt.Printf("EXAMPLE [%v]\n", example)

	if len(example) < 1 {
		panic("Usage: send_example.go -hook=<GUID> -adapter=glip -example=raygun")
	}

	sender := Sender{}
	if *adapterType == "glip" {
		if len(guid) < 1 {
			guid = os.Getenv(GLIP_WEBHOOK_ENV)
			fmt.Printf("GLIP_GUID_ENV [%v]\n", guid)
		}
		adapter, err := adapters.NewGlipAdapter(guid)
		if err != nil {
			panic("Incorrect Webhook GUID or URL")
		}
		sender.Adapter = &adapter
	} else if *adapterType == "slack" {
		if len(guid) < 1 {
			guid = os.Getenv(SLACK_WEBHOOK_ENV)
			fmt.Printf("SLACK_GUID_ENV [%v]\n", guid)
		}
		adapter, err := adapters.NewSlackAdapter(guid)
		if err != nil {
			panic("Incorrect Webhook GUID or URL")
		}
		sender.Adapter = &adapter
	} else {
		panic("Invalid Adapter")
	}

	exampleData, err := util.NewExampleData()
	if err != nil {
		panic(fmt.Sprintf("Invalid Example Data: %v\n", err))
	}
	fmtutil.PrintJSON(exampleData)

	switch example {
	case "appsignal":
		source := exampleData.Data[appsignal.HandlerKey]
		for _, eventSlug := range source.EventSlugs {
			sender.SendCcMessage(appsignal.ExampleMessage(exampleData, eventSlug))
		}
	case "circleci":
		sender.SendCcMessage(circleci.ExampleMessage(exampleData))
	case "codeship":
		sender.SendCcMessage(codeship.ExampleMessage(exampleData))
	case "confluence":
		source := exampleData.Data[confluence.HandlerKey]
		for _, eventSlug := range source.EventSlugs {
			sender.SendCcMessage(confluence.ExampleMessage(exampleData, eventSlug))
		}
	case "enchant":
		sender.SendCcMessage(enchant.ExampleMessage(exampleData))
	case "gosquared":
		source := exampleData.Data[gosquared.HandlerKey]
		for _, eventSlug := range source.EventSlugs {
			sender.SendCcMessage(gosquared.ExampleMessage(exampleData, eventSlug))
		}
	case "heroku":
		sender.SendCcMessage(heroku.ExampleMessage(exampleData))
	case "librato":
		source := exampleData.Data[librato.HandlerKey]
		for _, eventSlug := range source.EventSlugs {
			sender.SendCcMessage(librato.ExampleMessage(exampleData, eventSlug))
		}
	case "magnumci":
		sender.SendCcMessage(magnumci.ExampleMessage(exampleData))
	case "opsgenie":
		source := exampleData.Data[opsgenie.HandlerKey]
		for i, eventSlug := range source.EventSlugs {
			sender.SendCcMessage(opsgenie.ExampleMessage(exampleData, eventSlug))
			if i == 8 {
				time.Sleep(2000 * time.Millisecond)
			}
		}
	case "papertrail":
		source := exampleData.Data[papertrail.HandlerKey]
		for _, eventSlug := range source.EventSlugs {
			sender.SendCcMessage(papertrail.ExampleMessage(exampleData, eventSlug))
		}
	case "pingdom":
		source := exampleData.Data[pingdom.HandlerKey]
		for _, eventSlug := range source.EventSlugs {
			sender.SendCcMessage(pingdom.ExampleMessage(exampleData, eventSlug))
		}
	case "raygun":
		sender.SendCcMessage(raygun.ExampleMessage(exampleData))
	case "runscope":
		sender.SendCcMessage(runscope.ExampleMessage(exampleData))
	case "semaphore":
		source := exampleData.Data[semaphore.HandlerKey]
		for _, eventSlug := range source.EventSlugs {
			sender.SendCcMessage(semaphore.ExampleMessage(exampleData, eventSlug))
		}
	case "statuspage":
		source := exampleData.Data[statuspage.HandlerKey]
		for _, eventSlug := range source.EventSlugs {
			sender.SendCcMessage(statuspage.ExampleMessage(exampleData, eventSlug))
		}
	case "travisci":
		sender.SendCcMessage(travisci.ExampleMessage(exampleData))
	case "userlike":
		source := exampleData.Data[userlike.HandlerKey]
		for _, eventSlug := range source.EventSlugs {
			sender.SendCcMessage(userlike.ExampleMessage(exampleData, eventSlug))
		}
	default:
		panic(fmt.Sprintf("Unknown webhook source %v\n", example))
	}
}
