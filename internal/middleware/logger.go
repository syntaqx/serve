package middleware

import (
	"log"
	"net/http"
)

// Logger is a middleware that logs each request, along with some useful data
// about what was requested, and what the response was.
func Logger(log *log.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			sw := statusWriter{ResponseWriter: w}

			defer func() {
				log.Println(r.Method, r.URL.Path, sw.status, r.RemoteAddr, r.UserAgent())
			}()
			next.ServeHTTP(&sw, r)
		}
		return http.HandlerFunc(fn)
	}
}
