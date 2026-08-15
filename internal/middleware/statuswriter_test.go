package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatusWriterDefaultsToOK(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	sw := &statusWriter{ResponseWriter: rec}

	n, err := sw.Write([]byte("hello"))
	assert.NoError(t, err)
	assert.Equal(t, 5, n)
	assert.Equal(t, http.StatusOK, sw.status)
	assert.Equal(t, 5, sw.bytes)
}

func TestStatusWriterWriteHeaderOnce(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	sw := &statusWriter{ResponseWriter: rec}

	sw.WriteHeader(http.StatusTeapot)
	sw.WriteHeader(http.StatusOK)
	assert.Equal(t, http.StatusTeapot, sw.status)
}

func TestStatusWriterUnwrap(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	sw := &statusWriter{ResponseWriter: rec}
	assert.Same(t, rec, sw.Unwrap())
}
