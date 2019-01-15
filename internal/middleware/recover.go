package middleware

import (
	"fmt"
	"net/http"
)

// Recover is a middleware that recovers from panics that occur for a request.
func Recover() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					http.Error(w, fmt.Sprintf("[PANIC RECOVERED] %v", err), http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
