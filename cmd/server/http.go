// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/moov-io/base/log"

	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

var (
	routeHistogram = prometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
		Name: "http_response_duration_seconds",
		Help: "Histogram representing the http response durations",
	}, []string{"route"})
)

func wrapResponseWriter(logger log.Logger, w http.ResponseWriter, r *http.Request) http.ResponseWriter {
	route := fmt.Sprintf("%s%s", strings.ToLower(r.Method), strings.Replace(r.URL.Path, "/", "-", -1)) // TODO: filter out random ID's later
	// Use allowlisted CORS wrapper (not moovhttp.Wrap) until moov-io/base#509 is bumped.
	return wrapResponseWriterAllowlist(logger, routeHistogram.With("route", route), w, r)
}
