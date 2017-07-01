package victorops

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
	DisplayName = "VictorOps"
	HandlerKey  = "victorops"
	IconURL     = "https://victorops.com/wp-content/uploads/2015/04/download.png"
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

{
  "activity":"Event triggered",
  "title":"[Event ${{STATE.ENTITY_ID}}](https://portal.victorops.com/client/-/popoutIncident?incidentName=${{STATE.ENTITY_ID}}): ${{ALERT.entity_display_name}}",
  "body":"**State**\n${{ALERT.entity_state}}\n\n**Message**\n${{STATE.ACK_MSG}}\n\n**URL**\n${{ALERT.alert_url}}"}
}

Incident Alert
[${{ALERT.entity_display_name}}](${{ALERT.alert_url}})
Incident:${{STATE.INCIDENT_NAME}}
State: ${{STATE.CURRENT_STATE}}
Host:${{STATE.HOST}}

{
  "activity":"Event triggered",
  "title":"[Event ${{STATE.ENTITY_ID}}](https://portal.victorops.com/client/-/popoutIncident?incidentName=${{STATE.ENTITY_ID}}): ${{ALERT.entity_display_name}}",
  "body":"**State**\n${{ALERT.entity_state}}\n\n**Message**\n${{STATE.ACK_MSG}}"}

{
  "activity":"Incident alert",
  "body":"[${{ALERT.entity_display_name}}](https://portal.victorops.com/client/grokbase/popoutIncident?incidentName=${{STATE.INCIDENT_NAME}})\n**Incident**\n${{STATE.INCIDENT_NAME}}\n**State**\n${{STATE.CURRENT_STATE}}"
}

*/

func CcMessageFromBytes(bytes []byte) (cc.Message, error) {
	msg := cc.Message{}
	err := json.Unmarshal(bytes, &msg)
	return msg, err
}
