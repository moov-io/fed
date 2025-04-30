// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package fed

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/moov-io/base"
	"github.com/stretchr/testify/require"
)

// loadTestWireFiles returns two WIREDictionary, one from the JSON source file
// and other from the plaintext source file.
func loadTestWireFiles(t *testing.T) (*WIREDictionary, *WIREDictionary) {
	t.Helper()

	open := func(path string) *WIREDictionary {
		f, err := os.Open(path)
		if err != nil {
			t.Fatalf("%T: %s", err, err)
		}
		t.Cleanup(func() { f.Close() })

		dict := NewWIREDictionary()
		if err := dict.Read(f); err != nil {
			t.Fatalf("%T: %s", err, err)
		}
		return dict
	}

	jsonDict := open(filepath.Join("data", "fpddir.json"))
	plainDict := open(filepath.Join("data", "fpddir.txt"))

	return jsonDict, plainDict
}

func TestWIREParseParticipant(t *testing.T) {
	var line = "325280039MAC FCU           MAC FEDERAL CREDIT UNION            AKFAIRBANKS                Y Y20180629"

	f := NewWIREDictionary()
	f.Read(strings.NewReader(line))

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
		if fi.WIRELocation.State != "AK" {
			t.Errorf("Expected `AK` got ; %s", fi.State)
		}
		if fi.WIRELocation.City != "FAIRBANKS" {
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
	jsonDict, plainDict := loadTestWireFiles(t)

	check := func(t *testing.T, kind string, dict *WIREDictionary) {
		if fi, ok := dict.IndexWIRERoutingNumber["325280039"]; ok {
			if fi.TelegraphicName != "MAC FCU" {
				t.Errorf("Expected `MAC FCU` got : %s", fi.TelegraphicName)
			}
		} else {
			t.Errorf("routing number `325280039` not found")
		}
	}

	if len(jsonDict.WIREParticipants) != 10 {
		t.Errorf("Expected '7693' got: %v", len(jsonDict.WIREParticipants))
	}
	check(t, "json", jsonDict)

	if len(plainDict.WIREParticipants) != 7693 {
		t.Errorf("Expected '7693' got: %v", len(plainDict.WIREParticipants))
	}
	check(t, "plain", plainDict)
}

func TestWIREInvalidRecordLength(t *testing.T) {
	var line = "325280039MAC FCU           MAC FEDERAL CREDIT UNION"
	f := NewWIREDictionary()
	if err := f.Read(strings.NewReader(line)); err != nil {
		if !base.Has(err, NewRecordWrongLengthErr(101, 51)) {
			t.Errorf("%T: %s", err, err)
		}
	}
}

// TestWIRERoutingNumberSearch tests that a valid routing number defined in FedWIREDir returns the participant data
func TestWIRERoutingNumberSearchSingle(t *testing.T) {
	jsonDict, plainDict := loadTestWireFiles(t)

	check := func(t *testing.T, kind string, dict *WIREDictionary) {
		fi := dict.RoutingNumberSearchSingle("324172465")
		if fi == nil {
			t.Fatalf("wire routing number `324172465` not found")
		}
		if fi.CustomerName != "TRUGROCER FEDERAL CREDIT UNION" {
			t.Errorf("Expected `TRUGROCER FEDERAL CREDIT UNION` got : %s", fi.CustomerName)
		}
	}

	check(t, "json", jsonDict)
	check(t, "plain", plainDict)
}

// TestInvalidWIRERoutingNumberSearch tests that an invalid routing number returns nil
func TestInvalidWIRERoutingNumberSearchSingle(t *testing.T) {
	jsonDict, plainDict := loadTestWireFiles(t)

	check := func(t *testing.T, kind string, dict *WIREDictionary) {
		fi := dict.RoutingNumberSearchSingle("325183657")
		if fi != nil {
			t.Errorf("%s", "325183657 should have returned nil")
		}
	}

	check(t, "json", jsonDict)
	check(t, "plain", plainDict)
}

// TestWIREFinancialInstitutionSearch tests that a Financial Institution defined in FedWIREDir returns the participant
// data
func TestWIREFinancialInstitutionSearchSingle(t *testing.T) {
	jsonDict, plainDict := loadTestWireFiles(t)

	check := func(t *testing.T, kind string, dict *WIREDictionary) {
		fi := dict.FinancialInstitutionSearchSingle("TRUGROCER FEDERAL CREDIT UNION")
		if fi == nil {
			t.Fatalf("wire financial institution `TRUGROCER FEDERAL CREDIT UNION` not found")
		}
		for _, f := range fi {
			if f.CustomerName != "TRUGROCER FEDERAL CREDIT UNION" {
				t.Errorf("TRUGROCER FEDERAL CREDIT UNION` got : %v", f.CustomerName)
			}
		}
	}

	check(t, "json", jsonDict)
	check(t, "plain", plainDict)
}

// TestInvalidWIREFinancialInstitutionSearchSingle tests that a Financial Institution defined in FedWIREDir returns
// the participant data
func TestInvalidWIREFinancialInstitutionSearchSingle(t *testing.T) {
	jsonDict, plainDict := loadTestWireFiles(t)

	check := func(t *testing.T, kind string, dict *WIREDictionary) {
		fi := dict.FinancialInstitutionSearchSingle("XYZ")
		if fi != nil {
			t.Errorf("%s", "XYZ should have returned nil")
		}
	}

	check(t, "json", jsonDict)
	check(t, "plain", plainDict)
}

// TestWIRERoutingNumberSearch tests that routing number search returns nil or FEDWIRE participant data
func TestWIRERoutingNumberSearch(t *testing.T) {
	jsonDict, plainDict := loadTestWireFiles(t)

	check := func(t *testing.T, kind string, dict *WIREDictionary) {
		fi, err := dict.RoutingNumberSearch("325", 1)
		if err != nil {
			t.Fatalf("%T: %s", err, err)
		}
		if len(fi) == 0 {
			t.Errorf("%s", "325 should have returned values")
		}
	}

	check(t, "json", jsonDict)
	check(t, "plain", plainDict)
}

// TestWIRERoutingNumberSearch02 tests string `02` returns results
func TestWIRERoutingNumberSearch02(t *testing.T) {
	jsonDict, plainDict := loadTestWireFiles(t)

	check := func(t *testing.T, kind string, dict *WIREDictionary) {
		fi, err := dict.RoutingNumberSearch("02", 1)
		if err != nil {
			t.Fatalf("%T: %s", err, err)
		}
		if len(fi) == 0 {
			t.Fatalf("02 should have returned values")
		}
	}

	_ = jsonDict
	// check(t, "json", jsonDict)
	check(t, "plain", plainDict)
}

// TestWIRERoutingNumberSearchMinimumLength tests that routing number search returns a RecordWrongLengthErr if the
// length of the string passed in is less than 2.
func TestWIRERoutingNumberSearchMinimumLength(t *testing.T) {
	jsonDict, plainDict := loadTestWireFiles(t)

	check := func(t *testing.T, kind string, dict *WIREDictionary) {
		if _, err := dict.RoutingNumberSearch("0", 1); err != nil {
			if !base.Has(err, NewRecordWrongLengthErr(2, 1)) {
				t.Errorf("%T: %s", err, err)
			}
		}
	}

	check(t, "json", jsonDict)
	check(t, "plain", plainDict)
}

// TestInvalidWIRERoutingNumberSearch tests that routing number returns nil for an invalid RoutingNumber.
func TestInvalidWIRERoutingNumberSearch(t *testing.T) {
	jsonDict, plainDict := loadTestWireFiles(t)

	check := func(t *testing.T, kind string, dict *WIREDictionary) {
		fi, err := dict.RoutingNumberSearch("777777777", 1)
		if err != nil {
			t.Fatalf("%T: %s", err, err)
		}
		if len(fi) != 0 {
			t.Fatal("wire routing number search should have returned nil")
		}
	}

	check(t, "json", jsonDict)
	check(t, "plain", plainDict)
}

// TestWIRERoutingNumberMaximumLength tests that routing number search returns a RecordWrongLengthErr if the
// length of the string passed in is greater than 9.
func TestWIRERoutingNumberSearchMaximumLength(t *testing.T) {
	jsonDict, plainDict := loadTestWireFiles(t)

	check := func(t *testing.T, kind string, dict *WIREDictionary) {
		if _, err := dict.RoutingNumberSearch("1234567890", 1); err != nil {
			if !base.Has(err, NewRecordWrongLengthErr(9, 10)) {
				t.Errorf("%T: %s", err, err)
			}
		}

	}

	check(t, "json", jsonDict)
	check(t, "plain", plainDict)
}

// TestWIRERoutingNumberNumeric tests that routing number search returns an ErrRoutingNumberNumeric if the
// string passed in is not numeric.
func TestWIRERoutingNumberNumeric(t *testing.T) {
	jsonDict, plainDict := loadTestWireFiles(t)

	check := func(t *testing.T, kind string, dict *WIREDictionary) {
		if _, err := dict.RoutingNumberSearch("1  S5", 1); err != nil {
			if !base.Has(err, ErrRoutingNumberNumeric) {
				t.Errorf("%T: %s", err, err)
			}
		}
	}

	check(t, "json", jsonDict)
	check(t, "plain", plainDict)
}

func TestWIREParsingError(t *testing.T) {
	var line = "011000536FHLB BOSTON       FEDERAL HOME LOAN BANK              MABOSTON                   © Y20170818"
	f := NewWIREDictionary()
	if err := f.Read(strings.NewReader(line)); err != nil {
		if !base.Has(err, NewRecordWrongLengthErr(101, 51)) {
			t.Errorf("%T: %s", err, err)
		}
	}
}

func TestWIREFinancialInstitutionSearch__Examples(t *testing.T) {
	_, plainDict := loadTestWireFiles(t)

	cases := []struct {
		input    string
		expected *ACHParticipant
	}{
		{
			input: "Chase",
			expected: &ACHParticipant{
				RoutingNumber: "021000021",
				CustomerName:  "JPMORGAN CHASE BANK, NA",
			},
		},
		{
			input: "Wells",
			expected: &ACHParticipant{
				RoutingNumber: "101205940",
				CustomerName:  "WELLS BANK",
			},
		},
		{
			input: "Wells Fargo",
			expected: &ACHParticipant{
				RoutingNumber: "021052943",
				CustomerName:  "WELLS FARGO GNMA-P&I",
			},
		},
	}

	for i := range cases {
		// The plain dictionary has 18k records, so search is more realistic
		results := plainDict.FinancialInstitutionSearch(cases[i].input, 1)
		require.Equal(t, fmt.Sprintf("#%d = 1", i), fmt.Sprintf("#%d = %d", i, len(results)))

		require.Equal(t, cases[i].expected.RoutingNumber, results[0].RoutingNumber)
		require.Equal(t, cases[i].expected.CustomerName, results[0].CustomerName)
	}
}

// TestWIREFinancialInstitutionSearch tests search string `First Bank`
func TestWIREFinancialInstitutionSearch(t *testing.T) {
	jsonDict, plainDict := loadTestWireFiles(t)

	check := func(t *testing.T, kind string, dict *WIREDictionary) {
		fi := dict.FinancialInstitutionSearch("First Bank", 1)
		if len(fi) == 0 {
			t.Fatalf("No Financial Institutions matched your search query")
		}
	}

	check(t, "json", jsonDict)
	check(t, "plain", plainDict)
}

// TestWIREFinancialInstitutionFarmers tests search string `FaRmerS`
func TestWIREFinancialInstitutionFarmers(t *testing.T) {
	jsonDict, plainDict := loadTestWireFiles(t)

	check := func(t *testing.T, kind string, dict *WIREDictionary) {
		fi := dict.FinancialInstitutionSearch("FaRmerS", 1)
		if len(fi) == 0 {
			t.Fatalf("No Financial Institutions matched your search query")
		}
	}

	check(t, "json", jsonDict)
	check(t, "plain", plainDict)
}

// TestWIRESearchStateFilter tests search string `Farmers State Bank` and filters by the state of North Carolina, `NC`
func TestWIRESearchStateFilter(t *testing.T) {
	jsonDict, plainDict := loadTestWireFiles(t)

	check := func(t *testing.T, kind string, dict *WIREDictionary) {
		fi := dict.FinancialInstitutionSearch("Farmers State Bank", 100)
		if len(fi) == 0 {
			t.Fatalf("No Financial Institutions matched your search query")
		}

		filter := dict.WIREParticipantStateFilter(fi, "NC")
		if len(filter) == 0 {
			t.Fatalf("No Financial Institutions matched your search query")
		}
		for _, loc := range filter {
			if loc.WIRELocation.State != "NC" {
				t.Errorf("Expected `NC` got : %s", loc.WIRELocation.State)
			}
		}
	}

	check(t, "json", jsonDict)
	check(t, "plain", plainDict)
}

// TestWIRESearchCityFilter tests search string `Farmers State Bank` and filters by the city of `SALISBURY`
func TestWIRESearchCityFilter(t *testing.T) {
	jsonDict, plainDict := loadTestWireFiles(t)

	check := func(t *testing.T, kind string, dict *WIREDictionary) {
		fi := dict.FinancialInstitutionSearch("Farmers State Bank", 100)
		if len(fi) == 0 {
			t.Fatalf("No Financial Institutions matched your search query")
		}

		filter := dict.WIREParticipantCityFilter(fi, "SALISBURY")
		if len(filter) == 0 {
			t.Fatalf("No Financial Institutions matched your search query")
		}

		for _, loc := range filter {
			if loc.WIRELocation.City != "SALISBURY" {
				t.Errorf("Expected `SALISBURY` got : %s", loc.WIRELocation.City)
			}
		}
	}

	check(t, "json", jsonDict)
	check(t, "plain", plainDict)
}

// TestWIREDictionaryStateFilter tests filtering WIREDictionary.WIREParticipants by the state of `PA`
func TestWIREDictionaryStateFilter(t *testing.T) {
	jsonDict, plainDict := loadTestWireFiles(t)

	check := func(t *testing.T, kind string, dict *WIREDictionary) {
		filter := dict.StateFilter("pa")
		if len(filter) == 0 {
			t.Fatalf("No Financial Institutions matched your search query")
		}
		for _, loc := range filter {
			if loc.WIRELocation.State != "PA" {
				t.Errorf("Expected `PA` got : %s", loc.WIRELocation.State)
			}
		}
	}

	check(t, "json", jsonDict)
	check(t, "plain", plainDict)
}

// TestWIREDictionaryCityFilter tests filtering WIREDictionary.WIREParticipants by the city of `Reading`
func TestWIREDictionaryCityFilter(t *testing.T) {
	jsonDict, plainDict := loadTestWireFiles(t)

	check := func(t *testing.T, kind string, dict *WIREDictionary) {
		filter := dict.CityFilter("Reading")
		if len(filter) == 0 {
			t.Fatalf("No Financial Institutions matched your search query")
		}
		for _, loc := range filter {
			if loc.WIRELocation.City != "READING" {
				t.Errorf("Expected `READING` got : %s", loc.WIRELocation.City)
			}
		}
	}

	check(t, "json", jsonDict)
	check(t, "plain", plainDict)
}

func TestWireDictionary_EmptyRead(t *testing.T) {
	// Response for when headers are missing
	dict := NewWIREDictionary()
	input := strings.NewReader(`{ }`)
	err := dict.Read(input)
	require.NoError(t, err)
	require.Empty(t, dict.WIREParticipants)

	// Another response (when headers are set, but incorrect)
	dict = NewWIREDictionary()
	input = strings.NewReader(`{
  "fedwireParticipants" : {
    "response" : {
      "code" : 202
    },
    "fedwireParticipants" : [ ]
  }
}`)
	err = dict.Read(input)
	require.NoError(t, err)
	require.Empty(t, dict.WIREParticipants)
}
