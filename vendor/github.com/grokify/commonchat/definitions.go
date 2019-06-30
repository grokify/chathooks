package commonchat

type Message struct {
	Activity    string       `json:"activity,omitempty"`
	IconEmoji   string       `json:"icon_emoji,omitempty"`
	IconURL     string       `json:"icon_url,omitempty"`
	Title       string       `json:"title,omitempty"`
	Text        string       `json:"text,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

func NewMessage() Message {
	return Message{Attachments: []Attachment{}}
}

func (msg *Message) AddAttachment(att Attachment) {
	msg.Attachments = append(msg.Attachments, att)
}

type Attachment struct {
	AuthorIcon   string   `json:"author_icon,omitempty"`
	AuthorLink   string   `json:"author_link,omitempty"`
	AuthorName   string   `json:"author_name,omitempty"`
	Color        string   `json:"color,omitempty"`
	Fallback     string   `json:"fallback,omitempty"`
	Fields       []Field  `json:"fields,omitempty"`
	MarkdownIn   []string `json:"mrkdwn_in,omitempty"`
	Pretext      string   `json:"pretext,omitempty"`
	Text         string   `json:"text,omitempty"`
	ThumbnailURL string   `json:"thumbnail_url,omitempty"`
	Title        string   `json:"title,omitempty"`
}

func NewAttachment() Attachment {
	return Attachment{Fields: []Field{}, MarkdownIn: []string{"text"}}
}

type Field struct {
	Title string `json:"title,omitempty"`
	Value string `json:"value,omitempty"`
	Short bool   `json:"short,omitempty"`
}

func (attach *Attachment) AddField(field Field) {
	attach.Fields = append(attach.Fields, field)
}
