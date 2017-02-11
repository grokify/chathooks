package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	//log "github.com/Sirupsen/logrus"
	"github.com/grokify/glip-go-webhook"
	"github.com/grokify/glip-webhook-proxy-go/adapters/confluence"
	"github.com/grokify/glip-webhook-proxy-go/adapters/enchant"
	"github.com/grokify/glip-webhook-proxy-go/adapters/heroku"
	"github.com/grokify/glip-webhook-proxy-go/adapters/raygun"
	"github.com/grokify/glip-webhook-proxy-go/adapters/travisci"
	"github.com/grokify/glip-webhook-proxy-go/util"
)

// https://hooks.glip.com/webhook/848c88d9-d892-451a-9614-6045046d477a

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
	case "raygun":
		glipMsg, err := raygun.ExampleMessageGlip()
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
	default:
		fmt.Printf("Unknown webhook source %v\n", example)
	}
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
