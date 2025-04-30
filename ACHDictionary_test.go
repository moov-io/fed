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

func loadTestACHFiles(t *testing.T) (*ACHDictionary, *ACHDictionary) {
	t.Helper()

	open := func(path string) *ACHDictionary {
		f, err := os.Open(path)
		if err != nil {
			t.Fatalf("%T: %s", err, err)
		}
		t.Cleanup(func() { f.Close() })

		dict := NewACHDictionary()
		if err := dict.Read(f); err != nil {
			t.Fatalf("%T: %s", err, err)
		}
		return dict
	}

	jsonDict := open(filepath.Join("data", "fedachdir.json"))
	plainDict := open(filepath.Join("data", "FedACHdir.txt"))

	return jsonDict, plainDict
}

// Values within the tests can change if the FED ACH participants change (e.g. Number of participants, etc.)

func TestACHParseParticipant(t *testing.T) {
	var line = "073905527O0710003011012908000000000LINCOLN SAVINGS BANK                P O BOX E                           REINBECK            IA506690159319788644111     "
	f := NewACHDictionary()
	f.Read(strings.NewReader(line))
	if fi, ok := f.IndexACHRoutingNumber["073905527"]; ok {
		if fi.RoutingNumber != "073905527" {
			t.Errorf("CustomerName Expected '073905527' got: %v", fi.RoutingNumber)
		}
		if fi.OfficeCode != "O" {
			t.Errorf("OfficeCode Expected 'O' got: %v", fi.OfficeCode)
		}
		if fi.ServicingFRBNumber != "071000301" {
			t.Errorf("ServicingFrbNumber Expected '071000301' got: %v", fi.ServicingFRBNumber)
		}
		if fi.RecordTypeCode != "1" {
			t.Errorf("RecordTypeCode Expected '1' got: %v", fi.RecordTypeCode)
		}
		if fi.Revised != "012908" {
			t.Errorf("Revised Expected '012908' got: %v", fi.Revised)
		}
		if fi.NewRoutingNumber != "000000000" {
			t.Errorf("NewRoutingNumber Expected '000000000' got: %v", fi.NewRoutingNumber)
		}
		if fi.CustomerName != "LINCOLN SAVINGS BANK" {
			t.Errorf("CustomerName Expected 'LINCOLN SAVINGS BANK' got: %v", fi.CustomerName)
		}
		if fi.Address != "P O BOX E" {
			t.Errorf("Address Expected 'P O BOX E' got: %v", fi.Address)
		}
		if fi.City != "REINBECK" {
			t.Errorf("City Expected 'REINBECK' got: %v", fi.City)
		}
		if fi.State != "IA" {
			t.Errorf("State Expected 'REINBECK' got: %v", fi.State)
		}
		if fi.PostalCode != "50669" {
			t.Errorf("PostalCode Expected '50669' got: %v", fi.PostalCode)
		}
		if fi.PostalCodeExtension != "0159" {
			t.Errorf("PostalCodeExtension Expected '0159' got: %v", fi.PostalCodeExtension)
		}
		if fi.PhoneNumber != "3197886441" {
			t.Errorf("PhoneNumber Expected '3197886441' got: %v", fi.PhoneNumber)
		}
		if fi.StatusCode != "1" {
			t.Errorf("StatusCode Expected '1' got: %v", fi.StatusCode)
		}
		if fi.ViewCode != "1" {
			t.Errorf("ViewCode Expected '1' got: %v", fi.ViewCode)
		}
	} else {
		t.Errorf("routing number `073905527` not found")
	}
}

func TestACHDirectoryRead(t *testing.T) {
	check := func(t *testing.T, kind string, dict *ACHDictionary) {
		if fi, ok := dict.IndexACHRoutingNumber["073905527"]; ok {
			if fi.CustomerName != "LINCOLN SAVINGS BANK" {
				t.Errorf("%s: Expected `LINCOLN SAVINGS BANK` got : %v", kind, fi.CustomerName)
			}
		} else {
			t.Errorf("%s: ach routing number `073905527` not found", kind)
		}
	}

	jsonDict, plainDict := loadTestACHFiles(t)

	if len(jsonDict.ACHParticipants) != 6 {
		t.Errorf("got %d participants", len(jsonDict.ACHParticipants))
	}
	check(t, "json", jsonDict)

	if len(plainDict.ACHParticipants) != 18198 {
		t.Errorf("got %d participants", len(plainDict.ACHParticipants))
	}
	check(t, "plain", plainDict)
}

