package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoCache(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(err)
	res := httptest.NewRecorder()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	NoCache()(testHandler).ServeHTTP(res, req)

	headers := res.Header()

	assert.Equal("0", headers.Get("Expires"))
	assert.Equal("no-cache", headers.Get("Pragma"))
	assert.Equal("0", headers.Get("X-Accel-Expires"))
	assert.Contains(headers.Get("Cache-Control"), "no-cache")
}
