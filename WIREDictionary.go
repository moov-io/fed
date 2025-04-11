// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package fed

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"math"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/moov-io/base"
	"github.com/moov-io/fed/pkg/strcmp"
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
func NewWIREDictionary() *WIREDictionary {
	return &WIREDictionary{
		IndexWIRERoutingNumber: map[string]*WIREParticipant{},
		IndexWIRECustomerName:  map[string][]*WIREParticipant{},
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

	// CleanName is our cleaned up value of CustomerName
	CleanName string `json:"cleanName"`
}

// WIRELocation is the city and state
type WIRELocation struct {
	// City
	City string `json:"city"`
	// State
	State string `json:"state"`
}

// jsonWIREDictionary is a JSON representation of the FED Wire Participant directory
// Model generated with // https://mholt.github.io/json-to-go/
type jsonWIREDictionary struct {
	FedwireParticipants struct {
		Response struct {
			Code int `json:"code"`
		} `json:"response"`
		FedwireParticipants []struct {
			RoutingNumber             string `json:"routingNumber"`
			TelegraphicName           string `json:"telegraphicName"`
			CustomerName              string `json:"customerName"`
			CustomerState             string `json:"customerState"`
			CustomerCity              string `json:"customerCity"`
			FundsEligibility          string `json:"fundsEligibility"`
			FundsSettlementOnlyStatus string `json:"fundsSettlementOnlyStatus"`
			SecuritiesEligibility     string `json:"securitiesEligibility"`
			ChangeDate                string `json:"changeDate"`
		} `json:"fedwireParticipants"`
	} `json:"fedwireParticipants"`
}

// Read parses a single line or multiple lines of FedWIREdir text
func (f *WIREDictionary) Read(r io.Reader) error {
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

func (f *WIREDictionary) readJSON(r io.Reader) error {
	var wrapper jsonWIREDictionary
	if err := json.NewDecoder(r).Decode(&wrapper); err != nil {
		return err
	}
	ps := wrapper.FedwireParticipants.FedwireParticipants
	for i := range ps {
		p := &WIREParticipant{
			RoutingNumber:   ps[i].RoutingNumber,
			TelegraphicName: ps[i].TelegraphicName,
			CustomerName:    ps[i].CustomerName,
			WIRELocation: WIRELocation{
				City:  ps[i].CustomerCity,
				State: ps[i].CustomerState,
			},
			FundsTransferStatus:               ps[i].FundsEligibility,
			FundsSettlementOnlyStatus:         ps[i].FundsSettlementOnlyStatus,
			BookEntrySecuritiesTransferStatus: ps[i].SecuritiesEligibility,
			Date:                              ps[i].ChangeDate,

			// Our Custom Fields
			CleanName: Normalize(ps[i].CustomerName),
		}
		f.WIREParticipants = append(f.WIREParticipants, p)
		f.IndexWIRERoutingNumber[p.RoutingNumber] = p
	}
	f.createIndexWIRECustomerName()
	return nil
}

func (f *WIREDictionary) readPlaintext(r io.Reader) error {
	if f == nil || r == nil {
		return nil
	}
	// read each line and lift it into a ACHParticipant
	s := bufio.NewScanner(r)
	var line string
	for s.Scan() {
		line = s.Text()
		if utf8.RuneCountInString(line) != WIRELineLength {
			f.errors.Add(NewRecordWrongLengthErr(WIRELineLength, len(line)))
			// Return with error if the record length is incorrect as this file is a FED file
			return f.errors
		}
		if err := f.parseWIREParticipant(line); err != nil {
			f.errors.Add(err)
			return f.errors
		}
	}
	f.createIndexWIRECustomerName()
	return nil
}

// TODO return a parsing error if the format or file is wrong.
func (f *WIREDictionary) parseWIREParticipant(line string) error {
	p := new(WIREParticipant)

	//RoutingNumber (9): 011000015
	p.RoutingNumber = line[:9]
	// TelegraphicName (18): FED
	p.TelegraphicName = strings.Trim(line[9:27], " ")
	// CustomerName (36): FEDERAL RESERVE BANK
	p.CustomerName = strings.Trim(line[27:63], " ")
	p.WIRELocation = WIRELocation{
		// State (2): GA
		State: line[63:65],
		// City (25): ATLANTA
		City: strings.Trim(line[65:90], " "),
	}
	// FundsTransferStatus (1): Y or N
	p.FundsTransferStatus = line[90:91]
	// FundsSettlementOnlyStatus (1): " " or S - Settlement-Only
	p.FundsSettlementOnlyStatus = line[91:92]
	// BookEntrySecuritiesTransferStatus (1): Y or N
	p.BookEntrySecuritiesTransferStatus = line[92:93]
	// Date YYYYMMDD (8): 122415
	p.Date = line[93:101]

	// Our custom fields
	p.CleanName = Normalize(p.CustomerName)

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
func (f *WIREDictionary) RoutingNumberSearch(s string, limit int) ([]*WIREParticipant, error) {
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

	out := make([]*wireParticipantResult, 0)
	for _, wireP := range f.WIREParticipants {
		if exactMatch {
			if wireP.RoutingNumber == s {
				out = append(out, &wireParticipantResult{
					WIREParticipant: wireP,
					highestMatch:    1.0,
				})
			}
		} else {
			out = append(out, &wireParticipantResult{
				WIREParticipant: wireP,
				highestMatch:    strcmp.JaroWinkler(wireP.RoutingNumber, s),
			})
		}
	}
	return reduceWIREResults(out, limit), nil
}

// FinancialInstitutionSearch returns a FEDWIRE participant based on a WIREParticipant.CustomerName
func (f *WIREDictionary) FinancialInstitutionSearch(s string, limit int) []*WIREParticipant {
	s = strings.ToLower(s)

	out := make([]*wireParticipantResult, 0)

	for _, wireP := range f.WIREParticipants {
		// JaroWinkler is a more accurate version of the Jaro algorithm. It works by boosting the
		// score of exact matches at the beginning of the strings. By doing this, Winkler says that
		// typos are less common to happen at the beginning.
		jaroScore := strcmp.JaroWinkler(strings.ToLower(wireP.CleanName), s)

		// Levenshtein is the "edit distance" between two strings. This is the count of operations
		// (insert, delete, replace) needed for two strings to be equal.
		levenScore := strcmp.Levenshtein(strings.ToLower(wireP.CleanName), s)

		if jaroScore > ACHJaroWinklerSimilarity || levenScore > ACHLevenshteinSimilarity {
			out = append(out, &wireParticipantResult{
				WIREParticipant: wireP,
				highestMatch:    math.Max(jaroScore, levenScore),
			})
		}
	}

	return reduceWIREResults(out, limit)
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

type wireParticipantResult struct {
	*WIREParticipant

	highestMatch float64
}

func reduceWIREResults(in []*wireParticipantResult, limit int) []*WIREParticipant {
	sort.SliceStable(in, func(i, j int) bool { return in[i].highestMatch > in[j].highestMatch })

	out := make([]*WIREParticipant, 0)
	for i := 0; i < limit && i < len(in); i++ {
		out = append(out, in[i].WIREParticipant)
	}
	return out
}
