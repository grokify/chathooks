package heroku

import (
	"github.com/grokify/glip-go-webhook"
)

func ExampleMessageGlip() (glipwebhook.GlipWebhookMessage, error) {
	msg, err := ExampleMessageSource()
	if err != nil {
		return glipwebhook.GlipWebhookMessage{}, err
	}
	return Normalize(msg), nil
}

func ExampleMessageSource() (HerokuOutMessage, error) {
	return HerokuOutMessageFromQueryString(string(ExampleMessageBytes()))
}

func ExampleMessageBytes() []byte {
	return []byte(`app=secure-woodland-9775&user=example%40example.com&url=http%3A%2F%2Fsecure-woodland-9775.herokuapp.com&head=4f20bdd&head_long=4f20bdd&prev_head=&git_log=%20%20*%20Michael%20Friis%3A%20add%20bar&release=v7
`)
}
