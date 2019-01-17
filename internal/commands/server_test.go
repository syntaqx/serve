package commands

import (
	"bytes"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/syntaqx/serve/internal/config"
)

func TestServer(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	var b bytes.Buffer
	log := log.New(&b, "[test] ", 0)
	opt := config.Flags{Port: 0}

	go func() {
		assert.NoError(Server(log, opt, "."))
	}()

	time.Sleep(200 * time.Millisecond)
	assert.Contains(b.String(), "http server listening at")
}

func TestServerErr(t *testing.T) {
	assert := assert.New(t)

	var b bytes.Buffer
	log := log.New(&b, "[test] ", 8888)
	opt := config.Flags{Port: 8888}

	go func() {
		_ = http.ListenAndServe(":8888", nil)
	}()

	time.Sleep(200 * time.Millisecond)

	go func() {
		assert.Error(Server(log, opt, "."))
	}()

	time.Sleep(200 * time.Millisecond)
}

func TestServerHTTPS(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	var b bytes.Buffer
	log := log.New(&b, "[test] ", 0)

	opt := config.Flags{
		Port:      0,
		EnableSSL: true,
		CertFile:  "../../fixtures/cert.pem",
		KeyFile:   "../../fixtures/key.pem",
	}

	go func() {
		assert.NoError(Server(log, opt, "."))
	}()

	time.Sleep(200 * time.Millisecond)
	assert.Contains(b.String(), "https server listening at")
}
