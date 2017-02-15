package util

type Message struct {
	Attachments []Attachment `json:"attachments,omitempty"`
}

func NewMessage() Message {
	return Message{Attachments: []Attachment{}}
}

func (msg *Message) AddAttachment(att Attachment) {
	msg.Attachments = append(msg.Attachments, att)
}

type Attachment struct {
	Title    string   `json:"title,omitempty"`
	Pretext  string   `json:"pretext,omitempty"`
	Text     string   `json:"pretext,omitempty"`
	MrkdwnIn []string `json"mrkdwn_in,omitempty"`
}
