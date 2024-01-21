package wootric

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/grokify/mogo/time/timeutil"
	"github.com/grokify/mogo/type/maputil"
	"github.com/grokify/xgo/net/urlutil"
)

type WootricEvent struct {
	Response      WootricResponse `json:"response"`
	Decline       WootricDecline  `json:"decline"`
	EventName     string          `json:"event_name"`
	Timestamp     string          `json:"timestamp"`
	AccountToken  string          `json:"account_token"`
	SurveyMode    string          `json:"survey_mode"`
	TimestampTime time.Time
}

func (we *WootricEvent) IsResponse() bool {
	return len(strings.TrimSpace(we.Response.Email)) > 0
}

func (we *WootricEvent) IsDecline() bool {
	return len(strings.TrimSpace(we.Decline.Email)) > 0
}

func (we *WootricEvent) Activity() string {
	parts := []string{"NPS "}
	if we.IsResponse() {
		parts = append(parts, "Response")
	} else if we.IsDecline() {
		parts = append(parts, "Decline")
	}
	if 1 == 0 {
		evtName := strings.TrimSpace(we.EventName)
		if len(evtName) > 0 {
			parts = append(parts, evtName)
		}
	}
	return strings.Join(parts, " ")
}

func ParseQueryString(raw string) (WootricEvent, error) {
	var evt WootricEvent
	err := urlutil.UnmarshalRailsQS(raw, &evt)
	if err != nil {
		return evt, err
	}
	if len(strings.TrimSpace(evt.Timestamp)) > 0 {
		t1, err := time.Parse(timeutil.Ruby, strings.TrimSpace(evt.Timestamp))
		if err != nil {
			return evt, err
		}
		evt.TimestampTime = t1
	}
	if len(strings.TrimSpace(evt.Response.CreatedAt)) > 0 {
		t2, err := time.Parse(timeutil.Ruby, strings.TrimSpace(evt.Response.CreatedAt))
		if err != nil {
			return evt, err
		}
		evt.Response.CreatedAtTime = t2
	}
	if len(strings.TrimSpace(evt.Response.UpdatedAt)) > 0 {
		t3, err := time.Parse(timeutil.Ruby, strings.TrimSpace(evt.Response.UpdatedAt))
		if err != nil {
			return evt, err
		}
		evt.Response.UpdatedAtTime = t3
	}
	if len(strings.TrimSpace(evt.Decline.CreatedAt)) > 0 {
		t2, err := time.Parse(timeutil.Ruby, strings.TrimSpace(evt.Decline.CreatedAt))
		if err != nil {
			return evt, err
		}
		evt.Decline.CreatedAtTime = t2
	}
	if len(strings.TrimSpace(evt.Decline.UpdatedAt)) > 0 {
		t3, err := time.Parse(timeutil.Ruby, strings.TrimSpace(evt.Decline.UpdatedAt))
		if err != nil {
			return evt, err
		}
		evt.Decline.UpdatedAtTime = t3
	}
	return evt, err
}

type WootricResponse struct {
	ID                      string            `json:"id"`
	Email                   string            `json:"email"`
	ExternalID              string            `json:"external_id"`
	Score                   json.Number       `json:"score"`
	Text                    string            `json:"text"`
	IPAddress               string            `json:"ip_address"`
	OriginURL               string            `json:"origin_url"`
	EndUserID               string            `json:"end_user_id"`
	EndUserProperties       map[string]string `json:"end_user_properties"`
	SurveyID                string            `json:"survey_id"`
	CreatedAt               string            `json:"created_at"`
	UpdatedAt               string            `json:"updated_at"`
	ExcludeFromCalculations string            `json:"excluded_from_calculations"`
	CreatedAtTime           time.Time
	UpdatedAtTime           time.Time
}

func (resp *WootricResponse) Property(key string) string {
	if resp.EndUserProperties == nil {
		resp.EndUserProperties = map[string]string{}
	}
	return maputil.ValueStringOrDefault(resp.EndUserProperties, key, "")
}

type WootricDecline struct {
	ID                json.Number       `json:"id"`
	Email             string            `json:"email"`
	ExternalID        string            `json:"external_id"`
	IPAddress         string            `json:"ip_address"`
	OriginURL         string            `json:"origin_url"`
	EndUserID         string            `json:"end_user_id"`
	EndUserProperties map[string]string `json:"end_user_properties"`
	SurveyID          string            `json:"survey_id"`
	CreatedAt         string            `json:"created_at"`
	UpdatedAt         string            `json:"updated_at"`
	CreatedAtTime     time.Time
	UpdatedAtTime     time.Time
}

/*
decline[id]=19&
decline[email]=nps@example.com&
decline[external_id]=123abc&
decline[ip_address]=127.0.0.1&
decline[origin_url]=https%3A%2F%2Fwootric.com%2F&
decline[end_user_id]=31&
decline[end_user_properties][pricing_plan]=Pro&
decline[end_user_properties][product_plan]=Web%20App&
decline[survey_id]=1147&
decline[created_at]=2016-08-04%2013%3A58%3A21%20-0700&
decline[updated_at]=2016-08-04%2013%3A58%3


decline[id]=19&decline[email]=nps@example.com&decline[external_id]=123abc&decline[ip_address]=127.0.0.1&decline[origin_url]=https%3A%2F%2Fwootric.com%2F&decline[end_user_id]=31&decline[end_user_properties][pricing_plan]=Pro&decline[end_user_properties][product_plan]=Web%20App&decline[survey_id]=1147&decline[created_at]=2016-08-04%2013%3A58%3A21%20-0700&decline[updated_at]=2016-08-04%2013%3A58%3A21%20-0700&event_name=created&account_token=NPS-xxxxxxx&survey_mode=NPS&timestamp=2016-08-04%2013%3A58%3A23%20-0700
*/
