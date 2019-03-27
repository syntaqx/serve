# <img src="https://raw.githubusercontent.com/syntaqx/serve/master/docs/logo.svg?sanitize=true" width="250">

`serve` is a static http server anywhere you need one.

[homebrew]:   https://brew.sh/
[git]:        https://git-scm.com/
[golang]:     https://golang.org/
[releases]:   https://github.com/syntaqx/serve/releases
[modules]:    https://github.com/golang/go/wiki/Modules
[docker-hub]: https://hub.docker.com/r/syntaqx/serve

[![GoDoc](https://godoc.org/github.com/syntaqx/serve?status.svg)](https://godoc.org/github.com/syntaqx/serve)
[![Build Status](https://travis-ci.org/syntaqx/serve.svg?branch=master)](https://travis-ci.org/syntaqx/serve)
[![codecov](https://codecov.io/gh/syntaqx/serve/branch/master/graph/badge.svg)](https://codecov.io/gh/syntaqx/serve)
[![Go Report Card](https://goreportcard.com/badge/github.com/syntaqx/serve)](https://goreportcard.com/report/github.com/syntaqx/serve)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)

[![GitHub Release](https://img.shields.io/github/release-pre/syntaqx/serve.svg)][releases]
[![Docker Cloud Automated build](https://img.shields.io/docker/cloud/automated/syntaqx/serve.svg)][docker-hub]
[![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/syntaqx/serve.svg)][docker-hub]
[![Docker Pulls](https://img.shields.io/docker/pulls/syntaqx/serve.svg)][docker-hub]

[![Deploy to Heroku](https://www.herokucdn.com/deploy/button.svg)](https://dashboard.heroku.com/new?template=https://github.com/syntaqx/serve/tree/master)

## TL;DR

> It's basically `python -m SimpleHTTPServer 8080` written in Go, because who
> can remember that many letters?


### Features

* HTTPS (TLS)
* CORS support
* Request logging
* `net/http` compatible

## Installation

`serve` can be installed in a handful of ways:

### Homebrew on macOS

If you are using [Homebrew][] on macOS, you can install `serve` with the
following command:

```sh
brew install syntaqx/tap/serve
```

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

Currently, `serve` only supports using the `PORT` environment variable for
setting the listening port. All other configurations are available as CLI flags.

> In future releases, most configurations will be settable from both the CLI
> flag as well as a compatible environment variable, aligning with the
> expectations of a [12factor app][12-factor-config]. But, that will require a
> fair amount of work before the functionality is made available.

Here's an example using `docker-compose.yml` to configure `serve` to use HTTPS:

```yaml
version: '3'
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

The project repository provides an example [docker-compose](./docker-compose.yml)
that implements a variety of common use-cases for `serve`. Feel free to use
those to help you get started.

### Download the binary

Quickly download install the latest release:

```sh
curl -sfL https://install.goreleaser.com/github.com/syntaqx/serve.sh | sh
```

Or manually download the [latest release][releases] binary for your system and
architecture and install it into your `$PATH`.

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

* `--host` host address to bind to (defaults to `0.0.0.0`)
* `--port` listening port (defaults to `8080`)
* `--ssl` enable https (defaults to `false`)
* `--cert` path to the ssl cert file (defaults to `cert.pem`)
* `--key` path to the ssl key file (defaults to `key.pem`)
* `--dir` directory path to serve (defaults to `.`, also configurable by `arg[0]`)

## Development

To develop `serve` or interact with its source code in any meaningful way, be
sure you have the following installed:

### Prerequisites

* [Git][git]
* [Go 1.11][golang]+

You will need to activate [Modules][modules] for your version of Go, generally
by invoking `go` with the support `GO111MODULE=on` environment variable set.

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
serve version v0.0.6-8-g5074d63 windows/amd64
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
    fs := serve.NewFileServer()
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
