package glipwebhookproxy

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type EmojiToURL struct {
	EmojiURLPrefix string
	EmojiURLSuffix string
}

func (e2u *EmojiToURL) Convert(emoji string) (string, error) {
	emoji = strings.TrimSpace(emoji)
	if len(emoji) > 0 {
		rx := regexp.MustCompile(`^\s*:?([a-z_]+):?\s*`)
		rs := rx.FindStringSubmatch(emoji)
		if len(rs) > 1 {
			url := fmt.Sprintf("%v%v%v", e2u.EmojiURLPrefix, rs[1], e2u.EmojiURLSuffix)
			return url, nil
		}
	}
	return "", errors.New("No Emoji")
}
