package main

import (
	"github.com/grokify/gotilla/fmt/fmtutil"
	"github.com/grokify/webhookproxy/src/config"
)

func main() {
	swag := config.GetSwaggerSpec()
	fmtutil.PrintJSON(swag)
}
