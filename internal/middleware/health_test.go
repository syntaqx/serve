package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthOK(t *testing.T) {
	t.Parallel()

	called := false
	next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) { called = true })

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	res := httptest.NewRecorder()
	Health("/healthz")(next).ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, "ok\n", res.Body.String())
	assert.False(t, called, "health check should short-circuit")
}

func TestHealthHead(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodHead, "/healthz", nil)
	res := httptest.NewRecorder()
	Health("/healthz")(http.NotFoundHandler()).ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
}

func TestHealthPassthrough(t *testing.T) {
	t.Parallel()

	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { called = true })

	// Different path falls through.
	req := httptest.NewRequest(http.MethodGet, "/other", nil)
	res := httptest.NewRecorder()
	Health("/healthz")(next).ServeHTTP(res, req)
	assert.True(t, called)

	// Wrong method falls through.
	called = false
	req = httptest.NewRequest(http.MethodPost, "/healthz", nil)
	res = httptest.NewRecorder()
	Health("/healthz")(next).ServeHTTP(res, req)
	assert.True(t, called)
}

func TestHealthDisabled(t *testing.T) {
	t.Parallel()

	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { called = true })

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	res := httptest.NewRecorder()
	Health("")(next).ServeHTTP(res, req)

	assert.True(t, called, "empty path disables the health endpoint")
}
