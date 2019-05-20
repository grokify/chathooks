package bugsnag

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/grokify/chathooks/src/config"
	"github.com/grokify/chathooks/src/handlers"
	"github.com/grokify/chathooks/src/models"
	cc "github.com/grokify/commonchat"
	log "github.com/sirupsen/logrus"
)

const (
	DisplayName      = "Bugsnag"
	HandlerKey       = "bugsnag"
	MessageDirection = "out"
	MessageBodyType  = models.JSON
)

func NewHandler() handlers.Handler {
	return handlers.Handler{MessageBodyType: MessageBodyType, Normalize: Normalize}
}

/*
Bugsnag

**New error** in **production** from [Webpage]() in example:activity ([details]())

Error
ExampleError: An error occurred

Location
controllers/example.rb:11 - example_call
*/

func Normalize(cfg config.Configuration, bytes []byte) (cc.Message, error) {
	ccMsg := cc.NewMessage()
	iconURL, err := cfg.GetAppIconURL(HandlerKey)
	if err == nil {
		ccMsg.IconURL = iconURL.String()
	}

	ccMsg.Activity = fmt.Sprintf("%s alert", DisplayName)

	src, err := BugsnagOutMessageFromBytes(bytes)
	if err != nil {
		return ccMsg, err
	}

	parts := []string{}

	triggerMsg := strings.TrimSpace(src.Trigger.Message)
	if len(triggerMsg) > 0 {
		parts = append(parts, fmt.Sprintf("**%s**", triggerMsg))
	}

	stage := strings.TrimSpace(src.Release.ReleaseStage)
	if len(stage) > 0 {
		parts = append(parts, fmt.Sprintf("in **%s**", stage))
	}

	projLink := src.Project.MarkdownLink()
	if len(projLink) > 0 {
		parts = append(parts, fmt.Sprintf("from **%s**", projLink))
	}

	if len(src.Error.Context) > 0 {
		parts = append(parts, fmt.Sprintf("in %s", src.Error.Context))
	}

	if len(src.Error.Url) > 0 {
		details := fmt.Sprintf("([details](%s))", src.Error.Url)
		parts = append(parts, details)
	}

	if len(parts) > 0 {
		ccMsg.Title = strings.Join(parts, " ")
	}

	fields := []cc.Field{}

	if len(strings.TrimSpace(src.Error.Message)) > 0 {
		fields = append(fields, cc.Field{
			Title: "Error",
			Value: src.Error.Message})
	}

	for _, st := range src.Error.StackTrace {
		location := strings.TrimSpace(st.Location())
		if len(location) > 0 {
			fields = append(fields, cc.Field{
				Title: "Location",
				Value: location})
			break
		}
	}

	if len(strings.TrimSpace(src.Error.Status)) > 0 {
		status := src.Error.Status
		if strings.ToLower(strings.TrimSpace(status)) == "open" {
			unhandled := src.Error.Unhandled
			if unhandled {
				status += " - unhandled"
			} else {
				status += " - handled"
			}
		}

		fields = append(fields, cc.Field{
			Title: "Status",
			Value: status})
	}

	if len(fields) > 0 {
		attachment := cc.NewAttachment()
		for _, field := range fields {
			attachment.AddField(field)
		}
		ccMsg.AddAttachment(attachment)
	}

	/*
				if len(src.ActorName) > 0 {
					ccMsg.Activity = src.ActorName
				}

				attachment := cc.NewAttachment()

				if len(src.Model.Subject) > 0 {
					attachment.Text = src.Model.Subject
				}
				if len(src.Model.State) > 0 {
					attachment.AddField(cc.Field{
						Title: "State",
						Value: stringsutil.ToUpperFirst(src.Model.State)})
				}
				ccMsg.AddAttachment(attachment)

		   "stackTrace":[
		     {
		       "inProject":true,
		       "lineNumber":1234,
		       "columnNumber":123,
		       "file":"controllers/auth/session_controller.rb",
		       "method":"create",
		       "code":{
		         "1231":"  def a",
		         "1232":"",
		         "1233":"    if problem?",
		         "1234":"      raise 'something went wrong'",
		         "1235":"    end",
		         "1236":"",
		         "1237":"  end"
		       }
		     }
		   ]
	*/

	return ccMsg, nil
}

