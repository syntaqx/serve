// Package serve provides a static http server anywhere you need one.
package serve

import (
	"net/http"
	"os"
)

// FileServer is an http.Handler that serves static files from a directory tree.
// The zero value is not usable; construct one with NewFileServer.
type FileServer struct {
	handler http.Handler
}

// Option configures a FileServer.
type Option func(*options)

type options struct {
	directory  string
	prefix     string
	showHidden bool
	listDirs   bool
}

// WithDirectory sets the root directory from which files are served. When
// unset (or empty) the current working directory is used.
func WithDirectory(dir string) Option {
	return func(o *options) { o.directory = dir }
}

// WithPrefix strips the given URL path prefix before resolving files, mirroring
// http.StripPrefix.
func WithPrefix(prefix string) Option {
	return func(o *options) { o.prefix = prefix }
}

// WithHiddenFiles controls whether dotfiles (paths whose name begins with ".",
// such as .git or .env) are served. They are hidden by default; pass true to
// expose them.
func WithHiddenFiles(show bool) Option {
	return func(o *options) { o.showHidden = show }
}

// WithDirectoryListing controls whether directories without an index.html are
// rendered as an automatic file listing. Listings are disabled by default; pass
// true to enable them. Directories that contain an index.html always serve it.
func WithDirectoryListing(enabled bool) Option {
	return func(o *options) { o.listDirs = enabled }
}

// NewFileServer builds a FileServer from the provided options.
//
// Files are served through io/fs, which rejects path traversal at the boundary.
// By default dotfiles are hidden and directories without an index.html are not
// listed; use WithHiddenFiles and WithDirectoryListing to opt into either.
func NewFileServer(opts ...Option) *FileServer {
	o := options{directory: "."}
	for _, fn := range opts {
		fn(&o)
	}
	if o.directory == "" {
		o.directory = "."
	}

	fsys := os.DirFS(o.directory)
	if !o.showHidden {
		fsys = hiddenFS{fsys: fsys}
	}

	handler := http.FileServerFS(fsys)
	if !o.listDirs {
		handler = disableDirectoryListing(fsys, handler)
	}
	if o.prefix != "" {
		handler = http.StripPrefix(o.prefix, handler)
	}

	return &FileServer{handler: handler}
}

// Use wraps the handler with the given middleware. Middleware is applied in the
// order provided, so the last argument becomes the outermost layer and runs
// first for each request.
func (fs *FileServer) Use(mws ...func(http.Handler) http.Handler) {
	for _, mw := range mws {
		fs.handler = mw(fs.handler)
	}
}

// ServeHTTP implements the net/http.Handler interface.
func (fs *FileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fs.handler.ServeHTTP(w, r)
}
