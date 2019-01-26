package feddir

import (
	"os"
	"strings"
	"testing"
)

// Values within the tests can change if the FED WIRE participants change (e.g. Number of participants, etc.)

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
	f, err := os.Open("./data/fpddir.txt")
	if err != nil {
		t.Errorf("%T: %s", err, err)
	}
	defer f.Close()
	wireDir := NewWIREDictionary(f)
	err = wireDir.Read()
	if err != nil {
		t.Errorf("%T: %s", err, err)
	}
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
		if !Has(err, NewRecordWrongLengthErr(101, 51)) {
			t.Errorf("%T: %s", err, err)
		}
	}
}

// TestWIRERoutingNumberSearch tests that a valid routing number defined in FedWIREDir returns the participant data
func TestWIRERoutingNumberSearch(t *testing.T) {
	f, err := os.Open("./data/fpddir.txt")
	if err != nil {
		t.Errorf("%T: %s", err, err)
	}
	defer f.Close()
	wireDir := NewWIREDictionary(f)
	err = wireDir.Read()
	if err != nil {
		t.Errorf("%T: %s", err, err)
	}

	fi := wireDir.RoutingNumberSearch("324172465")

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
func TestInvalidWIRERoutingNumberSearch(t *testing.T) {
	f, err := os.Open("./data/fpddir.txt")
	if err != nil {
		t.Errorf("%T: %s", err, err)
	}
	defer f.Close()
	wireDir := NewWIREDictionary(f)
	err = wireDir.Read()
	if err != nil {
		t.Errorf("%T: %s", err, err)
	}

	fi := wireDir.RoutingNumberSearch("325183657")

	if fi != nil {
		t.Errorf("%s", "325183657 should have returned nil")
	}
}

// ToDo:  Add test for Wire Financial Institution
