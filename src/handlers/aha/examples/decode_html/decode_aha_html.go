package main

import (
	"fmt"
	"html"

	"github.com/grokify/chathooks/src/handlers/aha"
	"github.com/microcosm-cc/bluemonday"
)

func main() {
	//data := []byte(`{"event":"audit","audit":{"id":"1011112222333344445555666"}}`)
	data := []byte(`{"event":"audit","audit":{"id":"6514922586902400152","audit_action":"update","created_at":"2018-01-25T09:46:59.257Z","interesting":true,"user":{"id":"6355516420883588191","name":"John Wang","email":"john.wang@ringcentral.com","created_at":"2016-11-21T20:09:39.022Z","updated_at":"2018-01-24T23:50:12.097Z"},"auditable_type":"note","auditable_id":"6499593214164077189","description":"updated feature RCGG-112 User feedback collection","auditable_url":"https://ringcentral.aha.io/features/RCGG-112","changes":[{"field_name":"Description","value":"\u003cp\u003eBesides Google Chrome store reviews, we should provide user way to submit their feedback from within our Google app\u003c/p\u003e\u003cp\u003e\u003cb\u003eRequirement\u003c/b\u003e\u003c/p\u003e\u003col\u003e\n\u003cli\u003eThere should be a menu item 'Feedback' on Settings page\u003c/li\u003e\n\u003cli\u003eWhen user clicks 'Feedback', user shall be navigated to a new page with following content\u003cblockquote\u003e\n\u003cp\u003e\u003cb\u003eContact Customer Support\u003c/b\u003e\u003cbr\u003e Your feedback is valuable for us. If you have problems using the app, want to request a feature, or report a bug, we’re more than happy to help. Please fill in the form below and click \u003ci\u003eSend Your Feedback\u003c/i\u003e, or directly use your mailbox and send your request to integration.service@ringcentral.com.\u003c/p\u003e\n\u003cp\u003eYour email address (so we can reply to you)\u003cbr\u003e [Input box: ronald.app@ringcentral.com]\u003c/p\u003e\n\u003cp\u003eFeedback topic\u003cbr\u003e [Dropdown: Please select an option]\u003c/p\u003e\n\u003cp\u003eSubject\u003cbr\u003e [Input box: Let us know how we can help you]\u003c/p\u003e\n\u003cp\u003eFull description\u003cbr\u003e [Input box: Please include as much information as possible]\u003c/p\u003e\n\u003cp\u003e\u003cstrike\u003eAttachment\u003c/strike\u003e \u003cbr\u003e \u003cstrike\u003e[Drag and drop area]\u003c/strike\u003e\u003c/p\u003e\n\u003cp\u003e[Button: Send Your Feedback]\u003c/p\u003e\n\u003c/blockquote\u003e\n\u003c/li\u003e\n\u003cli\u003eFeedback topic options in the dropdown list: Please select an option (default) | Bug report | Feature request | Others\u003c/li\u003e\n\u003cli\u003eThere should be back icon on the page by clicking which user can be navigated back to Settings page.\u003c/li\u003e\n\u003cli\u003eWe should leverage an Email server and send the feedback including all the information/\u003cstrike\u003eattachment\u003c/strike\u003e user submitted to team alias \u003ci\u003eintegration.service@ringcentral.com\u003c/i\u003e with title 'Google User Feedback'.\u003cul\u003e\u003cli\u003eEmail content example\u003cblockquote\u003e\n\u003cp\u003eHi Integration Team,\u003c/p\u003e\n\u003cp\u003eYou've got feedback from customer on RingCentral for Google extension. This customer could be contacted via email [customer's email address].\u003c/p\u003e\n\u003cp\u003e\u003cb\u003eCustomer Feedback Topic\u003c/b\u003e\u003cbr\u003e *****\u003c/p\u003e\n\u003cp\u003e\u003cb\u003eSubject\u003c/b\u003e\u003cbr\u003e ******\u003c/p\u003e\n\u003cp\u003e\u003cb\u003eDescription\u003c/b\u003e\u003cbr\u003e ********\u003c/p\u003e\n\u003cp\u003eRegards,\u003cbr\u003e RingCentral for Google Extension\u003c/p\u003e\n\u003c/blockquote\u003e\n\u003c/li\u003e\u003c/ul\u003e\n\u003c/li\u003e\n\u003c/ol\u003e\u003cp\u003ePlease note that the feedback email should support rebranding, to avoid us using wrong words when reaching back to the customers.\u003c/p\u003e"}]}}`)

	msg, err := aha.AhaOutMessageFromBytes(data)
	if err != nil {
		panic(err)
	}

	val := msg.Audit.Changes[0].Value
	fmt.Println(val)

	fmt.Println("---")
	//p := bluemonday.UGCPolicy()
	p := bluemonday.StrictPolicy()
	val = p.Sanitize(val)
	fmt.Println(val)

	fmt.Println("===")
	val = html.UnescapeString(val)
	fmt.Println(val)

}
