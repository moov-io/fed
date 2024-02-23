// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package download

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

var MissingRoutingNumberErr = "missing routing number"
var MissingDownloadCodeErr = "missing download code"

type Client struct {
	httpClient *http.Client

	routingNumber string // X_FRB_EPAYMENTS_DIRECTORY_ORG_ID header
	downloadCode  string // X_FRB_EPAYMENTS_DIRECTORY_DOWNLOAD_CD
	downloadURL   string // defaults to "https://frbservices.org/EPaymentsDirectory/directories/%s?format=json" where %s is the list name

}

type ClientOpts struct {
	HTTPClient                               *http.Client
	RoutingNumber, DownloadCode, DownloadURL string
}

func NewClient(opts *ClientOpts) (*Client, error) {
	if opts == nil {
		opts = &ClientOpts{}
	}
	if opts.HTTPClient == nil {
		opts.HTTPClient = &http.Client{
			// These files can be fairly large to buffer to us and the FRB
			// services can be slow, so we default to a hefty timeout.
			Timeout: 90 * time.Second,
		}
	}

	// the client needs either routing number && download code OR download URL

	if opts.RoutingNumber == "" {
		return nil, errors.New(MissingRoutingNumberErr)
	}
	if opts.DownloadCode == "" {
		return nil, errors.New(MissingDownloadCodeErr)
	}

	if opts.DownloadURL == "" {
		opts.DownloadURL = "https://frbservices.org/EPaymentsDirectory/directories/%s?format=json"
	}

	return &Client{
		httpClient:    opts.HTTPClient,
		routingNumber: opts.RoutingNumber,
		downloadCode:  opts.DownloadCode,
		downloadURL:   opts.DownloadURL,
	}, nil
}

// GetList downloads an FRB list and saves it into an io.Reader.
// Example listName values: fedach, fedwire
func (c *Client) GetList(listName string) (io.Reader, error) {
	where, err := url.Parse(fmt.Sprintf(c.downloadURL, listName))
	if err != nil {
		return nil, fmt.Errorf("url: %v", err)
	}

	req, err := http.NewRequest("GET", where.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("building %s url: %v", listName, err)
	}

	if c.downloadCode != "" && c.routingNumber != "" {
		req.Header.Set("X_FRB_EPAYMENTS_DIRECTORY_ORG_ID", c.routingNumber)
		req.Header.Set("X_FRB_EPAYMENTS_DIRECTORY_DOWNLOAD_CD", c.downloadCode)
	}

	// perform our request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http get: %v", err)
	}
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	// Quit if we fail to download
	if resp.StatusCode >= 299 {
		return nil, fmt.Errorf("unexpected http status: %d", resp.StatusCode)
	}

	var out bytes.Buffer
	if n, err := io.Copy(&out, resp.Body); n == 0 || err != nil {
		return nil, fmt.Errorf("copying n=%d: %v", n, err)
	}
	if out.Len() > 0 {
		return &out, nil
	}
	return nil, nil
}
