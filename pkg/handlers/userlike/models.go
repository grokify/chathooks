package userlike

import (
	"encoding/json"
)

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
	ID              int64              `json:"id,omitempty"`
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
	var msg UserlikeOfflineMessageOutMessage
	return msg, json.Unmarshal(bytes, &msg)
}

type UserlikeChatWidget struct {
	ID   int64  `json:"id,omitempty"`
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
	ID                int64              `json:"id,omitempty"`
	InitialURL        string             `json:"initial_url,omitempty"`
	LocCity           string             `json:"loc_city,omitempty"`
	LocCountry        string             `json:"loc_country,omitempty"`
	LocLat            float64            `json:"loc_lat,omitempty"`
	LocLon            float64            `json:"loc_lon,omitempty"`
	MarkedRead        bool               `json:"marked_read,omitempty"`
	OperatorCurrentID int64              `json:"operator_current_id,omitempty"`
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
	var msg UserlikeChatMetaStartOutMessage
	return msg, json.Unmarshal(bytes, &msg)
}

type UserlikeOperatorOutMessage struct {
	UserlikeBaseOutMessage
	DashboardURL    string        `json:"dashboard_url,omitempty"`
	Email           string        `json:"email,omitempty"`
	FirstName       string        `json:"first_name,omitempty"`
	ID              int64         `json:"id,omitempty"`
	IsActive        bool          `json:"is_active,omitempty"`
	JID             string        `json:"jid,omitempty"`
	Lang            string        `json:"lang,omitempty"`
	LastName        string        `json:"last_name,omitempty"`
	Locale          string        `json:"locale,omitempty"`
	Name            string        `json:"name,omitempty"`
	OperatorGroup   OperatorGroup `json:"operator_group,omitempty"`
	OperatorGroupID int64         `json:"operator_group_id,omitempty"`
	Role            string        `json:"role,omitempty"`
	RoleName        string        `json:"role_name,omitempty"`
	Timezone        string        `json:"timezone,omitempty"`
	URLImage        string        `json:"url_image,omitempty"`
	Username        string        `json:"username,omitempty"`
}

func UserlikeOperatorOutMessageFromBytes(bytes []byte) (UserlikeOperatorOutMessage, error) {
	var msg UserlikeOperatorOutMessage
	return msg, json.Unmarshal(bytes, &msg)
}

type OperatorGroup struct {
	ID   int64  `json:"id,omitempty"`
	Name string `json:"string,omitempty"`
}

type UserlikeChatWidgetOutMessage struct {
	UserlikeBaseOutMessage
	CustomURL          string `json:"custom_url,omitempty"`
	Name               string `json:"name,omitempty"`
	TransitionDuration int64  `json:"transition_duration,omitempty"`
	StatusURL          string `json:"status_url,omitempty"`
	TestURL            string `json:"test_url,omitempty"`
	WidgetExternalType string `json:"widget_external_type,omitempty"`
	WidgetVersion      int64  `json:"widget_version,omitempty"`
}

func UserlikeChatWidgetOutMessageFromBytes(bytes []byte) (UserlikeChatWidgetOutMessage, error) {
	var msg UserlikeChatWidgetOutMessage
	return msg, json.Unmarshal(bytes, &msg)
}
