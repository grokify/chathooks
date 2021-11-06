package handlers

import (
	"fmt"

	"github.com/valyala/fasthttp"

	"github.com/grokify/chathooks/pkg/adapters"
	"github.com/grokify/chathooks/pkg/config"
)

const (
	QueryParamNamedOutputs    = config.ParamNameAdapters
	QueryParamInputType       = config.ParamNameInputType
	QueryParamOutputType      = config.ParamNameOutputType
	QueryParamOutputURL       = config.ParamNameOutputURL
	QueryParamURL             = config.ParamNameURL
	QueryParamToken           = config.ParamNameToken
	QueryParamDefaultActivity = config.ParamNameActivityDefault
	QueryParamDefaultIcon     = config.ParamNameIconDefault
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
