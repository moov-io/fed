// Copyright 2019 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"github.com/moov-io/base"
	"github.com/moov-io/fed"
	"sync"

	"github.com/go-kit/kit/log"
)

var (
	errNoSearchParams = errors.New("missing search parameter(s)")

	// ToDo: ?
	//softResultsLimit, hardResultsLimit = 10, 100
)

type searcher struct {
	ACHDictionary  *fed.ACHDictionary
	WIREDictionary *fed.WIREDictionary
	sync.RWMutex   // protects all above fields

	// errors holds each error encountered when attempting to parse the file
	errors base.ErrorList
	logger log.Logger
}

type searchResponse struct {
	ACHParticipants  []*fed.ACHParticipant  `json:"achParticipants"`
	WIREParticipants []*fed.WIREParticipant `json:"wireParticipants"`
}

func (s *searcher) ACHFindNameOnly(participantName string) ([]*fed.ACHParticipant, error) {
	s.RLock()
	defer s.RUnlock()
	fi, err := s.ACHDictionary.FinancialInstitutionSearch(participantName)
	if err != nil {
		return nil, err
	}
	return fi, nil
}

func (s *searcher) ACHFindRoutingNumberOnly(routingNumber string) ([]*fed.ACHParticipant, error) {
	s.RLock()
	defer s.RUnlock()
	fi, err := s.ACHDictionary.RoutingNumberSearch(routingNumber)
	if err != nil {
		return nil, err
	}
	return fi, nil
}

func (s *searcher) ACHFindStateOnly(state string) []*fed.ACHParticipant {
	s.RLock()
	defer s.RUnlock()
	fi := s.ACHDictionary.StateFilter(state)
	return fi
}

func (s *searcher) ACHFindCityOnly(city string) []*fed.ACHParticipant {
	s.RLock()
	defer s.RUnlock()
	fi := s.ACHDictionary.CityFilter(city)
	return fi
}

func (s *searcher) ACHFindPostalCodeOnly(postalCode string) []*fed.ACHParticipant {
	s.RLock()
	defer s.RUnlock()
	fi := s.ACHDictionary.PostalCodeFilter(postalCode)
	return fi
}

func (s *searcher) FindFEDACH(req FEDACHRequest) ([]*fed.ACHParticipant, error) {
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
