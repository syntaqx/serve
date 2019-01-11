package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCORS(t *testing.T) {
	assert := assert.New(t)

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(err)
	res := httptest.NewRecorder()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	CORS()(testHandler).ServeHTTP(res, req)

	assert.Equal("*", res.Header().Get("Access-Control-ALlow-Origin"))
	assert.Contains(http.MethodGet, res.Header().Get("Access-Control-ALlow-Methods"))
}

func TestNoCache(t *testing.T) {
	assert := assert.New(t)

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(err)
	res := httptest.NewRecorder()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	NoCache()(testHandler).ServeHTTP(res, req)

	assert.Contains("no-cache", res.Header().Get("Cache-Control"))
	assert.Equal("no-cache", res.Header().Get("Pragma"))
	assert.Equal("0", res.Header().Get("Expires"))
}
