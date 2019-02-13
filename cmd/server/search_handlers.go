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

func addSearchRoutes(logger log.Logger, r *mux.Router, searcher *searcher) {
	r.Methods("GET").Path("/FEDACH/search").HandlerFunc(searchFEDACH(logger, searcher))
}

type FEDACHRequest struct {
	Name          string `json:"name"`
	RoutingNumber string `json:"routingNumber"`
	City          string `json:"city"`
	State         string `json:"state"`
	PostalCode    string `json:"postalCode"`
}

func readFEDACHRequest(u *url.URL) FEDACHRequest {
	return FEDACHRequest{
		Name:          strings.ToLower(strings.TrimSpace(u.Query().Get("address"))),
		RoutingNumber: strings.ToLower(strings.TrimSpace(u.Query().Get("routingNumber"))),
		City:          strings.ToLower(strings.TrimSpace(u.Query().Get("city"))),
		State:         strings.ToLower(strings.TrimSpace(u.Query().Get("state"))),
		PostalCode:    strings.ToLower(strings.TrimSpace(u.Query().Get("postalCode"))),
	}
}

func (req FEDACHRequest) empty() bool {
	return req.Name == "" && req.RoutingNumber == "" && req.City == "" &&
		req.State == "" && req.PostalCode == ""
}

func (req FEDACHRequest) nameOnly() bool {
	return req.Name != "" && req.RoutingNumber == "" && req.City == "" &&
		req.State == "" && req.PostalCode == ""
}

func (req FEDACHRequest) routingNumberOnly() bool {
	return req.Name == "" && req.RoutingNumber != "" && req.City == "" &&
		req.State == "" && req.PostalCode == ""
}

func (req FEDACHRequest) cityOnly() bool {
	return req.Name == "" && req.RoutingNumber != "" && req.City != "" &&
		req.State == "" && req.PostalCode == ""
}

func (req FEDACHRequest) stateOnly() bool {
	return req.Name == "" && req.RoutingNumber == "" && req.City == "" &&
		req.State != "" && req.PostalCode == ""
}

func (req FEDACHRequest) isPostalCodeOnly() bool {
	return req.Name == "" && req.RoutingNumber == "" && req.City == "" &&
		req.State == "" && req.PostalCode != ""
}

func searchFEDACH(logger log.Logger, searcher *searcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w = wrapResponseWriter(logger, w, r)

		req := readFEDACHRequest(r.URL)

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
		} else if req.isPostalCodeOnly() {
			if logger != nil {
				logger.Log("searchFEDACH", fmt.Sprintf("searching FED ACH Dictionary by postal code only %s", req.PostalCode))
			}
			req.searchPostalCodeOnly(logger, searcher)(w, r)
			return
		} else {
			if logger != nil {
				logger.Log("searchFEDACH", fmt.Sprintf("searching FED ACH Dictionary by parameters %v", req.RoutingNumber))
			}
			req.search(logger, searcher)(w, r)
			return
		}

	}
}

func (req FEDACHRequest) searchNameOnly(logger log.Logger, searcher *searcher) http.HandlerFunc {
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

func (req FEDACHRequest) searchRoutingNumberOnly(logger log.Logger, searcher *searcher) http.HandlerFunc {
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

func (req FEDACHRequest) searchStateOnly(logger log.Logger, searcher *searcher) http.HandlerFunc {
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

func (req FEDACHRequest) searchCityOnly(logger log.Logger, searcher *searcher) http.HandlerFunc {
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

func (req FEDACHRequest) searchPostalCodeOnly(logger log.Logger, searcher *searcher) http.HandlerFunc {
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

func (req FEDACHRequest) search(logger log.Logger, searcher *searcher) http.HandlerFunc {
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
