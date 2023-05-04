package commands

import (
	"bufio"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
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
		Addr:              addr,
		Handler:           h,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		ReadHeaderTimeout: 60 * time.Second,
	}
}

// GetAuthUsers returns a map of users from a given io.Reader
func GetAuthUsers(r io.Reader) map[string]string {
	users := make(map[string]string)

	if r != nil {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			if line := strings.Split(scanner.Text(), ":"); len(line) == 2 { // use only if correct format
				users[line[0]] = line[1]
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatalf("error occurred during reading users file")
		}
	}

	return users
}

// Server implements the static http server command.
func Server(log *log.Logger, opt config.Flags, dir string) error {
	fs := serve.NewFileServer(serve.Options{
		Directory: dir,
	})

	// Authorization
	var f io.Reader
	if _, err := os.Stat(opt.UsersFile); !os.IsNotExist(err) {
		// Config file exists, load data
		f, err = os.Open(opt.UsersFile)
		if err != nil {
			log.Fatalf("unable to open users file %s", opt.UsersFile)
		}
	} else if opt.Debug {
		log.Printf("%s does not exist, authentication skipped", opt.UsersFile)
	}

	fs.Use(
		middleware.Logger(log),
		middleware.Recover(),
		middleware.CORS(),
		middleware.Auth(GetAuthUsers(f)),
	)

	addr := net.JoinHostPort(opt.Host, opt.Port)
	server := getHTTPServerFunc(addr, fs)

	if opt.EnableSSL {
		log.Printf("https server listening at %s", addr)
		return server.ListenAndServeTLS(opt.CertFile, opt.KeyFile)
	}

	log.Printf("http server listening at %s", addr)
	return server.ListenAndServe()
}
