// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"

	client "github.com/moov-io/fed/client"

	"github.com/antihax/optional"
)

func wireSearch(api *client.APIClient, requestID, routingNumber string) error {
	opts := &client.SearchFEDWIREOpts{
		XRequestID: optional.NewString(requestID),
	}
	if routingNumber != "" {
		opts.RoutingNumber = optional.NewString(routingNumber)
	}

	dict, resp, err := api.FEDApi.SearchFEDWIRE(context.Background(), opts)
	if err != nil {
		return fmt.Errorf("FED Wire error: %v", err)
	}
	defer resp.Body.Close()

	// Verify the requested routing number was found
	for i := range dict.WIREParticipants {
		if dict.WIREParticipants[i].RoutingNumber == routingNumber {
			return nil
		}
	}
	return fmt.Errorf("FED Wire no participant found for %s", routingNumber)
}
