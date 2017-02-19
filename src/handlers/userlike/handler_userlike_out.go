package userlike

import (
	"encoding/json"
	"errors"
	"fmt"

	log "github.com/Sirupsen/logrus"

	cc "github.com/grokify/commonchat"
	"github.com/grokify/glip-webhook-proxy-go/src/adapters"
	"github.com/grokify/glip-webhook-proxy-go/src/config"
	"github.com/grokify/glip-webhook-proxy-go/src/util"
	"github.com/valyala/fasthttp"
)

const (
	DisplayName = "Userlike"
	HandlerKey  = "userlike"
	IconURL     = "https://a.slack-edge.com/ae7f/img/services/userlike_512.png"
)

var (
	ChatMetaEvents = []string{"feedback", "forward", "rating", "receive", "start", "survey"}
	OperatorEvents = []string{"away", "back", "offline", "online"}
)

// FastHttp request handler for Userlike outbound webhook
type UserlikeOutToGlipHandler struct {
	Config  config.Configuration
	Adapter adapters.Adapter
}

// FastHttp request handler constructor for Userlike outbound webhook
func NewUserlikeOutToGlipHandler(cfg config.Configuration, adapter adapters.Adapter) UserlikeOutToGlipHandler {
	return UserlikeOutToGlipHandler{Config: cfg, Adapter: adapter}
}

