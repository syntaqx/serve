<div align="center">

![img](assets/logo.png)

`serve` is a static http server anywhere you need one.

[homebrew]: https://brew.sh/
[git]:      https://git-scm.com/
[golang]:   https://golang.org/
[releases]: https://github.com/syntaqx/serve/releases
[modules]:  https://github.com/golang/go/wiki/Modules

[![GoDoc](https://godoc.org/github.com/syntaqx/serve?status.svg)](https://godoc.org/github.com/syntaqx/serve)
[![Build Status](https://travis-ci.org/syntaqx/serve.svg?branch=master)](https://travis-ci.org/syntaqx/serve)
[![codecov](https://codecov.io/gh/syntaqx/serve/branch/master/graph/badge.svg)](https://codecov.io/gh/syntaqx/serve)
[![Go Report Card](https://goreportcard.com/badge/github.com/syntaqx/serve)](https://goreportcard.com/report/github.com/syntaqx/serve)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)

[![Pre-Release](https://img.shields.io/github/release-pre/syntaqx/serve.svg)][releases]

<br><br>

</div>

## TL;DR

> It's basically `python -m SimpleHTTPServer 8080` written in Go, because who
> can remember that many letters?

## Installation

`serve` can be installed in a handful of ways:

### Homebrew on macOS

If you are using [Homebrew][] on macOS, you can install `serve` with the
following command:

```sh
brew install syntaqx/tap/serve
```

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

- [Git][git]
- [Go 1.11][golang]+

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
