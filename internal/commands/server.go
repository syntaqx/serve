package commands

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/syntaqx/serve"
	"github.com/syntaqx/serve/internal/config"
	"github.com/syntaqx/serve/internal/middleware"
)

// shutdownTimeout bounds how long in-flight requests are given to drain when a
// termination signal is received.
const shutdownTimeout = 10 * time.Second

// GetAuthUsers parses BasicAuth credentials from r. Each non-empty, non-comment
// line has the form "user:secret", where secret may be a plaintext password or
// a bcrypt hash. Only the first colon is treated as the separator, so secrets
// may themselves contain colons.
func GetAuthUsers(r io.Reader) map[string]string {
	users := make(map[string]string)
	if r == nil {
		return users
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if name, secret, ok := strings.Cut(line, ":"); ok && name != "" {
			users[name] = secret
		}
	}

	return users
}

// loadAuthUsers reads Basic Auth credentials from the given path. A missing
// file is not an error: it simply means authentication is disabled.
func loadAuthUsers(path string, logger *slog.Logger, debug bool) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if debug {
				logger.Debug("users file not found, authentication disabled", "path", path)
			}
			return nil, nil
		}
		return nil, fmt.Errorf("open users file: %w", err)
	}
	defer f.Close()

	return GetAuthUsers(f), nil
}

// Handler builds the fully-wrapped static file handler for the given
// configuration and directory.
func Handler(logger *slog.Logger, cfg *config.Config, dir string) (http.Handler, error) {
	users, err := loadAuthUsers(cfg.UsersFile, logger, cfg.Debug)
	if err != nil {
		return nil, err
	}

	fs := serve.NewFileServer(
		serve.WithDirectory(dir),
		serve.WithHiddenFiles(cfg.ShowHidden),
		serve.WithDirectoryListing(cfg.DirectoryListing),
	)

	// Middleware is applied so the last entry is outermost. Health is outermost
	// so probes skip auth and logging; Logger records every other response
	// (including 401s from Auth and 500s from Recover); Recover wraps the rest
	// so panics anywhere are contained.
	fs.Use(
		middleware.SetContentType,
		middleware.Auth(users),
		middleware.CORS(cfg.CORSOrigin),
		middleware.Recover(logger),
		middleware.Logger(logger),
		middleware.Health(cfg.HealthPath),
	)

	return fs, nil
}

// Server builds the handler and runs an HTTP(S) server until ctx is cancelled,
// at which point it gracefully drains in-flight connections.
func Server(ctx context.Context, logger *slog.Logger, cfg *config.Config, dir string) error {
	handler, err := Handler(logger, cfg, dir)
	if err != nil {
		return err
	}

	addr := net.JoinHostPort(cfg.Host, cfg.Port)
	srv := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       120 * time.Second,
		BaseContext:       func(net.Listener) context.Context { return ctx },
	}

	start := srv.ListenAndServe
	if cfg.EnableSSL {
		logger.Info("https server listening", "addr", addr)
		start = func() error { return srv.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile) }
	} else {
		logger.Info("http server listening", "addr", addr)
	}

	return serveWithShutdown(ctx, logger, srv, start)
}

// serveWithShutdown runs start in the background and returns when either the
// server stops on its own or ctx is cancelled. On cancellation it drains
// in-flight connections within shutdownTimeout. A clean shutdown
// (http.ErrServerClosed) is reported as success.
func serveWithShutdown(ctx context.Context, logger *slog.Logger, srv *http.Server, start func() error) error {
	errCh := make(chan error, 1)
	go func() { errCh <- start() }()

	select {
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		logger.Info("shutdown signal received, draining connections")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	}
}
