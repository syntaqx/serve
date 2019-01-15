package main

import (
	"bytes"
	"log"
	"net/http"
	"testing"
)

func TestNewRouter(t *testing.T) {
	t.Parallel()

	var b bytes.Buffer
	log := log.New(&b, "[test] ", 0)
	opt := flags{Port: 0}

	r := NewRouter(log, opt)

	var _ http.Handler = r
}
