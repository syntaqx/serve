package commands

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/syntaqx/serve/internal/config"
	"github.com/syntaqx/serve/internal/router"
)

// Server implements the static http server command.
func Server(log *log.Logger, opt config.Flags) error {
	r := router.NewRouter(log, opt)

	server := &http.Server{
		Addr:         net.JoinHostPort(opt.Host, strconv.Itoa(opt.Port)),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Printf("http server listening at %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		return fmt.Errorf("http server closed unexpectedly: %v", err)
	}

	return nil
}
