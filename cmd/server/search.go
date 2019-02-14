// Copyright 2019 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"github.com/moov-io/fed"
	"sync"

	"github.com/go-kit/kit/log"
)

var (
	errNoSearchParams = errors.New("missing search parameter(s)")

	// ToDo: softResultsLimit, hardResultsLimit = 10, 499
)

// searcher defines a searcher struct
type searcher struct {
	ACHDictionary *fed.ACHDictionary
	// ToDo: WIREDictionary *fed.WIREDictionary
	sync.RWMutex // protects all above fields

	logger log.Logger
}

// searchResponse defines a FEDACH search response
type searchResponse struct {
	ACHParticipants []*fed.ACHParticipant `json:"achParticipants"`
	//ToDo: WIREParticipants []*fed.WIREParticipant `json:"wireParticipants"`
}

// ACHFindNameOnly finds ACH Participants by name only
func (s *searcher) ACHFindNameOnly(participantName string) ([]*fed.ACHParticipant, error) {
	s.RLock()
	defer s.RUnlock()
	fi, err := s.ACHDictionary.FinancialInstitutionSearch(participantName)
	if err != nil {
		return nil, err
	}
	return fi, nil
}

// ACHFindRoutingNumberOnly finds ACH Participants by routing number only
func (s *searcher) ACHFindRoutingNumberOnly(routingNumber string) ([]*fed.ACHParticipant, error) {
	s.RLock()
	defer s.RUnlock()
	fi, err := s.ACHDictionary.RoutingNumberSearch(routingNumber)
	if err != nil {
		return nil, err
	}
	return fi, nil
}

// ACHFindCityOnly finds ACH Participants by city only
func (s *searcher) ACHFindCityOnly(city string) []*fed.ACHParticipant {
	s.RLock()
	defer s.RUnlock()
	fi := s.ACHDictionary.CityFilter(city)
	return fi
}

// ACHFindSateOnly finds ACH Participants by state only
func (s *searcher) ACHFindStateOnly(state string) []*fed.ACHParticipant {
	s.RLock()
	defer s.RUnlock()
	fi := s.ACHDictionary.StateFilter(state)
	return fi
}

// ACHFindPostalCodeOnly finds ACH Participants by postal code only
func (s *searcher) ACHFindPostalCodeOnly(postalCode string) []*fed.ACHParticipant {
	s.RLock()
	defer s.RUnlock()
	fi := s.ACHDictionary.PostalCodeFilter(postalCode)
	return fi
}

// FindFEDACH finds ACH Participants based on multiple parameters
func (s *searcher) FindFEDACH(req fedachSearchRequest) ([]*fed.ACHParticipant, error) {
	s.RLock()
	defer s.RUnlock()

	fi, err := s.ACHDictionary.FinancialInstitutionSearch(req.Name)
	if err != nil {
		return nil, err
	}

	if req.RoutingNumber != "" {
		fi, err = s.ACHDictionary.ACHParticipantRoutingNumberFilter(fi, req.RoutingNumber)
		if err != nil {
			return nil, err
		}
	}

	if req.State != "" {
		fi = s.ACHDictionary.ACHParticipantStateFilter(fi, req.State)
	}

	if req.City != "" {
		fi = s.ACHDictionary.ACHParticipantCityFilter(fi, req.City)
	}

	if req.PostalCode != "" {
		fi = s.ACHDictionary.ACHParticipantPostalCodeFilter(fi, req.PostalCode)
	}
	return fi, nil
}
