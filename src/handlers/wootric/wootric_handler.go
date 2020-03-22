package wootric

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/grokify/chathooks/src/config"
	"github.com/grokify/chathooks/src/handlers"
	"github.com/grokify/chathooks/src/models"
	cc "github.com/grokify/commonchat"
	"github.com/grokify/gotilla/fmt/fmtutil"
	"github.com/grokify/gotilla/html/htmlutil"
)

const (
	DisplayName      = "Wootric"
	HandlerKey       = "wootric"
	MessageDirection = "out"
	MessageBodyType  = models.URLEncodedRails

	WootricResponseFormat = `score[Score],text[Why];company_name__c[Company Name],(rcAccountId__c[RC Account ID]);email[User email];directorySize__c[Number of users];brand[Brand];survey_id[Survey ID]`
)

func NewHandler() handlers.Handler {
	return handlers.Handler{MessageBodyType: MessageBodyType, Normalize: Normalize}
}

func Normalize(cfg config.Configuration, hReq handlers.HandlerRequest) (cc.Message, error) {
	ccMsg := cc.NewMessage()
	iconURL, err := cfg.GetAppIconURL(HandlerKey)
	if err == nil {
		ccMsg.IconURL = iconURL.String()
	}

	fmtutil.PrintJSON(ccMsg)

	src, err := ParseQueryString(string(hReq.Body))
	if err != nil {
		return ccMsg, err
	}

	ccMsg.Activity = src.Activity()
	fmtutil.PrintJSON(src)

	if src.IsResponse() {
		attachment := cc.NewAttachment()
		lines := ParseFields(WootricResponseFormat)
		fmtutil.PrintJSON(lines)

		scoreInt64, err := src.Response.Score.Int64()
		if err == nil {
			if scoreInt64 >= 9 {
				attachment.Color = htmlutil.Color2GreenHex
			} else if scoreInt64 >= 7 {
				attachment.Color = htmlutil.Color2YellowHex
			} else {
				attachment.Color = htmlutil.Color2RedHex
			}
		}

		for _, line := range lines {
			numFields := len(line.Fields)
			if numFields == 0 {
				continue
			} /*
				isShort := true
				if numFields == 1 {
					isShort = false
				}*/

			for _, field := range line.Fields {
				if field.Property == "score" {
					fmtutil.PrintJSON(src.Response)
					val := strings.TrimSpace(src.Response.Score.String())
					attachment.AddField(cc.Field{Title: field.Display, Value: val})
				} else if field.Property == "text" {
					attachment.AddField(cc.Field{
						Title: field.Display,
						Value: src.Response.Text})
				} else if field.Property == "email" {
					attachment.AddField(cc.Field{
						Title: field.Display,
						Value: src.Response.Email})
				} else if field.Property == "survey_id" {
					attachment.AddField(cc.Field{
						Title: field.Display,
						Value: src.Response.SurveyID})
				}
			}
		}
		if len(attachment.Fields) > 0 {
			ccMsg.AddAttachment(attachment)
		}
	}

	// score[Score],text(Why);company_name__c(Company Name),(rcAccountId__c[RC Account ID]);email[User email];directorySize[Number of users];brand[Brand];survey_id[Survey ID]

	return ccMsg, nil
}

type Line struct {
	Fields []Field
}

type Field struct {
	Property  string
	Display   string
	IsCustom  bool
	UseParens bool
}

var (
	rxParens   = regexp.MustCompile(`^\((.*)\)$`)
	rxBrackets = regexp.MustCompile(`^(.*)\[(.*)\]$`)
	rxCustom   = regexp.MustCompile(`^(.*)__c$`)
)

func ParseFields(fields string) []Line {
	lines := []Line{}
	parts := strings.Split(strings.TrimSpace(fields), ";")
	// Lines
	for _, part := range parts {
		line := Line{Fields: []Field{}}
		lineRaw := strings.TrimSpace(part)
		lineVars := strings.Split(lineRaw, ",")
		// Line Vars
		for _, lineVar := range lineVars {
			lineVar = strings.TrimSpace(lineVar)
			if len(lineVar) == 0 {
				continue
			}
			field := Field{}
			// Use parens
			m1 := rxParens.FindAllStringSubmatch(lineVar, -1)
			if len(m1) > 0 {
				field.UseParens = true
				lineVar = m1[0][1]
			}
			// Brackets
			m2 := rxBrackets.FindAllStringSubmatch(lineVar, -1)
			if len(m2) > 0 {
				field.Display = strings.TrimSpace(m2[0][2])
				property := strings.TrimSpace(m2[0][1])
				m3 := rxCustom.FindAllStringSubmatch(property, -1)
				if len(m3) > 0 {
					field.Property = strings.TrimSpace(m3[0][1])
					field.IsCustom = true
				} else {
					field.Property = property
				}
			}
			if len(field.Property) == 0 {
				fmt.Println(lineVar)
				fmtutil.PrintJSON(field)
				panic("Z")
			}
			line.Fields = append(line.Fields, field)
		}
		lines = append(lines, line)
	}
	return lines
}

//score[Score],text(Why);company_name__c(Company Name),(rcAccountId__c[RC Account ID]);email[User email];directorySize[Number of users];brand[Brand];survey_id[Survey ID]
