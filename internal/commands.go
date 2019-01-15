package internal

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

// VersionCommand implements the command `version` which outputs the current
// binary release version, if any.
func VersionCommand(version string, w io.Writer) error {
	fmt.Fprintf(w, fmt.Sprintf("serve version %s %s/%s\n", version, runtime.GOOS, runtime.GOARCH))
	return nil
}

// ServerCommand implements the static http server command.
func ServerCommand(log *log.Logger, opt Flags) error {
	r := NewRouter(log, opt)

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
