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

	// ToDo: ?
	//softResultsLimit, hardResultsLimit = 10, 100
)

type searcher struct {
	ACHDictionary  *fed.ACHDictionary
	WIREDictionary *fed.WIREDictionary
	sync.RWMutex   // protects all above fields

	logger log.Logger
}

type searchResponse struct {
	ACHParticipants  []*fed.ACHParticipant  `json:"achParticipants"`
	WIREParticipants []*fed.WIREParticipant `json:"wireParticipants"`
}

func (s *searcher) FindACHFinancialInstitution(participantName string) ([]*fed.ACHParticipant, error) {
	s.RLock()
	defer s.RUnlock()
	fi, err := s.ACHDictionary.FinancialInstitutionSearch(participantName)
	if err != nil {
		return nil, err
	}
	return fi, nil
}

func (s *searcher) FindACHRoutingNumber(routingNumber string) ([]*fed.ACHParticipant, error) {
	s.RLock()
	defer s.RUnlock()
	fi, err := s.ACHDictionary.RoutingNumberSearch(routingNumber)
	if err != nil {
		return nil, err
	}
	return fi, nil
}

func (s *searcher) FindWIREFinancialInstitution(participantName string) ([]*fed.WIREParticipant, error) {
	s.RLock()
	defer s.RUnlock()
	fi, err := s.WIREDictionary.FinancialInstitutionSearch(participantName)
	if err != nil {
		return nil, err
	}
	return fi, nil
}

func (s *searcher) FindWIRERoutingNumber(routingNumber string) ([]*fed.WIREParticipant, error) {
	s.RLock()
	defer s.RUnlock()
	fi, err := s.WIREDictionary.RoutingNumberSearch(routingNumber)
	if err != nil {
		return nil, err
	}
	return fi, nil
}
