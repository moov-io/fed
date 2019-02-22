// Copyright 2019 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package fed

import (
	"bufio"
	"github.com/moov-io/base"
	"github.com/moov-io/fed/pkg/strcmp"
	"io"
	"sort"
	"strings"
	"unicode/utf8"
)

var (
	// WIREJaroWinklerSimilarity is the search similarity percentage for strcmp.JaroWinkler for CustomerName
	// (Financial Institution Name)
	WIREJaroWinklerSimilarity = 0.85
	// WIRELevenshteinSimilarity is the search similarity percentage for strcmp.Levenshtein for CustomerName
	// (Financial Institution Name)
	WIRELevenshteinSimilarity = 0.85
)

// WIREDictionary of Participant records
type WIREDictionary struct {
	// Participants is a list of Participant structs
	WIREParticipants []*WIREParticipant
	//scanner provides a convenient interface for reading data
	scanner *bufio.Scanner
	// line being read
	line string
	// IndexWIRERoutingNumber creates an index of WIREParticipants keyed by WIREParticipant.RoutingNumber
	IndexWIRERoutingNumber map[string]*WIREParticipant
	// IndexWIRECustomerName creates an index of WIREParticipants keyed by WIREParticipant.CustomerName
	IndexWIRECustomerName map[string][]*WIREParticipant
	// errors holds each error encountered when attempting to parse the file
	errors base.ErrorList
	// validator is composed for data validation
	validator
}

// NewWIREDictionary creates a WIREDictionary
func NewWIREDictionary(r io.Reader) *WIREDictionary {
	return &WIREDictionary{
		IndexWIRERoutingNumber: map[string]*WIREParticipant{},
		IndexWIRECustomerName:  map[string][]*WIREParticipant{},
		scanner:                bufio.NewScanner(r),
	}
}

// WIREParticipant holds a FedWIRE dir routing record as defined by Fed WIRE Format
// https://frbservices.org/EPaymentsDirectory/fedwireFormat.html
type WIREParticipant struct {
	// RoutingNumber The institution's routing number
	RoutingNumber string `json:"routingNumber"`
	// TelegraphicName is the short name of financial institution  Wells Fargo
	TelegraphicName string `json:"telegraphicName"`
	// CustomerName (36): FEDERAL RESERVE BANK
	CustomerName string `json:"customerName"`
	// Location is the city and state
	WIRELocation `json:"wireLocation"`
	// FundsTransferStatus designates funds transfer status
	// Y - Eligible
	// N - Ineligible
	FundsTransferStatus string `json:"fundsTransferStatus"`
	// FundsSettlementOnlyStatus designates funds settlement only status
	// S - Settlement-Only
	FundsSettlementOnlyStatus string `json:"fundsSettlementOnlyStatus"`
	// BookEntrySecuritiesTransferStatus designates book entry securities transfer status
	BookEntrySecuritiesTransferStatus string `json:"bookEntrySecuritiesTransferStatus"`
	// Date of last revision: YYYYMMDD, or blank
	Date string `json:"date"`
}

// WIRELocation is the city and state
type WIRELocation struct {
	// City
	City string `json:"city"`
	// State
	State string `json:"state"`
}

// Read parses a single line or multiple lines of FedWIREdir text
func (f *WIREDictionary) Read() error {
	// read through the entire file
	for f.scanner.Scan() {
		f.line = f.scanner.Text()

		if utf8.RuneCountInString(f.line) != WIRELineLength {
			f.errors.Add(NewRecordWrongLengthErr(WIRELineLength, len(f.line)))
			// Return with error if the record length is incorrect as this file is a FED file
			return f.errors
		}

		if err := f.parseWIREParticipant(); err != nil {
			f.errors.Add(err)
			return f.errors
		}
	}
	f.createIndexWIRECustomerName()
	return nil
}

