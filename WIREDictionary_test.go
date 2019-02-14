// Copyright 2019 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package fed

import (
	"github.com/moov-io/base"
	"os"
	"strings"
	"testing"
)

// Values within the tests can change if the FED WIRE participants change (e.g. Number of participants, etc.)

func helperLoadFEDWIREFile(t *testing.T) *WIREDictionary {
	f, err := os.Open("./data/fpddir.txt")
	if err != nil {
		t.Fatalf("%T: %s", err, err)
	}
	defer f.Close()
	wireDir := NewWIREDictionary(f)
	err = wireDir.Read()
	if err != nil {
		t.Fatalf("%T: %s", err, err)
	}
	return wireDir
}

func TestWIREParseParticipant(t *testing.T) {
	var line = "325280039MAC FCU           MAC FEDERAL CREDIT UNION            AKFAIRBANKS                Y Y20180629"

	f := NewWIREDictionary(strings.NewReader(line))
	f.Read()

	if fi, ok := f.IndexWIRERoutingNumber["325280039"]; ok {
		if fi.RoutingNumber != "325280039" {
			t.Errorf("Expected `325280039` got : %s", fi.RoutingNumber)
		}
		if fi.TelegraphicName != "MAC FCU" {
			t.Errorf("Expected `MAC FCU` got : %s", fi.TelegraphicName)
		}
		if fi.CustomerName != "MAC FEDERAL CREDIT UNION" {
			t.Errorf("Expected `MAC FEDERAL CREDIT UNION` got : %s", fi.CustomerName)
		}
		if fi.State != "AK" {
			t.Errorf("Expected `AK` got ; %s", fi.State)
		}
		if fi.City != "FAIRBANKS" {
			t.Errorf("Expected `FAIRBANKS` got : %s", fi.City)
		}
		if fi.FundsTransferStatus != "Y" {
			t.Errorf("Expected `Y` got : %s", fi.FundsTransferStatus)
		}
		if fi.FundsSettlementOnlyStatus != " " {
			t.Errorf("Expected ` ` got : %s", fi.FundsSettlementOnlyStatus)
		}
		if fi.BookEntrySecuritiesTransferStatus != "Y" {
			t.Errorf("Expected `Y` got : %s", fi.BookEntrySecuritiesTransferStatus)
		}
		if fi.Date != "20180629" {
			t.Errorf("Expected `20180629` got : %s", fi.Date)
		}
	} else {
		t.Errorf("routing number `325280039` not found")
	}
}

func TestWIREDirectoryRead(t *testing.T) {
	wireDir := helperLoadFEDWIREFile(t)
	if len(wireDir.WIREParticipants) != 7693 {
		t.Errorf("Expected '7693' got: %v", len(wireDir.WIREParticipants))
	}
	if fi, ok := wireDir.IndexWIRERoutingNumber["325280039"]; ok {
		if fi.TelegraphicName != "MAC FCU" {
			t.Errorf("Expected `MAC FCU` got : %s", fi.TelegraphicName)
		}
	} else {
		t.Errorf("routing number `325280039` not found")
	}
}

func TestWIREInvalidRecordLength(t *testing.T) {
	var line = "325280039MAC FCU           MAC FEDERAL CREDIT UNION"
	f := NewWIREDictionary(strings.NewReader(line))
	if err := f.Read(); err != nil {
		if !base.Has(err, NewRecordWrongLengthErr(101, 51)) {
			t.Errorf("%T: %s", err, err)
		}
	}
}

// TestWIRERoutingNumberSearch tests that a valid routing number defined in FedWIREDir returns the participant data
func TestWIRERoutingNumberSearchSingle(t *testing.T) {
	wireDir := helperLoadFEDWIREFile(t)
	fi := wireDir.RoutingNumberSearchSingle("324172465")
	if fi == nil {
		t.Errorf("wire routing number `324172465` not found")
	}
	if fi != nil {
		if fi.CustomerName != "TRUGROCER FEDERAL CREDIT UNION" {
			t.Errorf("Expected `TRUGROCER FEDERAL CREDIT UNION` got : %s", fi.CustomerName)
		}
	}
}

