// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package fed

const (
	// ACHLineLength is the FedACH text file line length
	ACHLineLength = 155
	// WIRELineLength is the FedACH text file line length
	WIRELineLength = 101
	// MinimumRoutingNumberDigits is the minimum number of digits needed searching by routing numbers
	MinimumRoutingNumberDigits = 2
	// MaximumRoutingNumberDigits is the maximum number of digits allowed for searching by routing number
	// Based on https://www.frbservices.org/EPaymentsDirectory/search.html
	MaximumRoutingNumberDigits = 9
)
