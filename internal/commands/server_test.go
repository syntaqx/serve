package commands

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/syntaqx/serve/internal/config"
)

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestHandlerServesFiles(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "index.html"), []byte("<h1>hi</h1>"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".secret"), []byte("nope"), 0o600))

	cfg := &config.Config{CORSOrigin: "*", UsersFile: filepath.Join(dir, "missing")}
	handler, err := Handler(discardLogger(), cfg, dir)
	require.NoError(t, err)

	srv := httptest.NewServer(handler)
	defer srv.Close()

	res, err := http.Get(srv.URL + "/index.html")
	require.NoError(t, err)
	res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "*", res.Header.Get("Access-Control-Allow-Origin"))

	hidden, err := http.Get(srv.URL + "/.secret")
	require.NoError(t, err)
	hidden.Body.Close()
	assert.Equal(t, http.StatusNotFound, hidden.StatusCode)
}

func TestHandlerHealthAndDirListing(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "sub"), 0o700))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "sub", "file.txt"), []byte("ok"), 0o600))
	usersFile := filepath.Join(dir, "users.dat")
	require.NoError(t, os.WriteFile(usersFile, []byte("admin:secret\n"), 0o600))

	cfg := &config.Config{
		CORSOrigin:        "*",
		UsersFile:         usersFile,
		HealthPath:        "/healthz",
		DisableDirListing: true,
	}
	handler, err := Handler(discardLogger(), cfg, dir)
	require.NoError(t, err)

	srv := httptest.NewServer(handler)
	defer srv.Close()

	// Health endpoint bypasses auth and returns 200.
	health, err := http.Get(srv.URL + "/healthz")
	require.NoError(t, err)
	body := readBody(t, health)
	assert.Equal(t, http.StatusOK, health.StatusCode)
	assert.Equal(t, "ok\n", body)

	// Directory listing is disabled: an authed request to the directory 404s.
	dirReq, err := http.NewRequest(http.MethodGet, srv.URL+"/sub/", nil)
	require.NoError(t, err)
	dirReq.SetBasicAuth("admin", "secret")
	dirRes, err := http.DefaultClient.Do(dirReq)
	require.NoError(t, err)
	dirRes.Body.Close()
	assert.Equal(t, http.StatusNotFound, dirRes.StatusCode)
}

func readBody(t *testing.T, res *http.Response) string {
	t.Helper()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	return string(data)
}

func TestHandlerAuth(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "index.html"), []byte("ok"), 0o600))
	usersFile := filepath.Join(dir, "users.dat")
	require.NoError(t, os.WriteFile(usersFile, []byte("admin:hunter2\n"), 0o600))

	cfg := &config.Config{CORSOrigin: "*", UsersFile: usersFile}
	handler, err := Handler(discardLogger(), cfg, dir)
	require.NoError(t, err)

	srv := httptest.NewServer(handler)
	defer srv.Close()

	unauth, err := http.Get(srv.URL + "/index.html")
	require.NoError(t, err)
	unauth.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, unauth.StatusCode)

	req, err := http.NewRequest(http.MethodGet, srv.URL+"/index.html", nil)
	require.NoError(t, err)
	req.SetBasicAuth("admin", "hunter2")
	authed, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	authed.Body.Close()
	assert.Equal(t, http.StatusOK, authed.StatusCode)
}

func TestServerGracefulShutdown(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	cfg := &config.Config{Host: "127.0.0.1", Port: "0", UsersFile: filepath.Join(dir, "missing")}

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)
	go func() { done <- Server(ctx, discardLogger(), cfg, dir) }()

	// Give the server a moment to bind, then trigger shutdown.
	time.Sleep(50 * time.Millisecond)
	cancel()

	select {
	case err := <-done:
		assert.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("server did not shut down in time")
	}
}

func TestServerListenError(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	cfg := &config.Config{Host: "127.0.0.1", Port: "99999", UsersFile: filepath.Join(dir, "missing")}

	err := Server(context.Background(), discardLogger(), cfg, dir)
	assert.Error(t, err)
}

func TestServerHandlerError(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	cfg := &config.Config{Host: "127.0.0.1", Port: "0", UsersFile: "bad\x00name"}

	err := Server(context.Background(), discardLogger(), cfg, dir)
	assert.Error(t, err)
}

func TestServerHTTPSGracefulShutdown(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	cfg := &config.Config{
		Host:      "127.0.0.1",
		Port:      "0",
		EnableSSL: true,
		CertFile:  "../../fixtures/cert.pem",
		KeyFile:   "../../fixtures/key.pem",
		UsersFile: filepath.Join(dir, "missing"),
	}

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)
	go func() { done <- Server(ctx, discardLogger(), cfg, dir) }()

	time.Sleep(50 * time.Millisecond)
	cancel()

	select {
	case err := <-done:
		assert.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("server did not shut down in time")
	}
}

func TestServeWithShutdownServerClosed(t *testing.T) {
	t.Parallel()

	err := serveWithShutdown(context.Background(), discardLogger(), &http.Server{}, func() error {
		return http.ErrServerClosed
	})
	assert.NoError(t, err)
}

func TestServeWithShutdownStartError(t *testing.T) {
	t.Parallel()

	want := errors.New("boom")
	err := serveWithShutdown(context.Background(), discardLogger(), &http.Server{}, func() error {
		return want
	})
	assert.ErrorIs(t, err, want)
}

func TestServeWithShutdownGraceful(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	block := make(chan struct{})

	done := make(chan error, 1)
	go func() {
		done <- serveWithShutdown(ctx, discardLogger(), &http.Server{}, func() error {
			<-block
			return http.ErrServerClosed
		})
	}()

	cancel()
	err := <-done
	close(block)
	assert.NoError(t, err)
}

func TestHandlerDebugMissingUsers(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	cfg := &config.Config{CORSOrigin: "*", UsersFile: filepath.Join(dir, "missing"), Debug: true}

	_, err := Handler(discardLogger(), cfg, dir)
	require.NoError(t, err)
}

func TestHandlerUsersOpenError(t *testing.T) {
	t.Parallel()

	cfg := &config.Config{CORSOrigin: "*", UsersFile: "bad\x00name"}

	_, err := Handler(discardLogger(), cfg, t.TempDir())
	assert.Error(t, err)
}

func TestGetAuthUsers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  string
		output map[string]string
	}{
		{"single", "user1:pass1", map[string]string{"user1": "pass1"}},
		{"multiple", "user1:pass1\nuser2:pass2", map[string]string{"user1": "pass1", "user2": "pass2"}},
		{"empty", "", map[string]string{}},
		{"colon in secret", "user1:pa:ss1", map[string]string{"user1": "pa:ss1"}},
		{"comments and blanks", "# a comment\n\nuser1:pass1\n", map[string]string{"user1": "pass1"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.output, GetAuthUsers(strings.NewReader(tt.input)))
		})
	}
}

func TestGetAuthUsersNil(t *testing.T) {
	t.Parallel()
	assert.Equal(t, map[string]string{}, GetAuthUsers(nil))
}
