package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetContentType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		path     string
		expected string
	}{
		{"/app.js", "application/javascript"},
		{"/app.mjs", "application/javascript"},
		{"/style.css", "text/css; charset=utf-8"},
		{"/index.html", "text/html; charset=utf-8"},
		{"/data.json", "application/json"},
		{"/image.svg", "image/svg+xml"},
		{"/app.wasm", "application/wasm"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			t.Parallel()
			assert := assert.New(t)

			req, err := http.NewRequest(http.MethodGet, tt.path, nil)
			assert.NoError(err)

			res := httptest.NewRecorder()

			inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
			SetContentType(inner).ServeHTTP(res, req)

			assert.Equal(tt.expected, res.Header().Get("Content-Type"))
		})
	}
}

func TestSetContentType_NoExtension(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	req, err := http.NewRequest(http.MethodGet, "/no-extension", nil)
	assert.NoError(err)

	res := httptest.NewRecorder()

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	SetContentType(inner).ServeHTTP(res, req)

	assert.Empty(res.Header().Get("Content-Type"))
}
