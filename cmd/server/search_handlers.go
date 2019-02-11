// Copyright 2019 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	moovhttp "github.com/moov-io/base/http"
)

func addSearchRoutes(logger log.Logger, r *mux.Router, searcher *searcher) {
	r.Methods("GET").Path("/searchFEDACH").HandlerFunc(searchFEDACH(logger, searcher))
}

func searchFEDACH(logger log.Logger, searcher *searcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w = wrapResponseWriter(logger, w, r)

		// Search by Name
		if name := strings.TrimSpace(r.URL.Query().Get("name")); name != "" {
			if logger != nil {
				logger.Log("searchFEDACH", fmt.Sprintf("searching FED ACH Dictionary by name for %s", name))
			}
			searchByName(logger, searcher, name)(w, r)
			return
		}

		// Search By Routing Number
		if routingNumber := strings.TrimSpace(r.URL.Query().Get("routingNumber")); routingNumber != "" {
			if logger != nil {
				logger.Log("searchFEDACH", fmt.Sprintf("searching FED ACH Dictionary by routing number for %s", routingNumber))
			}
			searchByRoutingNumber(logger, searcher, routingNumber)(w, r)
			return
		}

		// Fallback if no search params were found
		moovhttp.Problem(w, errNoSearchParams)
	}
}

func searchByName(logger log.Logger, searcher *searcher, name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(name) == 0 {
			moovhttp.Problem(w, errNoSearchParams)
			return
		}
		if logger != nil {
			logger.Log("searchFEDACH", fmt.Sprintf("search by name for %s", name))
		}

		achP, err := searcher.FindACHFinancialInstitution(name)
		if err != nil {
			moovhttp.Problem(w, err)
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&searchResponse{ACHParticipants: achP}); err != nil {
			moovhttp.Problem(w, err)
			return
		}
	}
}

func searchByRoutingNumber(logger log.Logger, searcher *searcher, routingNumber string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(routingNumber) == 0 {
			moovhttp.Problem(w, errNoSearchParams)
			return
		}
		if logger != nil {
			logger.Log("searchFEDACH", fmt.Sprintf("search by routing number for %s", routingNumber))
		}

		achP, err := searcher.FindACHRoutingNumber(routingNumber)
		if err != nil {
			moovhttp.Problem(w, err)
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&searchResponse{ACHParticipants: achP}); err != nil {
			moovhttp.Problem(w, err)
			return
		}
	}
}
