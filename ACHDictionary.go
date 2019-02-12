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
	// ACHJaroWinklerSimilarity is the search similarity percentage for strcmp.JaroWinkler for CustomerName
	// (Financial Institution Name)
	ACHJaroWinklerSimilarity = 0.85
	// ACHLevenshteinSimilarity is the search similarity percentage for strcmp.Levenshtein for CustomerName
	// (Financial Institution Name)
	ACHLevenshteinSimilarity = 0.85
)

// ACHDictionary of Participant records
type ACHDictionary struct {
	// Participants is a list of Participant structs
	ACHParticipants []*ACHParticipant
	//scanner provides a convenient interface for reading data
	scanner *bufio.Scanner
	// line being read
	line string
	// IndexACHRoutingNumber creates an index of ACHParticipants keyed by ACHParticipant.RoutingNumber
	IndexACHRoutingNumber map[string]*ACHParticipant
	// IndexACHCustomerName creates an index of ACHParticipants keyed by ACHParticipant.CustomerName
	IndexACHCustomerName map[string][]*ACHParticipant
	// errors holds each error encountered when attempting to parse the file
	errors base.ErrorList
	// validator is composed for data validation
	validator
}

// NewACHDictionary creates a ACHDictionary
func NewACHDictionary(r io.Reader) *ACHDictionary {
	return &ACHDictionary{
		IndexACHRoutingNumber: make(map[string]*ACHParticipant),
		IndexACHCustomerName:  make(map[string][]*ACHParticipant),
		scanner:               bufio.NewScanner(r),
	}
}

// ACHParticipant holds a FedACH dir routing record as defined by Fed ACH Format
// https://www.frbservices.org/EPaymentsDirectory/achFormat.html
type ACHParticipant struct {
	// RoutingNumber The institution's routing number
	RoutingNumber string `json:"routingNumber"`
	// OfficeCode Main/Head Office or Branch. O=main B=branch
	OfficeCode string `json:"officeCode"`
	// ServicingFrbNumber Servicing Fed's main office routing number
	ServicingFrbNumber string `json:"servicingFrbNumber"`
	// RecordTypeCode The code indicating the ABA number to be used to route or send ACH items to the RDFI
	// 0 = Institution is a Federal Reserve Bank
	// 1 = Send items to customer routing number
	// 2 = Send items to customer using new routing number field
	RecordTypeCode string `json:"recordTypeCod"`
	// Revised Date of last revision: YYYYMMDD, or blank
	Revised string `json:"revised"`
	// NewRoutingNumber Institution's new routing number resulting from a merger or renumber
	NewRoutingNumber string `json:"newRoutingNumber"`
	// CustomerName (36): FEDERAL RESERVE BANK
	CustomerName string `json:"customerName"`
	// Location is the delivery address
	ACHLocation `json:"achlocation"`
	// PhoneNumber The institution's phone number
	PhoneNumber string `json:"phoneNumber"`
	// StatusCode Code is based on the customers receiver code
	// 1=Receives Gov/Comm
	StatusCode string `json:"statusCode"`
	// ViewCode is current view
	// 1 = Current view
	ViewCode string `json:"viewCode"`
}

// ACHLocation is the institution's delivery address
type ACHLocation struct {
	// Address
	Address string `json:"address"`
	// City
	City string `json:"city"`
	// State
	State string `json:"state"`
	// PostalCode
	PostalCode string `json:"postalCode"`
	// PostalCodeExtension
	PostalCodeExtension string `json:"postalCodeExtension"`
}

// Read parses a single line or multiple lines of FedACHdir text
func (f *ACHDictionary) Read() error {
	// read through the entire file
	for f.scanner.Scan() {
		f.line = f.scanner.Text()

		if utf8.RuneCountInString(f.line) != ACHLineLength {
			f.errors.Add(NewRecordWrongLengthErr(ACHLineLength, len(f.line)))
			// Return with error if the record length is incorrect as this file is a FED file
			return f.errors
		}
		if err := f.parseACHParticipant(); err != nil {
			f.errors.Add(err)
			return f.errors
		}
	}
	f.createIndexACHCustomerName()
	return nil
}

