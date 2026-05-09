package middleware

import (
	"mime"
	"net/http"
	"path/filepath"
)

func init() {
	// Ensure common MIME types are registered, as some operating systems
	// (notably Windows) may not have them in their registry.
	types := map[string]string{
		".css":  "text/css; charset=utf-8",
		".html": "text/html; charset=utf-8",
		".js":   "application/javascript",
		".json": "application/json",
		".mjs":  "application/javascript",
		".svg":  "image/svg+xml",
		".wasm": "application/wasm",
		".xml":  "text/xml; charset=utf-8",
	}

	for ext, ct := range types {
		_ = mime.AddExtensionType(ext, ct)
	}
}

// SetContentType overrides the Content-Type header based on file extension
// before the inner handler writes a response. This prevents the file server
// from falling back to "text/plain" on systems with incomplete MIME databases.
func SetContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ext := filepath.Ext(r.URL.Path)
		if ct := mime.TypeByExtension(ext); ct != "" {
			w.Header().Set("Content-Type", ct)
		}
		next.ServeHTTP(w, r)
	})
}
