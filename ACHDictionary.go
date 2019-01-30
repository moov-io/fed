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

// ACHDictionary of Participant records
type ACHDictionary struct {
	// Participants is a list of Participant structs
	ACHParticipants     []*ACHParticipant
	scanner             *bufio.Scanner
	line                string
	IndexACHParticipant map[string]*ACHParticipant
	// errors holds each error encountered when attempting to parse the file
	errors base.ErrorList
}

// NewACHDictionary creates a ACHDictionary
func NewACHDictionary(r io.Reader) *ACHDictionary {
	return &ACHDictionary{
		IndexACHParticipant: map[string]*ACHParticipant{},
		scanner:             bufio.NewScanner(r),
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
	// RecordTypeCode The code indicating the ABA number to be used to route or send ACH items to the RFI
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
	// ViewCode
	ViewCode string `json:"viewCode"`
}

// CustomerNameLabel returns a formatted string Title for displaying CustomerName
//ToDo: Review CU (Credit Union) which returns as Cu
func (p *ACHParticipant) CustomerNameLabel() string {
	return strings.Title(strings.ToLower(p.CustomerName))
}

// ACHLocation City name and state code in the institution's delivery address
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

		if utf8.RuneCountInString(f.line) != 155 {
			f.errors.Add(NewRecordWrongLengthErr(155, len(f.line)))
			// Return with error if the record length is incorrect as this file is a FED file
			return f.errors
		}
		if err := f.parseACHParticipant(); err != nil {
			f.errors.Add(err)
			return f.errors
		}
	}
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
	f.IndexACHParticipant[p.RoutingNumber] = p
	return nil
}

// RoutingNumberSearch returns a FEDACH participant based on a Participant.RoutingNumber.  Routing Number validation
// is only that it exists in IndexParticipant.  Expecting 9 digits, checksum needs to be included.
// ToDo: Should this remain exportable?
func (f *ACHDictionary) RoutingNumberSearch(s string) *ACHParticipant {
	if _, ok := f.IndexACHParticipant[s]; ok {
		return f.IndexACHParticipant[s]
	}
	return nil
}
