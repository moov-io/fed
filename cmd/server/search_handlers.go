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
		// Fallback if no search params were found
		moovhttp.Problem(w, errNoSearchParams)
	}
}

func searchByName(logger log.Logger, searcher *searcher, participantName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(participantName) == 0 {
			moovhttp.Problem(w, errNoSearchParams)
			return
		}
		if logger != nil {
			logger.Log("searchFEDACH", fmt.Sprintf("search by name for %s", participantName))
		}

		achP, err := searcher.FindACHFinancialInstitution(participantName)
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
