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
		dirs     []string
		expected string
	}{
		{[]string{"foo", "bar"}, "foo"},
		{[]string{"", "bar"}, "bar"},
		{[]string{"", ""}, cwd},
	}

	for _, tt := range tests {
		tt := tt
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
