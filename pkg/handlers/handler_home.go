package handlers

import (
	"fmt"

	"github.com/valyala/fasthttp"

	"github.com/grokify/chathooks/pkg/adapters"
	"github.com/grokify/chathooks/pkg/config"
)

const (
	QueryParamNamedOutputs    = "adapters"
	QueryParamInputType       = "inputType"
	QueryParamOutputType      = "outputType"
	QueryParamToken           = "token"
	QueryParamOutputURL       = "url"
	QueryParamDefaultActivity = "defaultActivity"
	QueryParamDefaultIcon     = "defaultIcon"
)

var (
	ShowDisplayName = false
)

// HomeHandler is a fasthttp handler for handling the webhoo proxy homepage.
func HomeHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "%s", []byte("Chathooks\nSource: https://github.com/grokify/chathooks"))
}

type Configuration struct {
	ConfigData config.Configuration
	AdapterSet adapters.AdapterSet
}

func IntegrationActivitySuffix(displayName string) string {
	if !ShowDisplayName || len(displayName) < 1 {
		return ""
	}
	return ""
}
