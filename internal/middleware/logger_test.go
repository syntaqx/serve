package middleware

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		handler    http.HandlerFunc
		wantStatus float64
	}{
		{
			name:       "ok",
			handler:    func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("hi")) },
			wantStatus: http.StatusOK,
		},
		{
			name:       "not found",
			handler:    func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNotFound) },
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "no write defaults to ok",
			handler:    func(http.ResponseWriter, *http.Request) {},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			logger := slog.New(slog.NewJSONHandler(&buf, nil))

			req := httptest.NewRequest(http.MethodGet, "/path", nil)
			res := httptest.NewRecorder()

			Logger(logger)(tt.handler).ServeHTTP(res, req)

			var record map[string]any
			require.NoError(t, json.Unmarshal(buf.Bytes(), &record))

			assert.Equal(t, "request", record["msg"])
			assert.Equal(t, http.MethodGet, record["method"])
			assert.Equal(t, "/path", record["path"])
			assert.Equal(t, tt.wantStatus, record["status"])
		})
	}
}
