package commands

import (
	"bytes"
	"log"
	"net/http"
	"strings"
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

func TestGetStdHTTPServer(t *testing.T) {
	_, ok := GetStdHTTPServer("", http.DefaultServeMux).(*http.Server)
	assert.True(t, ok)
}

func TestServer(t *testing.T) {
	getHTTPServerFunc = getMockHTTPServerFunc(false)

	assert := assert.New(t)

	var b bytes.Buffer
	log := log.New(&b, "[test] ", 0)
	opt := config.Flags{}

	assert.NoError(Server(log, opt, "."))
	assert.Contains(b.String(), "http server listening at")

	getHTTPServerFunc = GetStdHTTPServer
}

func TestServerErr(t *testing.T) {
	getHTTPServerFunc = getMockHTTPServerFunc(true)

	assert := assert.New(t)

	var b bytes.Buffer
	log := log.New(&b, "[test] ", 0)
	opt := config.Flags{}

	time.Sleep(200 * time.Millisecond)

	assert.Error(Server(log, opt, "."))
	time.Sleep(200 * time.Millisecond)

	getHTTPServerFunc = GetStdHTTPServer
}

func TestServerHTTPS(t *testing.T) {
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

	getHTTPServerFunc = GetStdHTTPServer
}

func TestGetAuthUsers(t *testing.T) {
	tests := []struct {
		input  string
		output map[string]string
	}{
		{ // Single user
			"user1:pass1", map[string]string{
				"user1": "pass1",
			},
		},
		{ // Multiple users
			"user1:pass1\nuser2:pass2", map[string]string{
				"user1": "pass1",
				"user2": "pass2",
			},
		},
		{ // Empty file
			"", map[string]string{},
		},
		{ // Incorrect structure
			"user1:pass1:field1", map[string]string{},
		},
	}

	for _, test := range tests {
		mockFile := strings.NewReader(test.input)
		assert.Equal(t, GetAuthUsers(mockFile), test.output)
	}

}
