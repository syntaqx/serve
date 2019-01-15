package internal

import (
	"bytes"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestVersionCommand(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	var b bytes.Buffer
	err := VersionCommand("mock", &b)

	assert.NoError(err)
	assert.Contains(b.String(), "version mock")
}

func TestServerCommand(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	var b bytes.Buffer
	log := log.New(&b, "[test] ", 0)
	opt := Flags{Port: 0}

	go func() {
		assert.NoError(ServerCommand(log, opt))
	}()

	time.Sleep(200 * time.Millisecond)
}

func TestServerCommandErr(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	var b bytes.Buffer
	log := log.New(&b, "[test] ", 8888)
	opt := Flags{Port: 8888}

	go func() {
		_ = http.ListenAndServe(":8888", nil)
	}()

	time.Sleep(200 * time.Millisecond)

	go func() {
		assert.Error(ServerCommand(log, opt))
	}()

	time.Sleep(200 * time.Millisecond)
}
