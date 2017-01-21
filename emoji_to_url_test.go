package glipwebhookproxy

import (
	"fmt"
	"testing"
)

var EmojiTests = []struct {
	v    string
	want string
}{
	{":ghost:", "https://grokify.github.io/emoji/assets/images/ghost.png"}}

func TestEmojiURL(t *testing.T) {
	converter := EmojiToURL{
		EmojiURLPrefix: "https://grokify.github.io/emoji/assets/images/",
		EmojiURLSuffix: ".png"}

	for _, tt := range EmojiTests {
		got, err := converter.Convert(tt.v)
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
	want string
}{
	{":ghXst:", "No Emoji"}}

func TestEmojiURLError(t *testing.T) {
	converter := EmojiToURL{
		EmojiURLPrefix: "https://grokify.github.io/emoji/assets/images/",
		EmojiURLSuffix: ".png"}

	for _, tt := range EmojiErrorTests {
		_, err := converter.Convert(tt.v)
		if err == nil {
			t.Errorf("EmojiToURL.Convert(%v): want %v, err %v", tt.v, tt.want, err)
		}
		if fmt.Sprintf("%v", err) != tt.want {
			t.Errorf("EmojiToURL.Convert(%v): want %v, got %v", tt.v, tt.want, err)
		}
	}
}
