package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCORS(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	CORS("*")(testHandler).ServeHTTP(res, req)

	assert.Equal("*", res.Header().Get("Access-Control-Allow-Origin"))
	assert.Contains(res.Header().Get("Access-Control-Allow-Methods"), http.MethodGet)
}

func TestCORSCustomOrigin(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	CORS("https://example.com")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(res, req)

	assert.Equal("https://example.com", res.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORSEmptyOriginDefaultsToWildcard(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	CORS("")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(res, req)

	assert.Equal("*", res.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORSPreflight(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	res := httptest.NewRecorder()

	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true })
	CORS("*")(next).ServeHTTP(res, req)

	assert.Equal(http.StatusNoContent, res.Code)
	assert.False(called, "preflight should short-circuit before the next handler")
}
