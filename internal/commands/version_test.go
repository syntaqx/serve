package commands

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	var b bytes.Buffer
	err := Version("mock", &b)

	assert.NoError(err)
	assert.Contains(b.String(), "version mock")
}
