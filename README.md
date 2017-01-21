Glip Webhook Proxy
==================

Proxy service to map different requests to Glip's inbound webhook service.

## Usage

Start the service with the following.

```go
package main

import (
	"github.com/grokify/glip-go-webhook-proxy"
)

func main() {
	config := glipwebhookproxy.Configuration{
		Port:           ":8080",
		EmojiURLPrefix: "https://grokify.github.io/emoji/assets/images/",
		EmojiURLSuffix: ".png"}

	glipwebhookproxy.StartServer(config)
}
```

## Example Request

```
curl -XPOST \
  -F 'payload={"username":"ghost-bot", "icon_emoji": ":ghost:", text":"BOO!"}' \
  "http://localhost:8080/slack/glip/00001111-2222-3333-4444-555566667777"
```
