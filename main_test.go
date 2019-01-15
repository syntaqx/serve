package main

import (
	"bytes"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionCommand(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	var b bytes.Buffer
	VersionCommand(&b)

	assert.Contains(fmt.Sprintf("version %s", version), b.String())
}

func TestServerCommand(t *testing.T) {
	t.Parallel()

	var b bytes.Buffer
	log := log.New(&b, "[test] ", 0)
	opt := flags{Port: 0}

	// @TODO - What do we even do here?
	go ServerCommand(log, opt)
}
