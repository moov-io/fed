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

// ToDo:  FED WIRE (write FED ACH tests first)

func addSearchRoutes(logger log.Logger, r *mux.Router, searcher *searcher) {
	r.Methods("GET").Path("/fed/ach/search").HandlerFunc(searchFEDACH(logger, searcher))
}

// fedachSearchRequest contains the properties for fed ach search request
type fedachSearchRequest struct {
	Name          string `json:"name"`
	RoutingNumber string `json:"routingNumber"`
	City          string `json:"city"`
	State         string `json:"state"`
	PostalCode    string `json:"postalCode"`
}

// readFEDACHSearchRequest returns a fedachSearchRequest based on url parameters for fed ach search
func readFEDACHSearchRequest(u *url.URL) fedachSearchRequest {
	return fedachSearchRequest{
		Name:          strings.ToUpper(strings.TrimSpace(u.Query().Get("name"))),
		RoutingNumber: strings.ToUpper(strings.TrimSpace(u.Query().Get("routingNumber"))),
		City:          strings.ToUpper(strings.TrimSpace(u.Query().Get("city"))),
		State:         strings.ToUpper(strings.TrimSpace(u.Query().Get("state"))),
		PostalCode:    strings.ToUpper(strings.TrimSpace(u.Query().Get("postalCode"))),
	}
}

// empty returns true if all of the properties in fedachSearchRequest are empty
func (req fedachSearchRequest) empty() bool {
	return req.Name == "" && req.RoutingNumber == "" && req.City == "" &&
		req.State == "" && req.PostalCode == ""
}

// nameOnly returns true if only Name is not ""
func (req fedachSearchRequest) nameOnly() bool {
	return req.Name != "" && req.RoutingNumber == "" && req.City == "" &&
		req.State == "" && req.PostalCode == ""
}

// routingNumberOnly returns true if only routingNumber is not ""
func (req fedachSearchRequest) routingNumberOnly() bool {
	return req.Name == "" && req.RoutingNumber != "" && req.City == "" &&
		req.State == "" && req.PostalCode == ""
}

// cityOnly returns true if only city is not ""
func (req fedachSearchRequest) cityOnly() bool {
	return req.Name == "" && req.RoutingNumber == "" && req.City != "" &&
		req.State == "" && req.PostalCode == ""
}

// stateOnly returns true if only state is not ""
func (req fedachSearchRequest) stateOnly() bool {
	return req.Name == "" && req.RoutingNumber == "" && req.City == "" &&
		req.State != "" && req.PostalCode == ""
}

// postalCodeOnly returns true if only postal code is not ""
func (req fedachSearchRequest) postalCodeOnly() bool {
	return req.Name == "" && req.RoutingNumber == "" && req.City == "" &&
		req.State == "" && req.PostalCode != ""
}

// searchFEDACH calls search functions based on the fed ach search request url parameters
func searchFEDACH(logger log.Logger, searcher *searcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w = wrapResponseWriter(logger, w, r)

		req := readFEDACHSearchRequest(r.URL)

		if req.empty() {
			moovhttp.Problem(w, errNoSearchParams)
		}

		// Search by Name Only
		if req.nameOnly() {
			if logger != nil {
				logger.Log("searchFEDACH", fmt.Sprintf("searching FED ACH Dictionary by name only %s", req.Name))
			}
			req.searchNameOnly(logger, searcher)(w, r)
			return
		} else if req.routingNumberOnly() {
			if logger != nil {
				logger.Log("searchFEDACH", fmt.Sprintf("searching FED ACH Dictionary by routing number only %s", req.RoutingNumber))
			}
			req.searchRoutingNumberOnly(logger, searcher)(w, r)
			return
		} else if req.stateOnly() {
			if logger != nil {
				logger.Log("searchFEDACH", fmt.Sprintf("searching FED ACH Dictionary by state only %s", req.State))
			}
			req.searchStateOnly(logger, searcher)(w, r)
			return
		} else if req.cityOnly() {
			if logger != nil {
				logger.Log("searchFEDACH", fmt.Sprintf("searching FED ACH Dictionary by city only %s", req.City))
			}
			req.searchCityOnly(logger, searcher)(w, r)
			return
		} else if req.postalCodeOnly() {
			if logger != nil {
				logger.Log("searchFEDACH", fmt.Sprintf("searching FED ACH Dictionary by postal code only %s", req.PostalCode))
			}
			req.searchPostalCodeOnly(logger, searcher)(w, r)
			return
		} else {
			if logger != nil {
				logger.Log("searchFEDACH", fmt.Sprintf("searching FED ACH Dictionary by parameters %v", req.RoutingNumber))
			}
			req.searchACH(logger, searcher)(w, r)
			return
		}

	}
}

// searchNameOnly searches FEDACH by name only
func (req fedachSearchRequest) searchNameOnly(logger log.Logger, searcher *searcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if logger != nil {
			logger.Log("searchFEDACH", fmt.Sprintf("search by name %s", req.Name))
		}

		achP, err := searcher.ACHFindNameOnly(req.Name)
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

// searchRoutingNumberOnly searches FEDACH by routing number only
func (req fedachSearchRequest) searchRoutingNumberOnly(logger log.Logger, searcher *searcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if logger != nil {
			logger.Log("searchFEDACH", fmt.Sprintf("search by routing number %s", req.RoutingNumber))
		}

		achP, err := searcher.ACHFindRoutingNumberOnly(req.RoutingNumber)
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

// searchStateOnly searches FEDACH by state only
func (req fedachSearchRequest) searchStateOnly(logger log.Logger, searcher *searcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if logger != nil {
			logger.Log("searchFEDACH", fmt.Sprintf("search by state %s", req.State))
		}

		achP := searcher.ACHFindStateOnly(req.RoutingNumber)

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&searchResponse{ACHParticipants: achP}); err != nil {
			moovhttp.Problem(w, err)
			return
		}
	}
}

// searchCityOnly searches FEDACH by city only
func (req fedachSearchRequest) searchCityOnly(logger log.Logger, searcher *searcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if logger != nil {
			logger.Log("searchFEDACH", fmt.Sprintf("search by city %s", req.City))
		}

		achP := searcher.ACHFindCityOnly(req.RoutingNumber)

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&searchResponse{ACHParticipants: achP}); err != nil {
			moovhttp.Problem(w, err)
			return
		}
	}
}

// searchPostalCodeOnly searches FEDACH by postal code only
func (req fedachSearchRequest) searchPostalCodeOnly(logger log.Logger, searcher *searcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if logger != nil {
			logger.Log("searchFEDACH", fmt.Sprintf("search by city %s", req.PostalCode))
		}

		achP := searcher.ACHFindPostalCodeOnly(req.RoutingNumber)

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&searchResponse{ACHParticipants: achP}); err != nil {
			moovhttp.Problem(w, err)
			return
		}
	}
}

// searchACH searches FEDACH by more than one url parameter
func (req fedachSearchRequest) searchACH(logger log.Logger, searcher *searcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		achP, err := searcher.FindFEDACH(req)
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
