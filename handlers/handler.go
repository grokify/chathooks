package handlers

import (
	"net/http"
)

// Result is returned by a handler to provide richer information than a plain
// error. Handlers can use it to control the HTTP status code returned to the
// webhook sender and to request that the message be skipped for filtering
// without treating the request as an error.
type Result struct {
	// StatusCode is the HTTP status code to return. If zero, the caller
	// should use http.StatusOK.
	StatusCode int

	// Skip indicates that the message should not be sent.
	Skip bool
}

// Handler is the interface implemented by all webhook handlers.
type Handler interface {
	// Handle processes the incoming webhook request. It returns a Result
	// describing how the request should be handled, and an error if the
	// request could not be processed.
	Handle(w http.ResponseWriter, r *http.Request) (Result, error)
}

// HandlerFunc adapts an ordinary function to the Handler interface.
type HandlerFunc func(w http.ResponseWriter, r *http.Request) (Result, error)

// Handle implements Handler.
func (f HandlerFunc) Handle(w http.ResponseWriter, r *http.Request) (Result, error) {
	return f(w, r)
}
