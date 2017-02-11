package util

import (
	"fmt"
	"strings"

	"github.com/grokify/glip-go-webhook"
	"github.com/valyala/fasthttp"
)

func SendGlipWebhookCtx(ctx *fasthttp.RequestCtx, glipClient glipwebhook.GlipWebhookClient, glipMsg glipwebhook.GlipWebhookMessage) error {
	glipWebhookGuid := fmt.Sprintf("%s", ctx.UserValue("glipguid"))
	glipWebhookGuid = strings.TrimSpace(glipWebhookGuid)

	req, resp, err := glipClient.PostWebhookGUIDFast(glipWebhookGuid, glipMsg)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
		return err
	}
	fmt.Fprintf(ctx, "%s", string(resp.Body()))
	fasthttp.ReleaseRequest(req)
	fasthttp.ReleaseResponse(resp)
	return nil
}

func SendGlipWebhook(glipClient glipwebhook.GlipWebhookClient, glipWebhookGuid string, glipMsg glipwebhook.GlipWebhookMessage) (int, []byte, error) {
	req, resp, err := glipClient.PostWebhookGUIDFast(glipWebhookGuid, glipMsg)

	if err != nil {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
		return -1, []byte(""), nil
	}
	status := resp.StatusCode()
	body := resp.Body()
	fasthttp.ReleaseRequest(req)
	fasthttp.ReleaseResponse(resp)
	return status, body, nil
}