func TestACHInvalidRecordLength(t *testing.T) {
	var line = "073905527O0710003011012908000000000LINCOLN SAVINGS BANK                P O BOX E"
	f := NewACHDictionary()
	if err := f.Read(strings.NewReader(line)); err != nil {
		if !base.Has(err, NewRecordWrongLengthErr(155, 80)) {
			t.Errorf("%T: %s", err, err)
		}
	}
}

func TestACHParticipantLabel(t *testing.T) {
	var line = "073905527O0710003011012908000000000LINCOLN SAVINGS BANK                P O BOX E                           REINBECK            IA506690159319788644111     "
	f := NewACHDictionary()
	f.Read(strings.NewReader(line))
	if fi, ok := f.IndexACHRoutingNumber["073905527"]; ok {
		if fi.CustomerNameLabel() != "Lincoln Savings Bank" {
			t.Errorf("CustomerNameLabel Expected 'Lincoln Savings Bank' got: %v", fi.CustomerNameLabel())
		}
	} else {
		t.Errorf("routing number `073905527` not found")
	}
}

// TestACHRoutingNumberSearchSingle tests that a valid routing number defined in FedACHDir returns participant data
func TestACHRoutingNumberSearchSingle(t *testing.T) {
	check := func(t *testing.T, kind string, dir *ACHDictionary) {
		fi := dir.RoutingNumberSearchSingle("325183657")
		if fi == nil {
			t.Fatalf("%s: ach routing number `325183657` not found", kind)
		}
		if fi.CustomerName != "LOWER VALLEY CU" {
			t.Errorf("%s: Expected `LOWER VALLEY CU` got : %v", kind, fi.CustomerName)
		}
	}
	jsonDict, plainDict := loadTestACHFiles(t)
	check(t, "json", jsonDict)
	check(t, "plain", plainDict)
}

// TestInvalidRoutingNumberSearchSingle tests that an invalid routing number returns nil
func TestInvalidACHRoutingNumberSearchSingle(t *testing.T) {
	jsonDict, plainDict := loadTestACHFiles(t)

	fi := jsonDict.RoutingNumberSearchSingle("433")
	if fi != nil {
		t.Error("json: 433 should have returned nil")
	}

	fi = plainDict.RoutingNumberSearchSingle("433")
	if fi != nil {
		t.Error("plain: 433 should have returned nil")
	}
}

// TestACHFinancialInstitutionSearchSingle tests that a Financial Institution defined in FedACHDir returns
// participant data
func TestACHFinancialInstitutionSearchSingle(t *testing.T) {
	check := func(t *testing.T, kind string, dict *ACHDictionary) {
		fi := dict.FinancialInstitutionSearchSingle("BANK OF AMERICA N.A")
		if len(fi) == 0 {
			t.Fatalf("%s: ach financial institution `BANK OF AMERICA N.A` not found", kind)
		}
		for _, f := range fi {
			if f.CustomerName != "BANK OF AMERICA N.A" {
				t.Errorf("%s: Expected `BANK OF AMERICA, N.A` got : %s", kind, f.CustomerName)
			}
		}
	}

	jsonDict, plainDict := loadTestACHFiles(t)
	check(t, "json", jsonDict)
	check(t, "plain", plainDict)
}

// TestInvalidACHFinancialInstitutionSearchSingle tests that a Financial Institution is not defined in FedACHDir
// returns nil
func TestInvalidACHFinancialInstitutionSearchSingle(t *testing.T) {
	jsonDict, plainDict := loadTestACHFiles(t)

	fi := jsonDict.FinancialInstitutionSearchSingle("XYZ")
	if len(fi) != 0 {
		t.Errorf("%s", "XYZ should have returned nil")
	}

	fi = plainDict.FinancialInstitutionSearchSingle("XYZ")
	if len(fi) != 0 {
		t.Errorf("%s", "XYZ should have returned nil")
	}
}

