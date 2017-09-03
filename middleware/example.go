package middleware

import (
	"net/http"
)

// Notify will implement the Adapter interface by returning a function
// taking and return a http.Handler. The Notify() adapter is an example
// showing how to create custom adapters.
func Notify() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// fmt.Println("This will be executed before routing")
			// defer fmt.Println("This will be executed after routing")

			h.ServeHTTP(w, r)
		})
	}
}
