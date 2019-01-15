// Package main implements the runtime for the serve binary.
// `serve` is a static http server anywhere you need one.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"
)

var version = "0.0.0-develop"

type flags struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	Dir  string `json:"dir"`
}

func main() {
	var opt flags
	flag.StringVar(&opt.Host, "host", "", "host address to bind to")
	flag.IntVar(&opt.Port, "port", 8080, "listening port")
	flag.StringVar(&opt.Dir, "dir", "", "directory to serve")
	flag.Parse()

	log := log.New(os.Stderr, "[serve] ", log.LstdFlags)

	// If an argument was provided, see if it's a command, or use it as opt.Dir
	cmd := flag.Arg(0)

	// If an argument is provided, use it as the root directory.
	if opt.Dir == "" {
		if len(cmd) == 0 {
			cwd, err := os.Getwd()
			if err != nil {
				log.Printf("unable to determine current working directory: %v\n", err)
				os.Exit(1)
			}
			opt.Dir = cwd
		} else {
			opt.Dir = cmd
		}
	}

	var err error

	// Execute the specified command
	switch cmd {
	case "version":
		err = VersionCommand(os.Stderr)
	default:
		err = ServerCommand(log, opt)
	}

	if err != nil {
		log.Printf("cmd error: %v", err)
		os.Exit(1)
	}
}

// VersionCommand implements the command `version` which outputs the current
// binary release version, if any.
func VersionCommand(w io.Writer) error {
	fmt.Fprintf(w, fmt.Sprintf("serve version %s %s/%s\n", version, runtime.GOOS, runtime.GOARCH))
	return nil
}

// ServerCommand implements the static http server command.
func ServerCommand(log *log.Logger, opt flags) error {
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

// NewRouter returns a new http.Handler for routing
func NewRouter(log *log.Logger, opt flags) http.Handler {
	r := http.NewServeMux()

	// Handler, wrapped with middleware
	handler := http.FileServer(http.Dir(opt.Dir))
	handler = Logger(log)(handler)
	handler = CORS()(handler)
	handler = NoCache()(handler)

	r.Handle("/", handler)

	return r
}
