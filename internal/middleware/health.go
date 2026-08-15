package middleware

import "net/http"

// Health returns middleware that responds to GET/HEAD requests for the given
// path with a 200 "ok", short-circuiting the rest of the chain. This is useful
// for load-balancer and container liveness probes. An empty path disables the
// endpoint, making the middleware a no-op.
func Health(path string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if path == "" {
			return next
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == path && (r.Method == http.MethodGet || r.Method == http.MethodHead) {
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("ok\n"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
