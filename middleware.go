package main

import (
	"log"
	"net/http"
	"strings"
)

// Logger is a middleware that logs each request, along with some useful data
// about what was requested, and what the response was.
func Logger(log *log.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				log.Println(r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
			}()
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// CORS sets permissive cross-origin resource sharing rules.
func CORS() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", strings.Join([]string{
				http.MethodHead,
				http.MethodOptions,
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodPatch,
				http.MethodDelete,
			}, ", "))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

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
