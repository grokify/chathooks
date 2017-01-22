package glipwebhookproxy

import (
	"fmt"
	"testing"
)

var EmojiTests = []struct {
	v    string
	v2   string
	want string
}{
	{"https://grokify.github.io/emoji/assets/images/%s.png",
		":ghost:", "https://grokify.github.io/emoji/assets/images/ghost.png"}}

func TestEmojiURL(t *testing.T) {
	for _, tt := range EmojiTests {
		got, err := EmojiToURL(tt.v, tt.v2)
		if err != nil {
			t.Errorf("EmojiToURL.Convert(%v): want %v, err %v", tt.v, tt.want, err)
		}
		if got != tt.want {
			t.Errorf("EmojiToURL.Convert(%v): want %v, got %v", tt.v, tt.want, got)
		}
	}
}

var EmojiErrorTests = []struct {
	v    string
	v2   string
	want string
}{
	{"%s", ":ghXst:", "No Emoji"}}

func TestEmojiURLError(t *testing.T) {
	for _, tt := range EmojiErrorTests {
		_, err := EmojiToURL(tt.v, tt.v2)
		if err == nil {
			t.Errorf("EmojiToURL.Convert(%v): want %v, err %v", tt.v, tt.want, err)
		}
		if fmt.Sprintf("%v", err) != tt.want {
			t.Errorf("EmojiToURL.Convert(%v): want %v, got %v", tt.v, tt.want, err)
		}
	}
}
