Glip Webhook Proxy
==================

[![Go Report Card][goreport-svg]][goreport-link]
[![License][license-svg]][license-link]

Proxy service to map different requests to Glip's inbound webhook service. This is useful because various chat services have similar, but slightly different inbound webhook services. With slight modifications, messages for one service can be converted to the Glip service.

This proxy currently supports converting Slack webhooks services to Glip webhooks. Setting up this service will allow you to use proxy URLs in services that support Slack to post into Glip. It currently:

* handles all request content types
* converts payload property names
* converts emoji to URL
* is tested with SDKs

Note: At this time, the proxy only supports the `text` body and not message attachments yet.

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
	"github.com/grokify/glip-webhook-proxy"
)

func main() {
	config := glipwebhookproxy.Configuration{
		Port:           8080,
		EmojiURLPrefix: "https://grokify.github.io/emoji/assets/images/",
		EmojiURLSuffix: ".png"}

	glipwebhookproxy.StartServer(config)
}
```

You can run the above by saving it to a file `start.go` and then running `$ go run start.go`.

### Creating the Glip Webhook

1. create a Glip webhook
2. use webhook URL's GUID to create the proxy URL as shown below
3. use the proxy URL in your outbound webhook service

| Name | Value |
|------|-------|
| Glip Webhook URL | `https://hooks.glip.com/webhook/11112222-3333-4444-5555-666677778888` |
| Proxy Webhook URL | `https://example.com/slack/glip/11112222-3333-4444-5555-666677778888` |

To create the Glip webhook and receive a webhook URL do the following:

#### Add the integration

![](glip_webhook_step-1_add-integration.png)

![](glip_webhook_step-2_add-webhook.png)

#### Get the Webhook URL

![](glip_webhook_step-3_details.png)

## Example Requests

Most of the time you will likely either:

* use the proxy URL in an outbound webhook service that supports the Slack format or
* use a client library

The following examples are provided for reference and testing.

### Using `application/json`

```
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"username":"ghost-bot", "icon_emoji": ":ghost:", "text":"BOO!"}' \
  "http://localhost:8080/slack/glip/11112222-3333-4444-5555-666677778888"
```

### Using `application/x-www-form-urlencoded`

```
curl -X POST \
  --data-urlencode 'payload={"username":"ghost-bot", "icon_emoji": ":ghost:", text":"BOO!"}' \
  "http://localhost:8080/slack/glip/11112222-3333-4444-5555-666677778888"
```

### Using `multipart/form-data`

```
curl -X POST \
  -F 'payload={"username":"ghost-bot", "icon_emoji": ":ghost:", text":"BOO!"}' \
  "http://localhost:8080/slack/glip/11112222-3333-4444-5555-666677778888"
```

### Using Community Ruby SDK

This has been tested using:

* [https://github.com/rikas/slack-poster](https://github.com/rikas/slack-poster)

```ruby
require 'slack/poster'

url = 'http://localhost:8080/slack/glip/11112222-3333-4444-5555-666677778888'

opts = {
	username: "MyBot [Bot]",
	icon_emoji: ':ghost:'
}

poster = Slack::Poster.new(url, opts)
poster.send_message('BOO!')
```

## Notes

Glip Webhook Proxy is built using:

* [https://github.com/valyala/fasthttp](https://github.com/valyala/fasthttp)
* [https://github.com/buaazp/fasthttprouter](https://github.com/buaazp/fasthttprouter)

 [build-status-svg]: https://api.travis-ci.org/grokify/glip-webhook-proxy.svg?branch=master
 [build-status-link]: https://travis-ci.org/grokify/glip-webhook-proxy
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/glip-webhook-proxy
 [goreport-link]: https://goreportcard.com/report/github.com/grokify/glip-webhook-proxy
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-link]: https://github.com/grokify/glip-webhook-proxy/blob/master/LICENSE.mds
