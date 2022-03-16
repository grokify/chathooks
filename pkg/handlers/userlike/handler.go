package userlike

import (
	"errors"
	"fmt"
	"strings"

	cc "github.com/grokify/commonchat"

	"github.com/grokify/chathooks/pkg/config"
	"github.com/grokify/chathooks/pkg/handlers"
	"github.com/grokify/chathooks/pkg/models"
)

const (
	DisplayName      = "Userlike"
	HandlerKey       = "userlike"
	MessageDirection = "out"
	MessageBodyType  = models.JSON
)

var (
	ChatMetaEvents = []string{"feedback", "forward", "rating", "receive", "start", "survey"}
	OperatorEvents = []string{"away", "back", "offline", "online"}
)

func NewHandler() handlers.Handler {
	return handlers.Handler{MessageBodyType: MessageBodyType, Normalize: Normalize}
}

func Normalize(cfg config.Configuration, hReq handlers.HandlerRequest) (cc.Message, error) {
	bodyBytes := hReq.Body
	srcMsgBase, err := UserlikeBaseOutMessageFromBytes(bodyBytes)
	if err != nil {
		return cc.Message{}, err
	}
	if srcMsgBase.Type == "offline_message" && srcMsgBase.Event == "receive" {
		srcMsg, err := UserlikeOfflineMessageOutMessageFromBytes(bodyBytes)
		if err != nil {
			return cc.Message{}, err
		}
		return NormalizeOfflineMessage(cfg, srcMsg), nil
	} else if srcMsgBase.Type == "chat_meta" {
		srcMsg, err := UserlikeChatMetaStartOutMessageFromBytes(bodyBytes)
		if err != nil {
			return cc.Message{}, err
		}
		return NormalizeChatMeta(cfg, srcMsg), nil
	} else if srcMsgBase.Type == "operator" {
		srcMsg, err := UserlikeOperatorOutMessageFromBytes(bodyBytes)
		if err != nil {
			return cc.Message{}, err
		}
		return NormalizeOperator(cfg, srcMsg), nil
	} else if srcMsgBase.Type == "chat_widget" {
		srcMsg, err := UserlikeChatWidgetOutMessageFromBytes(bodyBytes)
		if err != nil {
			return cc.Message{}, err
		}
		return NormalizeChatWidget(cfg, srcMsg), nil
	}

	return cc.Message{}, errors.New("type not supported")
}

func GlipActivityForChat(event string, feedback string) string {
	eventDisplay := event
	eventMap := map[string]string{
		"start":    "session started",
		"forward":  "session forwarded",
		"rating":   "rating received",
		"feedback": "feedback received",
		"survey":   "survey received",
		"receive":  "session ended",
		"goal":     "goal achieved"}
	if event == "rating" && len(feedback) > 0 {
		eventDisplay = eventMap["feedback"]
	} else {
		if displayTry, ok := eventMap[event]; ok {
			eventDisplay = displayTry
		}
	}
	return fmt.Sprintf("Chat %s", eventDisplay)
}

func NormalizeOfflineMessage(cfg config.Configuration, src UserlikeOfflineMessageOutMessage) cc.Message {
	ccMsg := cc.NewMessage()
	iconURL, err := cfg.GetAppIconURL(HandlerKey)
	if err == nil {
		ccMsg.IconURL = iconURL.String()
	}

	ccMsg.Activity = fmt.Sprintf("Offline message received%v", handlers.IntegrationActivitySuffix(DisplayName))

	attachment := cc.NewAttachment()

	if len(src.URL) > 0 {
		attachment.AddField(cc.Field{
			Title: "Message",
			Value: fmt.Sprintf("[%s](%v)", src.Message, src.URL)})
	}
	if len(src.ClientName) > 0 {
		attachment.AddField(cc.Field{
			Title: "Client Name",
			Value: src.ClientName})
	}

	if len(attachment.Fields) > 0 {
		ccMsg.AddAttachment(attachment)
	}
	return ccMsg
}

