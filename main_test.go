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
	assert.Contains(res.Header().Get("Access-Control-ALlow-Methods"), http.MethodGet)
}

func TestNoCache(t *testing.T) {
	assert := assert.New(t)

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(err)
	res := httptest.NewRecorder()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	NoCache()(testHandler).ServeHTTP(res, req)

	assert.Contains(res.Header().Get("Cache-Control"), "no-cache")
	assert.Equal("no-cache", res.Header().Get("Pragma"))
	assert.Equal("0", res.Header().Get("Expires"))
}
