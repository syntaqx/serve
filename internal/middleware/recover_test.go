package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecover(t *testing.T) {

	t.Parallel()
	assert := assert.New(t)

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(err)
	res := httptest.NewRecorder()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test")
	})

	Recover()(testHandler).ServeHTTP(res, req)

	assert.Equal(http.StatusInternalServerError, res.Code)
	assert.Equal(strings.TrimSpace(res.Body.String()), "[PANIC RECOVERED] test")
}
