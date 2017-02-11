package semaphoreci

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

func ExampleMessageSource() (SemaphoreciOutMessage, error) {
	return SemaphoreciOutMessageFromBytes(ExampleMessageBytes())
}

func ExampleMessageBytes() []byte {
	return []byte(`{
  "branch_name": "gem_updates",
  "branch_url": "https://semaphoreci.com/projects/44/branches/50",
  "project_name": "base-app",
  "project_hash_id": "123-aga-471-6a8",
  "build_url": "https://semaphoreci.com/projects/44/branches/50/builds/15",
  "build_number": 15,
  "result": "passed",
  "event": "build",
  "started_at": "2012-07-09T15:23:53Z",
  "finished_at": "2012-07-09T15:30:16Z",
  "commit": {
    "id": "dc395381e650f3bac18457909880829fc20e34ba",
    "url": "https://github.com/renderedtext/base-app/commit/dc395381e650f3bac18457909880829fc20e34ba",
    "author_name": "Vladimir Saric",
    "author_email": "vladimir@renderedtext.com",
    "message": "Update 'shoulda' gem.",
    "timestamp": "2012-07-04T18:14:08Z"
  }
}`)
}

/*

if strings.ToLower(msg.Event)== "build"
%v %v %v
msg.Commit.AuthorName
msg.ProjectName
Build #65 passed
fmt.Sprintf("%v %v #%v %v", msg.ProjectName, msg.BuildNumber, msg. Result)
%v", msg.COmmit.Message
[View Details](%v)


{
  "branch_name": "gem_updates",
  "branch_url": "https://semaphoreci.com/projects/44/branches/50",
  "project_name": "base-app",
  "project_hash_id": "123-aga-471-6a8",
  "build_url": "https://semaphoreci.com/projects/44/branches/50/builds/15",
  "build_number": 15,
  "result": "passed",
  "event": "build",
  "started_at": "2012-07-09T15:23:53Z",
  "finished_at": "2012-07-09T15:30:16Z",
  "commit": {
    "id": "dc395381e650f3bac18457909880829fc20e34ba",
    "url": "https://github.com/renderedtext/base-app/commit/dc395381e650f3bac18457909880829fc20e34ba",
    "author_name": "Vladimir Saric",
    "author_email": "vladimir@renderedtext.com",
    "message": "Update 'shoulda' gem.",
    "timestamp": "2012-07-04T18:14:08Z"
  }
}


Webhook Notification Reference

https://docs.travis-ci.com/user/notifications#Configuring-webhook-notifications

Format:

"Build <%{build_url}|#%{build_number}> (<%{compare_url}|%{commit}>) of %{repository}@%{branch} by %{author} %{result} in %{duration}"

Payload:

{
  "branch_name": "gem_updates",
  "branch_url": "https://semaphoreci.com/projects/44/branches/50",
  "project_name": "base-app",
  "project_hash_id": "123-aga-471-6a8",
  "build_url": "https://semaphoreci.com/projects/44/branches/50/builds/15",
  "build_number": 15,
  "result": "passed",
  "event": "build",
  "started_at": "2012-07-09T15:23:53Z",
  "finished_at": "2012-07-09T15:30:16Z",
  "commit": {
    "id": "dc395381e650f3bac18457909880829fc20e34ba",
    "url": "https://github.com/renderedtext/base-app/commit/dc395381e650f3bac18457909880829fc20e34ba",
    "author_name": "Vladimir Saric",
    "author_email": "vladimir@renderedtext.com",
    "message": "Update 'shoulda' gem.",
    "timestamp": "2012-07-04T18:14:08Z"
  }
}

*/
