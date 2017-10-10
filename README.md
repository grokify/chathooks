WebhookProxy - A webhook proxy
==============================

[![Build Status][build-status-svg]][build-status-link]
[![Go Report Card][goreport-svg]][goreport-link]
[![Code Climate][codeclimate-status-svg]][codeclimate-status-link]
[![Docs][docs-godoc-svg]][docs-godoc-link]
[![License][license-svg]][license-link]

WebhookProxy is a service that maps webhook posts from different services to message platforms such as Glip and Slack's inbound webhook service. It uses handlers to convert inbound messages to the [CommonChat](https://github.com/commonchat) canonical message format which are then sent via message platform adapters. This is useful because many services with outbound webhooks need to be formatted before they can be consumed by an inbound webhook. This proxy service does the conversion so you don't have to. Another use case is conversion of inbound messages so a message formatted for Slack inbound webhooks can be delivered to a Glip inbound webhook.

It is easy to add additional inbound webhook handlers and outbound webhook adapters by using the `adapters.Adapter` and `handlers.Handler` interfaces.

WebhookProxy currently supports two HTTP server engines.

* AWS API Gateway + AWS Lambda - [eawsy/aws-lambda-go-shim](https://github.com/eawsy/aws-lambda-go-shim)
* Locally - [valyala/fasthttp](https://github.com/valyala/fasthttp)

Conversion of the following webhook message formats to Glip inbound webhooks include:

Outbound Webhook Formats supported:

1. [AppSignal](http://docs.appsignal.com/application/integrations/webhooks.html)
1. [Apteligent/Crittercism]()
1. [Circle CI](https://circleci.com/docs/1.0/configuration/#notify)
1. [Codeship](https://documentation.codeship.com/basic/getting-started/webhooks/)
1. [Confluence](https://developer.atlassian.com/static/connect/docs/beta/modules/common/webhook.html)
1. [Datadog](http://docs.datadoghq.com/integrations/webhooks/)
1. [Desk.com](https://support.desk.com/customer/portal/articles/869334-configuring-webhooks-in-desk-com-apps)
1. [Enchant](https://dev.enchant.com/webhooks)
1. [GoSquared](https://www.gosquared.com/customer/portal/articles/1996494-webhooks)
1. [Heroku](https://devcenter.heroku.com/articles/deploy-hooks#http-post-hook)
1. [Librato]()
1. [Magnum CI](https://github.com/magnumci/documentation/blob/master/webhooks.md)
1. [Marketo]()
1. [OpsGenie]()
1. [Papertrail](http://help.papertrailapp.com/kb/how-it-works/web-hooks/)
1. [Pingdom]()
1. [Raygun](https://raygun.com/docs/integrations/webhooks)
1. [Runscope]()
1. [Semaphore CI](https://semaphoreci.com/docs/post-build-webhooks.html), [Deploy](https://semaphoreci.com/docs/post-deploy-webhooks.html)
1. [StatusPage]()
1. [Travis CI](https://docs.travis-ci.com/user/notifications#Configuring-webhook-notifications)
1. [Userlike](https://www.userlike.com/en/public/tutorial/addon/api)
1. [VictorOps]()

Inbound Webhook Format supported:

* Slack (inbound message format) - `text` only

**Note:** Slack inbound message formatting is for services sending outbound webhooks using Slack's inbound webhook message format, which can be directed to Glip via this proxy.

Example Webhook Message from Travis CI:

![](src/handlers/travisci/travisci_glip.png)

## Installation

```
$ go get github.com/grokify/webhookproxy
```

## Usage

### Starting the Service using FastHTTP

Start the service in `server.go`.

* To adjust supported handlers, edit server.go to add and remove handlers.

Start the service with the following.

Note: The emoji to URL is designed to take a `icon_emoji` value and convert it to a URL. `EmojiURLFormat` is a [`fmt`](https://golang.org/pkg/fmt/) `format` string with one `%s` verb to represent the emoji string without `:`. You can use any emoji image service. The example shows the emoji set from [github.com/wpeterson/emoji](https://github.com/wpeterson/emoji) forked and hosted at [grokify.github.io/emoji/](https://grokify.github.io/emoji/).

### Creating the Glip Webhook

1. create a Glip webhook
2. use webhook URL's GUID to create the proxy URL as shown below
3. use the proxy URL in your outbound webhook service

| Query Parameter | Required? | URL |
|-----------------|-----------|-----|
| `inputType` | required | An handler service like `marketo` |
| `outputType` | required | An adapter service like `glip` |
| `url` | required | A webhook URL or UID, e.g. `11112222-3333-4444-5555-666677778888` |
| `token` | optional | Must be included if service is configured to use auth tokens |

The webhook proxy URLs support both inbound and outbound formats. When available, these should be represented in the handler key.

To create the Glip webhook and receive a webhook URL do the following:

#### Add the Webhook Integration

At the top of any conversation page, click the Settings gear icon and then click `Add Integration`.

![](docs/images/glip_webhook_step-1_add-integration.png)

Select the `Glip Webhooks` integration.

![](docs/images/glip_webhook_step-2_add-webhook.png)

#### Get the Webhook URL

Once you get the URL, the proxy URL is created by appending the GUID (e.g. `1112222-3333-4444-5555-666677778888`) to the proxy URL base, `hooks?inputType=slack&outputType=glip` (e.g. `https://glip-proxy.example.com/hooks?inputType=slack&outputType=glip&url=1112222-3333-4444-5555-666677778888`). Use the proxy URL in the app that is posting the Slack webhook and the payload will be sent to Glip.

![](docs/images/glip_webhook_step-3_details.png)

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
  "http://localhost:8080/hooks?inputType=slack&outputType=glip&url=11112222-3333-4444-5555-666677778888"
```

### Using `application/x-www-form-urlencoded`

```bash
$ curl -X POST \
  --data-urlencode 'payload={"username":"ghost-bot", "icon_emoji": ":ghost:", text":"BOO!"}' \
  "http://localhost:8080/hooks?inputType=slack&outputType=glip&url=11112222-3333-4444-5555-666677778888"
```

### Using `multipart/form-data`

```bash
$ curl -X POST \
  -F 'payload={"username":"ghost-bot", "icon_emoji": ":ghost:", text":"BOO!"}' \
  "http://localhost:8080?hooks?inputType=slack&outputType=glip&url=11112222-3333-4444-5555-666677778888"
```

### Using Community Ruby SDK

This has been tested using:

* [https://github.com/rikas/slack-poster](https://github.com/rikas/slack-poster)

```ruby
require 'slack/poster'

url = 'http://localhost:8080?inputType=slack&outputType=glip&url=11112222-3333-4444-5555-666677778888'

opts = {
	username: 'Ghost Bot [Bot]',
	icon_emoji: ':ghost:'
}

poster = Slack::Poster.new url, opts
poster.send_message 'BOO!'
```

## Notes

WebhookProxy is built using:

* [valyala/fasthttp](https://github.com/valyala/fasthttp)
* [buaazp/fasthttprouter](https://github.com/buaazp/fasthttprouter)
* [sirupsen/logrus](https://github.com/sirupsen/logrus)
* [eawsy/aws-lambda-go-shim](https://github.com/eawsy/aws-lambda-go-shim)

 [build-status-svg]: https://api.travis-ci.org/grokify/webhookproxy.svg?branch=master
 [build-status-link]: https://travis-ci.org/grokify/webhookproxy
 [coverage-status-svg]: https://coveralls.io/repos/grokify/webhookproxy/badge.svg?branch=master
 [coverage-status-link]: https://coveralls.io/r/grokify/webhookproxy?branch=master
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/webhookproxy
 [goreport-link]: https://goreportcard.com/report/github.com/grokify/webhookproxy
 [codeclimate-status-svg]: https://codeclimate.com/github/grokify/webhookproxy/badges/gpa.svg
 [codeclimate-status-link]: https://codeclimate.com/github/grokify/webhookproxy
 [docs-godoc-svg]: https://img.shields.io/badge/docs-godoc-blue.svg
 [docs-godoc-link]: https://godoc.org/github.com/grokify/webhookproxy
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-link]: https://github.com/grokify/webhookproxy/blob/master/LICENSE.md
