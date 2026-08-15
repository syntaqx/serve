package commands

import (
	"bytes"
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersion(t *testing.T) {
	assert := assert.New(t)

	var b bytes.Buffer
	err := Version("mock", &b)

	assert.NoError(err)
	assert.Contains(b.String(), "version mock")
}

func TestVersionWithRevision(t *testing.T) {
	prev := readBuildInfo
	t.Cleanup(func() { readBuildInfo = prev })

	readBuildInfo = func() (*debug.BuildInfo, bool) {
		return &debug.BuildInfo{Settings: []debug.BuildSetting{
			{Key: "vcs.revision", Value: "0123456789abcdef0123"},
		}}, true
	}

	var b bytes.Buffer
	require.NoError(t, Version("mock", &b))
	assert.Contains(t, b.String(), "(0123456789ab)")
}

func TestVCSRevision(t *testing.T) {
	prev := readBuildInfo
	t.Cleanup(func() { readBuildInfo = prev })

	tests := []struct {
		name string
		info func() (*debug.BuildInfo, bool)
		want string
	}{
		{
			name: "unavailable",
			info: func() (*debug.BuildInfo, bool) { return nil, false },
			want: "",
		},
		{
			name: "no revision setting",
			info: func() (*debug.BuildInfo, bool) {
				return &debug.BuildInfo{Settings: []debug.BuildSetting{{Key: "vcs.time", Value: "now"}}}, true
			},
			want: "",
		},
		{
			name: "short revision",
			info: func() (*debug.BuildInfo, bool) {
				return &debug.BuildInfo{Settings: []debug.BuildSetting{{Key: "vcs.revision", Value: "abc123"}}}, true
			},
			want: "abc123",
		},
		{
			name: "long revision truncated",
			info: func() (*debug.BuildInfo, bool) {
				return &debug.BuildInfo{Settings: []debug.BuildSetting{{Key: "vcs.revision", Value: "abcdefabcdef123456"}}}, true
			},
			want: "abcdefabcdef",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			readBuildInfo = tt.info
			assert.Equal(t, tt.want, vcsRevision())
		})
	}
}
