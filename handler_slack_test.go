package glipwebhookproxy

import (
	"testing"
)

var SlackWebhookMessageFromBytesTests = []struct {
	v    []byte
	want SlackWebhookMessage
}{
	{[]byte(`{"username":"Ghost Bot [bot]"}`), SlackWebhookMessage{Username: "Ghost Bot [bot]"}}}

func TestSlackWebhookMessageFromBytes(t *testing.T) {
	for _, tt := range SlackWebhookMessageFromBytesTests {
		msg, err := SlackWebhookMessageFromBytes(tt.v)

		if err != nil {
			t.Errorf("SlackWebhookMessageFromBytes(%v): want %v, err %v", tt.v, tt.want, err)
		}

		if tt.want.Username != msg.Username {
			t.Errorf("SlackWebhookMessageFromBytes(%v): want %v, got %v", tt.v, tt.want, msg.Username)
		}
	}
}
