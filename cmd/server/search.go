// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"net/http"
	"strconv"
	"sync"

	"github.com/moov-io/base/log"
	"github.com/moov-io/fed"
)

var (
	errNoSearchParams                  = errors.New("missing search parameter(s)")
	softResultsLimit, hardResultsLimit = 100, 500
)

// searcher defines a searcher struct
type searcher struct {
	ACHDictionary  *fed.ACHDictionary
	WIREDictionary *fed.WIREDictionary
	sync.RWMutex   // protects all above fields

	logger log.Logger
}

// searchResponse defines a FEDACH search response
type searchResponse struct {
	ACHParticipants  []*fed.ACHParticipant  `json:"achParticipants,omitempty"`
	WIREParticipants []*fed.WIREParticipant `json:"wireParticipants,omitempty"`
}

// ACHFindNameOnly finds ACH Participants by name only
func (s *searcher) ACHFindNameOnly(limit int, participantName string) []*fed.ACHParticipant {
	s.RLock()
	defer s.RUnlock()

	return s.ACHDictionary.FinancialInstitutionSearch(participantName, limit)
}

// ACHFindRoutingNumberOnly finds ACH Participants by routing number only
func (s *searcher) ACHFindRoutingNumberOnly(limit int, routingNumber string) ([]*fed.ACHParticipant, error) {
	s.RLock()
	defer s.RUnlock()

	return s.ACHDictionary.RoutingNumberSearch(routingNumber, limit)
}

// ACHFindCityOnly finds ACH Participants by city only
func (s *searcher) ACHFindCityOnly(limit int, city string) []*fed.ACHParticipant {
	s.RLock()
	defer s.RUnlock()

	return achLimit(s.ACHDictionary.CityFilter(city), limit)
}

// ACHFindSateOnly finds ACH Participants by state only
func (s *searcher) ACHFindStateOnly(limit int, state string) []*fed.ACHParticipant {
	s.RLock()
	defer s.RUnlock()

	return achLimit(s.ACHDictionary.StateFilter(state), limit)
}

// ACHFindPostalCodeOnly finds ACH Participants by postal code only
func (s *searcher) ACHFindPostalCodeOnly(limit int, postalCode string) []*fed.ACHParticipant {
	s.RLock()
	defer s.RUnlock()

	return achLimit(s.ACHDictionary.PostalCodeFilter(postalCode), limit)
}

// ACHFind finds ACH Participants based on multiple parameters
func (s *searcher) ACHFind(limit int, req fedSearchRequest) ([]*fed.ACHParticipant, error) {
	s.RLock()
	defer s.RUnlock()
	var err error

	out := s.ACHDictionary.FinancialInstitutionSearch(req.Name, limit)
	if req.RoutingNumber != "" {
		out, err = s.ACHDictionary.ACHParticipantRoutingNumberFilter(out, req.RoutingNumber)
		if err != nil {
			return nil, err
		}
	}
	if req.State != "" {
		out = s.ACHDictionary.ACHParticipantStateFilter(out, req.State)
	}
	if req.City != "" {
		out = s.ACHDictionary.ACHParticipantCityFilter(out, req.City)
	}
	if req.PostalCode != "" {
		out = s.ACHDictionary.ACHParticipantPostalCodeFilter(out, req.PostalCode)
	}
	return out, nil
}

// WIRE Searches

// WIREFindNameOnly finds WIRE Participants by name only
func (s *searcher) WIREFindNameOnly(limit int, participantName string) []*fed.WIREParticipant {
	s.RLock()
	defer s.RUnlock()
	fi := s.WIREDictionary.FinancialInstitutionSearch(participantName, limit)
	out := wireLimit(fi, limit)
	return out
}

// WIREFindRoutingNumberOnly finds WIRE Participants by routing number only
func (s *searcher) WIREFindRoutingNumberOnly(limit int, routingNumber string) ([]*fed.WIREParticipant, error) {
	s.RLock()
	defer s.RUnlock()
	fi, err := s.WIREDictionary.RoutingNumberSearch(routingNumber, limit)
	if err != nil {
		return nil, err
	}
	out := wireLimit(fi, limit)
	return out, nil
}

// WIREFindCityOnly finds WIRE Participants by city only
func (s *searcher) WIREFindCityOnly(limit int, city string) []*fed.WIREParticipant {
	s.RLock()
	defer s.RUnlock()
	fi := s.WIREDictionary.CityFilter(city)
	out := wireLimit(fi, limit)
	return out
}

// WIREFindSateOnly finds WIRE Participants by state only
func (s *searcher) WIREFindStateOnly(limit int, state string) []*fed.WIREParticipant {
	s.RLock()
	defer s.RUnlock()
	fi := s.WIREDictionary.StateFilter(state)
	out := wireLimit(fi, limit)
	return out
}

// WIRE Find finds WIRE Participants based on multiple parameters
func (s *searcher) WIREFind(limit int, req fedSearchRequest) ([]*fed.WIREParticipant, error) {
	s.RLock()
	defer s.RUnlock()
	var err error
	fi := s.WIREDictionary.FinancialInstitutionSearch(req.Name, limit)

	if req.RoutingNumber != "" {
		fi, err = s.WIREDictionary.WIREParticipantRoutingNumberFilter(fi, req.RoutingNumber)
		if err != nil {
			return nil, err
		}
	}

	if req.State != "" {
		fi = s.WIREDictionary.WIREParticipantStateFilter(fi, req.State)
	}

	if req.City != "" {
		fi = s.WIREDictionary.WIREParticipantCityFilter(fi, req.City)
	}

	out := wireLimit(fi, limit)
	return out, nil
}

// extractSearchLimit extracts the search limit from url query parameters
func extractSearchLimit(r *http.Request) int {
	limit := softResultsLimit
	if v := r.URL.Query().Get("limit"); v != "" {
		n, _ := strconv.Atoi(v)
		if n > 0 {
			limit = n
		}
	}
	if limit > hardResultsLimit {
		limit = hardResultsLimit
	}
	return limit
}

// achLimit returns an FEDACH search result based on the search limit
func achLimit(fi []*fed.ACHParticipant, limit int) []*fed.ACHParticipant {
	var out []*fed.ACHParticipant
	for _, p := range fi {
		if len(out) == limit {
			break
		}
		out = append(out, p)
	}
	return out
}

// wireLimit returns a FEDWIRE search result based on the search limit
func wireLimit(fi []*fed.WIREParticipant, limit int) []*fed.WIREParticipant {
	var out []*fed.WIREParticipant
	for _, p := range fi {
		if len(out) == limit {
			break
		}
		out = append(out, p)
	}
	return out
}
