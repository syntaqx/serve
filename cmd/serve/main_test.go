package main

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/syntaqx/serve/internal/config"
)

func TestRunVersionSubcommand(t *testing.T) {
	var stdout, stderr bytes.Buffer

	err := run(context.Background(), []string{"version"}, &stdout, &stderr)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "serve version")
}

func TestRunVersionFlag(t *testing.T) {
	var stdout, stderr bytes.Buffer

	err := run(context.Background(), []string{"-version"}, &stdout, &stderr)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "serve version")
}

func TestRunParseError(t *testing.T) {
	var stdout, stderr bytes.Buffer

	err := run(context.Background(), []string{"-nonexistent"}, &stdout, &stderr)
	assert.Error(t, err)
}

func TestRunServeAndShutdown(t *testing.T) {
	t.Setenv("SERVE_HOST", "127.0.0.1")
	t.Setenv("SERVE_PORT", "0")

	dir := t.TempDir()

	ctx, cancel := context.WithCancel(context.Background())
	var stdout, stderr bytes.Buffer

	done := make(chan error, 1)
	go func() { done <- run(ctx, []string{dir}, &stdout, &stderr) }()

	time.Sleep(50 * time.Millisecond)
	cancel()

	select {
	case err := <-done:
		assert.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("run did not return in time")
	}
}

func TestRunServerError(t *testing.T) {
	t.Setenv("SERVE_HOST", "127.0.0.1")
	t.Setenv("SERVE_PORT", "99999")

	var stdout, stderr bytes.Buffer
	err := run(context.Background(), []string{t.TempDir()}, &stdout, &stderr)
	assert.Error(t, err)
}

func TestRunResolveDirError(t *testing.T) {
	prev := resolveDir
	t.Cleanup(func() { resolveDir = prev })
	resolveDir = func(...string) (string, error) { return "", errors.New("boom") }

	var stdout, stderr bytes.Buffer
	err := run(context.Background(), nil, &stdout, &stderr)
	assert.Error(t, err)
}

func TestRealMainSuccess(t *testing.T) {
	withArgs(t, []string{"serve", "-version"})
	assert.Equal(t, 0, realMain())
}

func TestRealMainHelp(t *testing.T) {
	withArgs(t, []string{"serve", "-h"})
	assert.Equal(t, 0, realMain())
}

func TestRealMainError(t *testing.T) {
	withArgs(t, []string{"serve", "-nonexistent"})
	assert.Equal(t, 1, realMain())
}

func TestMainInvokesExit(t *testing.T) {
	withArgs(t, []string{"serve", "-version"})

	prev := exit
	t.Cleanup(func() { exit = prev })

	var code int
	exit = func(c int) { code = c }
	main()

	assert.Equal(t, 0, code)
}

func TestPositionalDir(t *testing.T) {
	assert.Equal(t, "public", positionalDir([]string{"public"}))
	assert.Equal(t, "", positionalDir([]string{"version"}))
	assert.Equal(t, "", positionalDir(nil))
}

func TestNewLogger(t *testing.T) {
	var buf bytes.Buffer

	json := newLogger(&config.Config{LogFormat: "json", LogLevel: "info"}, &buf)
	json.Info("hello")
	assert.Contains(t, buf.String(), `"msg":"hello"`)

	buf.Reset()
	text := newLogger(&config.Config{LogFormat: "text", LogLevel: "warn"}, &buf)
	text.Info("skipped")
	text.Warn("kept")
	assert.NotContains(t, buf.String(), "skipped")
	assert.Contains(t, buf.String(), "kept")
}

func TestNewLoggerDebugOverride(t *testing.T) {
	var buf bytes.Buffer

	logger := newLogger(&config.Config{Debug: true, LogLevel: "error"}, &buf)
	logger.Debug("visible")
	assert.Contains(t, buf.String(), "visible")
}

func TestParseLevel(t *testing.T) {
	cases := map[string]slog.Level{
		"debug":   slog.LevelDebug,
		"INFO":    slog.LevelInfo,
		"warn":    slog.LevelWarn,
		"warning": slog.LevelWarn,
		"error":   slog.LevelError,
		"bogus":   slog.LevelInfo,
		"":        slog.LevelInfo,
	}
	for in, want := range cases {
		assert.Equal(t, want, parseLevel(in), in)
	}
}

func withArgs(t *testing.T, args []string) {
	t.Helper()
	prev := os.Args
	t.Cleanup(func() { os.Args = prev })
	os.Args = args
}
