// Package main implements the runtime for the serve binary.
package main

import (
	"context"
	"errors"
	"flag"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/syntaqx/serve/internal/commands"
	"github.com/syntaqx/serve/internal/config"
)

var version = "0.0.0-develop"

// Seams for tests.
var (
	exit       = os.Exit
	resolveDir = config.SanitizeDir
)

func main() {
	exit(realMain())
}

// realMain runs the program and returns the process exit code.
func realMain() int {
	if err := run(context.Background(), os.Args[1:], os.Stdout, os.Stderr); err != nil {
		// Help requests are not failures; usage was already printed.
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 1
	}
	return 0
}

func run(ctx context.Context, args []string, stdout, stderr io.Writer) error {
	cfg, positional, err := config.Load(args, stderr)
	if err != nil {
		return err
	}

	logger := newLogger(cfg, stderr)

	if cfg.ShowVersion || (len(positional) > 0 && positional[0] == "version") {
		return commands.Version(version, stdout)
	}

	dir, err := resolveDir(cfg.Directory, positionalDir(positional))
	if err != nil {
		logger.Error("resolve directory", "error", err)
		return err
	}

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := commands.Server(ctx, logger, cfg, dir); err != nil {
		logger.Error("server stopped", "error", err)
		return err
	}

	return nil
}

// positionalDir returns the first positional argument to use as the serve
// directory, ignoring the reserved "version" subcommand.
func positionalDir(args []string) string {
	if len(args) > 0 && args[0] != "version" {
		return args[0]
	}
	return ""
}

func newLogger(cfg *config.Config, w io.Writer) *slog.Logger {
	level := parseLevel(cfg.LogLevel)
	if cfg.Debug {
		level = slog.LevelDebug
	}

	opts := &slog.HandlerOptions{Level: level}

	var handler slog.Handler
	if strings.EqualFold(cfg.LogFormat, "json") {
		handler = slog.NewJSONHandler(w, opts)
	} else {
		handler = slog.NewTextHandler(w, opts)
	}

	return slog.New(handler)
}

func parseLevel(s string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
