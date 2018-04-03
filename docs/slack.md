# Example Requests

Most of the time you will likely either:

* use the proxy URL in an outbound webhook service that supports the Slack format or
* use a client library

The following examples are provided for reference and testing.

## Using `application/json`

```bash
$ curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"username":"ghost-bot", "icon_emoji": ":ghost:", "text":"BOO!"}' \
  "http://localhost:8080/hook?inputType=slack&outputType=glip&url=11112222-3333-4444-5555-666677778888"
```

## Using `application/x-www-form-urlencoded`

```bash
$ curl -X POST \
  --data-urlencode 'payload={"username":"ghost-bot", "icon_emoji": ":ghost:", text":"BOO!"}' \
  "http://localhost:8080/hook?inputType=slack&outputType=glip&url=11112222-3333-4444-5555-666677778888"
```

## Using `multipart/form-data`

```bash
$ curl -X POST \
  -F 'payload={"username":"ghost-bot", "icon_emoji": ":ghost:", text":"BOO!"}' \
  "http://localhost:8080/hook?inputType=slack&outputType=glip&url=11112222-3333-4444-5555-666677778888"
```

## Using Community Ruby SDK

This has been tested using:

* [https://github.com/rikas/slack-poster](https://github.com/rikas/slack-poster)

```ruby
require 'slack/poster'

url = 'http://localhost:8080/hook?inputType=slack&outputType=glip&url=11112222-3333-4444-5555-666677778888'

opts = {
    username: 'Ghost Bot [Bot]',
    icon_emoji: ':ghost:'
}

poster = Slack::Poster.new url, opts
poster.send_message 'BOO!'
```