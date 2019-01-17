package commands

import (
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/syntaqx/serve"
	"github.com/syntaqx/serve/internal/config"
	"github.com/syntaqx/serve/internal/middleware"
)

var getHTTPServerFunc = GetStdHTTPServer

// HTTPServer defines a returnable interface type for http.Server
type HTTPServer interface {
	ListenAndServe() error
	ListenAndServeTLS(certFile, keyFile string) error
}

// GetStdHTTPServer returns a standard net/http.Server configured for a given
// address and handler, and other sane defaults.
func GetStdHTTPServer(addr string, h http.Handler) HTTPServer {
	return &http.Server{
		Addr:         addr,
		Handler:      h,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
}

// Server implements the static http server command.
func Server(log *log.Logger, opt config.Flags, dir string) error {
	fs := serve.NewFileServer(serve.Options{
		Directory: dir,
	})

	fs.Use(
		middleware.Logger(log),
		middleware.Recover(),
		middleware.CORS(),
	)

	addr := net.JoinHostPort(opt.Host, strconv.Itoa(opt.Port))
	server := getHTTPServerFunc(addr, fs)

	if opt.EnableSSL {
		log.Printf("https server listening at %s", addr)
		return server.ListenAndServeTLS(opt.CertFile, opt.KeyFile)
	}

	log.Printf("http server listening at %s", addr)
	return server.ListenAndServe()
}
