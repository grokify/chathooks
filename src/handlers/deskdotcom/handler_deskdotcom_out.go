package deskdotcom

import (
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"

	cc "github.com/commonchat/commonchat-go"
	"github.com/grokify/chatmore/src/adapters"
	"github.com/grokify/chatmore/src/config"
	"github.com/grokify/chatmore/src/util"
	"github.com/valyala/fasthttp"
)

const (
	DisplayName = "Desk.com"
	HandlerKey  = "deskdotcom"
	IconURL     = "https://pbs.twimg.com/profile_images/529708455678840832/Cvy48n8B_400x400.png"
)

// FastHttp request handler for Travis CI outbound webhook
type Handler struct {
	Config  config.Configuration
	Adapter adapters.Adapter
}

// FastHttp request handler constructor for Travis CI outbound webhook
func NewHandler(cfg config.Configuration, adapter adapters.Adapter) Handler {
	return Handler{Config: cfg, Adapter: adapter}
}

// HandleFastHTTP is the method to respond to a fasthttp request.
func (h *Handler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	ccMsg, err := Normalize(ctx.PostBody())

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotAcceptable)
		log.WithFields(log.Fields{
			"type":   "http.response",
			"status": fasthttp.StatusNotAcceptable,
		}).Info(fmt.Sprintf("%v request is not acceptable.", DisplayName))
		return
	}

	util.SendWebhook(ctx, h.Adapter, ccMsg)
}

func Normalize(bytes []byte) (cc.Message, error) {
	message, err := CcMessageFromBytes(bytes)
	if err != nil {
		return message, err
	}
	message.IconURL = IconURL
	return message, nil
}

/*

Example template:

{
  "activity":"Case {{case.id}} updated",
  "body":"**Case Type**\n{{case.status}} {{case.type}}\n **Subject**\n[{{case.subject}}]({{case.direct_url}})"
}

*/

func CcMessageFromBytes(bytes []byte) (cc.Message, error) {
	msg := cc.Message{}
	err := json.Unmarshal(bytes, &msg)
	return msg, err
}
