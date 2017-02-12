package magnumci

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

func ExampleMessageSource() (MagnumciOutMessage, error) {
	return MagnumciOutMessageFromBytes(ExampleMessageBytes())
}

func ExampleMessageBytes() []byte {
	return []byte(`{
  "id": 1603,
  "project_id": 43,
  "title": "[PASS] project-name #130 (master - e91e132) by Dan Sosedoff",
  "number": 130,
  "commit": "e91e132612d263d95211aae6de2df9e503f22704",
  "author": "Dan Sosedoff",
  "committer": "Dan Sosedoff",
  "message": "Commit Message",
  "branch": "master",
  "state": "finished",
  "status": "pass",
  "result": 0,
  "duration": 158,
  "duration_string": "2m 38s",
  "commit_url": "http://domain.com/commit/e91e132612d263...",
  "compare_url": null,
  "build_url": "http://magnum-ci.com/projects/43/builds/1603",
  "started_at": "2013-02-14T00:09:01-06:00",
  "finished_at": "2013-02-14T00:11:39-06:00"
}`)
}
