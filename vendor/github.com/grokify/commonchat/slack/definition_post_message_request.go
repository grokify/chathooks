package slack

type CommonPost struct {
	Channel     string       `json:"channel,omitempty"`
	Username    string       `json:"username,omitempty"`
	IconEmoji   string       `json:"icon_emoji,omitempty"`
	LinkNames   int          `json:"link_names,omitempty"`
	Text        string       `json:"text,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

type CreatePostRequest struct {
	CommonPost
	AsUser         bool    `json:"as_user,omitempty"`
	IconUrl        string  `json:"icon_url,omitempty"`
	LinkNames      bool    `json:"link_names,omitempty"`
	Parse          bool    `json:"parse,omitempty"`
	ReplyBroadcast bool    `json:"reply_broadcast,omitempty"`
	UnfurlLinks    bool    `json:"unfurl_links,omitempty"`
	UnfurlMedia    bool    `json:"unfurl_media,omitempty"`
	ThreadTs       float64 `json:"thread_ts,omitempty"`
}

type WebhookRequest struct {
	CommonPost
}
