# Test Proxy Send

`proxy_send.go` is a test program which will send a test message to a live Chathooks URL.

Run it from the command line as follows:

```
$ go run proxy_send.go --url https://hooks.glip.com/webhook/<my_webhook_id> \
--input bugsnag --output glip --token <my_token> -c <my_chathooks_hook_url>
```

Using AWS API Gateway, a hook URL can look like the following: 

`https://0123456789.execute-api.us-west-1.amazonaws.com/prod/hook`
