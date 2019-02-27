// Copyright 2019 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	moovhttp "github.com/moov-io/base/http"
)

const (
	ACH  = "ACH"
	WIRE = "WIRE"
)

func addSearchRoutes(logger log.Logger, r *mux.Router, searcher *searcher) {
	r.Methods("GET").Path("/fed/ach/search").HandlerFunc(searchFEDACH(logger, searcher))
	r.Methods("GET").Path("/fed/wire/search").HandlerFunc(searchFEDWIRE(logger, searcher))
}

// fedSearchRequest contains the properties for fed ach search request
type fedSearchRequest struct {
	Name          string `json:"name"`
	RoutingNumber string `json:"routingNumber"`
	City          string `json:"city"`
	State         string `json:"state"`
	PostalCode    string `json:"postalCode"`
}

// readFEDSearchRequest returns a fedachSearchRequest based on url parameters for fed ach search
func readFEDSearchRequest(u *url.URL) fedSearchRequest {
	return fedSearchRequest{
		Name:          strings.ToUpper(strings.TrimSpace(u.Query().Get("name"))),
		RoutingNumber: strings.ToUpper(strings.TrimSpace(u.Query().Get("routingNumber"))),
		City:          strings.ToUpper(strings.TrimSpace(u.Query().Get("city"))),
		State:         strings.ToUpper(strings.TrimSpace(u.Query().Get("state"))),
		PostalCode:    strings.ToUpper(strings.TrimSpace(u.Query().Get("postalCode"))),
	}
}

// empty returns true if all of the properties in fedachSearchRequest are empty
func (req fedSearchRequest) empty() bool {
	return req.Name == "" && req.RoutingNumber == "" && req.City == "" &&
		req.State == "" && req.PostalCode == ""
}

// nameOnly returns true if only Name is not ""
func (req fedSearchRequest) nameOnly() bool {
	return req.Name != "" && req.RoutingNumber == "" && req.City == "" &&
		req.State == "" && req.PostalCode == ""
}

// routingNumberOnly returns true if only routingNumber is not ""
func (req fedSearchRequest) routingNumberOnly() bool {
	return req.Name == "" && req.RoutingNumber != "" && req.City == "" &&
		req.State == "" && req.PostalCode == ""
}

// cityOnly returns true if only city is not ""
func (req fedSearchRequest) cityOnly() bool {
	return req.Name == "" && req.RoutingNumber == "" && req.City != "" &&
		req.State == "" && req.PostalCode == ""
}

// stateOnly returns true if only state is not ""
func (req fedSearchRequest) stateOnly() bool {
	return req.Name == "" && req.RoutingNumber == "" && req.City == "" &&
		req.State != "" && req.PostalCode == ""
}

// postalCodeOnly returns true if only postal code is not ""
func (req fedSearchRequest) postalCodeOnly() bool {
	return req.Name == "" && req.RoutingNumber == "" && req.City == "" &&
		req.State == "" && req.PostalCode != ""
}

// setResponseHeader returns w with Content-Type and StatusOK
func setResponseHeader(w http.ResponseWriter) http.ResponseWriter {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return w
}

// searchFEDACH calls search functions based on the fed ach search request url parameters
func searchFEDACH(logger log.Logger, searcher *searcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w = wrapResponseWriter(logger, w, r)

		req := readFEDSearchRequest(r.URL)

		if req.empty() {
			moovhttp.Problem(w, errNoSearchParams)
		}

		if req.nameOnly() {
			if logger != nil {
				logger.Log("searchFEDACH", fmt.Sprintf("searching FED ACH Dictionary by name only %s", req.Name))
			}
			req.searchNameOnly(logger, searcher, ACH)(w, r)
			return
		} else if req.routingNumberOnly() {
			if logger != nil {
				logger.Log("searchFEDACH", fmt.Sprintf("searching FED ACH Dictionary by routing number only %s", req.RoutingNumber))
			}
			req.searchRoutingNumberOnly(logger, searcher, ACH)(w, r)
			return
		} else if req.stateOnly() {
			if logger != nil {
				logger.Log("searchFEDACH", fmt.Sprintf("searching FED ACH Dictionary by state only %s", req.State))
			}
			req.searchStateOnly(logger, searcher, ACH)(w, r)
			return
		} else if req.cityOnly() {
			if logger != nil {
				logger.Log("searchFEDACH", fmt.Sprintf("searching FED ACH Dictionary by city only %s", req.City))
			}
			req.searchCityOnly(logger, searcher, ACH)(w, r)
			return
		} else if req.postalCodeOnly() {
			if logger != nil {
				logger.Log("searchFEDACH", fmt.Sprintf("searching FED ACH Dictionary by postal code only %s", req.PostalCode))
			}
			req.searchPostalCodeOnly(logger, searcher, ACH)(w, r)
			return
		} else {
			if logger != nil {
				logger.Log("searchFEDACH", fmt.Sprintf("searching FED ACH Dictionary by parameters %v", req.RoutingNumber))
			}
			req.search(logger, searcher, ACH)(w, r)
			return
		}
	}
}

