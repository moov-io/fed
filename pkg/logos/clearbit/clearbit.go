// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package clearbit

import (
	"fmt"

	"github.com/moov-io/fed/pkg/logos"

	"github.com/clearbit/clearbit-go/clearbit"
)

type Client struct {
	underlying *clearbit.Client
}

func New(apiKey string) *Client {
	return &Client{
		underlying: clearbit.NewClient(clearbit.WithAPIKey(apiKey)),
	}
}

func (c *Client) Lookup(name string) (*logos.Logo, error) {
	result, resp, err := c.underlying.NameToDomain.Find(clearbit.NameToDomainFindParams{
		Name: name,
	})
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if result != nil {
		company, resp, err := c.underlying.Company.Find(clearbit.CompanyFindParams{
			Domain: result.Domain,
		})
		defer resp.Body.Close()
		if err != nil {
			return nil, err
		}

		return &logos.Logo{
			Name: company.Name,
			URL:  company.Logo,
		}, nil
	}

	return nil, fmt.Errorf("no results for %q found", name)
}