// TestInvalidWIRERoutingNumberSearch tests that an invalid routing number returns nil
func TestInvalidWIRERoutingNumberSearchSingle(t *testing.T) {
	wireDir := helperLoadFEDWIREFile(t)
	fi := wireDir.RoutingNumberSearchSingle("325183657")

	if fi != nil {
		t.Errorf("%s", "325183657 should have returned nil")
	}
}

// TestWIREFinancialInstitutionSearch tests that a Financial Institution defined in FedWIREDir returns the participant
// data
func TestWIREFinancialInstitutionSearchSingle(t *testing.T) {
	wireDir := helperLoadFEDWIREFile(t)
	fi := wireDir.FinancialInstitutionSearchSingle("TRUGROCER FEDERAL CREDIT UNION")
	if fi == nil {
		t.Fatalf("wire financial institution `TRUGROCER FEDERAL CREDIT UNION` not found")
	}
	for _, f := range fi {
		if f.CustomerName != "TRUGROCER FEDERAL CREDIT UNION" {
			t.Errorf("TRUGROCER FEDERAL CREDIT UNION` got : %v", f.CustomerName)
		}
	}
}

// TestInvalidWIREFinancialInstitutionSearchSingle tests that a Financial Institution defined in FedWIREDir returns
// the participant data
func TestInvalidWIREFinancialInstitutionSearchSingle(t *testing.T) {
	wireDir := helperLoadFEDWIREFile(t)
	fi := wireDir.FinancialInstitutionSearchSingle("XYZ")
	if fi != nil {
		t.Errorf("%s", "XYZ should have returned nil")
	}
}

// TestWIRERoutingNumberSearch tests that routing number search returns nil or FEDWIRE participant data
func TestWIRERoutingNumberSearch(t *testing.T) {
	wireDir := helperLoadFEDWIREFile(t)
	fi, err := wireDir.RoutingNumberSearch("325")
	if err != nil {
		t.Fatalf("%T: %s", err, err)
	}
	if len(fi) == 0 {
		t.Errorf("%s", "325 should have returned values")
	}
}

// TestWIRERoutingNumberSearch02 tests string `02` returns results
func TestWIRERoutingNumberSearch02(t *testing.T) {
	wireDir := helperLoadFEDWIREFile(t)
	fi, err := wireDir.RoutingNumberSearch("02")
	if err != nil {
		t.Fatalf("%T: %s", err, err)
	}
	if len(fi) == 0 {
		t.Fatalf("02 should have returned values")
	}

}

// TestWIRERoutingNumberSearchMinimumLength tests that routing number search returns a RecordWrongLengthErr if the
// length of the string passed in is less than 2.
func TestWIRERoutingNumberSearchMinimumLength(t *testing.T) {
	wireDir := helperLoadFEDWIREFile(t)
	if _, err := wireDir.RoutingNumberSearch("0"); err != nil {
		if !base.Has(err, NewRecordWrongLengthErr(2, 1)) {
			t.Errorf("%T: %s", err, err)
		}
	}
}

// TestInvalidWIRERoutingNumberSearch tests that routing number returns nil for an invalid RoutingNumber.
func TestInvalidWIRERoutingNumberSearch(t *testing.T) {
	wireDir := helperLoadFEDWIREFile(t)
	fi, err := wireDir.RoutingNumberSearch("777777777")
	if err != nil {
		t.Fatalf("%T: %s", err, err)
	}
	if len(fi) != 0 {
		t.Fatal("wire routing number search should have returned nil")
	}
}

// TestWIRERoutingNumberMaximumLength tests that routing number search returns a RecordWrongLengthErr if the
// length of the string passed in is greater than 9.
func TestWIRERoutingNumberSearchMaximumLength(t *testing.T) {
	wireDir := helperLoadFEDWIREFile(t)
	if _, err := wireDir.RoutingNumberSearch("1234567890"); err != nil {
		if !base.Has(err, NewRecordWrongLengthErr(9, 10)) {
			t.Errorf("%T: %s", err, err)
		}
	}
}

// TestWIRERoutingNumberNumeric tests that routing number search returns an ErrRoutingNumberNumeric if the
// string passed in is not numeric.
func TestWIRERoutingNumberNumeric(t *testing.T) {
	wireDir := helperLoadFEDWIREFile(t)
	if _, err := wireDir.RoutingNumberSearch("1  S5"); err != nil {
		if !base.Has(err, ErrRoutingNumberNumeric) {
			t.Errorf("%T: %s", err, err)
		}
	}
}

