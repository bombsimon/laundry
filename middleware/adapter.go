package middleware

import (
	"net/http"
)

// Adapter represents an adapter to use as middleware. The adapter takes
// a http.Handler and returns a http.Handler
type Adapter func(http.Handler) http.Handler

// Adapt will take a http.Handler and a list of adapters and adapt each adapter
func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}
