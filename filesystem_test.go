package serve

import (
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHiddenFSOpen(t *testing.T) {
	t.Parallel()

	// fstest.MapFS regular files do not implement fs.ReadDirFile, which
	// exercises the non-directory branch of hiddenFS.Open.
	base := fstest.MapFS{
		"file.txt":    {Data: []byte("ok")},
		".secret":     {Data: []byte("no")},
		"dir/visible": {Data: []byte("y")},
		"dir/.hidden": {Data: []byte("x")},
	}
	h := hiddenFS{fsys: base}

	t.Run("regular file", func(t *testing.T) {
		t.Parallel()
		f, err := h.Open("file.txt")
		require.NoError(t, err)
		require.NoError(t, f.Close())
	})

	t.Run("dotfile blocked", func(t *testing.T) {
		t.Parallel()
		_, err := h.Open(".secret")
		assert.ErrorIs(t, err, fs.ErrNotExist)
	})

	t.Run("missing file", func(t *testing.T) {
		t.Parallel()
		_, err := h.Open("nope.txt")
		assert.Error(t, err)
	})

	t.Run("directory listing filters dotfiles", func(t *testing.T) {
		t.Parallel()
		f, err := h.Open("dir")
		require.NoError(t, err)
		defer f.Close()

		rd, ok := f.(fs.ReadDirFile)
		require.True(t, ok)

		entries, err := rd.ReadDir(-1)
		require.NoError(t, err)

		names := make([]string, 0, len(entries))
		for _, e := range entries {
			names = append(names, e.Name())
		}
		assert.Equal(t, []string{"visible"}, names)
	})
}

func TestContainsDotFile(t *testing.T) {
	t.Parallel()

	cases := map[string]bool{
		".":           false,
		"..":          false,
		"index.html":  false,
		"dir/file":    false,
		".env":        true,
		"dir/.git":    true,
		".git/config": true,
	}
	for name, want := range cases {
		assert.Equal(t, want, containsDotFile(name), name)
	}
}
