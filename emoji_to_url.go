package glipwebhookproxy

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// EmojiToURL enables converting emoji strings Set up with a `EmojiURLPrefix`
// and `EmojiURLSuffix`.
type EmojiToURL struct {
	EmojiURLPrefix string
	EmojiURLSuffix string
}

// Convert function takes an emoji string with or without surrounding
// colons (`:`) and returns a URL. Emoji strings must satisfy `[a-z_]+` regexp.
func (e2u *EmojiToURL) Convert(emoji string) (string, error) {
	emoji = strings.TrimSpace(emoji)
	if len(emoji) > 0 {
		rx := regexp.MustCompile(`^:?([a-z_]+):?$`)
		rs := rx.FindStringSubmatch(emoji)
		if len(rs) > 1 {
			url := fmt.Sprintf("%v%v%v", e2u.EmojiURLPrefix, rs[1], e2u.EmojiURLSuffix)
			return url, nil
		}
	}
	return "", errors.New("No Emoji")
}
