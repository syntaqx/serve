package serve

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFileServerDefaults(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	fs := NewFileServer()
	assert.NotNil(fs)
	assert.NotNil(fs.handler)
}

func TestNewFileServerEmptyDirectory(t *testing.T) {
	t.Parallel()

	// An explicitly empty directory falls back to the current directory.
	fs := NewFileServer(WithDirectory(""))
	assert.NotNil(t, fs)
	assert.NotNil(t, fs.handler)
}

func TestFileServerMissingFile(t *testing.T) {
	t.Parallel()

	fs := NewFileServer(WithDirectory(t.TempDir()))
	srv := httptest.NewServer(fs)
	defer srv.Close()

	res, err := http.Get(srv.URL + "/missing.txt")
	require.NoError(t, err)
	res.Body.Close()
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestFileServerUse(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(err)
	res := httptest.NewRecorder()

	testMiddleware1 := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("start\n"))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}

	testMiddleware2 := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte("end\n"))
		}
		return http.HandlerFunc(fn)
	}

	fs := &FileServer{
		handler: http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			t.Fail()
		}),
	}

	fs.Use(testMiddleware2, testMiddleware1)

	fs.ServeHTTP(res, req)

	assert.Equal("start\nend\n", res.Body.String())
}

func TestFileServerServeHTTP(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(err)
	res := httptest.NewRecorder()

	fs := &FileServer{
		handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte("expected"))
		}),
	}

	fs.ServeHTTP(res, req)

	assert.Equal("expected", res.Body.String())
}

func TestFileServerServesFiles(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "index.html"), []byte("hello"), 0o600))

	fs := NewFileServer(WithDirectory(dir))
	srv := httptest.NewServer(fs)
	defer srv.Close()

	res, err := http.Get(srv.URL + "/index.html")
	require.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestFileServerHidesDotfiles(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "public.txt"), []byte("ok"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".env"), []byte("SECRET=1"), 0o600))
	require.NoError(t, os.MkdirAll(filepath.Join(dir, ".git"), 0o700))
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".git", "config"), []byte("[core]"), 0o600))

	t.Run("blocked by default", func(t *testing.T) {
		t.Parallel()
		fs := NewFileServer(WithDirectory(dir))
		srv := httptest.NewServer(fs)
		defer srv.Close()

		for _, path := range []string{"/.env", "/.git/config"} {
			res, err := http.Get(srv.URL + path)
			require.NoError(t, err)
			res.Body.Close()
			assert.Equal(t, http.StatusNotFound, res.StatusCode, path)
		}

		res, err := http.Get(srv.URL + "/public.txt")
		require.NoError(t, err)
		res.Body.Close()
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("listing omits dotfiles", func(t *testing.T) {
		t.Parallel()
		fs := NewFileServer(WithDirectory(dir))
		srv := httptest.NewServer(fs)
		defer srv.Close()

		res, err := http.Get(srv.URL + "/")
		require.NoError(t, err)
		body := readBody(t, res)
		assert.Contains(t, body, "public.txt")
		assert.NotContains(t, body, ".env")
		assert.NotContains(t, body, ".git")
	})

	t.Run("opt-in exposes dotfiles", func(t *testing.T) {
		t.Parallel()
		fs := NewFileServer(WithDirectory(dir), WithHiddenFiles(true))
		srv := httptest.NewServer(fs)
		defer srv.Close()

		res, err := http.Get(srv.URL + "/.env")
		require.NoError(t, err)
		res.Body.Close()
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}

func TestFileServerPrefix(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "file.txt"), []byte("ok"), 0o600))

	fs := NewFileServer(WithDirectory(dir), WithPrefix("/static"))
	srv := httptest.NewServer(fs)
	defer srv.Close()

	res, err := http.Get(srv.URL + "/static/file.txt")
	require.NoError(t, err)
	res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestFileServerDirectoryListing(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "sub"), 0o700))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "sub", "file.txt"), []byte("ok"), 0o600))
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "withindex"), 0o700))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "withindex", "index.html"), []byte("home"), 0o600))

	t.Run("enabled by default", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(NewFileServer(WithDirectory(dir)))
		defer srv.Close()

		res, err := http.Get(srv.URL + "/sub/")
		require.NoError(t, err)
		body := readBody(t, res)
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Contains(t, body, "file.txt")
	})

	t.Run("disabled returns 404 without index", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(NewFileServer(WithDirectory(dir), WithDirectoryListing(false)))
		defer srv.Close()

		res, err := http.Get(srv.URL + "/sub/")
		require.NoError(t, err)
		res.Body.Close()
		assert.Equal(t, http.StatusNotFound, res.StatusCode)

		// The root directory (no index.html) is also 404, not a listing.
		rootRes, err := http.Get(srv.URL + "/")
		require.NoError(t, err)
		rootRes.Body.Close()
		assert.Equal(t, http.StatusNotFound, rootRes.StatusCode)

		// A file inside is still served directly.
		fileRes, err := http.Get(srv.URL + "/sub/file.txt")
		require.NoError(t, err)
		fileRes.Body.Close()
		assert.Equal(t, http.StatusOK, fileRes.StatusCode)
	})

	t.Run("disabled still serves index.html", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(NewFileServer(WithDirectory(dir), WithDirectoryListing(false)))
		defer srv.Close()

		res, err := http.Get(srv.URL + "/withindex/")
		require.NoError(t, err)
		body := readBody(t, res)
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Contains(t, body, "home")
	})
}

func readBody(t *testing.T, res *http.Response) string {
	t.Helper()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	return string(data)
}
