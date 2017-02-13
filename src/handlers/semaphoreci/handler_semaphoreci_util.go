package semaphoreci

import (
	"github.com/grokify/glip-go-webhook"
)

func ExampleBuildMessageGlip() (glipwebhook.GlipWebhookMessage, error) {
	msg, err := ExampleBuildMessageSource()
	if err != nil {
		return glipwebhook.GlipWebhookMessage{}, err
	}
	return NormalizeSemaphoreciBuildOutMessage(msg), nil
}

func ExampleBuildMessageSource() (SemaphoreciBuildOutMessage, error) {
	return SemaphoreciBuildOutMessageFromBytes(ExampleBuildMessageBytes())
}

func ExampleBuildMessageBytes() []byte {
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

func ExampleDeployMessageGlip() (glipwebhook.GlipWebhookMessage, error) {
	msg, err := ExampleDeployMessageSource()
	if err != nil {
		return glipwebhook.GlipWebhookMessage{}, err
	}
	return NormalizeSemaphoreciDeployOutMessage(msg), nil
}

func ExampleDeployMessageSource() (SemaphoreciDeployOutMessage, error) {
	return SemaphoreciDeployOutMessageFromBytes(ExampleDeployMessageBytes())
}

func ExampleDeployMessageBytes() []byte {
	return []byte(`{
  "project_name": "heroku-deploy-test",
  "project_hash_id": "123-aga-471-6a8",
  "result": "passed",
  "event": "deploy",
  "server_name": "server-heroku-master-automatic-2",
  "number": 2,
  "created_at": "2013-07-30T13:52:33Z",
  "updated_at": "2013-07-30T13:53:21Z",
  "started_at": "2013-07-30T13:52:38Z",
  "finished_at": "2013-07-30T13:53:21Z",
  "html_url": "https://semaphoreci.com/projects/2420/servers/81/deploys/2",
  "build_number": 10,
  "branch_name": "master",
  "branch_html_url": "https://semaphoreci.com/projects/2420/branches/58394",
  "build_html_url": "https://semaphoreci.com/projects/2420/branches/58394/builds/7",
  "commit": {
    "author_email": "rastasheep3@gmail.com",
    "author_name": "Aleksandar Diklic",
    "id": "43ddb7516ecc743f0563abd7418f0bd3617348c4",
    "message": "One more time",
    "timestamp": "2013-07-19T12:56:25Z",
    "url": "https://github.com/rastasheep/heroku-deploy-test/commit/43ddb7516ecc743f0563abd7418f0bd3617348c4"
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

{
  "project_name": "heroku-deploy-test",
  "project_hash_id": "123-aga-471-6a8",
  "result": "passed",
  "event": "deploy",
  "server_name": "server-heroku-master-automatic-2",
  "number": 2,
  "created_at": "2013-07-30T13:52:33Z",
  "updated_at": "2013-07-30T13:53:21Z",
  "started_at": "2013-07-30T13:52:38Z",
  "finished_at": "2013-07-30T13:53:21Z",
  "html_url": "https://semaphoreci.com/projects/2420/servers/81/deploys/2",
  "build_number": 10,
  "branch_name": "master",
  "branch_html_url": "https://semaphoreci.com/projects/2420/branches/58394",
  "build_html_url": "https://semaphoreci.com/projects/2420/branches/58394/builds/7",
  "commit": {
    "author_email": "rastasheep3@gmail.com",
    "author_name": "Aleksandar Diklic",
    "id": "43ddb7516ecc743f0563abd7418f0bd3617348c4",
    "message": "One more time",
    "timestamp": "2013-07-19T12:56:25Z",
    "url": "https://github.com/rastasheep/heroku-deploy-test/commit/43ddb7516ecc743f0563abd7418f0bd3617348c4"
  }
}


*/
