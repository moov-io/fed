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
	"os"
	"time"

	"github.com/moov-io/base/strx"
)

type Client struct {
	httpClient *http.Client

	routingNumber string // X_FRB_EPAYMENTS_DIRECTORY_ORG_ID header
	downloadCode  string // X_FRB_EPAYMENTS_DIRECTORY_DOWNLOAD_CD
}

type ClientOpts struct {
	HTTPClient                  *http.Client
	RoutingNumber, DownloadCode string
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

	if opts.RoutingNumber == "" {
		return nil, errors.New("missing routing number")
	}
	if opts.DownloadCode == "" {
		return nil, errors.New("missing download code")
	}

	return &Client{
		httpClient:    opts.HTTPClient,
		routingNumber: opts.RoutingNumber,
		downloadCode:  opts.DownloadCode,
	}, nil
}

var (
	downloadDirectory = strx.Or(os.Getenv("DOWNLOAD_DIRECTORY"), os.TempDir())
)

func init() {
	if _, err := os.Stat(downloadDirectory); os.IsNotExist(err) {
		os.MkdirAll(downloadDirectory, 0777)
	}
}

// GetList downloads an FRB list and saves it into an io.Reader.
// Example listName values: fedach, fedwire
func (c *Client) GetList(listName string) (io.Reader, error) {
	where, err := url.Parse(fmt.Sprintf("https://frbservices.org/EPaymentsDirectory/directories/%s?format=json", listName))
	if err != nil {
		return nil, fmt.Errorf("url: %v", err)
	}

	req, err := http.NewRequest("GET", where.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("building %s url: %v", listName, err)
	}
	req.Header.Set("X_FRB_EPAYMENTS_DIRECTORY_ORG_ID", c.routingNumber)
	req.Header.Set("X_FRB_EPAYMENTS_DIRECTORY_DOWNLOAD_CD", c.downloadCode)

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
