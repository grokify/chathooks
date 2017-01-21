Glip Webhook Proxy
==================

[![Go Report Card][goreport-svg]][goreport-link]
[![License][license-svg]][license-link]

Proxy service to map different requests to Glip's inbound webhook service.

## Usage

Start the service with the following.

```go
package main

import (
	"github.com/grokify/glip-webhook-proxy"
)

func main() {
	config := glipwebhookproxy.Configuration{
		Port:           ":8080",
		EmojiURLPrefix: "https://grokify.github.io/emoji/assets/images/",
		EmojiURLSuffix: ".png"}

	glipwebhookproxy.StartServer(config)
}
```

Then, to use this service:

1. create a Glip webhook
2. use webhook URL's GUID to create the proxy URL as shown below
3. use the proxy URL in your outbound webhook service

| Name | Value |
|------|-------|
| Glip Webhook URL | `https://hooks.glip.com/webhook/11112222-3333-4444-5555-666677778888` |
| Proxy Webhook URL | `https://example.com/slack/glip/11112222-3333-4444-5555-666677778888` |

## Example Requests

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

## Functionality

* handles all request content types
* converts payload property names
* converts emoji to URL
* tested with SDKs

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
