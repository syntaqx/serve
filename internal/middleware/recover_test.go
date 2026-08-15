package middleware

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecover(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	})

	Recover(logger)(testHandler).ServeHTTP(res, req)

	assert.Equal(http.StatusInternalServerError, res.Code)
	assert.Equal("Internal Server Error", strings.TrimSpace(res.Body.String()))
	assert.Contains(buf.String(), "panic recovered")
	assert.NotContains(res.Body.String(), "boom", "panic detail must not leak to the client")
}

func TestRecoverNilLogger(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	})

	Recover(nil)(testHandler).ServeHTTP(res, req)

	assert.Equal(http.StatusInternalServerError, res.Code)
}
