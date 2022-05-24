// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package fed

import (
	"regexp"
	"strings"
)

var (
	symbolReplacer = strings.NewReplacer(
		".", "",
		",", "",
		":", "",
		"/", " ",
		"(", "",
		")", "",
		"-", " ", // 'CITIBANK-NEW YORK' => 'CITIBANK NEW YORK'
	)
	wasteReplacer = strings.NewReplacer(
		"N.A.", " ",
	)
	spaceTrimming = regexp.MustCompile(`\s{2,}`)

	// Replacements for banks that have lots of similar names.
	// The big banks will have '$name - Arizona' which confuses logo search tools, so just replace them with known-good names.
	//
	// Based on https://github.com/wealthsimple/frb-participants/blob/master/data/manually-normalized-institution-names.yml
	nameReplacements = map[string]string{
		"ALLY BANK":           "Ally Bank",
		"AMERICAN EXPRESS":    "American Express",
		"BANK OF AMERICA":     "Bank of America",
		"CAPITAL ONE":         "Capital One",
		"CHARLES SCHWAB BANK": "Charles Schwab",
		"CITIBANK":            "Citibank",
		"FIDELITY BANK":       "Fidelity",
		"HSBC":                "HSBC Bank",
		"JPMORGAN CHASE":      "Chase",
		"PNC BANK":            "PNC Bank",
		"SUNTRUST":            "SunTrust",
		"TD BANK":             "TD Bank",
		"WELLS FARGO":         "Wells Fargo",
		"US BANK":             "US Bank",
		"USAA":                "USAA",
	}
)

func Normalize(name string) string {
	for sub, answer := range nameReplacements {
		if strings.Contains(name, sub) {
			return answer
		}
	}
	return RemoveDuplicatedSpaces(StripSymbols(StripWaste(name)))
}

func StripSymbols(name string) string {
	return symbolReplacer.Replace(name)
}

func StripWaste(name string) string {
	return wasteReplacer.Replace(name)
}

func RemoveDuplicatedSpaces(name string) string {
	return spaceTrimming.ReplaceAllString(name, " ")
}