func BugsnagOutMessageFromBytes(bytes []byte) (BugsnagOutMessage, error) {
	log.WithFields(log.Fields{
		"type":    "message.raw",
		"message": string(bytes),
	}).Debug(fmt.Sprintf("%v message.", DisplayName))
	msg := BugsnagOutMessage{}
	err := json.Unmarshal(bytes, &msg)
	if err != nil {
		log.WithFields(log.Fields{
			"type":  "message.json.unmarshal",
			"error": fmt.Sprintf("%v\n", err),
		}).Warn(fmt.Sprintf("%v request unmarshal failure.", DisplayName))
	}
	return msg, err
}

type BugsnagOutMessage struct {
	Account BugsnagAccount `json:"account,omitempty"`
	Project BugsnagProject `json:"project,omitempty"`
	Trigger BugsnagTrigger `json:"trigger,omitempty"`
	Comment BugsnagComment `json:"created_at,omitempty"`
	User    BugsnagUser    `json:"user,omitempty"`
	Error   BugsnagError   `json:"error,omitempty"`
	Release BugsnagRelease `json:"release,omitempty"`
}

func (msg *BugsnagOutMessage) ReleaseStage() string {
	relRelStage := strings.TrimSpace(msg.Release.ReleaseStage)
	appRelStage := strings.TrimSpace(msg.Error.App.ReleaseStage)
	if len(appRelStage) > 0 {
		return appRelStage
	}
	return relRelStage
}

type BugsnagAccount struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

type BugsnagProject struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

func (bp *BugsnagProject) MarkdownLink() string {
	if len(bp.Name) > 0 && len(bp.URL) > 0 {
		return fmt.Sprintf("[%s](%s)", bp.Name, bp.URL)
	} else if len(bp.URL) > 0 {
		return fmt.Sprintf("[%s](%s)", bp.URL, bp.URL)
	}
	return ""
}

type BugsnagTrigger struct {
	Type        string            `json:"type,omitempty"`
	Message     string            `json:"message,omitempty"`
	SnoozeRule  BugsnagSnoozeRule `json:"snoozeRule,omitempty"`
	Rate        int32             `json:"rate,omitempty"`
	StateChange string            `json:"stateChange,omitempty"`
}

type BugsnagSnoozeRule struct {
	Type      string `json:"type,omitempty"`
	RuleValue int32  `json:"ruleValue,omitempty"`
}

type BugsnagComment struct{}

