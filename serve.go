// Package serve provides a static http server anywhere you need one.
package serve

import "net/http"

// Options is a struct for specifying configuration options for a FileServer.
type Options struct {
	// Directory is the root directory from which to serve files.
	Directory string

	// Prefix is a filepath prefix that should be ignored by the FileServer.
	Prefix string
}

// FileServer wraps an http.FileServer.
type FileServer struct {
	opt     Options
	handler http.Handler
}

// NewFileServer initializes a FileServer.
func NewFileServer(options ...Options) *FileServer {
	var opt Options
	if len(options) == 0 {
		opt = Options{}
	} else {
		opt = options[0]
	}

	fs := &FileServer{
		opt: opt,
	}

	fs.handler = http.StripPrefix(opt.Prefix, http.FileServer(http.Dir(opt.Directory)))
	return fs
}

// Use wraps the Handler with middleware(s).
func (fs *FileServer) Use(mws ...func(http.Handler) http.Handler) {
	for _, h := range mws {
		fs.handler = h(fs.handler)
	}
}

// ServeHTTP implements the net/http.Handler interface.
func (fs *FileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fs.handler.ServeHTTP(w, r)
}
