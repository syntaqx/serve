package middleware

import (
	"net/http"
	"strings"
)

// CORS sets cross-origin resource sharing headers using the given allowed
// origin (use "*" to allow any) and short-circuits preflight (OPTIONS)
// requests with a 204 response.
func CORS(origin string) func(next http.Handler) http.Handler {
	if origin == "" {
		origin = "*"
	}

	methods := strings.Join([]string{
		http.MethodHead,
		http.MethodOptions,
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
	}, ", ")

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", methods)
			w.Header().Set("Access-Control-Allow-Headers", "*")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
