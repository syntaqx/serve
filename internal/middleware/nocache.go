package middleware

import (
	"net/http"
)

// NoCache sets a number of HTTP Headers instructing clients not to cache a
// given response, or use an existing cache.
func NoCache() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			headers := w.Header()
			headers.Set("Expires", "0")
			headers.Set("Cache-Control", "no-cache, no-store, no-transform, must-revalidate, private, max-age=0")
			headers.Set("Pragma", "no-cache")
			headers.Set("X-Accel-Expires", "0")
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
