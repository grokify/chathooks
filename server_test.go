package glipwebhookproxy

import (
	"testing"
)

var ConfigurationTests = []struct {
	v    int
	want string
}{
	{8080, ":8080"}}

func TestConfigurationAddress(t *testing.T) {
	for _, tt := range ConfigurationTests {
		config := Configuration{
			Port: tt.v}

		addr := config.Address()
		if tt.want != addr {
			t.Errorf("Configuration.Address(%v): want %v, got %v", tt.v, tt.want, addr)
		}
	}
}

var ConfigurationEmojiTests = []struct {
	test1 string
	true1 string
	test2 string
	true2 string
}{
	{"https://grokify.gitub.io/", "https://grokify.gitub.io/", ".png", ".png"}}

func TestConfigurationEmoji(t *testing.T) {
	for _, tt := range ConfigurationEmojiTests {
		config := Configuration{
			EmojiURLPrefix: tt.test1,
			EmojiURLSuffix: tt.test1}

		if tt.true1 != config.EmojiURLPrefix {
			t.Errorf("Configuration.EmojiURLPrefix.%v: want %v, got %v", tt.test1, tt.true1, config.EmojiURLPrefix)
		}
		if tt.true1 != config.EmojiURLSuffix {
			t.Errorf("Configuration.EmojiURLSuffix.%v: want %v, got %v", tt.test2, tt.true2, config.EmojiURLSuffix)
		}
	}
}
