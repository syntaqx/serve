<div align="center">

![img](assets/logo.png)

`serve` is a static http server anywhere you need one.

[homebrew]: https://brew.sh/
[git]:      https://git-scm.com/
[golang]:   https://golang.org/
[releases]: https://github.com/syntaqx/serve/releases

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

```
brew install syntaqx/tap/serve
```

### Download the prebuilt binary

You can visit [github releases][releases] to download the latest binary release
for your operating system and architecture.

### Manually build from source

To manually build from source, check out instructions on getting started with
[development](#development).

## Usage

Run a server from the current directory:

```
serve
```

Or, specify the directory. Paths can be both relative and absolute:

```
serve /var/www
```

## Development

To develop `serve` or interact with its source code in any meaningful way, be
sure you have the following installed:

### Prerequisites

- [Git][git]
- [Go 1.11][golang]+

Additionally, you will need to set `GO111MODULES=on` when interacting with the
project.

### Install

You can download and install the project from GitHub by simply running:

```
git clone git@github.com:syntaqx/serve.git && cd $(basename $_ .git)
go install ./...
```

This will install `serve` into your `$GOPATH/bin` directory, which if appended
to your `$PATH`, can now be used.

```
serve version
```

## Using `serve` manually

Serve also exposes a reusable `FileServer` convenience struct, letting you
easily create your own static file server:

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
