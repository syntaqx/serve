<div align="center">

![img](assets/logo.png)

`serve` is a static http server anywhere you need one.

[homebrew]: https://brew.sh/
[git]:      https://git-scm.com/
[golang]:   https://golang.org/
[releases]: https://github.com/syntaqx/serve/releases
[modules]:  https://github.com/golang/go/wiki/Modules

[![Release](https://img.shields.io/github/release/syntaqx/serve.svg)][releases]
[![Build Status](https://travis-ci.org/syntaqx/serve.svg?branch=master)](https://travis-ci.org/syntaqx/serve)
[![codecov](https://codecov.io/gh/syntaqx/serve/branch/master/graph/badge.svg)](https://codecov.io/gh/syntaqx/serve)
[![Go Report Card](https://goreportcard.com/badge/github.com/syntaqx/serve)](https://goreportcard.com/report/github.com/syntaqx/serve)
[![GoDoc](https://godoc.org/github.com/syntaqx/serve?status.svg)](https://godoc.org/github.com/syntaqx/serve)

<br><br>

</div>

> TL;DR: `python -m SimpleHTTPServer 8080` because who can remember that many
> letters?

## Installation

`serve` can be installed in a handeful of ways:

### Homebrew on macOS

If you are using [Homebrew][] on macOS, you can install `serve` with the
following command:

```sh
brew install syntaqx/tap/serve
```

### Download the prebuilt binary

You can visit [github releases][releases] to download the latest binary release
for your operating system and architecture.

### From source

To build from source, check out the instructions on getting started with
[development](#development).

## Usage

Run a server from the current directory:

```sh
serve
```

Or, specify the directory. Paths can be both relative and absolute:

```sh
serve /var/www # or serve -dir=/var/www
```

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

This will install the `serve` command into your `$GOPATH/bin` directory, which
if properly appended to your `$PATH`, can now be used.

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
	fs := serve.NewFileServer(".")
	log.Fatal(http.ListenAndServe(":8080", fs))
}
```

## License

[MIT]: https://opensource.org/licenses/MIT

`serve` is open source software released under the [MIT license][MIT].
