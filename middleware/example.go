package middleware

import (
	"fmt"
	"net/http"
)

func Notify() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("This is an example adapter - calling before handler!")
			defer fmt.Println("This is an example adapter - calling after handler!")

			h.ServeHTTP(w, r)
		})
	}
}
