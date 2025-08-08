// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/moov-io/base/admin"
	moovhttp "github.com/moov-io/base/http"
	"github.com/moov-io/base/http/bind"
	"github.com/moov-io/base/log"
	"github.com/moov-io/fed"
	"github.com/moov-io/fed/webui"

	"github.com/gorilla/mux"
)

var (
	httpAddr  = flag.String("http.addr", bind.HTTP("fed"), "HTTP listen address")
	adminAddr = flag.String("admin.addr", bind.Admin("fed"), "Admin HTTP listen address")

	flagLogFormat = flag.String("log.format", "", "Format for log lines (Options: json, plain")
)

func main() {
	flag.Parse()

	var logger log.Logger
	if v := os.Getenv("LOG_FORMAT"); v != "" {
		*flagLogFormat = v
	}
	if strings.ToLower(*flagLogFormat) == "json" {
		logger = log.NewJSONLogger()
	} else {
		logger = log.NewDefaultLogger()
	}
	logger.Info().Logf("Starting fed server version %s", fed.Version)

	// Channel for errors
	errs := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	// Setup business HTTP routes
	router := mux.NewRouter()
	moovhttp.AddCORSHandler(router)
	addPingRoute(router)

	// Start business HTTP server
	readTimeout, _ := time.ParseDuration("30s")
	writTimeout, _ := time.ParseDuration("30s")
	idleTimeout, _ := time.ParseDuration("60s")

	// Check to see if our -http.addr flag has been overridden
	if v := os.Getenv("HTTP_BIND_ADDRESS"); v != "" {
		*httpAddr = v
	}

	serve := &http.Server{
		Addr:    *httpAddr,
		Handler: router,
		TLSConfig: &tls.Config{
			InsecureSkipVerify:       false,
			PreferServerCipherSuites: true,
			MinVersion:               tls.VersionTLS12,
		},
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readTimeout,
		WriteTimeout:      writTimeout,
		IdleTimeout:       idleTimeout,
	}
	shutdownServer := func() {
		if err := serve.Shutdown(context.TODO()); err != nil {
			logger.Logf("shutting down: %v", err)
		}
	}

	// Check to see if our -admin.addr flag has been overridden
	if v := os.Getenv("HTTP_ADMIN_BIND_ADDRESS"); v != "" {
		*adminAddr = v
	}

	// Start Admin server (with Prometheus metrics)
	adminServer, err := admin.New(admin.Opts{
		Addr: *adminAddr,
	})
	if err != nil {
		logger.LogErrorf("problem creating admin server: %v", err)
		os.Exit(1)
	}
	adminServer.AddVersionHandler(fed.Version) // Setup 'GET /version'
	go func() {
		logger.Info().Logf(fmt.Sprintf("listening on %s", adminServer.BindAddr()))
		if err := adminServer.Listen(); err != nil {
			err = fmt.Errorf("problem starting admin http: %v", err)
			logger.Logf("admin: %v", err)
			errs <- err
		}
	}()
	defer adminServer.Shutdown()

	// Start our searcher
	searcher := &searcher{logger: logger}

	fedACHData, err := fedACHDataFile(logger)
	if err != nil {
		logger.LogErrorf("problem downloading FedACH: %v", err)
		os.Exit(1)
	}
	fedWireData, err := fedWireDataFile(logger)
	if err != nil {
		logger.LogErrorf("problem downloading FedWire: %v", err)
		os.Exit(1)
	}

	if err := setupSearcher(logger, searcher, fedACHData, fedWireData); err != nil {
		logger.Logf("read: %v", err)
		os.Exit(1)
	}

	// Add searcher for HTTP routes
	addSearchRoutes(logger, router, searcher)

	// Add webui routes
	webuiController := webui.NewController(logger)
	webuiController.AppendRoutes(router)

	// Start business logic HTTP server
	go func() {
		if certFile, keyFile := os.Getenv("HTTPS_CERT_FILE"), os.Getenv("HTTPS_KEY_FILE"); certFile != "" && keyFile != "" {
			logger.Logf("binding to %s for secure HTTP server", *httpAddr)
			if err := serve.ListenAndServeTLS(certFile, keyFile); err != nil {
				logger.Logf("listen: %v", err)
			}
		} else {
			logger.Logf("binding to %s for HTTP server", *httpAddr)
			if err := serve.ListenAndServe(); err != nil {
				logger.Logf("listen: %v", err)
			}
		}
	}()

	// Block/Wait for an error
	if err := <-errs; err != nil {
		shutdownServer()
		logger.Logf("exit: %v", err)
		os.Exit(1)
	}
}

func addPingRoute(r *mux.Router) {
	r.Methods("GET").Path("/ping").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		moovhttp.SetAccessControlAllowHeaders(w, r.Header.Get("Origin"))
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("PONG"))
	})
}

func setupSearcher(logger log.Logger, s *searcher, achFile, wireFile io.Reader) error {
	if achFile == nil {
		return errors.New("missing fedach data file")
	}
	if wireFile == nil {
		return errors.New("missing fedwire data file")
	}

	if err := s.readFEDACHData(achFile); err != nil {
		return fmt.Errorf("error reading ACH data: %v", err)
	}
	if err := s.readFEDWIREData(wireFile); err != nil {
		return fmt.Errorf("error reading wire data: %v", err)
	}

	return s.precompute()
}
