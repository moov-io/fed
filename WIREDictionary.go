// Copyright 2019 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package fed

import (
	"bufio"
	"github.com/moov-io/base"
	"io"
	"strings"
	"unicode/utf8"
)

// WIREDictionary of Participant records
type WIREDictionary struct {
	// Participants is a list of Participant structs
	WIREParticipants     []*WIREParticipant
	scanner              *bufio.Scanner
	line                 string
	IndexWIREParticipant map[string]*WIREParticipant
	// errors holds each error encountered when attempting to parse the file
	errors base.ErrorList
}

// NewWIREDictionary creates a WIREDictionary
func NewWIREDictionary(r io.Reader) *WIREDictionary {
	return &WIREDictionary{
		IndexWIREParticipant: map[string]*WIREParticipant{},
		scanner:              bufio.NewScanner(r),
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

		if utf8.RuneCountInString(f.line) != 101 {
			f.errors.Add(NewRecordWrongLengthErr(101, len(f.line)))
			// Return with error if the record length is incorrect as this file is a FED file
			return f.errors
		}

		if err := f.parseWIREParticipant(); err != nil {
			f.errors.Add(err)
			return f.errors
		}
	}
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
	f.IndexWIREParticipant[p.RoutingNumber] = p
	return nil
}

// RoutingNumberSearch returns a FEDWIRE participant based on a Participant.RoutingNumber
func (f *WIREDictionary) RoutingNumberSearch(s string) *WIREParticipant {
	if _, ok := f.IndexWIREParticipant[s]; ok {
		return f.IndexWIREParticipant[s]
	}
	return nil
}
