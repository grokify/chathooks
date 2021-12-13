package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/grokify/mogo/config"
	"github.com/grokify/mogo/net/http/httpsimple"

	"github.com/grokify/chathooks/pkg/service"
)

/*

Use the `CHATHOOKS_TOKENS` environment variable to load secret
tokens as a comma delimited string.

*/

// CHATHOOKS_URL=http://localhost:8080/hook CHATHOOKS_HOME_URL=http://localhost:8080 go run main.go

func portAddress(port int) string { return ":" + strconv.Itoa(port) }

func main() {
	if err := config.LoadDotEnvSkipEmpty(
		os.Getenv("ENV_PATH"), "./.env"); err != nil {
		panic(err)
	}

	svc := service.NewService()
	fmt.Printf("Starting on port [%d] with engine [%s].\n",
		svc.PortInt(), svc.HttpEngine())
	httpsimple.Serve(svc)
}
