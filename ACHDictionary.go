// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package fed

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/moov-io/base"
	"github.com/moov-io/fed/pkg/strcmp"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
func NewACHDictionary() *ACHDictionary {
	return &ACHDictionary{
		IndexACHRoutingNumber: make(map[string]*ACHParticipant),
		IndexACHCustomerName:  make(map[string][]*ACHParticipant),
	}
}

// ACHParticipant holds a FedACH dir routing record as defined by Fed ACH Format
// https://www.frbservices.org/EPaymentsDirectory/achFormat.html
type ACHParticipant struct {
	// RoutingNumber The institution's routing number
	RoutingNumber string `json:"routingNumber"`
	// OfficeCode Main/Head Office or Branch. O=main B=branch
	OfficeCode string `json:"officeCode"`
	// ServicingFRBNumber Servicing Fed's main office routing number
	ServicingFRBNumber string `json:"servicingFRBNumber"`
	// RecordTypeCode The code indicating the ABA number to be used to route or send ACH items to the RDFI
	// 0 = Institution is a Federal Reserve Bank
	// 1 = Send items to customer routing number
	// 2 = Send items to customer using new routing number field
	RecordTypeCode string `json:"recordTypeCode"`
	// Revised Date of last revision: YYYYMMDD, or blank
	Revised string `json:"revised"`
	// NewRoutingNumber Institution's new routing number resulting from a merger or renumber
	NewRoutingNumber string `json:"newRoutingNumber"`
	// CustomerName (36): FEDERAL RESERVE BANK
	CustomerName string `json:"customerName"`
	// Location is the delivery address
	ACHLocation `json:"achLocation"`
	// PhoneNumber The institution's phone number
	PhoneNumber string `json:"phoneNumber"`
	// StatusCode Code is based on the customers receiver code
	// 1 = Receives Gov/Comm
	StatusCode string `json:"statusCode"`
	// ViewCode is current view
	// 1 = Current view
	ViewCode string `json:"viewCode"`

	// CleanName is our cleaned up value of CustomerName
	CleanName string `json:"cleanName"`
}

