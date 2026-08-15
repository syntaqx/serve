package serve

import (
	"io/fs"
	"net/http"
	"path"
	"strings"
)

// hiddenFS wraps an fs.FS and hides dotfiles: any path containing a segment
// that begins with "." is reported as non-existent, and directory listings
// omit such entries. This prevents accidental exposure of files like .env or
// the contents of a .git directory.
type hiddenFS struct {
	fsys fs.FS
}

func (h hiddenFS) Open(name string) (fs.File, error) {
	if containsDotFile(name) {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}

	f, err := h.fsys.Open(name)
	if err != nil {
		return nil, err
	}

	if dir, ok := f.(fs.ReadDirFile); ok {
		return &hiddenDir{ReadDirFile: dir}, nil
	}
	return f, nil
}

// hiddenDir filters dotfiles out of directory listings.
type hiddenDir struct {
	fs.ReadDirFile
}

func (d *hiddenDir) ReadDir(n int) ([]fs.DirEntry, error) {
	entries, err := d.ReadDirFile.ReadDir(n)
	visible := entries[:0]
	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name(), ".") {
			visible = append(visible, entry)
		}
	}
	return visible, err
}

// containsDotFile reports whether any path segment begins with a dot, ignoring
// the "." and ".." navigation segments.
func containsDotFile(name string) bool {
	for _, part := range strings.Split(name, "/") {
		if part != "." && part != ".." && strings.HasPrefix(part, ".") {
			return true
		}
	}
	return false
}

// disableDirectoryListing wraps next so that a request resolving to a directory
// without an index.html returns 404 instead of an automatic file listing.
// Directories that do contain an index.html continue to serve it.
func disableDirectoryListing(fsys fs.FS, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
		if name == "" {
			name = "."
		}

		if info, err := fs.Stat(fsys, name); err == nil && info.IsDir() {
			if _, err := fs.Stat(fsys, path.Join(name, "index.html")); err != nil {
				http.NotFound(w, r)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
