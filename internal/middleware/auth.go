package middleware

import "net/http"

// Auth sets basic HTTP authorization
func Auth(users map[string]string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// Only require auth if we have any users
			if len(users) > 0 {
				authUser, authPass, ok := r.BasicAuth()
				if !ok {
					// No username/password received
					w.Header().Set("WWW-Authenticate", "Basic realm=Authenticate")
					w.WriteHeader(http.StatusUnauthorized)
				} else {
					if pass, ok := users[authUser]; ok {
						// User exists
						if pass == authPass {
							// Authentication successful
							next.ServeHTTP(w, r)
						} else {
							http.Error(w, "Incorrect login details", http.StatusUnauthorized)
							return
						}
					} else {
						http.Error(w, "Incorrect login details", http.StatusUnauthorized)
						return
					}
				}
			} else {
				next.ServeHTTP(w, r)
			}
		}

		return http.HandlerFunc(fn)
	}
}
