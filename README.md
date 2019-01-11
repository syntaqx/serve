# serve

[homebrew]: https://brew.sh/
[git]:      https://git-scm.com/
[golang]:   https://golang.org/
[releases]: https://github.com/syntaqx/serve/releases

[![Release](https://img.shields.io/github/release/syntaqx/serve.svg)][releases]
[![Build Status](https://travis-ci.org/syntaqx/serve.svg?branch=master)](https://travis-ci.org/syntaqx/serve)
[![codecov](https://codecov.io/gh/syntaqx/serve/branch/master/graph/badge.svg)](https://codecov.io/gh/syntaqx/serve)
[![Go Report Card](https://goreportcard.com/badge/github.com/syntaqx/protokit)](https://goreportcard.com/report/github.com/syntaqx/protokit)
[![GoDoc](https://godoc.org/github.com/syntaqx/serve?status.svg)](https://godoc.org/github.com/syntaqx/serve)

`serve` runs a static http server.

## Installation

`serve` can be installed in a handeful of ways:

### Homebrew on macOS

If you are using [Homebrew][] on macOS, you can install `serve` with the
following command:

```
<coming soon>
```

### Download the prebuilt binary

You can visit [github releases][releases] to download the latest binary release
for your operating system and architecture.

### Manually build from source

To manually build from source, check out the [development][#development]
section.

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
go install .
```

This will install `serve` into your `$GOPATH/bin` directory, which should be
appended to your `$PATH`, and available for immediate use.

## License

[MIT]: https://opensource.org/licenses/MIT

`serve` is open source software released under the [MIT license][MIT].
