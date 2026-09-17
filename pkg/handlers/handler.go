package handlers

import "net/http"

// Result represents the outcome of a handler processing a webhook event.
// It allows handlers to specify a custom status code and to skip
// sending a notification for desired filtering.
type Result struct {
	// Message is the formatted notification message to be sent.
	// An empty message means no notification is sent.
	Message string
	// StatusCode is the HTTP status code returned to the webhook caller.
	// If zero, http.StatusOK is used.
	StatusCode int
	// Skip indicates that the notification should not be sent.
	// This enables handlers to filter out undesired events.
	Skip bool
}

// EffectiveStatusCode returns the HTTP status code to use,
// defaulting to http.StatusOK when StatusCode is zero.
func (r *Result) EffectiveStatusCode() int {
	if r.StatusCode == 0 {
		return http.StatusOK
	}
	return r.StatusCode
}

// Handler is the function type for processing webhook events.
// It receives the raw request payload and returns a *Result
// describing the desired response. A non-nil error signals a
// processing failure and results in an HTTP 500 response.
type Handler func(payload []byte) (*Result, error)
