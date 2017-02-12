package appsignal

import (
	"github.com/grokify/glip-go-webhook"
)

func ExampleMarkerMessageGlip() (glipwebhook.GlipWebhookMessage, error) {
	msg, err := ExampleMarkerMessageSource()
	if err != nil {
		return glipwebhook.GlipWebhookMessage{}, err
	}
	return Normalize(msg), nil
}

func ExampleMarkerMessageSource() (AppsignalOutMessage, error) {
	return AppsignalOutMessageFromBytes(ExampleMarkerMessageBytes())
}

func ExampleMarkerMessageBytes() []byte {
	return []byte(`{
  "marker":{
    "user": "thijs",
    "site": "AppSignal",
    "environment": "test",
    "revision": "3107ddc4bb053d570083b4e3e425b8d62532ddc9",
    "repository": "git@github.com:appsignal/appsignal.git",
    "url": "https://appsignal.com/test/sites/1385f7e38c5ce90000000000/web/exceptions"
  }
}`)
}

func ExampleExceptionMessageGlip() (glipwebhook.GlipWebhookMessage, error) {
	msg, err := ExampleExceptionMessageSource()
	if err != nil {
		return glipwebhook.GlipWebhookMessage{}, err
	}
	return Normalize(msg), nil
}

func ExampleExceptionMessageSource() (AppsignalOutMessage, error) {
	return AppsignalOutMessageFromBytes(ExampleExceptionMessageBytes())
}

func ExampleExceptionMessageBytes() []byte {
	return []byte(`{
  "exception":{
    "exception": "ActionView::Template::Error",
    "site": "AppSignal",
    "message": "undefined method 'encoding' for nil:NilClass",
    "action": "App::ErrorController#show",
    "path": "/errors",
    "revision": "3107ddc4bb053d570083b4e3e425b8d62532ddc9",
    "user": "thijs",
    "first_backtrace_line": "/usr/local/rbenv/versions/2.0.0-p353/lib/ruby/2.0.0/cgi/util.rb:7:in 'escape'",
    "url": "https://appsignal.com/test/sites/1385f7e38c5ce90000000000/web/exceptions/App::SnapshotsController-show/ActionView::Template::Error",
    "environment": "test"
  }
}`)
}

func ExamplePerformanceMessageGlip() (glipwebhook.GlipWebhookMessage, error) {
	msg, err := ExamplePerformanceMessageSource()
	if err != nil {
		return glipwebhook.GlipWebhookMessage{}, err
	}
	return Normalize(msg), nil
}

func ExamplePerformanceMessageSource() (AppsignalOutMessage, error) {
	return AppsignalOutMessageFromBytes(ExamplePerformanceMessageBytes())
}

func ExamplePerformanceMessageBytes() []byte {
	return []byte(`{
  "performance":{
    "site": "AppSignal",
    "action": "App::ExceptionsController#index",
    "path": "/slow",
    "duration": 552.7897429999999,
    "status": 200,
    "hostname": "frontend.appsignal.com",
    "revision": "3107ddc4bb053d570083b4e3e425b8d62532ddc9",
    "user": "thijs",
    "url": "https://appsignal.com/test/sites/1385f7e38c5ce90000000000/web/performance/App::ExceptionsController-index",
    "environment": "test"
  }
}`)
}