// searchNameOnly searches FEDACH / FEDWIRE by name only
func (req fedSearchRequest) searchNameOnly(logger log.Logger, searcher *searcher, searchType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if logger != nil {
			logger.Log("searchFEDACH", fmt.Sprintf("search by name %s", req.Name))
		}

		switch searchType {
		case ACH:
			achP, err := searcher.ACHFindNameOnly(extractSearchLimit(r), req.Name)
			if err != nil {
				moovhttp.Problem(w, err)
			}
			w := setResponseHeader(w)
			if err := json.NewEncoder(w).Encode(&searchResponse{ACHParticipants: achP}); err != nil {
				moovhttp.Problem(w, err)
				return
			}
		case WIRE:
			wireP, err := searcher.WIREFindNameOnly(extractSearchLimit(r), req.Name)
			if err != nil {
				moovhttp.Problem(w, err)
			}
			w := setResponseHeader(w)
			if err := json.NewEncoder(w).Encode(&searchResponse{WIREParticipants: wireP}); err != nil {
				moovhttp.Problem(w, err)
				return
			}
		}
	}
}

// searchRoutingNumberOnly searches FEDACH / FEDWIRE by routing number only
func (req fedSearchRequest) searchRoutingNumberOnly(logger log.Logger, searcher *searcher, searchType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if logger != nil {
			logger.Log("searchFEDACH", fmt.Sprintf("search by routing number %s", req.RoutingNumber))
		}
		switch searchType {
		case ACH:
			achP, err := searcher.ACHFindRoutingNumberOnly(extractSearchLimit(r), req.RoutingNumber)
			if err != nil {
				moovhttp.Problem(w, err)
			}
			w := setResponseHeader(w)
			if err := json.NewEncoder(w).Encode(&searchResponse{ACHParticipants: achP}); err != nil {
				moovhttp.Problem(w, err)
				return
			}
		case WIRE:
			wireP, err := searcher.WIREFindRoutingNumberOnly(extractSearchLimit(r), req.RoutingNumber)
			if err != nil {
				moovhttp.Problem(w, err)
			}
			w := setResponseHeader(w)
			if err := json.NewEncoder(w).Encode(&searchResponse{WIREParticipants: wireP}); err != nil {
				moovhttp.Problem(w, err)
				return
			}
		}
	}
}

// searchStateOnly searches FEDACH / FEDWIRE by state only
func (req fedSearchRequest) searchStateOnly(logger log.Logger, searcher *searcher, searchType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if logger != nil {
			logger.Log("searchFEDACH", fmt.Sprintf("search by state %s", req.State))
		}

		switch searchType {
		case ACH:
			achP := searcher.ACHFindStateOnly(extractSearchLimit(r), req.State)
			w := setResponseHeader(w)
			if err := json.NewEncoder(w).Encode(&searchResponse{ACHParticipants: achP}); err != nil {
				moovhttp.Problem(w, err)
				return
			}
		case WIRE:
			wireP := searcher.WIREFindStateOnly(extractSearchLimit(r), req.State)
			w := setResponseHeader(w)
			if err := json.NewEncoder(w).Encode(&searchResponse{WIREParticipants: wireP}); err != nil {
				moovhttp.Problem(w, err)
				return
			}

		}
	}
}