// HandleFastHTTP is the method to respond to a fasthttp request.
func (h *UserlikeOutToGlipHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
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

func Normalize(bodyBytes []byte) (cc.Message, error) {
	srcMsgBase, err := UserlikeBaseOutMessageFromBytes(bodyBytes)
	if err != nil {
		return cc.Message{}, err
	}
	if srcMsgBase.Type == "offline_message" && srcMsgBase.Event == "receive" {
		srcMsg, err := UserlikeOfflineMessageOutMessageFromBytes(bodyBytes)
		if err != nil {
			return cc.Message{}, err
		}
		return NormalizeOfflineMessage(srcMsg), nil
	} else if srcMsgBase.Type == "chat_meta" {
		srcMsg, err := UserlikeChatMetaStartOutMessageFromBytes(bodyBytes)
		if err != nil {
			return cc.Message{}, err
		}
		return NormalizeChatMeta(srcMsg), nil
	} else if srcMsgBase.Type == "operator" {
		srcMsg, err := UserlikeOperatorOutMessageFromBytes(bodyBytes)
		if err != nil {
			return cc.Message{}, err
		}
		return NormalizeOperator(srcMsg), nil
	} else if srcMsgBase.Type == "chat_widget" {
		srcMsg, err := UserlikeChatWidgetOutMessageFromBytes(bodyBytes)
		if err != nil {
			return cc.Message{}, err
		}
		return NormalizeChatWidget(srcMsg), nil
	}

	return cc.Message{}, errors.New("Type Not Supported")
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
	return fmt.Sprintf("%s chat %s", DisplayName, eventDisplay)
}

func NormalizeOfflineMessage(src UserlikeOfflineMessageOutMessage) cc.Message {
	message := cc.NewMessage()
	message.IconURL = IconURL

	clientName := src.ClientName
	if len(clientName) < 1 {
		clientName = "Website Visitor"
	}
	message.Activity = fmt.Sprintf("%s sent a new Offline Message%v", clientName, adapters.IntegrationActivitySuffix(DisplayName))

	attachment := cc.NewAttachment()
	attachment.ThumbnailURL = IconURL

	if len(src.URL) > 0 {
		attachment.AddField(cc.Field{
			Value: fmt.Sprintf("[View message](%v)", src.URL)})
	}

	if len(attachment.Fields) > 0 {
		message.AddAttachment(attachment)
	}
	return message
}

func NormalizeChatMeta(src UserlikeChatMetaStartOutMessage) cc.Message {
	message := cc.NewMessage()
	message.IconURL = IconURL
	message.Activity = fmt.Sprintf("%s%s",
		GlipActivityForChat(src.Event, src.FeedbackMessage), adapters.IntegrationActivitySuffix(DisplayName))

	attachment := cc.NewAttachment()
	attachment.ThumbnailURL = IconURL

	if src.Event == "rating" || src.Event == "survey" { // includes feedback
		if len(src.FeedbackMessage) > 0 {
			attachment.AddField(cc.Field{
				Title: "Feedback",
				Value: src.FeedbackMessage,
				Short: false})
		}
		if len(src.PostSurveyOption) > 0 {
			attachment.AddField(cc.Field{
				Title: "Rating",
				Value: src.PostSurveyOption,
				Short: true})
		}
	}
	if len(src.ClientName) > 0 {
		attachment.AddField(cc.Field{
			Title: "Client Name",
			Value: src.ClientName,
			Short: true})
	} else {
		attachment.AddField(cc.Field{
			Title: "Client Name",
			Value: "Unknown",
			Short: true})
	}

	if len(src.URL) > 0 {
		attachment.AddField(cc.Field{
			Value: fmt.Sprintf("[View details](%v)", src.URL)})
	}

	if len(attachment.Fields) > 0 {
		message.AddAttachment(attachment)
	}
	return message
}

func NormalizeChatWidget(src UserlikeChatWidgetOutMessage) cc.Message {
	message := cc.NewMessage()
	message.IconURL = IconURL
	message.Activity = fmt.Sprintf("Chat widget configuration updated%s", adapters.IntegrationActivitySuffix(DisplayName))

	attachment := cc.NewAttachment()
	attachment.ThumbnailURL = IconURL

	if len(src.Name) > 0 {
		attachment.AddField(cc.Field{
			Title: "Widget Name",
			Value: fmt.Sprintf("[%s](%s)", src.Name, src.CustomUrl),
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
	if len(src.StatusUrl) > 0 {
		attachment.AddField(cc.Field{
			Value: fmt.Sprintf("[Widget Status](%v)", src.StatusUrl),
			Short: false})
	}
	if len(src.TestUrl) > 0 {
		attachment.AddField(cc.Field{
			Value: fmt.Sprintf("[Test Widget](%v)", src.TestUrl),
			Short: false})
	}

	message.AddAttachment(attachment)
	return message
}

func NormalizeOperator(src UserlikeOperatorOutMessage) cc.Message {
	message := cc.NewMessage()
	message.IconURL = IconURL
	message.Activity = fmt.Sprintf("%s (operator) is %s%s",
		src.Name, src.Event, adapters.IntegrationActivitySuffix(DisplayName))
	/*
		gmsg := glipwebhook.GlipWebhookMessage{
			Activity: fmt.Sprintf("%s is %s as operator%s",
				src.Name, src.Event, adapters.IntegrationActivitySuffix(DISPLAY_NAME)),
			Icon: ICON_URL}
	*/
	//message := util.NewMessage()

	attachment := cc.NewAttachment()

	if len(src.DashboardUrl) > 0 {
		attachment.AddField(cc.Field{
			Value: fmt.Sprintf("[Operator Details](%v)", src.DashboardUrl)})
	}

	//gmsg.Body = adapters.RenderMessage(message)
	message.AddAttachment(attachment)
	return message
}

type UserlikeBaseOutMessage struct {
	Event string `json:"_event,omitempty"`
	Type  string `json:"_type,omitempty"`
}

func UserlikeBaseOutMessageFromBytes(bytes []byte) (UserlikeBaseOutMessage, error) {
	msg := UserlikeBaseOutMessage{}
	err := json.Unmarshal(bytes, &msg)
	return msg, err
}

type UserlikeOfflineMessageOutMessage struct {
	UserlikeBaseOutMessage
	BrowserName     string             `json:"browser_name,omitempty"`
	BrowserOS       string             `json:"browser_os,omitempty"`
	BrowserVersion  string             `json:"browser_version,omitempty"`
	ChatWidget      UserlikeChatWidget `json:"chat_widget,omitempty"`
	ClientEmail     string             `json:"client_email,omitempty"`
	ClientName      string             `json:"client_name,omitempty"`
	CreatedAt       string             `json:"created_at,omitempty"`
	Custom          interface{}        `json:"custom,omitempty"`
	DataPrivacy     interface{}        `json:"data_privacy,omitempty"`
	Id              int64              `json:"id,omitempty"`
	LocCity         string             `json:"loc_city,omitempty"`
	LocCountry      string             `json:"loc_country,omitempty"`
	LocLat          float64            `json:"loc_lat,omitempty"`
	LocLon          float64            `json:"loc_lon,omitempty"`
	MarkedRead      bool               `json:"marked_read,omitempty"`
	Message         string             `json:"message,omitempty"`
	PageImpressions int64              `json:"page_impresions,omitempty"`
	ScreenshotOID   string             `json:"screenshot_oid,omitempty"`
	ScreenshotURL   string             `json:"screenshot_url,omitempty"`
	Status          string             `json:"status,omitempty"`
	Topic           string             `json:"topic,omitempty"`
	URL             string             `json:"url,omitempty"`
	Visits          int64              `json:"visits,omitempty"`
}

func UserlikeOfflineMessageOutMessageFromBytes(bytes []byte) (UserlikeOfflineMessageOutMessage, error) {
	msg := UserlikeOfflineMessageOutMessage{}
	err := json.Unmarshal(bytes, &msg)
	return msg, err
}

type UserlikeChatWidget struct {
	Id   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type UserlikeChatMetaStartOutMessage struct {
	UserlikeBaseOutMessage
	BrowserName       string             `json:"browser_name,omitempty"`
	BrowserOS         string             `json:"browser_os,omitempty"`
	BrowserVersion    string             `json:"browser_version,omitempty"`
	ChatWidget        UserlikeChatWidget `json:"chat_widget,omitempty"`
	ClientEmail       string             `json:"client_email,omitempty"`
	ClientName        string             `json:"client_name,omitempty"`
	ClientUUID        string             `json:"client_uuid,omitempty"`
	CreatedAt         string             `json:"created_at,omitempty"`
	DataPrivacy       bool               `json:"data_privacy,omitempty"`
	Duration          string             `json:"duration,omitempty"`
	EndedAt           string             `json:"ended_at,omitempty"`
	FeedbackMessage   string             `json:"feedback_message,omitempty"`
	Id                int64              `json:"id,omitempty"`
	InitialURL        string             `json:"initial_url,omitempty"`
	LocCity           string             `json:"browser_version,omitempty"`
	LocCountry        string             `json:"browser_version,omitempty"`
	LocLat            float64            `json:"loc_lat,omitempty"`
	LocLon            float64            `json:"loc_lon,omitempty"`
	MarkedRead        bool               `json:"marked_read,omitempty"`
	OperatorCurrentId int64              `json:"operator_current_id,omitempty"`
	PageImpressions   int64              `json:"page_impressions,omitempty"`
	PostSurveyOption  string             `json:"post_survey_option,omitempty"`
	Rate              int64              `json:"rate,omitempty"`
	Referrer          string             `json:"referrer,omitempty"`
	Status            string             `json:"status,omitempty"`
	Topic             string             `json:"topic,omitempty"`
	URL               string             `json:"url,omitempty"`
	Visits            int64              `json:"visits,omitempty"`
	WasProactive      bool               `json:"was_proactive,omitempty"`
}

func UserlikeChatMetaStartOutMessageFromBytes(bytes []byte) (UserlikeChatMetaStartOutMessage, error) {
	msg := UserlikeChatMetaStartOutMessage{}
	err := json.Unmarshal(bytes, &msg)
	return msg, err
}

type UserlikeOperatorOutMessage struct {
	UserlikeBaseOutMessage
	DashboardUrl    string        `json:"dashboard_url,omitempty"`
	Email           string        `json:"email,omitempty"`
	FirstName       string        `json:"first_name,omitempty"`
	Id              int64         `json:"id,omitempty"`
	IsActive        bool          `json:"is_active,omitempty"`
	JID             string        `json:"jid,omitempty"`
	Lang            string        `json:"lang,omitempty"`
	LastName        string        `json:"last_name,omitempty"`
	Locale          string        `json:"locale,omitempty"`
	Name            string        `json:"name,omitempty"`
	OperatorGroup   OperatorGroup `json:"operator_group,omitempty"`
	OperatorGroupId int64         `json:"operator_group_id,omitempty"`
	Role            string        `json:"role,omitempty"`
	RoleName        string        `json:"role_name,omitempty"`
	Timezone        string        `json:"timezone,omitempty"`
	UrlImage        string        `json:"url_image,omitempty"`
	Username        string        `json:"username,omitempty"`
}

func UserlikeOperatorOutMessageFromBytes(bytes []byte) (UserlikeOperatorOutMessage, error) {
	msg := UserlikeOperatorOutMessage{}
	err := json.Unmarshal(bytes, &msg)
	return msg, err
}

type OperatorGroup struct {
	Id   int64  `json:"id,omitempty"`
	Name string `json:"string,omitempty"`
}

type UserlikeChatWidgetOutMessage struct {
	UserlikeBaseOutMessage
	CustomUrl          string `json:"custom_url,omitempty"`
	Name               string `json:"name,omitempty"`
	TransitionDuration int64  `json:"transition_duration,omitempty"`
	StatusUrl          string `json:"status_url,omitempty"`
	TestUrl            string `json:"test_url,omitempty"`
	WidgetExternalType string `json:"widget_external_type,omitempty"`
	WidgetVersion      int64  `json:"widget_version,omitempty"`
}

func UserlikeChatWidgetOutMessageFromBytes(bytes []byte) (UserlikeChatWidgetOutMessage, error) {
	msg := UserlikeChatWidgetOutMessage{}
	err := json.Unmarshal(bytes, &msg)
	return msg, err
}

/*
{
    "_event":"config",
    "_type":"chat_widget",
    "audio_message_received":true,
    "audio_message_sent":false,
    "audio_only_inactive":true,
    "chat_inactive_timeout":600,
    "chat_is_typing":"{{name}} is typing....",
    "chat_me":"Me",
    "chat_stopped_typing":"{{name}} stopped typing",
    "colorscheme_hi":"rgba(0,174,239,0.67)",
    "colorscheme_lo":"rgba(0,174,239,0.21)",
    "cookie_expire":365,
    "custom_css":"",
    "custom_header":"",
    "custom_url":"https://devel.userlike.local/custom/9",
    "customer_id":2,
    "data_privacy":false,
    "data_privacy_link":"https://userlike.com",
    "data_privacy_name":"Data Privacy",
    "default_message":"Welcome: {{name}}",
    "default_textarea":"Type your message here...",
    "delete_empty_transcripts":true,
    "disclaimer_headline":"Disclaimer",
    "disclaimer_show":false,
    "disclaimer_text":"This website uses Userlike, a live chat software provided by Userlike UG. Userlike uses \"Cookies\", which are text files placed on your computer and that enable having personal chats on this website. The data collected will not be used to identify a visitor personally and is it not aggregated with any personal data of this user.",
    "drag_enabled":true,
    "emit_chat_state":true,
    "facebook_app_id":"505255409521592",
    "facebook_connect":"disabled",
    "facebook_integration":"load",
    "facebook_like_enabled":true,
    "facebook_like_headline":"Like us on Facebook",
    "facebook_like_href":"https://www.userlike.com",
    "facebook_like_layout":"standard",
    "facebook_like_show_faces":false,
    "facebook_like_verb":"like",
    "favicon_enabled":true,
    "feedback_download_link":"Download Link: ",
    "feedback_error":"An error has occured. Please try again later.",
    "feedback_expired":"This account expired.",
    "feedback_invite_video_meeting":"You are invited to join a online video meeting. Please go to the meeting url:",
    "feedback_no_cookies":"To use this chat tool, cookies must be enabled in your browser.",
    "feedback_offline":"Sorry, we are offline right now.",
    "feedback_quota":"This free account reached the monthly quota limit.",
    "feedback_transfer":"Your chat window has been transferred to another browser window.",
    "font":"sans-serif",
    "force_ssl":false,
    "forward_message":"Your chat has been successfully forwarded. {{name}} is happy to provide you with further assistance.",
    "geo_type":"ip",
    "goals":[

    ],
    "group_select_enabled":false,
    "group_select_headline":"Select a group you want to chat to",
    "hide_button":false,
    "hide_poweredby":false,
    "id":9,
    "inactivity_action":false,
    "inactivity_message":"Sorry to keep you waiting, I will be with you in a minute.",
    "inactivity_mode":"none",
    "inactivity_timeout":60,
    "is_default":false,
    "javascript_loader_snippet":"\n<script type=\"text/javascript\" src=\"//devel.userlike.local/static/chat/widgets/4d70b54030ded794aac80a14ea886d61f998ece75d78273940b3eb64bbece074.js\"></script>\n",
    "javascript_loader_url":"devel.userlike.local/static/chat/widgets/4d70b54030ded794aac80a14ea886d61f998ece75d78273940b3eb64bbece074.js",
    "lang":"en",
    "links_offsite":"new",
    "links_onsite":"same",
    "locale":"en_US",
    "logo_mode":"userlike",
    "logo_upload_url":"devel-cdn-logos.s3-eu-west-1.amazonaws.com/4d70b54030ded794aac80a14ea886d61f998ece75d78273940b3eb64bbece074.png",
    "logo_url":"https://www.google.com/a/cpanel/userlike.com/images/logo.gif",
    "mobile_mode":"enabled",
    "mobile_reduced_size":true,
    "mode_proactive":false,
    "mode_registration":false,
    "mode_remote":false,
    "name":"Testing David",
    "offline_message_body":"Please leave a message and we contact you very soon.",
    "offline_message_default_textarea":"Leave a message...",
    "offline_message_enter_email":"Enter your email here",
    "offline_message_enter_name":"Enter your name here",
    "offline_message_header":"We are not here to chat right now",
    "offline_message_response":"Thanks for your request. We get back to you soon.",
    "offline_message_send_screenshot":"Send Screenshot",
    "offline_mode":"message",
    "operator_group_chat_meta":false,
    "operator_group_id":14,
    "operator_group_offline_message":false,
    "operator_picture_style":"square",
    "optional_registration":false,
    "orientation":"bottomRight",
    "passive_connect":false,
    "post_survey_enabled":true,
    "post_survey_option01":"Not at all satisfied",
    "post_survey_option02":"Somewhat Satisfied",
    "post_survey_option03":"Satisfied",
    "post_survey_option04":"Very Satisfied",
    "post_survey_question":"Are you satisfied with our service?",
    "pre_survey_enabled":true,
    "pre_survey_option01":"Question about a product",
    "pre_survey_option02":"Check shipping and delivery status",
    "pre_survey_option03":"Technical support",
    "pre_survey_option04":"",
    "pre_survey_question":"What is the topic of your chat?",
    "proactive_message":"Proactive: {{name}}",
    "proactive_passive_connect":true,
    "proactive_timeout":5,
    "quit_message":"{{name}} left the chat.",
    "rating_enabled":true,
    "rating_question":"Was this chat helpful?",
    "register_additional01_default":"Enter your phone number here",
    "register_additional01_enabled":false,
    "register_additional01_name":"Phone number",
    "register_additional01_optional":true,
    "register_additional02_default":"Enter your company name here",
    "register_additional02_enabled":false,
    "register_additional02_name":"Company name",
    "register_additional02_optional":true,
    "register_additional03_default":"Enter your customer id here",
    "register_additional03_enabled":false,
    "register_additional03_name":"Customer ID",
    "register_additional03_optional":true,
    "register_body":"Please enter your name and email",
    "register_enter_email":"Enter your email here",
    "register_enter_name":"Enter your name here",
    "register_header":"Welcome to our live-chat",
    "screenshot_command_enabled":true,
    "screenshot_enabled":true,
    "show_data_privacy_link":false,
    "status_url":"https://devel.userlike.local/status/9",
    "tab_color_css":"#009dd6",
    "tab_color_png":"green",
    "tab_form_css":"big",
    "tab_icon":"comments-alt",
    "tab_text":"live",
    "tab_text_label":"Do you have questions?",
    "tab_type":"css",
    "test_url":"https://devel.userlike.local/test/9",
    "text_color":"rgba(169,176,183,1)",
    "tracking":"disabled",
    "transition_duration":100,
    "twitter_follow_show_username":true,
    "twitter_follow_user":"userlike",
    "twitter_headline":"Tweet about us",
    "twitter_integration":"load",
    "twitter_mention_text":"Great support",
    "twitter_mention_user":"userlike",
    "twitter_share_hashtag":"livechat",
    "twitter_share_show_count":true,
    "twitter_share_text":"I love Userlike live chat software!",
    "twitter_share_url":"",
    "twitter_share_via":"userlike",
    "twitter_type":"share",
    "widget_external_type":"web",
    "widget_id":9,
    "widget_key":"4d70b54030ded794aac80a14ea886d61f998ece75d78273940b3eb64bbece074",
    "widget_type":"d",
    "widget_version":2,
    "wildcard_cookie":false
}


*/
