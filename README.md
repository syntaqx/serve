# <img src="https://raw.githubusercontent.com/syntaqx/serve/main/docs/logo.svg?sanitize=true" width="250">

`serve` is a static http server anywhere you need one.

[git]:        https://git-scm.com/
[golang]:     https://golang.org/
[releases]:   https://github.com/syntaqx/serve/releases
[docker-hub]: https://hub.docker.com/r/syntaqx/serve

[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)

[![codecov](https://codecov.io/gh/syntaqx/serve/graph/badge.svg?token=FGkU1ntp8z)](https://codecov.io/gh/syntaqx/serve)
[![Go Reference](https://pkg.go.dev/badge/github.com/syntaqx/serve.svg)](https://pkg.go.dev/github.com/syntaqx/serve)

[![GitHub Release](https://img.shields.io/github/release-pre/syntaqx/serve.svg)][releases]
[![Docker Pulls](https://img.shields.io/docker/pulls/syntaqx/serve.svg)][docker-hub]

## TL;DR

> It's basically `python -m SimpleHTTPServer 8080` written in Go, because who
> can remember that many letters?

### Features

* HTTPS (TLS)
* Configurable CORS support
* Structured request logging (text or JSON)
* Graceful shutdown on `SIGINT`/`SIGTERM`
* Health-check endpoint at `/healthz` for load balancers and container probes
* Dotfiles (`.env`, `.git`, ...) hidden by default
* Optional directory-listing control via `--dirlisting`
* `net/http` compatible
* [BasicAuth](https://en.wikipedia.org/wiki/Basic_access_authentication) via a users file (plaintext or bcrypt)
* Configurable via flags or `SERVE_*` environment variables

## Installation

`serve` can be installed in a handful of ways:

### Docker

The official [syntaqx/serve][docker-hub] image is available on Docker Hub.

To get started, try hosting a directory from your docker host:

```sh
docker run -v .:/var/www:ro -d syntaqx/serve
```

Alternatively, a simple `Dockerfile` can be used to generate a new image that
includes the necessary content:

```dockerfile
FROM syntaqx/serve
COPY . /var/www
```

Place this in the same directory as your content, then `build` and `run` the
container:

```sh
docker build -t some-content-serve .
docker run --name some-serve -d some-content-serve
```

#### Exposing an external port

```sh
docker run --name some-serve -d -p 8080:8080 some-content-serve
```

Then you can navigate to http://localhost:8080/ or http://host-ip:8080/ in your
browser.

#### Using environment variables for configuration

[12-factor-config]: https://12factor.net/config

Every flag can also be set from an environment variable, following the
expectations of a [12factor app][12-factor-config]. Flags take precedence over
environment variables, which take precedence over the built-in defaults.

| Flag | Environment variable | Default |
| --- | --- | --- |
| `--host` | `SERVE_HOST` | `0.0.0.0` |
| `--port` | `SERVE_PORT` (or `PORT`) | `8080` |
| `--dir` | `SERVE_DIR` | current directory |
| `--ssl` | `SERVE_SSL` | `false` |
| `--cert` | `SERVE_CERT` | `cert.pem` |
| `--key` | `SERVE_KEY` | `key.pem` |
| `--users` | `SERVE_USERS` | `users.dat` |
| `--all` | `SERVE_ALL` | `false` |
| `--dirlisting` | `SERVE_DIRLISTING` | `false` |
| `--health-path` | `SERVE_HEALTH_PATH` | `/healthz` |
| `--cors-origin` | `SERVE_CORS_ORIGIN` | `*` |
| `--log-format` | `SERVE_LOG_FORMAT` | `text` |
| `--log-level` | `SERVE_LOG_LEVEL` | `info` |
| `--debug` | `SERVE_DEBUG` | `false` |

Here's an example using `compose.yml` to configure `serve` to use HTTPS:

```yaml
services:
  web:
    image: syntaqx/serve
    volumes:
      - ./static:/var/www
      - ./fixtures:/etc/ssl
    environment:
      - PORT=1234
    ports:
      - 1234
    command: serve -ssl -cert=/etc/ssl/cert.pem -key=/etc/ssl/key.pem -dir=/var/www
```

The project repository provides an example [compose](./compose.yml) that
implements a variety of common use-cases for `serve`. Feel free to use those to
help you get started.

### Download the binary

Download the [latest release][releases] binary for your system and architecture
and install it into your `$PATH`. If you have the Go toolchain installed, you can
also install directly from source:

```sh
go install github.com/syntaqx/serve/cmd/serve@latest
```

### From source

To build from source, check out the instructions on getting started with
[development](#development).

## Usage

```sh
serve [options] [path]
```

> `[path]` defaults to `.` (relative path to the current directory)

Then simply open your browser to http://localhost:8080 to view your server.

### Options

The following configuration options are available:

* `--host` host address to bind to (defaults to `0.0.0.0` so the server is reachable from other machines and from inside Docker; use `127.0.0.1` to restrict to localhost)
* `--port` listening port (defaults to `8080`, also honors `PORT`)
* `--ssl` enable https (defaults to `false`)
* `--cert` path to the TLS cert file (defaults to `cert.pem`)
* `--key` path to the TLS key file (defaults to `key.pem`)
* `--dir` directory path to serve (defaults to the first argument or the current directory)
* `--users` path to the users file (defaults to `users.dat`); each line is `username:password` or `username:bcrypt-hash`
* `--all` serve dotfiles, which are hidden by default
* `--dirlisting` enable an automatic listing for directories without an `index.html` (disabled by default)
* `--health-path` path for the health-check endpoint (defaults to `/healthz`, empty to disable)
* `--cors-origin` value for the `Access-Control-Allow-Origin` header (defaults to `*`)
* `--log-format` log output format, `text` or `json` (defaults to `text`)
* `--log-level` log level: `debug`, `info`, `warn`, or `error` (defaults to `info`)
* `--debug` enable debug logging
* `--version` print version information and exit

## Development

To develop `serve` or interact with its source code in any meaningful way, be
sure you have the following installed:

### Prerequisites

* [Git][git]
* [Go][golang]

### Install

You can download and install the project from GitHub by simply running:

```sh
git clone git@github.com:syntaqx/serve.git && cd $(basename $_ .git)
make install
```

This will install `serve` into your `$GOPATH/bin` directory, which assuming is
properly appended to your `$PATH`, can now be used:

```sh
$ serve version
serve version v0.8.0 windows/amd64 go1.26.5
```

## Using `serve` manually

Besides running `serve` using the provided binary, you can also embed a
`serve.FileServer` into your own Go program:

```go
package main

import (
    "log"
    "net/http"

    "github.com/syntaqx/serve"
)

func main() {
    fs := serve.NewFileServer(
        serve.WithDirectory("."),
    )
    log.Fatal(http.ListenAndServe(":8080", fs))
}
```

## License

[MIT]: https://opensource.org/licenses/MIT

`serve` is open source software released under the [MIT license][MIT].

As with all Docker images, these likely also contain other software which may be
under other licenses (such as Bash, etc from the base distribution, along with
any direct or indirect dependencies of the primary software being contained).

As for any pre-built image usage, it is the image user's responsibility to
ensure that any use of this image complies with any relevant licenses for all
software contained within.
