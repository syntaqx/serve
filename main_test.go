package main

import (
	"bytes"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestVersionCommand(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	var b bytes.Buffer
	VersionCommand(&b)

	assert.Contains(b.String(), fmt.Sprintf("version %s", version))
}

func TestServerCommand(t *testing.T) {
	t.Skip()

	t.Parallel()
	assert := assert.New(t)

	var b bytes.Buffer
	log := log.New(&b, "[test] ", 0)
	opt := flags{Port: 0}

	go ServerCommand(log, opt)

	// @TODO: Better way of giving the ServerCommand a chance to start?
	time.Sleep(time.Millisecond * 200)
	assert.Contains(b.String(), "http server listening at")
}
