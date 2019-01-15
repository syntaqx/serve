package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartHTTPServer(t *testing.T) {
	t.Skip()

	t.Parallel()

	assert := assert.New(t)

	opt := flags{Port: 0, Dir: "."}

	var b bytes.Buffer
	log := log.New(&b, "[test] ", 0)

	go startHTTPServer(opt, log)

	assert.Contains(b.String(), "http server listening at")
}

func TestLogger(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	var b bytes.Buffer
	log := log.New(&b, "[test] ", 0)

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(err)
	res := httptest.NewRecorder()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	Logger(log)(testHandler).ServeHTTP(res, req)

	assert.Equal("[test] GET /", strings.TrimSpace(b.String()))
}
func TestCORS(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(err)
	res := httptest.NewRecorder()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	CORS()(testHandler).ServeHTTP(res, req)

	assert.Equal("*", res.Header().Get("Access-Control-ALlow-Origin"))
	assert.Contains(res.Header().Get("Access-Control-ALlow-Methods"), http.MethodGet)
}

func TestNoCache(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(err)
	res := httptest.NewRecorder()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	NoCache()(testHandler).ServeHTTP(res, req)

	headers := res.Header()

	assert.Equal(headers.Get("Expires"), "0")
	assert.Equal(headers.Get("Pragma"), "no-cache")
	assert.Equal(headers.Get("X-Accel-Expires"), "0")
	assert.Contains(headers.Get("Cache-Control"), "no-cache")
}
