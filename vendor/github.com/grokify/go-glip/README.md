Glip Webhook Client in Go
=========================

[![Build Status][build-status-svg]][build-status-link]
[![Go Report Card][goreport-svg]][goreport-link]
[![Docs][docs-godoc-svg]][docs-godoc-link]
[![License][license-svg]][license-link]
[![Chat][chat-svg]][chat-url]

## Installation

```bash
$ go get github.com/grokify/go-glip
```

## Usage

1. Create a webhook URL for a conversation in Glip
2. Use the code below to send a message to the webhook URL

```go
import (
    "fmt"
    "github.com/grokify/go-glip"
)

func sendMessage() {
    // Can instantiate webhook client with full URL or GUID only
    url := "https://hooks.glip.com/webhook/00001111-2222-3333-4444-555566667777"
    client, err := glipwebhook.NewGlipWebhookClient(url)
    if err != nil {
        panic("BAD URL")
    }

    msg := glipwebhook.GlipWebhookMessage{
        Icon:     "https://raw.githubusercontent.com/grokify/glip-go-webhook/master/glip_gopher_600x600xfff.png",
        Activity: "Gopher [Bot]",
        Title:    "Test Message Title",
        Body:     "Test Message Body"}

    resp, err := client.PostMessage(msg)

    respBodyBytes, err := client.SendMessage(msg)
    if err == nil {
        fmt.Printf("%v\n", string(respBodyBytes))
    }
}
```

### Using `fasthttp` client

Posts can be made using [`fasthttp`](https://github.com/valyala/fasthttp).

```go
import (
    "fmt"
    "github.com/grokify/go-glip"
)

func sendMessage() {
    // Can instantiate webhook client with full URL or GUID only
    url := "https://hooks.glip.com/webhook/00001111-2222-3333-4444-555566667777"
    client, err := glipwebhook.NewGlipWebhookClientFast(url)
    if err != nil {
        panic("BAD URL")
    }

    msg := glipwebhook.GlipWebhookMessage{
        Body: "Test Message Body"}

    req, resp, err := client.PostMessageFast(msg)
    if err == nil {
        fmt.Println(string(resp.Body()))
    }
    fasthttp.ReleaseRequest(req)
    fasthttp.ReleaseResponse(resp)
}
```

You can reuse the client for different Webhook URLs or GUIDs as follows:

```go
// Webhook URL
res, resp, err := client.PostWebhookFast(url, msg)

// Webhook GUID
res, resp, err := client.PostWebhookGUIDFast(guid, msg)
```

 [build-status-svg]: https://api.travis-ci.org/grokify/go-glip.svg?branch=master
 [build-status-link]: https://travis-ci.org/grokify/go-glip
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/go-glip
 [goreport-link]: https://goreportcard.com/report/github.com/grokify/go-glip
 [docs-godoc-svg]: https://img.shields.io/badge/docs-godoc-blue.svg
 [docs-godoc-link]: https://godoc.org/github.com/grokify/go-glip
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-link]: https://github.com/grokify/go-glip/blob/master/LICENSE.md
 [chat-svg]: https://img.shields.io/badge/chat-on%20glip-orange.svg
 [chat-url]: https://glipped.herokuapp.com/