func NormalizeChatMeta(cfg config.Configuration, src UserlikeChatMetaStartOutMessage) cc.Message {
	ccMsg := cc.NewMessage()
	iconURL, err := cfg.GetAppIconURL(HandlerKey)
	if err == nil {
		ccMsg.IconURL = iconURL.String()
	}

	ccMsg.Activity = fmt.Sprintf("%s%s",
		GlipActivityForChat(src.Event, src.FeedbackMessage), handlers.IntegrationActivitySuffix(DisplayName))

	attachment := cc.NewAttachment()

	displayedURL := false

	if src.Event == "rating" || src.Event == "survey" { // includes feedback
		if len(src.FeedbackMessage) > 0 {
			url, linked := LinkifyURL(src.FeedbackMessage, src.URL, displayedURL)
			displayedURL = linked
			attachment.AddField(cc.Field{
				Title: "Feedback",
				Value: url,
				Short: false})
		}
		if len(src.PostSurveyOption) > 0 {
			url, linked := LinkifyURL(src.PostSurveyOption, src.URL, displayedURL)
			displayedURL = linked
			attachment.AddField(cc.Field{
				Title: "Rating",
				Value: url,
				Short: true})
		}
	}
	if len(src.ClientName) > 0 {
		url, _ := LinkifyURL(src.ClientName, src.URL, displayedURL)
		attachment.AddField(cc.Field{
			Title: "Client Name",
			Value: url,
			Short: true})
	} else {
		url, _ := LinkifyURL("Unknown", src.URL, displayedURL)
		attachment.AddField(cc.Field{
			Title: "Client Name",
			Value: url,
			Short: true})
	}

	if len(attachment.Fields) > 0 {
		ccMsg.AddAttachment(attachment)
	}
	return ccMsg
}

func LinkifyURL(innerHTML string, url string, skipLinking bool) (string, bool) {
	if len(innerHTML) == 0 && len(url) > 0 {
		innerHTML = url
	}
	if skipLinking {
		return innerHTML, skipLinking
	}
	if len(url) < 1 {
		return innerHTML, false
	}
	if len(innerHTML) > 0 {
		return fmt.Sprintf("[%s](%s)", innerHTML, url), true
	}
	return fmt.Sprintf("[%s](%s)", url, url), true
}

func NormalizeChatWidget(cfg config.Configuration, src UserlikeChatWidgetOutMessage) cc.Message {
	ccMsg := cc.NewMessage()
	iconURL, err := cfg.GetAppIconURL(HandlerKey)
	if err == nil {
		ccMsg.IconURL = iconURL.String()
	}

	ccMsg.Activity = fmt.Sprintf("Chat widget configuration updated%s", handlers.IntegrationActivitySuffix(DisplayName))

	titleParts := []string{}
	if len(src.StatusURL) > 0 {
		titleParts = append(titleParts, fmt.Sprintf("[Check status](%s)", src.StatusURL))
	}
	if len(src.TestURL) > 0 {
		titleParts = append(titleParts, fmt.Sprintf("[test widget](%s)", src.TestURL))
	}
	if len(titleParts) > 0 {
		ccMsg.Title = strings.Join(titleParts, " and ")
	}

	attachment := cc.NewAttachment()

	if len(src.Name) > 0 {
		attachment.AddField(cc.Field{
			Title: "Widget Name",
			Value: fmt.Sprintf("[%s](%s)", src.Name, src.CustomURL),
			Short: true})
	}
	attachment.AddField(cc.Field{
		Title: "Widget Version",
		Value: fmt.Sprintf("%v", src.WidgetVersion),
		Short: true})
	if len(src.WidgetExternalType) > 0 {
		attachment.AddField(cc.Field{
			Title: "Widget Type",
			Value: fmt.Sprintf("%v", src.WidgetExternalType),
			Short: true})
	}

	ccMsg.AddAttachment(attachment)
	return ccMsg
}

func NormalizeOperator(cfg config.Configuration, src UserlikeOperatorOutMessage) cc.Message {
	ccMsg := cc.NewMessage()
	iconURL, err := cfg.GetAppIconURL(HandlerKey)
	if err == nil {
		ccMsg.IconURL = iconURL.String()
	}

	ccMsg.Activity = fmt.Sprintf("Operator is %s%s",
		src.Event, handlers.IntegrationActivitySuffix(DisplayName))

	attachment := cc.NewAttachment()

	if len(src.DashboardURL) > 0 {
		attachment.AddField(cc.Field{
			Title: "Operator",
			Value: fmt.Sprintf("[%s](%s)", src.Name, src.DashboardURL)})
	}

	ccMsg.AddAttachment(attachment)
	return ccMsg
}