// TODO return a parsing error if the format or file is wrong.
func (f *ACHDictionary) parseACHParticipant() error {
	p := new(ACHParticipant)

	//RoutingNumber (9): 011000015
	p.RoutingNumber = f.line[:9]
	// OfficeCode (1): O
	p.OfficeCode = f.line[9:10]
	// ServicingFrbNumber (9): 011000015
	p.ServicingFrbNumber = f.line[10:19]
	// RecordTypeCode (1): 0
	p.RecordTypeCode = f.line[19:20]
	// ChangeDate (6): 122415
	p.Revised = f.line[20:26]
	// NewRoutingNumber (9): 000000000
	p.NewRoutingNumber = f.line[26:35]
	// CustomerName (36): FEDERAL RESERVE BANK
	p.CustomerName = strings.Trim(f.line[35:71], " ")
	// Address (36): 1000 PEACHTREE ST N.E.
	p.Address = strings.Trim(f.line[71:107], " ")
	// City (20): ATLANTA
	p.City = strings.Trim(f.line[107:127], " ")
	// State (2): GA
	p.State = f.line[127:129]
	// PostalCode (5): 30309
	p.PostalCode = f.line[129:134]
	// PostalCodeExtension (4): 4470
	p.PostalCodeExtension = f.line[134:138]
	// PhoneNumber(10): 8773722457
	p.PhoneNumber = f.line[138:148]
	// StatusCode (1): 1
	p.StatusCode = f.line[148:149]
	// ViewCode (1): 1
	p.ViewCode = f.line[149:150]

	f.ACHParticipants = append(f.ACHParticipants, p)
	f.IndexACHRoutingNumber[p.RoutingNumber] = p
	return nil
}

// createIndexACHCustomerName creates an index of Financial Institutions keyed by ACHParticipant.CustomerName
func (f *ACHDictionary) createIndexACHCustomerName() {
	for _, achP := range f.ACHParticipants {
		f.IndexACHCustomerName[achP.CustomerName] = append(f.IndexACHCustomerName[achP.CustomerName], achP)
	}
}

// CustomerNameLabel returns a formatted string Title for displaying ACHParticipant.CustomerName
func (p *ACHParticipant) CustomerNameLabel() string {
	s := strings.Title(strings.ToLower(p.CustomerName))
	return s
}

// RoutingNumberSearchSingle returns a FEDACH participant based on a ACHParticipant.RoutingNumber.  Routing Number
// validation is only that it exists in IndexParticipant.  Expecting a valid 9 digit routing number.
func (f *ACHDictionary) RoutingNumberSearchSingle(s string) *ACHParticipant {
	if _, ok := f.IndexACHRoutingNumber[s]; ok {
		return f.IndexACHRoutingNumber[s]
	}
	return nil
}

// FinancialInstitutionSearchSingle returns FEDACH participants based on a ACHParticipant.CustomerName
func (f *ACHDictionary) FinancialInstitutionSearchSingle(s string) []*ACHParticipant {
	if _, ok := f.IndexACHCustomerName[s]; ok {
		return f.IndexACHCustomerName[s]
	}
	return nil
}

// RoutingNumberSearch returns FEDACH participants if ACHParticipant.RoutingNumber begins with prefix string s.
// The first 2 digits of the routing number are required.
// Based on https://www.frbservices.org/EPaymentsDirectory/search.html
func (f *ACHDictionary) RoutingNumberSearch(s string) ([]*ACHParticipant, error) {
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

	Participants := make([]*ACHParticipant, 0)

	for _, achP := range f.ACHParticipants {
		if strings.HasPrefix(achP.RoutingNumber, s) {
			Participants = append(Participants, achP)
		}
	}

	return Participants, nil
}

