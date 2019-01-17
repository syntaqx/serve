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

	server := &http.Server{
		Addr:         net.JoinHostPort(opt.Host, strconv.Itoa(opt.Port)),
		Handler:      fs,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	if opt.EnableSSL {
		log.Printf("https server listening at %s", server.Addr)
		return server.ListenAndServeTLS(opt.CertFile, opt.KeyFile)
	}

	log.Printf("http server listening at %s", server.Addr)
	return server.ListenAndServe()
}
