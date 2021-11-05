package slack

import (
	"testing"

	ccslack "github.com/grokify/commonchat/slack"
)

var SlackWebhookMessageFromBytesTests = []struct {
	v    []byte
	want ccslack.Message
}{
	{[]byte(`{"username":"Ghost Bot [bot]"}`), ccslack.Message{Username: "Ghost Bot [bot]"}}}

func TestSlackWebhookMessageFromBytes(t *testing.T) {
	for _, tt := range SlackWebhookMessageFromBytesTests {
		msg, err := ccslack.ParseMessageJSON(tt.v)

		if err != nil {
			t.Errorf("ParseMessageJSON(%v): want %v, err %v", tt.v, tt.want, err)
		}

		if tt.want.Username != msg.Username {
			t.Errorf("ParseMessageJSON(%v): want %v, got %v", tt.v, tt.want, msg.Username)
		}
	}
}
