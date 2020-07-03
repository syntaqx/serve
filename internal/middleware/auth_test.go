package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(err)
	res := httptest.NewRecorder()

	// No users
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	Auth(nil)(testHandler).ServeHTTP(res, req)
	assert.Equal("", res.Header().Get("WWW-Authenticate"))

	// Some users
	testUsers := map[string]string{
		"user1": "pass1",
		"user2": "pass2",
	}
	Auth(testUsers)(testHandler).ServeHTTP(res, req)
	assert.Equal("Basic realm=Authenticate", res.Header().Get("WWW-Authenticate"))
	assert.Equal(http.StatusUnauthorized, res.Result().StatusCode)

	// Correct password
	// Recreate new environment
	testHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	req, err = http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(err)
	res = httptest.NewRecorder()

	req.SetBasicAuth("user1", "pass1")
	Auth(testUsers)(testHandler).ServeHTTP(res, req)
	assert.Equal("", res.Header().Get("WWW-Authenticate"))
	assert.Equal(http.StatusOK, res.Result().StatusCode)

	// Incorrect password
	// Recreate new environment
	testHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	req, err = http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(err)
	res = httptest.NewRecorder()

	req.SetBasicAuth("user1", "pass2")
	Auth(testUsers)(testHandler).ServeHTTP(res, req)
	assert.Equal(http.StatusUnauthorized, res.Result().StatusCode)
}
