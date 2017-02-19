package raygun

import (
	"github.com/grokify/commonchat"
)

func ExampleMessage() (commonchat.Message, error) {
	bytes, err := ExampleMessageBytes()
	if err != nil {
		return commonchat.Message{}, err
	}
	return Normalize(bytes)
}

func ExampleMessageBytes() ([]byte, error) {
	return []byte(`{
  "event":"error_notification",
  "eventType":"NewErrorOccurred",
  "error": {
    "url":"http://app.raygun.io/error-url",
    "message":"",
    "firstOccurredOn":"1970-01-28T01:49:36Z",
    "lastOccurredOn":"1970-01-28T01:49:36Z",
    "usersAffected":1,
    "totalOccurrences":1
  },
  "application": {
    "name":"application name",
    "url":"http://app.raygun.io/application-url"
  }
}`), nil
}
