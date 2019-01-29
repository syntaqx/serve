package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var logTests = []struct {
	in  func(w http.ResponseWriter, r *http.Request)
	out string
}{
	{
		in: func(w http.ResponseWriter, _ *http.Request) {
			w.Write([]byte{})
		},
		out: "[test] GET / 200",
	},
	{
		in: func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		},
		out: "[test] GET / 404",
	},
}

func TestLogger(t *testing.T) {
	t.Parallel()

	for _, tt := range logTests {
		assert := assert.New(t)

		var b bytes.Buffer
		log := log.New(&b, "[test] ", 0)

		req, err := http.NewRequest(http.MethodGet, "/", nil)
		assert.NoError(err)
		res := httptest.NewRecorder()

		testHandler := http.HandlerFunc(tt.in)
		Logger(log)(testHandler).ServeHTTP(res, req)

		assert.Equal(tt.out, strings.TrimSpace(b.String()))
	}
}
