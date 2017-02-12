package travisci

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

func ExampleMessageSource() (TravisciOutMessage, error) {
	return TravisciOutMessageFromBytes(ExampleMessageBytes())
}

func ExampleMessageBytes() []byte {
	return []byte(`{
  "id": 1,
  "number": "1",
  "status": null,
  "started_at": null,
  "finished_at": null,
  "status_message": "Passed",
  "commit": "62aae5f70ceee39123ef",
  "branch": "master",
  "message": "the commit message",
  "compare_url": "https://github.com/svenfuchs/minimal/compare/master...develop",
  "committed_at": "2011-11-11T11: 11: 11Z",
  "committer_name": "Sven Fuchs",
  "committer_email": "svenfuchs@artweb-design.de",
  "author_name": "Sven Fuchs",
  "author_email": "svenfuchs@artweb-design.de",
  "type": "push",
  "build_url": "https://travis-ci.org/svenfuchs/minimal/builds/1",
  "repository": {
    "id": 1,
    "name": "minimal",
    "owner_name": "svenfuchs",
    "url": "http://github.com/svenfuchs/minimal"
   },
  "config": {
    "notifications": {
      "webhooks": ["http://evome.fr/notifications", "http://example.com/"]
    }
  },
  "matrix": [
    {
      "id": 2,
      "repository_id": 1,
      "number": "1.1",
      "state": "created",
      "started_at": null,
      "finished_at": null,
      "config": {
        "notifications": {
          "webhooks": ["http://evome.fr/notifications", "http://example.com/"]
        }
      },
      "status": null,
      "log": "",
      "result": null,
      "parent_id": 1,
      "commit": "62aae5f70ceee39123ef",
      "branch": "master",
      "message": "the commit message",
      "committed_at": "2011-11-11T11: 11: 11Z",
      "committer_name": "Sven Fuchs",
      "committer_email": "svenfuchs@artweb-design.de",
      "author_name": "Sven Fuchs",
      "author_email": "svenfuchs@artweb-design.de",
      "compare_url": "https://github.com/svenfuchs/minimal/compare/master...develop"
    }
  ]
}`)
}

/*

Webhook Notification Reference

https://docs.travis-ci.com/user/notifications#Configuring-webhook-notifications

Format:

"Build <%{build_url}|#%{build_number}> (<%{compare_url}|%{commit}>) of %{repository}@%{branch} by %{author} %{result} in %{duration}"

Payload:



*/
