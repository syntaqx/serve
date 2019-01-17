// Package main implements the runtime for the serve binary.
package main

import (
	"flag"
	"log"
	"os"

	"github.com/syntaqx/serve/internal/commands"
	"github.com/syntaqx/serve/internal/config"
)

var version = "0.0.0-develop"

func main() {
	var opt config.Flags
	flag.StringVar(&opt.Host, "host", "0.0.0.0", "host address to bind to")
	flag.IntVar(&opt.Port, "port", 8080, "listening port")
	flag.BoolVar(&opt.EnableSSL, "ssl", false, "enable https")
	flag.StringVar(&opt.CertFile, "cert", "cert.pem", "path to the ssl cert file")
	flag.StringVar(&opt.KeyFile, "key", "key.pem", "path to the ssl key file")
	flag.StringVar(&opt.Directory, "dir", "", "directory path to serve")
	flag.Parse()

	log := log.New(os.Stderr, "[serve] ", log.LstdFlags)

	cmd := flag.Arg(0)

	dir, err := config.SanitizeDir(opt.Directory, cmd)
	if err != nil {
		log.Printf("sanitize directory: %v", err)
		os.Exit(1)
	}

	switch cmd {
	case "version":
		err = commands.Version(version, os.Stderr)
	default:
		err = commands.Server(log, opt, dir)
	}

	if err != nil {
		log.Printf("cmd.%s: %v", cmd, err)
		os.Exit(1)
	}
}
