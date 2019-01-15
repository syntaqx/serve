// Package main implements the runtime for the serve binary.
// `serve` is a static http server anywhere you need one.
package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type flags struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func main() {
	var opt flags
	flag.StringVar(&opt.Host, "host", "", "host address to bind to")
	flag.IntVar(&opt.Port, "port", 8080, "listening port")
	flag.Parse()

	log := log.New(os.Stderr, "[serve] ", log.LstdFlags)

	// If an argument is provided, use it as the root directory.
	dir := flag.Arg(0)
	if len(dir) == 0 {
		cwd, err := os.Getwd()
		if err != nil {
			log.Printf("unable to determine current working directory: %v\n", err)
			os.Exit(1)
		}
		dir = cwd
	}

	r := http.NewServeMux()

	// Handler, wrapped with middleware
	handler := http.FileServer(http.Dir(dir))
	handler = Logger(log)(handler)
	handler = CORS()(handler)
	handler = NoCache()(handler)

	r.Handle("/", handler)

	server := &http.Server{
		Addr:         net.JoinHostPort(opt.Host, strconv.Itoa(opt.Port)),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Printf("http server listening at %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("http server closed unexpectedly: %v", err)
	}
}

// Logger is a middleware that logs each request, along with some useful data
// about what was requested, and what the response was.
func Logger(log *log.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				log.Println(r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
			}()
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// CORS sets permissive cross-origin resource sharing rules.
func CORS() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", strings.Join([]string{
				http.MethodHead,
				http.MethodOptions,
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodPatch,
				http.MethodDelete,
			}, ", "))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// NoCache sets a number of HTTP Headers instructing clients not to cache a
// given response, or use an existing cache.
func NoCache() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			headers := w.Header()
			headers.Set("Expires", "0")
			headers.Set("Cache-Control", "no-cache, no-store, no-transform, must-revalidate, private, max-age=0")
			headers.Set("Pragma", "no-cache")
			headers.Set("X-Accel-Expires", "0")
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
