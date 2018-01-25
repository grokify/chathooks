package main

import (
	"fmt"

	"github.com/google/go-querystring/query"
)

type Options struct {
	InputType  string `url:"inputType"`
	OutputType string `url:"outputType"`
	URL        string `url:"url"`
	Token      string `url:"token"`
}

func BuildURL(baseUrl string, opts Options) string {
	v, _ := query.Values(opts)
	return fmt.Sprintf("%v?%v", baseUrl, v.Encode())
}

func main() {
	baseUrl := "https://12345678.ngrok.io/hook"
	opts := Options{
		InputType:  "aha",
		OutputType: "glip",
		URL:        "https://hooks.glip.com/webhook/11112222-3333-4444-5555-666677778888",
		Token:      "deadbeefdeadbeefdeadbeefdeadbeef",
	}
	fmt.Println(BuildURL(baseUrl, opts))
}
