package serve

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileServerUse(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(err)
	res := httptest.NewRecorder()

	testMiddleware1 := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("start\n"))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}

	testMiddleware2 := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("end\n"))
		}
		return http.HandlerFunc(fn)
	}

	fs := &FileServer{
		handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fail()
		}),
	}

	fs.Use(testMiddleware2, testMiddleware1)

	fs.ServeHTTP(res, req)

	assert.Equal("start\nend\n", res.Body.String())
}

func TestFileServerServeHTTP(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(err)
	res := httptest.NewRecorder()

	fs := &FileServer{
		handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("expected"))
		}),
	}

	fs.ServeHTTP(res, req)

	assert.Equal("expected", res.Body.String())
}
