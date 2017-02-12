package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	//log "github.com/Sirupsen/logrus"
	"github.com/grokify/glip-go-webhook"
	"github.com/grokify/glip-webhook-proxy-go/handlers/appsignal"
	"github.com/grokify/glip-webhook-proxy-go/handlers/confluence"
	"github.com/grokify/glip-webhook-proxy-go/handlers/enchant"
	"github.com/grokify/glip-webhook-proxy-go/handlers/heroku"
	"github.com/grokify/glip-webhook-proxy-go/handlers/magnumci"
	"github.com/grokify/glip-webhook-proxy-go/handlers/raygun"
	"github.com/grokify/glip-webhook-proxy-go/handlers/semaphoreci"
	"github.com/grokify/glip-webhook-proxy-go/handlers/travisci"
	"github.com/grokify/glip-webhook-proxy-go/handlers/userlike"
	"github.com/grokify/glip-webhook-proxy-go/util"
)

const (
	GLIP_WEBHOOK_ENV = "GLIP_WEBHOOK"
)

func main() {
	guidPointer := flag.String("guid", "", "Glip webhook GUID or URL")
	examplePointer := flag.String("example", "", "Example message type")
	flag.Parse()
	guid := strings.TrimSpace(*guidPointer)
	example := strings.ToLower(strings.TrimSpace(*examplePointer))

	fmt.Printf("LENGUID[%v]\n", len(guid))
	if len(guid) < 1 {
		guid = os.Getenv(GLIP_WEBHOOK_ENV)
		fmt.Printf("HERE [%v]\n", guid)
	}

	//glip, _ := glipwebhook.NewGlipWebhookClient(guid)
	fmt.Printf("GUID [%v]\n", guid)
	fmt.Printf("EXAMPLE [%v]\n", example)

	if len(example) < 1 {
		panic("Usage: send_example.go -hook=<GUID> -example=raygun")
	}

	glipClient, err := glipwebhook.NewGlipWebhookClient(guid)
	if err != nil {
		panic("Incorrect Webhook GUID or URL")
	}

	switch example {
	case "appsignal":
		SendAppsignal(glipClient, guid)
	case "confluence":
		SendConfluence(glipClient, guid)
	case "enchant":
		glipMsg, err := enchant.ExampleMessageGlip()
		if err != nil {
			panic("Bad Test Message")
		}
		util.SendGlipWebhook(glipClient, guid, glipMsg)
	case "heroku":
		glipMsg, err := heroku.ExampleMessageGlip()
		if err != nil {
			panic("Bad Test Message")
		}
		util.SendGlipWebhook(glipClient, guid, glipMsg)
	case "magnumci":
		glipMsg, err := magnumci.ExampleMessageGlip()
		if err != nil {
			panic(fmt.Sprintf("Bad Test Message [%v]", err))
		}
		util.SendGlipWebhook(glipClient, guid, glipMsg)
	case "raygun":
		glipMsg, err := raygun.ExampleMessageGlip()
		if err != nil {
			panic("Bad Test Message")
		}
		util.SendGlipWebhook(glipClient, guid, glipMsg)
	case "semaphoreci":
		glipMsg, err := semaphoreci.ExampleMessageGlip()
		if err != nil {
			panic("Bad Test Message")
		}
		util.SendGlipWebhook(glipClient, guid, glipMsg)
	case "travisci":
		glipMsg, err := travisci.ExampleMessageGlip()
		if err != nil {
			panic("Bad Test Message")
		}
		util.SendGlipWebhook(glipClient, guid, glipMsg)
	case "userlike":
		SendUserlike(glipClient, guid)
	default:
		fmt.Printf("Unknown webhook source %v\n", example)
	}
}

func SendAppsignal(glipClient glipwebhook.GlipWebhookClient, guid string) {
	glipMsg, err := appsignal.ExampleMarkerMessageGlip()
	if err != nil {
		panic("Bad Test Message")
	}
	util.SendGlipWebhook(glipClient, guid, glipMsg)
	glipMsg, err = appsignal.ExampleExceptionMessageGlip()
	if err != nil {
		panic("Bad Test Message")
	}
	util.SendGlipWebhook(glipClient, guid, glipMsg)
	glipMsg, err = appsignal.ExamplePerformanceMessageGlip()
	if err != nil {
		panic("Bad Test Message")
	}
	util.SendGlipWebhook(glipClient, guid, glipMsg)
}

func SendConfluence(glipClient glipwebhook.GlipWebhookClient, guid string) {
	glipMsg, err := confluence.ExamplePageCreatedMessageGlip()
	if err != nil {
		panic("Bad Test Message")
	}
	util.SendGlipWebhook(glipClient, guid, glipMsg)
	glipMsg, err = confluence.ExampleCommentCreatedMessageGlip()
	if err != nil {
		panic("Bad Test Message")
	}
	util.SendGlipWebhook(glipClient, guid, glipMsg)
}

func SendUserlike(glipClient glipwebhook.GlipWebhookClient, guid string) {
	glipMsg, err := userlike.ExampleOfflineMessageReceiveMessageGlip()
	if err != nil {
		panic(fmt.Sprintf("Bad Test Message [%v]", err))
	}
	util.SendGlipWebhook(glipClient, guid, glipMsg)
	glipMsg, err = userlike.ExampleUserlikeChatMetaStartOutMessageGlip()
	if err != nil {
		panic("Bad Test Message")
	}
	util.SendGlipWebhook(glipClient, guid, glipMsg)
}
