package glipwebhookproxy

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

func HomeHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "%s", []byte("Webhook Proxy"))
}
