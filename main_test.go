package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionCommand(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	var b bytes.Buffer
	VersionCommand(&b)

	assert.Contains(b.String(), fmt.Sprintf("version %s", version))
}

func TestNewRouter(t *testing.T) {
	t.Parallel()

	var b bytes.Buffer
	log := log.New(&b, "[test] ", 0)
	opt := flags{Port: 0}

	r := NewRouter(log, opt)

	var _ http.Handler = r
}
