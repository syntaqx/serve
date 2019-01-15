// Package main implements the runtime for the serve binary.
package main

import (
	"flag"
	"log"
	"os"
)

var version = "0.0.0-develop"

type flags struct {
	Host string
	Port int
	Dir  string
}

func main() {
	var opt flags
	flag.StringVar(&opt.Host, "host", "", "host address to bind to")
	flag.IntVar(&opt.Port, "port", 8080, "listening port")
	flag.StringVar(&opt.Dir, "dir", "", "directory to serve")
	flag.Parse()

	log := log.New(os.Stderr, "[serve] ", log.LstdFlags)

	var err error
	switch cmd := flag.Arg(0); cmd {
	case "version":
		err = VersionCommand(os.Stderr)
	default:
		opt.Dir = sanitizeDirFlagArg(opt.Dir, cmd)
		err = ServerCommand(log, opt)
	}

	if err != nil {
		log.Printf("cmd error: %v", err)
		os.Exit(1)
	}
}

func sanitizeDirFlagArg(opt, cmd string) string {
	if opt != "" {
		return opt
	} else if len(cmd) != 0 {
		return cmd
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("unable to determine current working directory: %v\n", err)
		os.Exit(1)
	}

	return cwd
}
