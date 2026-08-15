package middleware

import (
	"crypto/subtle"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// Auth returns middleware that enforces HTTP Basic authentication whenever one
// or more users are configured. When no users are provided the middleware is a
// no-op, allowing all requests through.
//
// Stored secrets may be plaintext or bcrypt hashes. Plaintext secrets are
// compared with crypto/subtle and bcrypt hashes with bcrypt.CompareHashAndPassword,
// both of which run in constant time to avoid leaking information through
// response timing.
func Auth(users map[string]string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if len(users) == 0 {
				next.ServeHTTP(w, r)
				return
			}

			user, pass, ok := r.BasicAuth()
			if !ok {
				w.Header().Set("WWW-Authenticate", `Basic realm="serve"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			stored, known := users[user]
			if !known || !credentialsMatch(stored, pass) {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func credentialsMatch(stored, provided string) bool {
	if isBcryptHash(stored) {
		return bcrypt.CompareHashAndPassword([]byte(stored), []byte(provided)) == nil
	}
	return subtle.ConstantTimeCompare([]byte(stored), []byte(provided)) == 1
}

func isBcryptHash(s string) bool {
	return strings.HasPrefix(s, "$2a$") ||
		strings.HasPrefix(s, "$2b$") ||
		strings.HasPrefix(s, "$2y$")
}
