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
	ACHDictionary *fed.ACHDictionary
	sync.RWMutex  // protects all above fields

	logger log.Logger
}

type searchResponse struct {
	ACHParticipants []*fed.ACHParticipant `json:"achParticipants"`
}

func (s *searcher) FindACHFinancialInstitution(participantName string) ([]*fed.ACHParticipant, error) {
	s.RLock()
	defer s.RUnlock()
	// ToDo: Routing Number Search
	fi, err := s.ACHDictionary.FinancialInstitutionSearch(participantName)
	if err != nil {
		return nil, err
	}
	return fi, nil
}
