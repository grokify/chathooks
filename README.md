Glip Webhook Proxy
==================

[![Build Status][build-status-svg]][build-status-link]
[![Go Report Card][goreport-svg]][goreport-link]
[![Code Climate][codeclimate-status-svg]][codeclimate-status-link]
[![Docs][docs-godoc-svg]][docs-godoc-link]
[![License][license-svg]][license-link]

Proxy service to map different requests to Glip's inbound webhook service. This is useful because various chat services have similar, but slightly different inbound webhook services. This proxy service does the conversion so you don't have to. Applications already integrated with Slack's inbound webhooks can create messages on Glip simply by using the proxy URL.

Conversion of the following webhook message formats to Glip inbound webhooks include:

* Slack (inbound message format) - `text` only
* Travis CI (outbound message format)

**Note:** Slack inbound message formatting is for services sending outbound webhooks using Slack's inbound webhook message format, which can be directed to Glip via this proxy.

Example Webhook Message from Travis CI:

![](adapters/travisci/travisci_glip.png)

## Installation

```
$ go get github.com/grokify/glip-webhook-proxy
```

## Usage

### Starting the Service

Start the service with the following.

```go
package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/grokify/glip-webhook-proxy"
)

func main() {
	config := glipwebhookproxy.Configuration{
		Port:           8080,
		EmojiURLFormat: "https://grokify.github.io/emoji/assets/images/%s.png",
		LogLevel:       log.DebugLevel}

	glipwebhookproxy.StartServer(config)
}
```

You can run the above by saving it to a file `start.go` and then running `$ go run start.go`.

Note: The emoji to URL is designed to take a `icon_emoji` value and convert it to a URL. `EmojiURLFormat` is a [`fmt`](https://golang.org/pkg/fmt/) `format` string with one `%s` verb to represent the emoji string without `:`. You can use any emoji image service. The example shows the emoji set from [github.com/wpeterson/emoji](https://github.com/wpeterson/emoji) forked and hosted at [grokify.github.io/emoji/](https://grokify.github.io/emoji/).

### Creating the Glip Webhook

1. create a Glip webhook
2. use webhook URL's GUID to create the proxy URL as shown below
3. use the proxy URL in your outbound webhook service

| Service | URL |
|------|-------|
| Glip | `https://hooks.glip.com/webhook/11112222-3333-4444-5555-666677778888` |
| Slack Inbound | `https://example.com/webhook/slack/in/glip/11112222-3333-4444-5555-666677778888` |
| Travis CI Outbound | `https://example.com/webhook/travisci/out/glip/11112222-3333-4444-5555-666677778888` |

The webhook proxy URLs support both inbound and outbound formats. For example:

* when using Travis CI's webhook format use `travisci/out/glip` to indicate converting a Travis CI outbound webhook format message to Glip.
* when using a service that supports Slack inbound webhook format use `slack/in/glip` to indicate converting a Slack inbound webhook format message to Glip.

To create the Glip webhook and receive a webhook URL do the following:

#### Add the Webhook Integration

At the top of any conversation page, click the Settings gear icon and then click `Add Integration`.

![](images/glip_webhook_step-1_add-integration.png)

Select the `Glip Webhooks` integration.

![](images/glip_webhook_step-2_add-webhook.png)

#### Get the Webhook URL

Once you get the URL, the proxy URL is created by appending the GUID (e.g. `1112222-3333-4444-5555-666677778888`) to the proxy URL base, `/webhook/slack/glip` (e.g. `https://glip-proxy.example.com/webhook/slack/glip/1112222-3333-4444-5555-666677778888`). Use the proxy URL in the app that is posting the Slack webhook and the payload will be sent to Glip.

![](images/glip_webhook_step-3_details.png)

## Example Requests

Most of the time you will likely either:

* use the proxy URL in an outbound webhook service that supports the Slack format or
* use a client library

The following examples are provided for reference and testing.

### Using `application/json`

```bash
$ curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"username":"ghost-bot", "icon_emoji": ":ghost:", "text":"BOO!"}' \
  "http://localhost:8080/webhook/slack/glip/11112222-3333-4444-5555-666677778888"
```

### Using `application/x-www-form-urlencoded`

```bash
$ curl -X POST \
  --data-urlencode 'payload={"username":"ghost-bot", "icon_emoji": ":ghost:", text":"BOO!"}' \
  "http://localhost:8080/webhook/slack/glip/11112222-3333-4444-5555-666677778888"
```

### Using `multipart/form-data`

```bash
$ curl -X POST \
  -F 'payload={"username":"ghost-bot", "icon_emoji": ":ghost:", text":"BOO!"}' \
  "http://localhost:8080/webhook/slack/glip/11112222-3333-4444-5555-666677778888"
```

### Using Community Ruby SDK

This has been tested using:

* [https://github.com/rikas/slack-poster](https://github.com/rikas/slack-poster)

```ruby
require 'slack/poster'

url = 'http://localhost:8080/webhook/slack/glip/11112222-3333-4444-5555-666677778888'

opts = {
	username: 'Ghost Bot [Bot]',
	icon_emoji: ':ghost:'
}

poster = Slack::Poster.new url, opts
poster.send_message 'BOO!'
```

## Links

1. [Confluence webhooks](https://developer.atlassian.com/static/connect/docs/beta/modules/common/webhook.html)
1. [Enchant webhooks](https://dev.enchant.com/webhooks)
1. [Heroku webhooks](https://devcenter.heroku.com/articles/deploy-hooks#http-post-hook)
1. [Magnum CI webhooks](https://github.com/magnumci/documentation/blob/master/webhooks.md)
1. [Raygun webhooks](https://raygun.com/docs/integrations/webhooks)
1. [Semaphore CI webhooks](https://semaphoreci.com/docs/post-build-webhooks.html)
1. [Travis CI webhooks](https://docs.travis-ci.com/user/notifications#Configuring-webhook-notifications)
1. [Userlike webhooks](https://www.userlike.com/en/public/tutorial/addon/api)

## Notes

Glip Webhook Proxy is built using:

* [fasthttp](https://github.com/valyala/fasthttp)
* [fasthttprouter](https://github.com/buaazp/fasthttprouter)
* [logrus](https://github.com/sirupsen/logrus)

 [build-status-svg]: https://api.travis-ci.org/grokify/glip-webhook-proxy-go.svg?branch=master
 [build-status-link]: https://travis-ci.org/grokify/glip-webhook-proxy-go
 [coverage-status-svg]: https://coveralls.io/repos/grokify/glip-webhook-proxy-go/badge.svg?branch=master
 [coverage-status-link]: https://coveralls.io/r/grokify/glip-webhook-proxy-go?branch=master
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/glip-webhook-proxy-go
 [goreport-link]: https://goreportcard.com/report/github.com/grokify/glip-webhook-proxy-go
 [codeclimate-status-svg]: https://codeclimate.com/github/grokify/glip-webhook-proxy-go/badges/gpa.svg
 [codeclimate-status-link]: https://codeclimate.com/github/grokify/glip-webhook-proxy-go
 [docs-godoc-svg]: https://img.shields.io/badge/docs-godoc-blue.svg
 [docs-godoc-link]: https://godoc.org/github.com/grokify/glip-webhook-proxy-go
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-link]: https://github.com/grokify/glip-webhook-proxy-go/blob/master/LICENSE.md
