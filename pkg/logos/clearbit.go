// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package logos

import (
	"fmt"
	"os"
	"strconv"

	"github.com/clearbit/clearbit-go/clearbit"
	hashilru "github.com/hashicorp/golang-lru"
)

type Client struct {
	underlying *clearbit.Client
	lruCache   *hashilru.Cache
}

func newClearbit(apiKey string) *Client {
	client := &Client{
		underlying: clearbit.NewClient(clearbit.WithAPIKey(apiKey)),
	}

	maxSize, _ := strconv.ParseInt(os.Getenv("LOGO_CACHE_SIZE"), 10, 32)
	if maxSize > 0 {
		client.lruCache, _ = hashilru.New(int(maxSize))
	}

	return client
}

func (c *Client) Lookup(name string) (*Logo, error) {
	if c.lruCache != nil {
		item, found := c.lruCache.Get(name)
		if found {
			return item.(*Logo), nil
		}
	}

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

		logo := &Logo{
			Name: company.Name,
			URL:  company.Logo,
		}
		if c.lruCache != nil && logo.Name != "" && logo.URL != "" {
			c.lruCache.Add(name, logo)
		}
		return logo, nil
	}

	return nil, fmt.Errorf("no results for %q found", name)
}
