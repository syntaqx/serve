# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.8.0] - 2026-08-14

### Added
- Full environment-variable configuration for every flag via `SERVE_*` (with
  `PORT` still honored for Heroku-style deployments).
- Structured logging with `log/slog`, configurable via `--log-format`
  (`text`/`json`) and `--log-level` (`debug`/`info`/`warn`/`error`).
- Health-check endpoint at `/healthz` (configurable with `--health-path`, empty
  to disable); it bypasses auth and logging for probes.
- `--dirlisting` to opt into an automatic file listing for directories without
  an `index.html`; directory listing is disabled by default.
- Dotfiles (e.g. `.env`, `.git`) are now hidden from access and directory
  listings by default; opt in with `--all`.
- BasicAuth users file now supports bcrypt password hashes in addition to
  plaintext, and secrets may contain colons.
- Configurable CORS origin via `--cors-origin`.
- Graceful shutdown: the server drains in-flight connections on `SIGINT`/`SIGTERM`.
- Functional options for the library API (`WithDirectory`, `WithPrefix`,
  `WithHiddenFiles`, `WithDirectoryListing`).
- `--version` flag, richer version output (Go toolchain and VCS revision), and a
  Docker `HEALTHCHECK`.

### Changed
- Requires Go 1.26.
- Request logging now captures every response, including `401` and `500`s.
- Panics return a generic `500` to clients; details are logged server-side only.
- BasicAuth passwords are compared in constant time.
- Static files are served through `io/fs` (`os.DirFS`) for safer path handling.

### Removed
- The internal `mock` HTTP server package and the `HTTPServer` interface
  indirection, replaced by real integration tests.
- The `serve.Options` struct, replaced by functional options.

[Unreleased]: https://github.com/syntaqx/serve/compare/v0.8.0...HEAD
[0.8.0]: https://github.com/syntaqx/serve/compare/v0.7.1...v0.8.0
