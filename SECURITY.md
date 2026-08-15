# Security Policy

## Supported Versions

Only the latest released `0.x` version receives security fixes while the project
is pre-1.0.

## Reporting a Vulnerability

Please report security issues privately via GitHub's
[private vulnerability reporting](https://github.com/syntaqx/serve/security/advisories/new)
rather than opening a public issue. We aim to acknowledge reports within a few
days.

## Security Model

`serve` is a static file server. Keep the following in mind when deploying it:

- **Dotfiles are hidden by default.** Files and directories whose name begins
  with `.` (for example `.env` or `.git`) are not served and do not appear in
  directory listings. Pass `--all` (or `SERVE_ALL=true`) only when you
  explicitly intend to expose them.
- **Path traversal** is prevented by serving through `io/fs` (`os.DirFS`), which
  rejects `..` and rooted paths at the filesystem boundary.
- **BasicAuth** credentials in the users file may be plaintext or bcrypt hashes.
  Prefer bcrypt hashes; plaintext passwords are compared in constant time, but
  they are still stored in the clear on disk. Always serve behind TLS (`--ssl`)
  when using authentication so credentials are not sent over the wire in the
  clear.
- **Panics** are recovered and logged server-side; clients receive a generic
  `500` with no internal detail.

`serve` is intended for development, internal tooling, and serving trusted
static assets. Review your directory contents before exposing it publicly.
