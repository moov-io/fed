// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package fed

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalize(t *testing.T) {
	cases := []struct {
		input, expected string
	}{
		{input: "ALLY BANK", expected: "Ally Bank"},
		{input: "BANK OF AMERICA, N.A. - ARIZONA", expected: "Bank of America"},
		{input: "CITIBANK FSB", expected: "Citibank"},
		{input: "CITIBANK /FLORIDA/BR", expected: "Citibank"},
		{input: "CITIBANK-NEW YORK STATE", expected: "Citibank"},
		{input: "CITIBANK (SOUTH DAKOTA) NA", expected: "Citibank"},
		{input: "PNC BANK INC. - BALTIMORE", expected: "PNC Bank"},
		{input: "SUNTRUST BANK/ SKYLIGHT", expected: "SunTrust"},
		{input: "TD BANK NA (PC)", expected: "TD Bank"},
		{input: "WELLS FARGO BANK, N.A.", expected: "Wells Fargo"},
	}
	for i := range cases {
		output := Normalize(cases[i].input)
		require.Equal(t, cases[i].expected, output)
	}
}

func TestStripSymbols(t *testing.T) {
	cases := []struct {
		input, expected string
	}{
		{input: "ALLY BANK", expected: "ALLY BANK"},
		{input: "BANK OF AMERICA, N.A. - ARIZONA", expected: "BANK OF AMERICA NA   ARIZONA"},
		{input: "CITIBANK FSB", expected: "CITIBANK FSB"},
		{input: "CITIBANK /FLORIDA/BR", expected: "CITIBANK  FLORIDA BR"},
		{input: "CITIBANK-NEW YORK STATE", expected: "CITIBANK NEW YORK STATE"},
		{input: "CITIBANK (SOUTH DAKOTA) NA", expected: "CITIBANK SOUTH DAKOTA NA"},
		{input: "PNC BANK INC. - BALTIMORE", expected: "PNC BANK INC   BALTIMORE"},
		{input: "SUNTRUST BANK/ SKYLIGHT", expected: "SUNTRUST BANK  SKYLIGHT"},
		{input: "TD BANK NA (PC)", expected: "TD BANK NA PC"},
		{input: "WELLS FARGO BANK, N.A.", expected: "WELLS FARGO BANK NA"},
	}
	for i := range cases {
		output := StripSymbols(cases[i].input)
		require.Equal(t, cases[i].expected, output)
	}
}

func TestStripWaste(t *testing.T) {
	cases := []struct {
		input, expected string
	}{
		{input: "BANK OF AMERICA, N.A. - ARIZONA", expected: "BANK OF AMERICA,   - ARIZONA"},
		{input: "WELLS FARGO BANK, N.A.", expected: "WELLS FARGO BANK,  "},
	}
	for i := range cases {
		output := StripWaste(cases[i].input)
		require.Equal(t, cases[i].expected, output)
	}
}

func TestRemoveDuplicatedSpaces(t *testing.T) {
	cases := []struct {
		input, expected string
	}{
		{input: "ALLY BANK", expected: "ALLY BANK"},
		{input: "BANK OF AMERICA NA   ARIZONA", expected: "BANK OF AMERICA NA ARIZONA"},
		{input: "CITIBANK  FLORIDA BR", expected: "CITIBANK FLORIDA BR"},
		{input: "PNC BANK INC   BALTIMORE", expected: "PNC BANK INC BALTIMORE"},
		{input: "PNC  BANK   INC   BALTIMORE", expected: "PNC BANK INC BALTIMORE"},
	}
	for i := range cases {
		output := RemoveDuplicatedSpaces(cases[i].input)
		require.Equal(t, cases[i].expected, output)
	}
}
