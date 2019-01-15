package router

import (
	"bytes"
	"log"
	"net/http"
	"testing"

	"github.com/syntaqx/serve/internal/config"
)

func TestNewRouter(t *testing.T) {
	t.Parallel()

	var b bytes.Buffer
	log := log.New(&b, "[test] ", 0)
	opt := config.Flags{Port: 0}

	r := NewRouter(log, opt)

	var _ http.Handler = r
}