func TestWIREParsingError(t *testing.T) {
	var line = "011000536FHLB BOSTON       FEDERAL HOME LOAN BANK              MABOSTON                   © Y20170818"
	f := NewWIREDictionary(strings.NewReader(line))
	if err := f.Read(); err != nil {
		if !base.Has(err, NewRecordWrongLengthErr(101, 51)) {
			t.Errorf("%T: %s", err, err)
		}
	}
}

// TestWIREFinancialInstitutionSearch tests search string `First Bank`
func TestWIREFinancialInstitutionSearch(t *testing.T) {
	wireDir := helperLoadFEDWIREFile(t)
	fi, err := wireDir.FinancialInstitutionSearch("First Bank")
	if err != nil {
		t.Fatalf("%T: %s", err, err)
	}
	if len(fi) == 0 {
		t.Fatalf("No Financial Institutions matched your search query")
	}
}

// TestWIREFinancialInstitutionFarmers tests search string `FaRmerS`
func TestWIREFinancialInstitutionFarmers(t *testing.T) {
	wireDir := helperLoadFEDWIREFile(t)
	fi, err := wireDir.FinancialInstitutionSearch("FaRmerS")
	if err != nil {
		t.Fatalf("%T: %s", err, err)
	}
	if len(fi) == 0 {
		t.Fatalf("No Financial Institutions matched your search query")
	}
}

// TestWIRESearchStateFilter tests search string `Farmers State Bank` and filters by the state of North Carolina, `NC`
func TestWIRESearchStateFilter(t *testing.T) {
	wireDir := helperLoadFEDWIREFile(t)
	fi, err := wireDir.FinancialInstitutionSearch("Farmers State Bank")
	if err != nil {
		t.Fatalf("%T: %s", err, err)
	}
	if len(fi) == 0 {
		t.Fatalf("No Financial Institutions matched your search query")
	}

	filter := wireDir.WIREParticipantStateFilter(fi, "NC")
	if len(filter) == 0 {
		t.Fatalf("No Financial Institutions matched your search query")
	}
	for _, loc := range filter {
		if loc.WIRELocation.State != "NC" {
			t.Errorf("Expected `NC` got : %s", loc.WIRELocation.State)
		}
	}
}

// TestWIRESearchCityFilter tests search string `Farmers State Bank` and filters by the city of `SALISBURY`
func TestWIRESearchCityFilter(t *testing.T) {
	wireDir := helperLoadFEDWIREFile(t)
	fi, err := wireDir.FinancialInstitutionSearch("Farmers State Bank")
	if err != nil {
		t.Fatalf("%T: %s", err, err)
	}
	if len(fi) == 0 {
		t.Fatalf("No Financial Institutions matched your search query")
	}

	filter := wireDir.WIREParticipantCityFilter(fi, "SALISBURY")
	if len(filter) == 0 {
		t.Fatalf("No Financial Institutions matched your search query")
	}
	for _, loc := range filter {
		if loc.WIRELocation.City != "SALISBURY" {
			t.Errorf("Expected `SALISBURY` got : %s", loc.WIRELocation.City)
		}
	}
}

// TestWIREDictionaryStateFilter tests filtering WIREDictionary.WIREParticipants by the state of `PA`
func TestWIREDictionaryStateFilter(t *testing.T) {
	wireDir := helperLoadFEDWIREFile(t)

	filter := wireDir.StateFilter("pa")
	if len(filter) == 0 {
		t.Fatalf("No Financial Institutions matched your search query")
	}
	for _, loc := range filter {
		if loc.WIRELocation.State != "PA" {
			t.Errorf("Expected `PA` got : %s", loc.WIRELocation.State)
		}
	}
}

// TestWIREDictionaryCityFilter tests filtering WIREDictionary.WIREParticipants by the city of `Reading`
func TestWIREDictionaryCityFilter(t *testing.T) {
	wireDir := helperLoadFEDWIREFile(t)

	filter := wireDir.CityFilter("Reading")
	if len(filter) == 0 {
		t.Fatalf("No Financial Institutions matched your search query")
	}
	for _, loc := range filter {
		if loc.WIRELocation.City != "READING" {
			t.Errorf("Expected `READING` got : %s", loc.WIRELocation.City)
		}
	}
}
