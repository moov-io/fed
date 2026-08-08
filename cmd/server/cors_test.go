// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func withCORSAllowOrigins(t *testing.T, origins string) {
	t.Helper()
	t.Setenv(CORSAllowedOriginsEnv, origins)
	resetCORSAllowlistForTest()
	t.Cleanup(resetCORSAllowlistForTest)
}

func TestCORSRejectsUnlistedOrigin(t *testing.T) {
	withCORSAllowOrigins(t, "https://moov.io")
	w := httptest.NewRecorder()
	setAccessControlAllowHeaders(w, "https://evil.example")
	if v := w.Header().Get("Access-Control-Allow-Origin"); v != "" {
		t.Errorf("expected no ACAO, got %q", v)
	}
}

func TestCORSAllowsListedOrigin(t *testing.T) {
	withCORSAllowOrigins(t, "https://moov.io")
	w := httptest.NewRecorder()
	setAccessControlAllowHeaders(w, "https://moov.io")
	if v := w.Header().Get("Access-Control-Allow-Origin"); v != "https://moov.io" {
		t.Errorf("got %q", v)
	}
	if v := w.Header().Get("Access-Control-Allow-Credentials"); v != "true" {
		t.Errorf("got credentials %q", v)
	}
}

func TestCORSAllowsLocalhost(t *testing.T) {
	withCORSAllowOrigins(t, "")
	w := httptest.NewRecorder()
	setAccessControlAllowHeaders(w, "http://localhost:3000")
	if v := w.Header().Get("Access-Control-Allow-Origin"); v != "http://localhost:3000" {
		t.Errorf("got %q", v)
	}
}

func TestPingCORSUsesAllowlist(t *testing.T) {
	withCORSAllowOrigins(t, "https://moov.io")
	req := httptest.NewRequest("GET", "/ping", nil)
	req.Header.Set("Origin", "https://moov.io")
	w := httptest.NewRecorder()
	setAccessControlAllowHeaders(w, req.Header.Get("Origin"))
	if w.Header().Get("Access-Control-Allow-Origin") != "https://moov.io" {
		t.Fatal("expected listed origin")
	}
	_ = http.StatusOK
}
