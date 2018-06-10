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
	Title        string   `json:"title,omitempty"`
	Pretext      string   `json:"pretext,omitempty"`
	Text         string   `json:"text,omitempty"`
	Color        string   `json:"color,omitempty"`
	ThumbnailURL string   `json:"thumbnail_url,omitempty"`
	MrkdwnIn     []string `json:"mrkdwn_in,omitempty"`
	Fields       []Field  `json:"fields,omitempty"`
}

func NewAttachment() Attachment {
	return Attachment{Fields: []Field{}, MrkdwnIn: []string{"text"}}
}

type Field struct {
	Title string `json:"title,omitempty"`
	Value string `json:"value,omitempty"`
	Short bool   `json:"short,omitempty"`
}

func (attach *Attachment) AddField(field Field) {
	attach.Fields = append(attach.Fields, field)
}
