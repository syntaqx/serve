package middleware

import "net/http"

// statusWriter wraps http.ResponseWriter to capture the response status code
// and number of bytes written for logging.
type statusWriter struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (w *statusWriter) WriteHeader(status int) {
	if w.status == 0 {
		w.status = status
	}
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	n, err := w.ResponseWriter.Write(b)
	w.bytes += n
	return n, err
}

// Unwrap exposes the underlying writer so http.ResponseController can reach
// optional interfaces such as http.Flusher and http.Hijacker.
func (w *statusWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}
