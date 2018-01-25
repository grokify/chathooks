package main

import (
	"fmt"
	"html"

	"github.com/grokify/chathooks/src/handlers/aha"
	"github.com/microcosm-cc/bluemonday"
)

func main() {
	data := []byte(`{"event":"audit","audit":{"id":"1011112222333344445555666"}}`)

	msg, err := aha.AhaOutMessageFromBytes(data)
	if err != nil {
		panic(err)
	}

	val := msg.Audit.Changes[0].Value
	fmt.Println(val)

	fmt.Println("---")
	val2 := html.UnescapeString(val)
	fmt.Println(val2)

	fmt.Println("---")
	//p := bluemonday.UGCPolicy()
	p := bluemonday.StrictPolicy()
	val3 := p.Sanitize(val2)
	fmt.Println(val4)

}
