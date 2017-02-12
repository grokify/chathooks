package heroku

import (
	"testing"
)

const (
	TEST_APP = "secure-woodland-9775"
)

var ConfigurationTests = []struct {
	v    string
	want HerokuOutMessage
}{
	{HookData(), HerokuOutMessage{App: "secure-woodland-9775"}}}

func TestConfluence(t *testing.T) {
	for _, tt := range ConfigurationTests {
		msg, err := HerokuOutMessageFromQueryString(tt.v)
		if err != nil {
			t.Errorf("Error %v\n", err)
			continue
		}
		if msg.App != TEST_APP {
			t.Errorf("HerokuOutMessageFromQueryString(%v): want %v, got %v", tt.v, "secure-woodland-9775", tt.want.App)
		}
	}
}

func HookData() string {
	return `app=secure-woodland-9775&user=example%40example.com&url=http%3A%2F%2Fsecure-woodland-9775.herokuapp.com&head=4f20bdd&head_long=4f20bdd&prev_head=&git_log=%20%20*%20Michael%20Friis%3A%20add%20bar&release=v7
`
}
