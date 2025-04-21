// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
	moovhttp "github.com/moov-io/base/http"
	"github.com/moov-io/base/log"
	"github.com/moov-io/fed"
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

// searchFEDACH calls search functions based on the fed ach search request url parameters
func searchFEDACH(logger log.Logger, searcher *searcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if logger == nil {
			logger = log.NewDefaultLogger()
		}

		w = wrapResponseWriter(logger, w, r)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		requestID, userID := moovhttp.GetRequestID(r), moovhttp.GetUserID(r)
		logger = logger.With(log.Fields{
			"requestID": log.String(requestID),
			"userID":    log.String(userID),
		})

		req := readFEDSearchRequest(r.URL)
		if req.empty() {
			logger.Error().Logf("searchFedACH", log.String(errNoSearchParams.Error()))
			moovhttp.Problem(w, errNoSearchParams)
			return
		}

		searchLimit := extractSearchLimit(r)

		var achParticipants []*fed.ACHParticipant
		var err error

		switch {
		case req.nameOnly():
			logger.Logf("searching FED ACH Dictionary by name only %s", req.Name)
			achParticipants = searcher.ACHFindNameOnly(searchLimit, req.Name)

		case req.routingNumberOnly():
			logger.Logf("searching FED ACH Dictionary by routing number only %s", req.RoutingNumber)
			achParticipants, err = searcher.ACHFindRoutingNumberOnly(searchLimit, req.RoutingNumber)
			if err != nil {
				moovhttp.Problem(w, err)
				return
			}

		case req.stateOnly():
			logger.Logf("searching FED ACH Dictionary by state only %s", req.State)
			achParticipants = searcher.ACHFindStateOnly(searchLimit, req.State)

		case req.cityOnly():
			logger.Logf("searching FED ACH Dictionary by city only %s", req.City)
			achParticipants = searcher.ACHFindCityOnly(searchLimit, req.City)

		case req.postalCodeOnly():
			logger.Logf("searching FED ACH Dictionary by postal code only %s", req.PostalCode)
			achParticipants = searcher.ACHFindPostalCodeOnly(searchLimit, req.PostalCode)

		default:
			logger.Logf("searching FED ACH Dictionary by parameters %v", req.RoutingNumber)
			achParticipants, err = searcher.ACHFind(searchLimit, req)
			if err != nil {
				moovhttp.Problem(w, err)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&searchResponse{
			ACHParticipants: achParticipants,
			Stats:           &searcher.achStats,
		})
	}
}

// searchFEDWIRE calls search functions based on the fed wire search request url parameters
func searchFEDWIRE(logger log.Logger, searcher *searcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if logger == nil {
			logger = log.NewDefaultLogger()
		}

		w = wrapResponseWriter(logger, w, r)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		requestID, userID := moovhttp.GetRequestID(r), moovhttp.GetUserID(r)
		logger = logger.With(log.Fields{
			"requestID": log.String(requestID),
			"userID":    log.String(userID),
		})

		req := readFEDSearchRequest(r.URL)
		if req.empty() {
			logger.Error().Logf("searchFEDWIRE: %v", errNoSearchParams)
			moovhttp.Problem(w, errNoSearchParams)
			return
		}

		searchLimit := extractSearchLimit(r)

		var wireParticipants []*fed.WIREParticipant
		var err error

		switch {
		case req.nameOnly():
			logger.Logf("searchFEDWIRE: searching FED WIRE Dictionary by name only %s", req.Name)
			wireParticipants = searcher.WIREFindNameOnly(searchLimit, req.Name)

		case req.routingNumberOnly():
			logger.Logf("searchFEDWIRE: searching FED WIRE Dictionary by routing number only %s", req.RoutingNumber)
			wireParticipants, err = searcher.WIREFindRoutingNumberOnly(searchLimit, req.RoutingNumber)
			if err != nil {
				moovhttp.Problem(w, err)
				return
			}

		case req.stateOnly():
			logger.Logf("searchFEDWIRE: searching FED WIRE Dictionary by state only %s", req.State)
			wireParticipants = searcher.WIREFindStateOnly(searchLimit, req.State)

		case req.cityOnly():
			logger.Logf("searchFEDWIRE: searching FED WIRE Dictionary by city only %s", req.City)
			wireParticipants = searcher.WIREFindCityOnly(searchLimit, req.City)

		default:
			logger.Logf("searchFEDWIRE: searching FED WIRE Dictionary by parameters %v", req.RoutingNumber)
			wireParticipants, err = searcher.WIREFind(searchLimit, req)
			if err != nil {
				moovhttp.Problem(w, err)
			}
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&searchResponse{
			WIREParticipants: wireParticipants,
			Stats:            &searcher.wireStats,
		})
	}
}
