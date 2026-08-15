package config

import (
	"errors"
	"flag"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadDefaults(t *testing.T) {
	cfg, args, err := Load(nil, io.Discard)
	require.NoError(t, err)

	assert.Empty(t, args)
	assert.Equal(t, "0.0.0.0", cfg.Host)
	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, "users.dat", cfg.UsersFile)
	assert.Equal(t, "*", cfg.CORSOrigin)
	assert.Equal(t, "text", cfg.LogFormat)
	assert.Equal(t, "info", cfg.LogLevel)
	assert.Equal(t, "/healthz", cfg.HealthPath)
	assert.False(t, cfg.ShowHidden)
	assert.False(t, cfg.DisableDirListing)
}

func TestLoadNewToggles(t *testing.T) {
	cfg, _, err := Load([]string{"-no-dirlisting", "-health-path", "/status"}, io.Discard)
	require.NoError(t, err)

	assert.True(t, cfg.DisableDirListing)
	assert.Equal(t, "/status", cfg.HealthPath)
}

func TestLoadFlags(t *testing.T) {
	cfg, args, err := Load([]string{"-port", "9000", "-all", "-log-format", "json", "public"}, io.Discard)
	require.NoError(t, err)

	assert.Equal(t, "9000", cfg.Port)
	assert.True(t, cfg.ShowHidden)
	assert.Equal(t, "json", cfg.LogFormat)
	assert.Equal(t, []string{"public"}, args)
}

func TestLoadEnv(t *testing.T) {
	t.Setenv("SERVE_PORT", "7000")
	t.Setenv("SERVE_HOST", "127.0.0.1")
	t.Setenv("SERVE_ALL", "true")

	cfg, _, err := Load(nil, io.Discard)
	require.NoError(t, err)

	assert.Equal(t, "7000", cfg.Port)
	assert.Equal(t, "127.0.0.1", cfg.Host)
	assert.True(t, cfg.ShowHidden)
}

func TestLoadFlagOverridesEnv(t *testing.T) {
	t.Setenv("SERVE_PORT", "7000")

	cfg, _, err := Load([]string{"-port", "8081"}, io.Discard)
	require.NoError(t, err)

	assert.Equal(t, "8081", cfg.Port)
}

func TestLoadPortFallback(t *testing.T) {
	t.Setenv("PORT", "6000")

	cfg, _, err := Load(nil, io.Discard)
	require.NoError(t, err)

	assert.Equal(t, "6000", cfg.Port)
}

func TestLoadHelp(t *testing.T) {
	_, _, err := Load([]string{"-h"}, io.Discard)
	assert.ErrorIs(t, err, flag.ErrHelp)
}

func TestSanitizeDir(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	cwd, err := os.Getwd()
	assert.NoError(err)

	var tests = []struct {
		dirs     []string
		expected string
	}{
		{[]string{"foo", "bar"}, "foo"},
		{[]string{"", "bar"}, "bar"},
		{[]string{"", ""}, cwd},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			t.Parallel()
			dir, err := SanitizeDir(tt.dirs...)
			assert.Equal(tt.expected, dir)
			assert.NoError(err)
		})
	}
}

func TestSanitizeDirCwdErr(t *testing.T) {
	assert := assert.New(t)

	getwd = func() (string, error) {
		return "", errors.New("mock")
	}

	dir, err := SanitizeDir()
	assert.Empty(dir)
	assert.Error(err)

	getwd = os.Getwd
}
