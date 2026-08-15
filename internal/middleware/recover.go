package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

// Recover returns middleware that recovers from panics raised while handling a
// request. The panic value and stack trace are logged server-side, and the
// client receives a generic 500 response so internal details are not leaked.
func Recover(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					if logger != nil {
						logger.LogAttrs(r.Context(), slog.LevelError, "panic recovered",
							slog.Any("panic", rec),
							slog.String("method", r.Method),
							slog.String("path", r.URL.Path),
							slog.String("stack", string(debug.Stack())),
						)
					}
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
