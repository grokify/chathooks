package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	cfg "github.com/grokify/gotilla/config"
	log "github.com/sirupsen/logrus"

	"github.com/grokify/chathooks/src/service"
)

/*

Use the `CHATHOOKS_TOKENS` environment variable to load secret
tokens as a comma delimited string.

*/

// CHATHOOKS_URL=http://localhost:8080/hook CHATHOOKS_HOME_URL=http://localhost:8080 go run main.go

func portAddress(port int) string { return ":" + strconv.Itoa(port) }

func main() {
	if err := cfg.LoadDotEnvSkipEmpty(os.Getenv("ENV_PATH"), "./.env"); err != nil {
		panic(err)
	}

	svc := service.NewService()

	if strings.ToLower(strings.TrimSpace(svc.Config.LogFormat)) == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{})
	}

	engine := svc.Config.Engine

	switch engine {
	case "awslambda":
		service.ServeAwsLambda(svc)
	case "nethttp":
		service.ServeNetHttp(svc)
	case "fasthttp":
		service.ServeFastHttp(svc)
	default:
		log.Fatal(fmt.Sprintf("Engine Not Supported [%v]", engine))
	}
}