// TODO return a parsing error if the format or file is wrong.
func (f *WIREDictionary) parseWIREParticipant() error {
	p := new(WIREParticipant)

	//RoutingNumber (9): 011000015
	p.RoutingNumber = f.line[:9]
	// TelegraphicName (18): FED
	p.TelegraphicName = strings.Trim(f.line[9:27], " ")
	// CustomerName (36): FEDERAL RESERVE BANK
	p.CustomerName = strings.Trim(f.line[27:63], " ")
	// State (2): GA
	p.State = f.line[63:65]
	// City (25): ATLANTA
	p.City = strings.Trim(f.line[65:90], " ")
	// FundsTransferStatus (1): Y or N
	p.FundsTransferStatus = f.line[90:91]
	// FundsSettlementOnlyStatus (1): " " or S - Settlement-Only
	p.FundsSettlementOnlyStatus = f.line[91:92]
	// BookEntrySecuritiesTransferStatus (1): Y or N
	p.BookEntrySecuritiesTransferStatus = f.line[92:93]
	// Date YYYYMMDD (8): 122415
	p.Date = f.line[93:101]
	f.WIREParticipants = append(f.WIREParticipants, p)
	f.IndexWIRERoutingNumber[p.RoutingNumber] = p
	return nil
}

// createIndexWIRECustomerName creates an index of Financial Institutions keyed by WIREParticipant.CustomerName
func (f *WIREDictionary) createIndexWIRECustomerName() {
	for _, wireP := range f.WIREParticipants {
		f.IndexWIRECustomerName[wireP.CustomerName] = append(f.IndexWIRECustomerName[wireP.CustomerName], wireP)
	}
}

// RoutingNumberSearchSingle returns a FEDWIRE participant based on a WIREParticipant.RoutingNumber.  Routing Number
// validation is only that it exists in IndexParticipant.  Expecting 9 digits, checksum needs to be included.
func (f *WIREDictionary) RoutingNumberSearchSingle(s string) *WIREParticipant {
	if _, ok := f.IndexWIRERoutingNumber[s]; ok {
		return f.IndexWIRERoutingNumber[s]
	}
	return nil
}

// FinancialInstitutionSearchSingle returns a FEDWIRE participant based on a WIREParticipant.CustomerName
func (f *WIREDictionary) FinancialInstitutionSearchSingle(s string) []*WIREParticipant {
	if _, ok := f.IndexWIRECustomerName[s]; ok {
		return f.IndexWIRECustomerName[s]
	}
	return nil
}

// RoutingNumberSearch returns FEDWIRE participants if WIREParticipant.RoutingNumber begins with prefix string s.
// The first 2 digits of the routing number are required.
// Based on https://www.frbservices.org/EPaymentsDirectory/search.html
func (f *WIREDictionary) RoutingNumberSearch(s string) ([]*WIREParticipant, error) {
	s = strings.TrimSpace(s)

	if utf8.RuneCountInString(s) < MinimumRoutingNumberDigits {
		// The first 2 digits (characters) are required
		f.errors.Add(NewRecordWrongLengthErr(2, len(s)))
		return nil, f.errors
	}
	if utf8.RuneCountInString(s) > MaximumRoutingNumberDigits {
		f.errors.Add(NewRecordWrongLengthErr(9, len(s)))
		// Routing Number cannot be greater than 10 digits (characters)
		return nil, f.errors
	}
	if err := f.isNumeric(s); err != nil {
		// Routing Number is not numeric
		f.errors.Add(ErrRoutingNumberNumeric)
		return nil, f.errors
	}

	Participants := make([]*WIREParticipant, 0)

	for _, wireP := range f.WIREParticipants {
		if strings.HasPrefix(wireP.RoutingNumber, s) {
			Participants = append(Participants, wireP)
		}
	}

	return Participants, nil
}

