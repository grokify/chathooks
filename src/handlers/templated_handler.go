package handlers

import (
	"strings"
	"regexp"
    "encoding/json"
    "strconv"

	"github.com/grokify/chathooks/src/config"
    cc "github.com/grokify/commonchat"
    "github.com/tidwall/gjson"
)

func NewTemplatedHandler(tmpl string) Handler {
	return Handler {
        Normalize: getTemplatedNormalizer(tmpl),
    }
}

func getTemplatedNormalizer(tmpl string) func(cfg config.Configuration, bytes []byte) (cc.Message, error) {
    return func(cfg config.Configuration, bytes []byte) (cc.Message, error) {
    	ccMsg := cc.NewMessage()
    	src := string(bytes)

        tokenPattern := regexp.MustCompile(`\${.+?}`)
        keyPattern := regexp.MustCompile(`\${(.+?)}`)
        formattedJson := tokenPattern.ReplaceAllStringFunc(tmpl, func(match string) string {
            matches := keyPattern.FindStringSubmatch(match)
            result := gjson.Get(src, strings.TrimSpace(matches[1]))
            switch result.Type {
            case gjson.String:
                return result.Str
            case gjson.Number:
                return strconv.FormatFloat(result.Num, 'f', -1, 64)
            case gjson.JSON:
                return result.Raw
            default:
                return result.Type.String()
            }
        })

        err := json.Unmarshal([]byte(formattedJson), &ccMsg)
        return ccMsg, err
    }
}