// TestACHRoutingNumberSearch tests that routing number search returns nil or FEDACH participant data
func TestACHRoutingNumberSearch(t *testing.T) {
	jsonDict, plainDict := loadTestACHFiles(t)

	fi, err := jsonDict.RoutingNumberSearch("325", 10)
	if err != nil {
		t.Fatalf("%T: %s", err, err)
	}
	if len(fi) == 0 {
		t.Errorf("%s", "325 should have returned values")
	}

	fi, err = plainDict.RoutingNumberSearch("325", 10)
	if err != nil {
		t.Fatalf("%T: %s", err, err)
	}
	if len(fi) == 0 {
		t.Errorf("%s", "325 should have returned values")
	}
}

// TestACHRoutingNumberSearch02 tests string `02` returns results
func TestACHRoutingNumberSearch02(t *testing.T) {
	jsonDict, plainDict := loadTestACHFiles(t)

	fi, err := jsonDict.RoutingNumberSearch("03", 10)
	if err != nil {
		t.Fatalf("%T: %s", err, err)
	}
	if len(fi) == 0 {
		t.Fatalf("02 should have returned values")
	}

	fi, err = plainDict.RoutingNumberSearch("02", 10)
	if err != nil {
		t.Fatalf("%T: %s", err, err)
	}
	if len(fi) == 0 {
		t.Fatalf("02 should have returned values")
	}
}

// TestACHRoutingNumberSearchMinimumLength tests that routing number search returns a RecordWrongLengthErr if the
// length of the string passed in is less than 2.
func TestACHRoutingNumberSearchMinimumLength(t *testing.T) {
	jsonDict, plainDict := loadTestACHFiles(t)

	if _, err := jsonDict.RoutingNumberSearch("0", 10); err != nil {
		if !base.Has(err, NewRecordWrongLengthErr(2, 1)) {
			t.Errorf("%T: %s", err, err)
		}
	}

	if _, err := plainDict.RoutingNumberSearch("0", 10); err != nil {
		if !base.Has(err, NewRecordWrongLengthErr(2, 1)) {
			t.Errorf("%T: %s", err, err)
		}
	}
}

// TestInvalidACHRoutingNumberSearch tests that routing number returns nil for an invalid RoutingNumber.
func TestInvalidACHRoutingNumberSearch(t *testing.T) {
	check := func(t *testing.T, kind string, dict *ACHDictionary) {
		fi, err := dict.RoutingNumberSearch("777777777", 10)
		if err != nil {
			t.Fatalf("%s: %T: %s", kind, err, err)
		}
		for i := range fi {
			fmt.Printf("fi[%d]=%#v\n", i, fi[i])
		}
		if len(fi) != 0 {
			t.Fatalf("%s: ach routing number search should have returned nil", kind)
		}
	}

	jsonDict, plainDict := loadTestACHFiles(t)
	check(t, "json", jsonDict)
	check(t, "plain", plainDict)
}

// TestACHRoutingNumberMaximumLength tests that routing number search returns a RecordWrongLengthErr if the
// length of the string passed in is greater than 9.
func TestACHRoutingNumberSearchMaximumLength(t *testing.T) {
	jsonDict, plainDict := loadTestACHFiles(t)

	if _, err := jsonDict.RoutingNumberSearch("1234567890", 10); err != nil {
		if !base.Has(err, NewRecordWrongLengthErr(9, 10)) {
			t.Errorf("%T: %s", err, err)
		}
	}
	if _, err := plainDict.RoutingNumberSearch("1234567890", 10); err != nil {
		if !base.Has(err, NewRecordWrongLengthErr(9, 10)) {
			t.Errorf("%T: %s", err, err)
		}
	}
}

// TestACHRoutingNumberNumeric tests that routing number search returns an ErrRoutingNumberNumeric if the
// string passed in is not numeric.
func TestACHRoutingNumberNumeric(t *testing.T) {
	jsonDict, plainDict := loadTestACHFiles(t)

	if _, err := jsonDict.RoutingNumberSearch("1  S5", 10); err != nil {
		if !base.Has(err, ErrRoutingNumberNumeric) {
			t.Errorf("%T: %s", err, err)
		}
	}
	if _, err := plainDict.RoutingNumberSearch("1  S5", 10); err != nil {
		if !base.Has(err, ErrRoutingNumberNumeric) {
			t.Errorf("%T: %s", err, err)
		}
	}
}

