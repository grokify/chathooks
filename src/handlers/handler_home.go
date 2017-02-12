package handlers

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

// HomeHandler is a fasthttp handler for handling the webhoo proxy homepage.
func HomeHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "%s", []byte("Webhook Proxy"))
}
