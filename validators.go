// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package fed

import (
	"errors"
	"regexp"
)

var (
	numericRegex = regexp.MustCompile(`[^0-9]`)
	msgNumeric   = "is not 0-9"
)

// validator is common validation and formatting of golang types to fed type strings
type validator struct{}

// isNumeric checks if a string only contains ASCII numeric (0-9) characters
func (v *validator) isNumeric(s string) error {
	if numericRegex.MatchString(s) {
		// [^0-9]
		return errors.New(msgNumeric)
	}
	return nil
}
