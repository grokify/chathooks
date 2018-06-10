package groupbot

import (
	"regexp"
	"strings"

	"github.com/grokify/gotilla/encoding/jsonutil"
	hum "github.com/grokify/gotilla/net/httputilmore"
	"github.com/grokify/gotilla/strings/stringsutil"
	log "github.com/sirupsen/logrus"
)

/*
type EventResponse struct {
	StatusCode int               `json:"statusCode,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Message    string            `json:"message,omitempty"`
}

func (er *EventResponse) ToJson() []byte {
	if len(er.Message) == 0 {
		er.Message = ""
	}
	msgJson, err := json.Marshal(er)
	if err != nil {
		return []byte(`{"statusCode":500,"message":"Cannot Marshal to JSON"}`)
	}
	return msgJson
}
*/
type IntentRouter struct {
	Intents []Intent
}

func NewIntentRouter() IntentRouter {
	return IntentRouter{Intents: []Intent{}}
}

func (ir *IntentRouter) ProcessRequest(bot *Groupbot, glipPostEventInfo *GlipPostEventInfo) (*hum.ResponseInfo, error) {

	tryCmdsNotMatched := []string{}
	intentResponses := []*hum.ResponseInfo{}

	regexps := []*regexp.Regexp{
		regexp.MustCompile(`[^a-zA-Z0-9\-]+`),
		regexp.MustCompile(`\s+`)}

	tryCmdsLc := stringsutil.SliceCondenseRegexps(
		glipPostEventInfo.TryCommandsLc,
		regexps,
		" ",
	)

	for _, tryCmdLc := range tryCmdsLc {
		matched := false
		for _, intent := range ir.Intents {
			if intent.Type == MatchStringLowerCase {
				for _, try := range intent.Strings {
					if try == tryCmdLc {
						matched = true
						evtResp, err := intent.HandleIntent(bot, glipPostEventInfo)
						if err == nil {
							intentResponses = append(intentResponses, evtResp)
						}
					}
				}
			}
		}
		if !matched {
			tryCmdsNotMatched = append(tryCmdsNotMatched, tryCmdLc)
		}
	}

	tryCmdsNotMatched = stringsutil.SliceCondenseRegexps(
		tryCmdsNotMatched,
		regexps,
		" ",
	)

	if len(tryCmdsNotMatched) > 0 {
		log.Info("TRY_CMDS_NOT_MATCHED " + jsonutil.MustMarshalString(tryCmdsNotMatched, true))
		glipPostEventInfo.TryCommandsLc = tryCmdsNotMatched
		for _, intent := range ir.Intents {
			if intent.Type == MatchAny {
				return intent.HandleIntent(bot, glipPostEventInfo)
			}
		}
	}

	return &hum.ResponseInfo{}, nil
}

func (ir *IntentRouter) ProcessRequestSingle(bot *Groupbot, textNoBotMention string, glipPostEventInfo *GlipPostEventInfo) (*hum.ResponseInfo, error) {
	textNoBotMention = strings.TrimSpace(textNoBotMention)
	textNoBotMentionLc := strings.ToLower(textNoBotMention)
	for _, intent := range ir.Intents {
		if intent.Type == MatchStringLowerCase {
			for _, try := range intent.Strings {
				if try == textNoBotMentionLc {
					return intent.HandleIntent(bot, glipPostEventInfo)
				}
			}
		} else if intent.Type == MatchAny {
			return intent.HandleIntent(bot, glipPostEventInfo)
		}
	}
	return &hum.ResponseInfo{}, nil
}

type IntentType int

const (
	MatchString IntentType = iota
	MatchStringLowerCase
	MatchRegexp
	MatchAny
)

type Intent struct {
	Type         IntentType
	Strings      []string
	Regexps      []*regexp.Regexp
	HandleIntent func(bot *Groupbot, glipPostEventInfo *GlipPostEventInfo) (*hum.ResponseInfo, error)
}