type achParticipantResult struct {
	*ACHParticipant

	highestMatch float64
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
func (f *ACHDictionary) Read(r io.Reader) error {
	if f == nil {
		return nil
	}

	bs, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	// Try validating the file as JSON and if that fails read as plaintext
	if json.Valid(bs) {
		err = f.readJSON(bytes.NewReader(bs))
		if err != nil {
			return err
		}
		return nil
	}

	return f.readPlaintext(bytes.NewReader(bs))
}

// jsonACHDictionary is a JSON representation of the FED ACH Participant directory
// Model generated with // https://mholt.github.io/json-to-go/
type jsonACHDictionary struct {
	FedACHParticipants struct {
		Response struct {
			Code int `json:"code"`
		} `json:"response"`
		FedACHParticipants []struct {
			RoutingNumber         string `json:"routingNumber"`
			OfficeCode            string `json:"officeCode"`
			ServicingFRBNumber    string `json:"servicingFRBNumber"`
			RecordTypeCode        string `json:"recordTypeCode"`
			ChangeDate            string `json:"changeDate"`
			NewRoutingNumber      string `json:"newRoutingNumber"`
			CustomerName          string `json:"customerName"`
			CustomerAddress       string `json:"customerAddress"`
			CustomerCity          string `json:"customerCity"`
			CustomerState         string `json:"customerState"`
			CustomerZip           string `json:"customerZip"`
			CustomerZipExt        string `json:"customerZipExt"`
			CustomerAreaCode      string `json:"customerAreaCode"`
			CustomerPhonePrefix   string `json:"customerPhonePrefix"`
			CustomerPhoneSuffix   string `json:"customerPhoneSuffix"`
			InstitutionStatusCode string `json:"institutionStatusCode"`
			DataViewCode          string `json:"dataViewCode"`
		} `json:"fedACHParticipants"`
	} `json:"fedACHParticipants"`
}

func (f *ACHDictionary) readJSON(r io.Reader) error {
	var wrapper jsonACHDictionary
	if err := json.NewDecoder(r).Decode(&wrapper); err != nil {
		return err
	}
	ps := wrapper.FedACHParticipants.FedACHParticipants
	for i := range ps {
		p := &ACHParticipant{
			RoutingNumber:      ps[i].RoutingNumber,
			OfficeCode:         ps[i].OfficeCode,
			ServicingFRBNumber: ps[i].ServicingFRBNumber,
			RecordTypeCode:     ps[i].RecordTypeCode,
			Revised:            ps[i].ChangeDate,
			NewRoutingNumber:   ps[i].NewRoutingNumber,
			CustomerName:       ps[i].CustomerName,
			ACHLocation: ACHLocation{
				Address:             ps[i].CustomerAddress,
				City:                ps[i].CustomerCity,
				State:               ps[i].CustomerState,
				PostalCode:          ps[i].CustomerZip,
				PostalCodeExtension: ps[i].CustomerZipExt,
			},
			PhoneNumber: fmt.Sprintf("%s%s%s", ps[i].CustomerAreaCode, ps[i].CustomerPhonePrefix, ps[i].CustomerPhoneSuffix),
			StatusCode:  ps[i].InstitutionStatusCode,
			ViewCode:    ps[i].DataViewCode,

			// Our Custom Fields
			CleanName: Normalize(ps[i].CustomerName),
		}
		f.IndexACHRoutingNumber[ps[i].RoutingNumber] = p
		f.ACHParticipants = append(f.ACHParticipants, p)
	}
	f.createIndexACHCustomerName()
	return nil
}

func (f *ACHDictionary) readPlaintext(r io.Reader) error {
	if f == nil || r == nil {
		return nil
	}
	// read each line and lift it into a ACHParticipant
	s := bufio.NewScanner(r)
	var line string
	for s.Scan() {
		line = s.Text()

		if utf8.RuneCountInString(line) != ACHLineLength {
			f.errors.Add(NewRecordWrongLengthErr(ACHLineLength, len(line)))
			// Return with error if the record length is incorrect as this file is a FED file
			return f.errors
		}
		if err := f.parseACHParticipant(line); err != nil {
			f.errors.Add(err)
			return f.errors
		}
	}
	f.createIndexACHCustomerName()
	return nil
}

// TODO return a parsing error if the format or file is wrong.
func (f *ACHDictionary) parseACHParticipant(line string) error {
	p := new(ACHParticipant)

	//RoutingNumber (9): 011000015
	p.RoutingNumber = line[:9]
	// OfficeCode (1): O
	p.OfficeCode = line[9:10]
	// ServicingFrbNumber (9): 011000015
	p.ServicingFRBNumber = line[10:19]
	// RecordTypeCode (1): 0
	p.RecordTypeCode = line[19:20]
	// ChangeDate (6): 122415
	p.Revised = line[20:26]
	// NewRoutingNumber (9): 000000000
	p.NewRoutingNumber = line[26:35]
	// CustomerName (36): FEDERAL RESERVE BANK
	p.CustomerName = strings.Trim(line[35:71], " ")
	// Address (36): 1000 PEACHTREE ST N.E.
	p.Address = strings.Trim(line[71:107], " ")
	// City (20): ATLANTA
	p.City = strings.Trim(line[107:127], " ")
	// State (2): GA
	p.State = line[127:129]
	// PostalCode (5): 30309
	p.PostalCode = line[129:134]
	// PostalCodeExtension (4): 4470
	p.PostalCodeExtension = line[134:138]
	// PhoneNumber(10): 8773722457
	p.PhoneNumber = line[138:148]
	// StatusCode (1): 1
	p.StatusCode = line[148:149]
	// ViewCode (1): 1
	p.ViewCode = line[149:150]

	// Our custom fields
	p.CleanName = Normalize(p.CustomerName)

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
	s := cases.Title(language.AmericanEnglish).String(strings.ToLower(p.CustomerName))
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
func (f *ACHDictionary) RoutingNumberSearch(s string, limit int) ([]*ACHParticipant, error) {
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
	exactMatch := len(s) == 9

	out := make([]*achParticipantResult, 0)
	for _, achP := range f.ACHParticipants {
		if exactMatch {
			if achP.RoutingNumber == s {
				out = append(out, &achParticipantResult{
					ACHParticipant: achP,
					highestMatch:   1.0,
				})
			}
		} else {
			out = append(out, &achParticipantResult{
				ACHParticipant: achP,
				highestMatch:   strcmp.JaroWinkler(achP.RoutingNumber, s),
			})
		}
	}
	return reduceACHResults(out, limit), nil
}

// FinancialInstitutionSearch returns a FEDACH participant based on a ACHParticipant.CustomerName
func (f *ACHDictionary) FinancialInstitutionSearch(s string, limit int) []*ACHParticipant {
	s = strings.ToLower(s)

	out := make([]*achParticipantResult, 0)

	for _, achP := range f.ACHParticipants {
		// JaroWinkler is a more accurate version of the Jaro algorithm. It works by boosting the
		// score of exact matches at the beginning of the strings. By doing this, Winkler says that
		// typos are less common to happen at the beginning.
		jaroScore := strcmp.JaroWinkler(strings.ToLower(achP.CleanName), s)

		// Levenshtein is the "edit distance" between two strings. This is the count of operations
		// (insert, delete, replace) needed for two strings to be equal.
		levenScore := strcmp.Levenshtein(strings.ToLower(achP.CleanName), s)

		if jaroScore > ACHJaroWinklerSimilarity || levenScore > ACHLevenshteinSimilarity {
			out = append(out, &achParticipantResult{
				ACHParticipant: achP,
				highestMatch:   math.Max(jaroScore, levenScore),
			})
		}
	}

	return reduceACHResults(out, limit)
}

// ACHParticipantStateFilter filters ACHParticipant by State.
func (f *ACHDictionary) ACHParticipantStateFilter(achParticipants []*ACHParticipant, s string) []*ACHParticipant {
	nsl := make([]*ACHParticipant, 0)
	for _, achP := range achParticipants {
		if strings.EqualFold(achP.ACHLocation.State, s) {
			nsl = append(nsl, achP)
		}
	}
	return nsl
}

// ACHParticipantCityFilter filters ACHParticipant by City
func (f *ACHDictionary) ACHParticipantCityFilter(achParticipants []*ACHParticipant, s string) []*ACHParticipant {
	nsl := make([]*ACHParticipant, 0)
	for _, achP := range achParticipants {
		if strings.EqualFold(achP.ACHLocation.City, s) {
			nsl = append(nsl, achP)
		}
	}
	return nsl
}

// ACHParticipantPostalCodeFilter filters ACHParticipant by Postal Code.
func (f *ACHDictionary) ACHParticipantPostalCodeFilter(achParticipants []*ACHParticipant, s string) []*ACHParticipant {
	nsl := make([]*ACHParticipant, 0)
	for _, achP := range achParticipants {
		if strings.EqualFold(achP.ACHLocation.PostalCode, s) {
			nsl = append(nsl, achP)
		}
	}
	return nsl
}

// ACHParticipantRoutingNumberFilter filters ACHParticipant by Routing Number
func (f *ACHDictionary) ACHParticipantRoutingNumberFilter(achParticipants []*ACHParticipant, s string) ([]*ACHParticipant, error) {
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
func (f *ACHDictionary) StateFilter(s string) []*ACHParticipant {
	nsl := make([]*ACHParticipant, 0)
	for _, achP := range f.ACHParticipants {
		if strings.EqualFold(achP.ACHLocation.State, s) {

			nsl = append(nsl, achP)
		}
	}
	return nsl
}

// CityFilter filters ACHDictionary.ACHParticipant by city
func (f *ACHDictionary) CityFilter(s string) []*ACHParticipant {
	nsl := make([]*ACHParticipant, 0)
	for _, achP := range f.ACHParticipants {
		if strings.EqualFold(achP.ACHLocation.City, s) {
			nsl = append(nsl, achP)
		}
	}
	return nsl
}

// PostalCodeFilter filters ACHParticipant by postal code
func (f *ACHDictionary) PostalCodeFilter(s string) []*ACHParticipant {
	nsl := make([]*ACHParticipant, 0)
	for _, achP := range f.ACHParticipants {
		if strings.EqualFold(achP.ACHLocation.PostalCode, s) {
			nsl = append(nsl, achP)
		}
	}
	return nsl
}

func reduceACHResults(in []*achParticipantResult, limit int) []*ACHParticipant {
	sort.SliceStable(in, func(i, j int) bool { return in[i].highestMatch > in[j].highestMatch })

	out := make([]*ACHParticipant, 0)
	for i := 0; i < limit && i < len(in); i++ {
		out = append(out, in[i].ACHParticipant)
	}
	return out
}
