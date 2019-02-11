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
	r.Methods("GET").Path("/searchFEDWIRE").HandlerFunc(searchFEDWIRE(logger, searcher))
}

// ToDo:  Expand Searches to include RoutingNumber, State, City, (Zipcode for ACH)

func searchFEDACH(logger log.Logger, searcher *searcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w = wrapResponseWriter(logger, w, r)

		// Search by Name
		if name := strings.TrimSpace(r.URL.Query().Get("name")); name != "" {
			if logger != nil {
				logger.Log("searchFEDACH", fmt.Sprintf("searching FED ACH Dictionary by name for %s", name))
			}
			achSearchByName(logger, searcher, name)(w, r)
			return
		}

		// Search By Routing Number
		if routingNumber := strings.TrimSpace(r.URL.Query().Get("routingNumber")); routingNumber != "" {
			if logger != nil {
				logger.Log("searchFEDACH", fmt.Sprintf("searching FED ACH Dictionary by routing number for %s", routingNumber))
			}
			achSearchByRoutingNumber(logger, searcher, routingNumber)(w, r)
			return
		}

		// Fallback if no search params were found
		moovhttp.Problem(w, errNoSearchParams)
	}
}

func achSearchByName(logger log.Logger, searcher *searcher, name string) http.HandlerFunc {
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

func achSearchByRoutingNumber(logger log.Logger, searcher *searcher, routingNumber string) http.HandlerFunc {
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

func searchFEDWIRE(logger log.Logger, searcher *searcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w = wrapResponseWriter(logger, w, r)

		// Search by Name
		if name := strings.TrimSpace(r.URL.Query().Get("name")); name != "" {
			if logger != nil {
				logger.Log("searchFEDWIRE", fmt.Sprintf("searching FED WIRE Dictionary by name for %s", name))
			}
			wireSearchByName(logger, searcher, name)(w, r)
			return
		}

		// Search By Routing Number
		if routingNumber := strings.TrimSpace(r.URL.Query().Get("routingNumber")); routingNumber != "" {
			if logger != nil {
				logger.Log("searchFEDACH", fmt.Sprintf("searching FED WIRE Dictionary by routing number for %s", routingNumber))
			}
			wireSearchByRoutingNumber(logger, searcher, routingNumber)(w, r)
			return
		}

		// Fallback if no search params were found
		moovhttp.Problem(w, errNoSearchParams)
	}
}

func wireSearchByName(logger log.Logger, searcher *searcher, name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(name) == 0 {
			moovhttp.Problem(w, errNoSearchParams)
			return
		}
		if logger != nil {
			logger.Log("searchFEDWIRE", fmt.Sprintf("search by name for %s", name))
		}

		wireP, err := searcher.FindWIREFinancialInstitution(name)
		if err != nil {
			moovhttp.Problem(w, err)
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&searchResponse{WIREParticipants: wireP}); err != nil {
			moovhttp.Problem(w, err)
			return
		}
	}
}

func wireSearchByRoutingNumber(logger log.Logger, searcher *searcher, routingNumber string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(routingNumber) == 0 {
			moovhttp.Problem(w, errNoSearchParams)
			return
		}
		if logger != nil {
			logger.Log("searchFEDWIRE", fmt.Sprintf("search by routing number for %s", routingNumber))
		}

		wireP, err := searcher.FindWIRERoutingNumber(routingNumber)
		if err != nil {
			moovhttp.Problem(w, err)
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&searchResponse{WIREParticipants: wireP}); err != nil {
			moovhttp.Problem(w, err)
			return
		}
	}
}
