package handlers

import (
	"fmt"

	"github.com/valyala/fasthttp"
	//"github.com/grokify/chatmore/src/adapters"
	//"github.com/grokify/chatmore/src/config"
)

const (
	handlerKey       = "base"
	messageDirection = "out"
)

// HomeHandler is a fasthttp handler for handling the webhoo proxy homepage.
func HomeHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "%s", []byte("Webhook Proxy"))
}

/*
type BaseHandler struct {
	Config  config.Configuration
	Adapter adapters.Adapter
}

func (h BaseHandler) HandlerKey() string {
	return handlerKey
}

func (h BaseHandler) MessageDirection() string {
	return messageDirection
}
*/
