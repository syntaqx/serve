package internal

import (
	"log"
	"net/http"
)

// NewRouter returns a new http.Handler for routing
func NewRouter(log *log.Logger, opt Flags) http.Handler {
	r := http.NewServeMux()

	// Handler, wrapped with middleware
	handler := http.FileServer(http.Dir(opt.Dir))
	handler = Logger(log)(handler)
	handler = CORS()(handler)
	handler = NoCache()(handler)

	r.Handle("/", handler)

	return r
}