func TestACHFinancialInstitutionSearch__Examples(t *testing.T) {
	_, plainDict := loadTestACHFiles(t)

	cases := []struct {
		input    string
		expected *ACHParticipant
	}{
		{
			input: "Chase",
			expected: &ACHParticipant{
				RoutingNumber: "021000021",
				CustomerName:  "JPMORGAN CHASE",
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
			input: "Fargo",
			expected: &ACHParticipant{
				RoutingNumber: "291378392",
				CustomerName:  "FARGO VA FEDERAL CU",
			},
		},
		{
			input: "Wells Fargo",
			expected: &ACHParticipant{
				RoutingNumber: "011100106",
				CustomerName:  "WELLS FARGO BANK",
			},
		},
	}

	for i := range cases {
		// The plain dictionary has 18k records, so search is more realistic
		results := plainDict.FinancialInstitutionSearch(cases[i].input, 1)
		require.Len(t, results, 1)

		require.Equal(t, cases[i].expected.RoutingNumber, results[0].RoutingNumber)
		require.Equal(t, cases[i].expected.CustomerName, results[0].CustomerName)
	}
}

// TestACHFinancialInstitutionSearch tests search string `First Bank`
func TestACHFinancialInstitutionSearch(t *testing.T) {
	jsonDict, plainDict := loadTestACHFiles(t)

	fi := jsonDict.FinancialInstitutionSearch("First Bank", 10)
	if len(fi) == 0 {
		t.Fatalf("json: No Financial Institutions matched your search query")
	}

	fi = plainDict.FinancialInstitutionSearch("First Bank", 10)
	if len(fi) == 0 {
		t.Fatalf("plain: No Financial Institutions matched your search query")
	}
}

// TestACHFinancialInstitutionFarmers tests search string `FaRmerS`
func TestACHFinancialInstitutionFarmers(t *testing.T) {
	jsonDict, plainDict := loadTestACHFiles(t)

	fi := jsonDict.FinancialInstitutionSearch("FaRmerS", 10)
	if len(fi) == 0 {
		t.Fatalf("json: No Financial Institutions matched your search query")
	}

	fi = plainDict.FinancialInstitutionSearch("FaRmerS", 10)
	if len(fi) == 0 {
		t.Fatalf("plain: No Financial Institutions matched your search query")
	}
}

// TestACHSearchStateFilter tests search string `Farmers State Bank` and filters by the state of Ohio, `OH`
func TestACHSearchStateFilter(t *testing.T) {
	check := func(t *testing.T, kind string, dict *ACHDictionary) {
		fi := dict.FinancialInstitutionSearch("Farmers State Bank", 100)
		if len(fi) == 0 {
			t.Fatalf("%s: No Financial Institutions matched your search query", kind)
		}

		filter := dict.ACHParticipantStateFilter(fi, "MO")
		if len(filter) == 0 {
			t.Fatalf("%s: No Financial Institutions matched your search query", kind)
		}
		for _, loc := range filter {
			if loc.ACHLocation.State != "MO" {
				t.Errorf("%s: Expected `MO` got : %s", kind, loc.ACHLocation.State)
			}
		}
	}

	jsonDict, plainDict := loadTestACHFiles(t)
	check(t, "json", jsonDict)
	check(t, "plain", plainDict)
}

// TestACHSearchCityFilter tests search string `Farmers State Bank` and filters by the city of `ARCHBOLD`
func TestACHSearchCityFilter(t *testing.T) {
	check := func(t *testing.T, kind string, dict *ACHDictionary) {
		fi := dict.FinancialInstitutionSearch("Farmers State Bank", 100)
		if len(fi) == 0 {
			t.Fatalf("%s: No Financial Institutions matched your search query", kind)
		}

		filter := dict.ACHParticipantCityFilter(fi, "CAMERON")
		if len(filter) == 0 {
			t.Fatalf("%s: No Financial Institutions matched your search query", kind)
		}
		for _, loc := range filter {
			if loc.ACHLocation.City != "CAMERON" {
				t.Errorf("%s: Expected `CAMERON` got : %s", kind, loc.ACHLocation.City)
			}
		}
	}

	jsonDict, plainDict := loadTestACHFiles(t)
	check(t, "json", jsonDict)
	check(t, "plain", plainDict)
}

// TestACHSearchPostalCodeFilter tests search string `Farmers State Bank` and filters by the postal code of
func TestACHSearchPostalCodeFilter(t *testing.T) {
	check := func(t *testing.T, kind string, dict *ACHDictionary) {
		fi := dict.FinancialInstitutionSearch("Farmers State Bank", 100)
		if len(fi) == 0 {
			t.Fatalf("%s: No Financial Institutions matched your search query", kind)
		}

		filter := dict.ACHParticipantPostalCodeFilter(fi, "64429")
		if len(filter) == 0 {
			t.Fatalf("%s: No Financial Institutions matched your search query", kind)
		}
		for _, loc := range filter {
			if loc.ACHLocation.PostalCode != "64429" {
				t.Errorf("%s: Expected `64429` got : %s", kind, loc.ACHLocation.PostalCode)
			}
		}
	}

	jsonDict, plainDict := loadTestACHFiles(t)
	check(t, "json", jsonDict)
	check(t, "plain", plainDict)
}

// TestACHDictionaryStateFilter tests filtering ACHDictionary.ACHParticipants by the state of `PA`
func TestACHDictionaryStateFilter(t *testing.T) {
	check := func(t *testing.T, kind string, dict *ACHDictionary) {
		filter := dict.StateFilter("nj")
		if len(filter) == 0 {
			t.Fatalf("%s: No Financial Institutions matched your search query", kind)
		}
		for _, loc := range filter {
			if loc.ACHLocation.State != "NJ" {
				t.Errorf("%s: Expected `NJ` got : %s", kind, loc.ACHLocation.State)
			}
		}
	}

	jsonDict, plainDict := loadTestACHFiles(t)
	check(t, "json", jsonDict)
	check(t, "plain", plainDict)
}

// TestACHDictionaryCityFilter tests filtering ACHDictionary.ACHParticipants by the city of `Reading`
func TestACHDictionaryCityFilter(t *testing.T) {
	check := func(t *testing.T, kind string, dict *ACHDictionary) {
		filter := dict.CityFilter("Hamilton")
		if len(filter) == 0 {
			t.Fatalf("%s: No Financial Institutions matched your search query", kind)
		}
		for _, loc := range filter {
			if loc.ACHLocation.City != "HAMILTON" {
				t.Errorf("%s: Expected `HAMILTON` got : %s", kind, loc.ACHLocation.City)
			}
		}
	}

	jsonDict, plainDict := loadTestACHFiles(t)
	check(t, "json", jsonDict)
	check(t, "plain", plainDict)
}

// TestACHDictionaryPostalCodeFilter tests filtering ACHDictionary.ACHParticipants by the postal code of `19468`
func TestACHDictionaryPostalCodeFilter(t *testing.T) {
	check := func(t *testing.T, kind string, dict *ACHDictionary) {
		filter := dict.PostalCodeFilter("08690")
		if len(filter) == 0 {
			t.Fatalf("%s: No Financial Institutions matched your search query", kind)
		}
		for _, loc := range filter {
			if loc.ACHLocation.PostalCode != "08690" {
				t.Errorf("%s: Expected `08690` got : %s", kind, loc.ACHLocation.PostalCode)
			}
		}
	}

	jsonDict, plainDict := loadTestACHFiles(t)
	check(t, "json", jsonDict)
	check(t, "plain", plainDict)
}

func TestACHDictionary_EmptyRead(t *testing.T) {
	// Response for when headers are missing
	dict := NewACHDictionary()
	input := strings.NewReader(`{ }`)
	err := dict.Read(input)
	require.NoError(t, err)
	require.Empty(t, dict.ACHParticipants)

	// Another response (when headers are set, but incorrect)
	dict = NewACHDictionary()
	input = strings.NewReader(`{
  "fedACHParticipants" : {
    "response" : {
      "code" : 202
    },
    "fedACHParticipants" : [ ]
  }
}`)
	err = dict.Read(input)
	require.NoError(t, err)
	require.Empty(t, dict.ACHParticipants)
}
