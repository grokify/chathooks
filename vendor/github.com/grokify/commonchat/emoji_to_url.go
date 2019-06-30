package commonchat

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// EmojiToURL function takes a `Sprintf` `format` string and emoji string with
// or without surrounding colons (`:`) and returns a URL. Emoji strings must
// satisfy `[a-z_]+` regexp.
func EmojiToURL(urlFormat string, emoji string) (string, error) {
	emoji = strings.TrimSpace(emoji)
	if len(emoji) > 0 {
		rx := regexp.MustCompile(`^:?([a-z_]+):?$`)
		rs := rx.FindStringSubmatch(emoji)
		if len(rs) > 1 {
			return fmt.Sprintf(urlFormat, rs[1]), nil
		}
	}
	return "", errors.New("No Emoji")
}