// FinancialInstitutionSearch returns a FEDACH participant based on a ACHParticipant.CustomerName
func (f *ACHDictionary) FinancialInstitutionSearch(s string) ([]*ACHParticipant, error) {
	s = strings.ToLower(s)

	// Participants is a subset ACHDictionary.ACHParticipants that match the search based on JaroWinkler similarity
	// and Levenshtein similarity
	Participants := make([]*ACHParticipant, 0)

	// JaroWinkler is a more accurate version of the Jaro algorithm. It works by boosting the
	// score of exact matches at the beginning of the strings. By doing this, Winkler says that
	// typos are less common to happen at the beginning.
	for _, achP := range f.ACHParticipants {
		if strcmp.JaroWinkler(strings.ToLower(achP.CustomerName), s) > ACHJaroWinklerSimilarity {
			Participants = append(Participants, achP)
		}
	}

	// Levenshtein is the "edit distance" between two strings. This is the count of operations
	// (insert, delete, replace) needed for two strings to be equal.
	for _, achP := range f.ACHParticipants {
		if strcmp.Levenshtein(strings.ToLower(achP.CustomerName), s) > ACHLevenshteinSimilarity {

			// Only append if the not included in the Participant sub-set
			if len(Participants) != 0 {
				for _, p := range Participants {
					if p.CustomerName == achP.CustomerName && p.RoutingNumber == achP.RoutingNumber {
						break
					}
				}
				Participants = append(Participants, achP)

			} else {
				Participants = append(Participants, achP)
			}
		}
	}
	// Sort the result
	sort.SliceStable(Participants, func(i, j int) bool { return Participants[i].CustomerName < Participants[j].CustomerName })

	return Participants, nil
}

// ACHParticipantStateFilter filters ACHParticipant by State.
func (f ACHDictionary) ACHParticipantStateFilter(achParticipants []*ACHParticipant, s string) []*ACHParticipant {
	nsl := make([]*ACHParticipant, 0)
	for _, achP := range achParticipants {
		if strings.EqualFold(achP.ACHLocation.State, s) {
			nsl = append(nsl, achP)
		}
	}
	return nsl
}

// ACHParticipantCityFilter filters ACHParticipant by City
func (f ACHDictionary) ACHParticipantCityFilter(achParticipants []*ACHParticipant, s string) []*ACHParticipant {
	nsl := make([]*ACHParticipant, 0)
	for _, achP := range achParticipants {
		if strings.EqualFold(achP.ACHLocation.City, s) {
			nsl = append(nsl, achP)
		}
	}
	return nsl
}

// ACHParticipantPostalCodeFilter filters ACHParticipant by Postal Code.
func (f ACHDictionary) ACHParticipantPostalCodeFilter(achParticipants []*ACHParticipant, s string) []*ACHParticipant {
	nsl := make([]*ACHParticipant, 0)
	for _, achP := range achParticipants {
		if strings.EqualFold(achP.ACHLocation.PostalCode, s) {
			nsl = append(nsl, achP)
		}
	}
	return nsl
}

// ACHParticipantRoutingNumberFilter filters ACHParticipant by Routing Number
func (f ACHDictionary) ACHParticipantRoutingNumberFilter(achParticipants []*ACHParticipant, s string) ([]*ACHParticipant, error) {
	s = strings.TrimSpace(s)

	if len(s) < MinimumRoutingNumberDigits {
		// The first 2 digits (characters) are required
		f.errors.Add(NewRecordWrongLengthErr(2, len(s)))
		return nil, f.errors
	}
	nsl := make([]*ACHParticipant, 0)
	for _, achP := range achParticipants {
		if strings.HasPrefix(achP.RoutingNumber, s) {
			nsl = append(nsl, achP)
		}
	}
	return nsl, nil
}

// StateFilter filters ACHDictionary.ACHParticipant by state
func (f ACHDictionary) StateFilter(s string) []*ACHParticipant {
	nsl := make([]*ACHParticipant, 0)
	for _, achP := range f.ACHParticipants {
		if strings.EqualFold(achP.ACHLocation.State, s) {

			nsl = append(nsl, achP)
		}
	}
	return nsl
}

// CityFilter filters ACHDictionary.ACHParticipant by city
func (f ACHDictionary) CityFilter(s string) []*ACHParticipant {
	nsl := make([]*ACHParticipant, 0)
	for _, achP := range f.ACHParticipants {
		if strings.EqualFold(achP.ACHLocation.City, s) {
			nsl = append(nsl, achP)
		}
	}
	return nsl
}

// PostalCodeFilter filters ACHParticipant by postal code
func (f ACHDictionary) PostalCodeFilter(s string) []*ACHParticipant {
	nsl := make([]*ACHParticipant, 0)
	for _, achP := range f.ACHParticipants {
		if strings.EqualFold(achP.ACHLocation.PostalCode, s) {
			nsl = append(nsl, achP)
		}
	}
	return nsl
}
