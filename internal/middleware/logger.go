package middleware

import (
	"log"
	"net/http"
)

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	return w.ResponseWriter.Write(b)
}

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