type BugsnagUser struct {
	Id    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

type BugsnagRelease struct {
	Id            string               `json:"id,omitempty"`
	Version       string               `json:"version,omitempty"`
	VersionCode   string               `json:"versionCode,omitempty"`
	BundleVersion string               `json:"bundleVersion,omitempty"`
	ReleaseStage  string               `json:"releaseStage,omitempty"`
	URL           string               `json:"url,omitempty"`
	ReleaseTime   time.Time            `json:"releaseTime,omitempty"`
	ReleasedBy    string               `json:"releasedBy,omitempty"`
	SourceControl BugsnagSourceControl `json:"sourceControl,omitempty"`
	Metadata      map[string]string    `json:"metadata,omitempty"`
}

type BugsnagSourceControl struct {
	Provider    string `json:"provider,omitempty"`
	Revision    string `json:"revision,omitempty"`
	RevisionURL string `json:"revisionUrl,omitempty"`
	DiffURL     string `json:"diffUrl,omitempty"`
}

type BugsnagError struct {
	Id             string                   `json:"id,omitempty"`
	ErrorId        string                   `json:"errorId,omitempty"`
	ExceptionClass string                   `json:"exceptionClass,omitempty"`
	Message        string                   `json:"message,omitempty"`
	Context        string                   `json:"context,omitempty"`
	FirstReceived  time.Time                `json:"firstReceived,omitempty"`
	ReceivedAt     time.Time                `json:"receivedAt,omitempty"`
	RequestUrl     string                   `json:"requestUrl,omitempty"`
	AssignedUserId string                   `json:"assignedUserId,omitempty"`
	Url            string                   `json:"url,omitempty"`
	Severity       string                   `json:"severity,omitempty"`
	Status         string                   `json:"status,omitempty"`
	Unhandled      bool                     `json:"unhandled,omitempty"`
	CreatedIssue   BugsnagErrorIssue        `json:"createdIssue,omitempty"`
	User           BugsnagUser              `json:"user,omitempty"`
	App            BugsnaggErrorApp         `json:"app,omitempty"`
	Device         BugsnagErrorDevice       `json:"device,omitempty"`
	StackTrace     []BugsnagErrorStackTrace `json:"stackTrace,omitempty"`
}

type BugsnagErrorIssue struct {
	Id     string `json:"id,omitempty"`
	Number int    `json:"name,omitempty"`
	Type   string `json:"type,omitempty"`
	URL    string `json:"url,omitempty"`
}

type BugsnaggErrorApp struct {
	Id                   string   `json:"id,omitempty"`
	Version              string   `json:"version,omitempty"`
	VersionCode          string   `json:"versionCode,omitempty"`
	BundleVersion        string   `json:"bundleVersion,omitempty"`
	CodeBundleId         string   `json:"codeBundleId,omitempty"`
	BuildUUID            string   `json:"buildUUID,omitempty"`
	ReleaseStage         string   `json:"releaseStage,omitempty"`
	Type                 string   `json:"type,omitempty"`
	DsymUUIDs            []string `json:"dsymUUIDs,omitempty"`
	Duration             int      `json:"duration,omitempty"`
	durationInForeground int      `json:"durationInForeground,omitempty"`
	inForeground         bool     `json:"inForeground,omitempty"`
}

type BugsnagErrorDevice struct {
	Id             string    `json:"id,omitempty"`
	Manufacturer   string    `json:"manufacturer,omitempty"`
	Model          string    `json:"model,omitempty"`
	ModelNumber    string    `json:"modelNumber,omitempty"`
	OsName         string    `json:"osName,omitempty"`
	OsVersion      string    `json:"osVersion,omitempty"`
	FreeMemory     int       `json:"freeMemory,omitempty"`
	TotalMemory    int       `json:"totalMemory,omitempty"`
	FreeDisk       int       `json:"freeDisk,omitempty"`
	BrowserName    string    `json:"browserName,omitempty"`
	BrowserVersion string    `json:"browserVersion,omitempty"`
	Jailbroken     bool      `json:"jailbroken,omitempty"`
	Orientation    string    `json:"orientation,omitempty"`
	Locale         string    `json:"locale,omitempty"`
	Charging       bool      `json:"charging,omitempty"`
	BatteryLevel   float32   `json:"batteryLevel,omitempty"`
	Time           time.Time `json:"time,omitempty"`
	Timezone       string    `json:"timezone,omitempty"`
}

type BugsnagErrorStackTrace struct {
	InProject    bool              `json:"inProject,omitempty"`
	LineNumber   int               `json:"lineNumber,omitempty"`
	ColumnNumber int               `json:"columnNumber,omitempty"`
	File         string            `json:"file,omitempty"`
	Method       string            `json:"method,omitempty"`
	Code         map[string]string `json:"code,omitempty"`
}

// Location returns a string per the Slack integration
func (st *BugsnagErrorStackTrace) Location() string {
	location := ""
	st.File = strings.TrimSpace(st.File)
	st.Method = strings.TrimSpace(st.Method)
	if len(st.File) > 0 {
		location = st.File
		if st.LineNumber > 0 {
			location += ":" + strconv.Itoa(st.LineNumber)
		}
		if len(st.Method) > 0 {
			location += " - " + st.Method
		}
	}
	return location
}

/*

controllers/example.rb:11 - example_call

   "stackTrace":[
     {
       "inProject":true,
       "lineNumber":1234,
       "columnNumber":123,
       "file":"controllers/auth/session_controller.rb",
       "method":"create",
       "code":{
         "1231":"  def a",
         "1232":"",
         "1233":"    if problem?",
         "1234":"      raise 'something went wrong'",
         "1235":"    end",
         "1236":"",
         "1237":"  end"
       }
     }
   ]
*/
