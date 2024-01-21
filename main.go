package main

import (
	"fmt"
	"os"

	"github.com/grokify/mogo/config"
	"github.com/grokify/xgo/net/http/httpsimple"

	"github.com/grokify/chathooks/pkg/service"
)

/*

Use the `CHATHOOKS_TOKENS` environment variable to load secret
tokens as a comma delimited string.

*/

// CHATHOOKS_URL=http://localhost:8080/hook CHATHOOKS_HOME_URL=http://localhost:8080 go run main.go

func main() {
	if _, err := config.LoadDotEnv(
		[]string{os.Getenv("ENV_PATH"), "./.env"}, 1); err != nil {
		panic(err)
	}

	svc := service.NewService()
	fmt.Printf("Starting on port [%d] with engine [%s].\n",
		svc.PortInt(), svc.HTTPEngine())
	httpsimple.Serve(svc)
}