// FinancialInstitutionSearch returns a FEDWIRE participant based on a WIREParticipant.CustomerName
func (f *WIREDictionary) FinancialInstitutionSearch(s string) ([]*WIREParticipant, error) {
	s = strings.ToLower(s)

	// Participants is a subset WIREDictionary.WIREParticipants that match the search based on JaroWinkler similarity
	// and Levenshtein similarity
	Participants := make([]*WIREParticipant, 0)

	// JaroWinkler is a more accurate version of the Jaro algorithm. It works by boosting the
	// score of exact matches at the beginning of the strings. By doing this, Winkler says that
	// typos are less common to happen at the beginning.
	for _, wireP := range f.WIREParticipants {
		if strcmp.JaroWinkler(strings.ToLower(wireP.CustomerName), s) > WIREJaroWinklerSimilarity {
			Participants = append(Participants, wireP)
		}
	}

	// Levenshtein is the "edit distance" between two strings. This is the count of operations
	// (insert, delete, replace) needed for two strings to be equal.
	for _, wireP := range f.WIREParticipants {
		if strcmp.Levenshtein(strings.ToLower(wireP.CustomerName), s) > WIRELevenshteinSimilarity {

			// Only append if the not included in the Participant sub-set
			if len(Participants) != 0 {
				for _, p := range Participants {
					if p.CustomerName == wireP.CustomerName && p.RoutingNumber == wireP.RoutingNumber {
						break
					}
				}
				Participants = append(Participants, wireP)

			} else {
				Participants = append(Participants, wireP)
			}
		}
	}
	// Sort the result
	sort.SliceStable(Participants, func(i, j int) bool { return Participants[i].CustomerName < Participants[j].CustomerName })

	return Participants, nil
}

// WIREParticipantRoutingNumberFilter filters WIREParticipant by Routing Number
func (f *WIREDictionary) WIREParticipantRoutingNumberFilter(wireParticipants []*WIREParticipant, s string) ([]*WIREParticipant, error) {
	s = strings.TrimSpace(s)

	if len(s) < MinimumRoutingNumberDigits {
		// The first 2 digits (characters) are required
		f.errors.Add(NewRecordWrongLengthErr(2, len(s)))
		return nil, f.errors
	}
	nsl := make([]*WIREParticipant, 0)
	for _, wireP := range wireParticipants {
		if strings.HasPrefix(wireP.RoutingNumber, s) {
			nsl = append(nsl, wireP)
		}
	}
	return nsl, nil
}

// WIREParticipantStateFilter filters WIREParticipant by State.
func (f *WIREDictionary) WIREParticipantStateFilter(wireParticipants []*WIREParticipant, s string) []*WIREParticipant {
	nsl := make([]*WIREParticipant, 0)
	for _, wireP := range wireParticipants {
		if strings.EqualFold(wireP.WIRELocation.State, s) {
			nsl = append(nsl, wireP)
		}
	}
	return nsl
}

// WIREParticipantCityFilter filters WIREParticipant by City
func (f *WIREDictionary) WIREParticipantCityFilter(wireParticipants []*WIREParticipant, s string) []*WIREParticipant {
	nsl := make([]*WIREParticipant, 0)
	for _, wireP := range wireParticipants {
		if strings.EqualFold(wireP.WIRELocation.City, s) {
			nsl = append(nsl, wireP)
		}
	}
	return nsl
}

// StateFilter filters WIREDictionary.WIREParticipant by state
func (f *WIREDictionary) StateFilter(s string) []*WIREParticipant {
	nsl := make([]*WIREParticipant, 0)
	for _, wireP := range f.WIREParticipants {
		if strings.EqualFold(wireP.WIRELocation.State, s) {
			nsl = append(nsl, wireP)
		}
	}
	return nsl
}

// CityFilter filters WIREDictionary.WIREParticipant by city
func (f *WIREDictionary) CityFilter(s string) []*WIREParticipant {
	nsl := make([]*WIREParticipant, 0)
	for _, wireP := range f.WIREParticipants {
		if strings.EqualFold(wireP.WIRELocation.City, s) {
			nsl = append(nsl, wireP)
		}
	}
	return nsl
}
