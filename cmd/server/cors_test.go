// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"net/http/httptest"
	"testing"

	moovhttp "github.com/moov-io/base/http"
)

func withCORSAllowOrigins(t *testing.T, origins ...string) {
	t.Helper()
	moovhttp.SetCORSAllowedOrigins(origins)
	t.Cleanup(moovhttp.ResetCORSAllowlistForTest)
}

func TestCORSRejectsUnlistedOrigin(t *testing.T) {
	withCORSAllowOrigins(t, "https://moov.io")
	w := httptest.NewRecorder()
	moovhttp.SetAccessControlAllowHeaders(w, "https://evil.example")
	if v := w.Header().Get("Access-Control-Allow-Origin"); v != "" {
		t.Errorf("expected no ACAO, got %q", v)
	}
}

func TestCORSAllowsListedOrigin(t *testing.T) {
	withCORSAllowOrigins(t, "https://moov.io")
	w := httptest.NewRecorder()
	moovhttp.SetAccessControlAllowHeaders(w, "https://moov.io")
	if v := w.Header().Get("Access-Control-Allow-Origin"); v != "https://moov.io" {
		t.Errorf("got %q", v)
	}
	if v := w.Header().Get("Access-Control-Allow-Credentials"); v != "true" {
		t.Errorf("got credentials %q", v)
	}
}

func TestCORSAllowsLocalhost(t *testing.T) {
	withCORSAllowOrigins(t)
	w := httptest.NewRecorder()
	moovhttp.SetAccessControlAllowHeaders(w, "http://localhost:3000")
	if v := w.Header().Get("Access-Control-Allow-Origin"); v != "http://localhost:3000" {
		t.Errorf("got %q", v)
	}
}
