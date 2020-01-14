// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package fed

import (
	"errors"
	"fmt"
)

// ErrFileTooLong is the error given when a file exceeds the maximum possible length
var (
	ErrFileTooLong = errors.New("file exceeds maximum possible number of lines")
	// Similar to FEDACH site
	ErrRoutingNumberNumeric = errors.New("the routing number entered is not numeric")
)

// RecordWrongLengthErr is the error given when a record is the wrong length
type RecordWrongLengthErr struct {
	Message        string
	LengthRequired int
	Length         int
}

// NewRecordWrongLengthErr creates a new error of the RecordWrongLengthErr type
func NewRecordWrongLengthErr(lengthRequired int, length int) RecordWrongLengthErr {
	return RecordWrongLengthErr{
		Message:        fmt.Sprintf("must be %d characters and found %d", lengthRequired, length),
		LengthRequired: lengthRequired,
		Length:         length,
	}
}

func (e RecordWrongLengthErr) Error() string {
	return e.Message
}
