package config

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeDir(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	cwd, err := os.Getwd()
	assert.NoError(err)

	var tests = []struct {
		opt      string
		cmd      string
		expected string
	}{
		{"foo", "bar", "foo"},
		{"", "bar", "bar"},
		{"", "", cwd},
	}

	for _, tt := range tests {
		tt := tt
		t.Run("", func(t *testing.T) {
			t.Parallel()
			dir, err := SanitizeDir(tt.opt, tt.cmd)
			assert.Equal(tt.expected, dir)
			assert.NoError(err)
		})
	}
}

func TestSanitizeDirCwdErr(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	getwd = func() (string, error) {
		return "", errors.New("mock")
	}

	dir, err := SanitizeDir("", "")
	assert.Empty(dir)
	assert.Error(err)
}
