package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func serveAuth(t *testing.T, users map[string]string, setup func(*http.Request)) *http.Response {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if setup != nil {
		setup(req)
	}
	res := httptest.NewRecorder()

	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	Auth(users)(next).ServeHTTP(res, req)
	return res.Result()
}

func TestAuthNoUsers(t *testing.T) {
	t.Parallel()
	res := serveAuth(t, nil, nil)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Empty(t, res.Header.Get("WWW-Authenticate"))
}

func TestAuthChallenge(t *testing.T) {
	t.Parallel()
	res := serveAuth(t, map[string]string{"user": "pass"}, nil)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	assert.Equal(t, `Basic realm="serve"`, res.Header.Get("WWW-Authenticate"))
}

func TestAuthPlaintext(t *testing.T) {
	t.Parallel()
	users := map[string]string{"user": "pass"}

	ok := serveAuth(t, users, func(r *http.Request) { r.SetBasicAuth("user", "pass") })
	assert.Equal(t, http.StatusOK, ok.StatusCode)

	wrong := serveAuth(t, users, func(r *http.Request) { r.SetBasicAuth("user", "nope") })
	assert.Equal(t, http.StatusUnauthorized, wrong.StatusCode)

	unknown := serveAuth(t, users, func(r *http.Request) { r.SetBasicAuth("ghost", "pass") })
	assert.Equal(t, http.StatusUnauthorized, unknown.StatusCode)
}

func TestAuthBcrypt(t *testing.T) {
	t.Parallel()

	hash, err := bcrypt.GenerateFromPassword([]byte("s3cret"), bcrypt.MinCost)
	require.NoError(t, err)
	users := map[string]string{"user": string(hash)}

	ok := serveAuth(t, users, func(r *http.Request) { r.SetBasicAuth("user", "s3cret") })
	assert.Equal(t, http.StatusOK, ok.StatusCode)

	wrong := serveAuth(t, users, func(r *http.Request) { r.SetBasicAuth("user", "wrong") })
	assert.Equal(t, http.StatusUnauthorized, wrong.StatusCode)
}
