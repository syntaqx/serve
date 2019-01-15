package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeDirFlagArg(t *testing.T) {
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
			assert.Equal(tt.expected, sanitizeDirFlagArg(tt.opt, tt.cmd))
		})
	}
}
