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

func getMockHTTPServer(addr string, h http.Handler) HTTPServer {
	return &mock.HTTPServer{}
}

func getMockErrHTTPServer(addr string, h http.Handler) HTTPServer {
	return &mock.HTTPServer{ShouldError: true}
}

func TestGetStdHTTPServer(t *testing.T) {
	_, ok := GetStdHTTPServer("", http.DefaultServeMux).(*http.Server)
	assert.True(t, ok)
}

func TestServer(t *testing.T) {
	getHTTPServerFunc = getMockHTTPServer

	assert := assert.New(t)

	var b bytes.Buffer
	log := log.New(&b, "[test] ", 0)
	opt := config.Flags{}

	go func() {
		assert.NoError(Server(log, opt, "."))
	}()

	time.Sleep(200 * time.Millisecond)
	assert.Contains(b.String(), "http server listening at")
}

func TestServerErr(t *testing.T) {
	getHTTPServerFunc = getMockErrHTTPServer

	assert := assert.New(t)

	var b bytes.Buffer
	log := log.New(&b, "[test] ", 0)
	opt := config.Flags{}

	time.Sleep(200 * time.Millisecond)

	go func() {
		assert.Error(Server(log, opt, "."))
	}()

	time.Sleep(200 * time.Millisecond)
}

func TestServerHTTPS(t *testing.T) {
	getHTTPServerFunc = getMockHTTPServer

	t.Parallel()
	assert := assert.New(t)

	var b bytes.Buffer
	log := log.New(&b, "[test] ", 0)

	opt := config.Flags{
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
