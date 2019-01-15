// Package serve provides a static http server anywhere you need one.
package serve

import "net/http"

// FileServer wraps an http.FileServer.
type FileServer struct {
	Handler http.Handler
}

// NewFileServer initializes a FileServer.
func NewFileServer(dir string) *FileServer {
	fs := &FileServer{
		Handler: http.FileServer(http.Dir(dir)),
	}

	return fs
}

// Use wraps the handler with another, middleware style.
func (fs *FileServer) Use(mws ...func(http.Handler) http.Handler) {
	for _, h := range mws {
		fs.Handler = h(fs.Handler)
	}
}

// ServeHTTP implements the net/http.Handler interface.
func (fs *FileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fs.Handler.ServeHTTP(w, r)
}
