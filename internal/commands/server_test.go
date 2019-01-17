package commands

import (
	"bytes"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/syntaqx/serve/internal/config"
	"github.com/syntaqx/serve/mock"
)

func getMockHTTPServerFunc(shouldError bool) func(addr string, h http.Handler) HTTPServer {
	return func(addr string, h http.Handler) HTTPServer {
		return &mock.HTTPServer{ShouldError: shouldError}
	}
}

func TestServer(t *testing.T) {
	t.Parallel()
	getHTTPServerFunc = getMockHTTPServerFunc(false)
	assert := assert.New(t)

	var b bytes.Buffer
	log := log.New(&b, "[test] ", 0)
	opt := config.Flags{}

	assert.NoError(Server(log, opt, "."))
	assert.Contains(b.String(), "http server listening at")
}

func TestServerErr(t *testing.T) {
	t.Parallel()
	getHTTPServerFunc = getMockHTTPServerFunc(true)
	assert := assert.New(t)

	var b bytes.Buffer
	log := log.New(&b, "[test] ", 0)
	opt := config.Flags{}

	time.Sleep(200 * time.Millisecond)

	assert.Error(Server(log, opt, "."))
	time.Sleep(200 * time.Millisecond)
}

func TestServerHTTPS(t *testing.T) {
	t.Parallel()
	getHTTPServerFunc = getMockHTTPServerFunc(false)
	assert := assert.New(t)

	var b bytes.Buffer
	log := log.New(&b, "[test] ", 0)

	opt := config.Flags{
		EnableSSL: true,
		CertFile:  "../../fixtures/cert.pem",
		KeyFile:   "../../fixtures/key.pem",
	}

	assert.NoError(Server(log, opt, "."))
	assert.Contains(b.String(), "https server listening at")
}