// searchCityOnly searches FEDACH by city only
func (req fedSearchRequest) searchCityOnly(logger log.Logger, searcher *searcher, searchType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if logger != nil {
			logger.Log("searchFEDACH", fmt.Sprintf("search by city %s", req.City))
		}

		switch searchType {
		case ACH:
			achP := searcher.ACHFindCityOnly(extractSearchLimit(r), req.City)
			w := setResponseHeader(w)
			if err := json.NewEncoder(w).Encode(&searchResponse{ACHParticipants: achP}); err != nil {
				moovhttp.Problem(w, err)
				return
			}
		case WIRE:
			wireP := searcher.WIREFindCityOnly(extractSearchLimit(r), req.City)
			w := setResponseHeader(w)
			if err := json.NewEncoder(w).Encode(&searchResponse{WIREParticipants: wireP}); err != nil {
				moovhttp.Problem(w, err)
				return
			}
		}
	}
}

// searchPostalCodeOnly searches FEDACH by postal code only
func (req fedSearchRequest) searchPostalCodeOnly(logger log.Logger, searcher *searcher, searchType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if logger != nil {
			logger.Log("searchFEDACH", fmt.Sprintf("search by city %s", req.PostalCode))
		}
		switch searchType {
		case ACH:
			achP := searcher.ACHFindPostalCodeOnly(extractSearchLimit(r), req.PostalCode)
			w := setResponseHeader(w)
			if err := json.NewEncoder(w).Encode(&searchResponse{ACHParticipants: achP}); err != nil {
				moovhttp.Problem(w, err)
				return
			}
		}
	}
}

// search searches FEDACH / FEDWIRE by more than one url parameter
func (req fedSearchRequest) search(logger log.Logger, searcher *searcher, searchType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch searchType {
		case ACH:
			achP, err := searcher.ACHFind(extractSearchLimit(r), req)
			if err != nil {
				moovhttp.Problem(w, err)
			}
			w = setResponseHeader(w)
			if err := json.NewEncoder(w).Encode(&searchResponse{ACHParticipants: achP}); err != nil {
				moovhttp.Problem(w, err)
				return
			}
		case WIRE:
			wireP, err := searcher.WIREFind(extractSearchLimit(r), req)
			if err != nil {
				moovhttp.Problem(w, err)
			}
			w = setResponseHeader(w)
			if err := json.NewEncoder(w).Encode(&searchResponse{WIREParticipants: wireP}); err != nil {
				moovhttp.Problem(w, err)
				return
			}

		}

	}
}

// searchFEDWIRE calls search functions based on the fed wire search request url parameters
func searchFEDWIRE(logger log.Logger, searcher *searcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w = wrapResponseWriter(logger, w, r)

		req := readFEDSearchRequest(r.URL)
		if req.empty() {
			moovhttp.Problem(w, errNoSearchParams)
		}

		if req.nameOnly() {
			if logger != nil {
				logger.Log("searchFEDWIRE", fmt.Sprintf("searching FED WIRE Dictionary by name only %s", req.Name))
			}
			req.searchNameOnly(logger, searcher, WIRE)(w, r)
			return
		} else if req.routingNumberOnly() {
			if logger != nil {
				logger.Log("searchFEDWIRE", fmt.Sprintf("searching FED WIRE Dictionary by routing number only %s", req.RoutingNumber))
			}
			req.searchRoutingNumberOnly(logger, searcher, WIRE)(w, r)
			return
		} else if req.stateOnly() {
			if logger != nil {
				logger.Log("searchFEDWIRE", fmt.Sprintf("searching FED WIRE Dictionary by state only %s", req.State))
			}
			req.searchStateOnly(logger, searcher, WIRE)(w, r)
			return
		} else if req.cityOnly() {
			if logger != nil {
				logger.Log("searchFEDWIRE", fmt.Sprintf("searching FED WIRE Dictionary by city only %s", req.City))
			}
			req.searchCityOnly(logger, searcher, WIRE)(w, r)
			return
		} else {
			if logger != nil {
				logger.Log("searchFEDWIRE", fmt.Sprintf("searching FED WIRE Dictionary by parameters %v", req.RoutingNumber))
			}
			req.search(logger, searcher, WIRE)(w, r)
			return
		}
	}
}
