// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/moov-io/base/log"

	"github.com/go-kit/kit/metrics"
)

// CORSAllowedOriginsEnv is a comma-separated list of exact Origins permitted for
// credentialed CORS. Example: "https://moov.io,https://dashboard.moov.io".
// Localhost (http://localhost:<port>) remains allowed for development.
//
// Temporary FED-local allowlist until moov-io/base#509 lands and is bumped here.
const CORSAllowedOriginsEnv = "MOOV_CORS_ALLOW_ORIGINS"

var (
	corsAllowlistOnce sync.Once
	corsAllowlist     map[string]struct{}
)

func loadCORSAllowlist() map[string]struct{} {
	corsAllowlistOnce.Do(func() {
		corsAllowlist = make(map[string]struct{})
		for _, part := range strings.Split(os.Getenv(CORSAllowedOriginsEnv), ",") {
			origin := strings.TrimSpace(part)
			if origin != "" {
				corsAllowlist[origin] = struct{}{}
			}
		}
	})
	return corsAllowlist
}

func resetCORSAllowlistForTest() {
	corsAllowlistOnce = sync.Once{}
	corsAllowlist = nil
}

func originAllowedForCORS(origin string) bool {
	if origin == "" {
		return false
	}
	if strings.HasPrefix(origin, "http://localhost:") {
		return true
	}
	_, ok := loadCORSAllowlist()[origin]
	return ok
}

func setAccessControlAllowHeaders(w http.ResponseWriter, origin string) {
	if !originAllowedForCORS(origin) {
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PATCH,DELETE,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Cookie,X-User-Id,X-Request-Id,Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

func addCORSHandler(r *mux.Router) {
	r.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		origin := req.Header.Get("Origin")
		if origin == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		setAccessControlAllowHeaders(w, origin)
		w.WriteHeader(http.StatusOK)
	})
}

// responseWriter mirrors moov-io/base Wrap metrics but uses the local CORS
// allowlist instead of reflecting arbitrary https Origins.
type responseWriter struct {
	http.ResponseWriter
	start          time.Time
	request        *http.Request
	metric         metrics.Histogram
	headersWritten bool
	log            log.Logger
}

func (w *responseWriter) WriteHeader(code int) {
	if w == nil || w.headersWritten {
		return
	}
	w.headersWritten = true
	setAccessControlAllowHeaders(w, w.request.Header.Get("Origin"))
	defer w.ResponseWriter.WriteHeader(code)

	diff := time.Since(w.start)
	if w.metric != nil {
		w.metric.Observe(diff.Seconds())
	}
	if w.ResponseWriter.Header().Get("Content-Type") == "" {
		w.ResponseWriter.Header().Set("Content-Type", "text/plain")
		w.ResponseWriter.Header().Set("X-Content-Type-Options", "nosniff")
	}
}

func wrapResponseWriterAllowlist(logger log.Logger, m metrics.Histogram, w http.ResponseWriter, r *http.Request) http.ResponseWriter {
	return &responseWriter{
		ResponseWriter: w,
		start:          time.Now(),
		request:        r,
		metric:         m,
		log:            logger,
	}
}
