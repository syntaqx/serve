package commands

import (
	"fmt"
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
	fs := serve.NewFileServer(dir)

	fs.Use(middleware.Logger(log))
	fs.Use(middleware.Recover())
	fs.Use(middleware.CORS())
	fs.Use(middleware.NoCache())

	server := &http.Server{
		Addr:         net.JoinHostPort(opt.Host, strconv.Itoa(opt.Port)),
		Handler:      fs,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Printf("http server listening at %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		return fmt.Errorf("http server closed unexpectedly: %v", err)
	}

	return nil
}
