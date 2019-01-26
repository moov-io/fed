// Copyright 2019 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package feddir

import (
	"os"
	"strings"
	"testing"
)

// Values within the tests can change if the FED ACH participants change (e.g. Number of participants, etc.)

func TestACHParseParticipant(t *testing.T) {
	var line = "073905527O0710003011012908000000000LINCOLN SAVINGS BANK                P O BOX E                           REINBECK            IA506690159319788644111     "

	f := NewACHDictionary(strings.NewReader(line))
	f.Read()

	if fi, ok := f.IndexACHRoutingNumber["073905527"]; ok {
		if fi.RoutingNumber != "073905527" {
			t.Errorf("CustomerName Expected '073905527' got: %v", fi.RoutingNumber)
		}
		if fi.OfficeCode != "O" {
			t.Errorf("OfficeCode Expected 'O' got: %v", fi.OfficeCode)
		}
		if fi.ServicingFrbNumber != "071000301" {
			t.Errorf("ServicingFrbNumber Expected '071000301' got: %v", fi.ServicingFrbNumber)
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
	f, err := os.Open("./data/FedACHdir.txt")
	if err != nil {
		t.Errorf("%T: %s", err, err)
	}
	defer f.Close()
	achDir := NewACHDictionary(f)
	err = achDir.Read()
	if err != nil {
		t.Fatalf("%T: %s", err, err)
	}
	if len(achDir.ACHParticipants) != 18198 {
		t.Errorf("Expected '18198' got: %v", len(achDir.ACHParticipants))
	}
	if fi, ok := achDir.IndexACHRoutingNumber["073905527"]; ok {
		if fi.CustomerName != "LINCOLN SAVINGS BANK" {
			t.Errorf("Expected `LINCOLN SAVINGS BANK` got : %v", fi.CustomerName)
		}
	} else {
		t.Errorf("ach routing number `073905527` not found")
	}
}

func TestACHInvalidRecordLength(t *testing.T) {
	var line = "073905527O0710003011012908000000000LINCOLN SAVINGS BANK                P O BOX E"
	f := NewACHDictionary(strings.NewReader(line))
	if err := f.Read(); err != nil {
		if !Has(err, NewRecordWrongLengthErr(155, 80)) {
			t.Errorf("%T: %s", err, err)
		}
	}
}

func TestACHParticipantLabel(t *testing.T) {
	var line = "073905527O0710003011012908000000000LINCOLN SAVINGS BANK                P O BOX E                           REINBECK            IA506690159319788644111     "

	f := NewACHDictionary(strings.NewReader(line))
	f.Read()

	if fi, ok := f.IndexACHRoutingNumber["073905527"]; ok {
		if fi.CustomerNameLabel() != "Lincoln Savings Bank" {
			t.Errorf("CustomerNameLabel Expected 'Lincoln Savings Bank' got: %v", fi.CustomerNameLabel())
		}
	} else {
		t.Errorf("routing number `073905527` not found")
	}
}

// TestACHRoutingNumberSearch tests that a valid routing number defined in FedACHDir returns the participant data
func TestACHRoutingNumberSearch(t *testing.T) {
	f, err := os.Open("./data/FedACHdir.txt")
	if err != nil {
		t.Errorf("%T: %s", err, err)
	}
	defer f.Close()
	achDir := NewACHDictionary(f)
	err = achDir.Read()
	if err != nil {
		t.Fatalf("%T: %s", err, err)
	}

	fi := achDir.RoutingNumberSearch("325183657")

	if fi == nil {
		t.Errorf("ach routing number `325183657` not found")
	}

	if fi.CustomerName != "LOWER VALLEY CU" {
		t.Errorf("Expected `LOWER VALLEY CU` got : %v", fi.CustomerName)
	}
}

// TestInvalidRoutingNumberSearch tests that an invalid routing number returns nil
func TestInvalidACHRoutingNumberSearch(t *testing.T) {
	f, err := os.Open("./data/FedACHdir.txt")
	if err != nil {
		t.Errorf("%T: %s", err, err)
	}
	defer f.Close()
	achDir := NewACHDictionary(f)
	err = achDir.Read()
	if err != nil {
		t.Fatalf("%T: %s", err, err)
	}

	fi := achDir.RoutingNumberSearch("433")

	if fi != nil {
		t.Errorf("%s", "433 should have returned nil")
	}
}

// TestACHFinancialInstitutionSearch tests that a Financial Institution defined in FedACHDir returns the participant data
func TestACHFinancialInstitutionSearch(t *testing.T) {
	f, err := os.Open("./data/FedACHdir.txt")
	if err != nil {
		t.Errorf("%T: %s", err, err)
	}
	defer f.Close()
	achDir := NewACHDictionary(f)
	err = achDir.Read()
	if err != nil {
		t.Errorf("%T: %s", err, err)
	}

	fi := achDir.FinancialInstitutionSearch("BANK OF AMERICA N.A")

	if fi == nil {
		t.Fatalf("ach financial institution `BANK OF AMERICA N.A` not found")
	}

	for _, f := range fi {
		if f.CustomerName != "BANK OF AMERICA N.A" {
			t.Errorf("Expected `BANK OF AMERICA, N.A.` got : %v", f.CustomerName)
		}
	}
}
