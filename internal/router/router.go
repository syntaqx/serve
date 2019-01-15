package router

import (
	"log"
	"net/http"

	"github.com/syntaqx/serve/internal/config"
	"github.com/syntaqx/serve/internal/middleware"
)

// NewRouter returns a new http.Handler for routing
func NewRouter(log *log.Logger, opt config.Flags) http.Handler {
	r := http.NewServeMux()

	// Handler, wrapped with middleware
	handler := http.FileServer(http.Dir(opt.Dir))
	handler = middleware.Logger(log)(handler)
	handler = middleware.CORS()(handler)
	handler = middleware.NoCache()(handler)

	r.Handle("/", handler)

	return r
}
