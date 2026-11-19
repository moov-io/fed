// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package download

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

const DefaultFRBDownloadURLTemplate = "https://frbservices.org/EPaymentsDirectory/directories/%s?format=json"

// FRBSuccessResponseCode is the response.code FRB eServices returns when a
// directory download succeeds. Sample files in this repo use 100. Other codes
// (for example 202) are returned with an empty participant list when the
// routing number or download code is missing or incorrect.
const FRBSuccessResponseCode = 100

var (
	ErrMissingConfigValue   = errors.New("missing config value")
	ErrMissingRoutingNumber = errors.New("missing routing number")
	ErrMissingDownloadCD    = errors.New("missing download code")
	ErrEmptyDirectory       = errors.New("FRB directory contained no participants")
)

type Client struct {
	httpClient *http.Client

	routingNumber string // X_FRB_EPAYMENTS_DIRECTORY_ORG_ID header
	downloadCode  string // X_FRB_EPAYMENTS_DIRECTORY_DOWNLOAD_CD
	downloadURL   string // defaults to "https://frbservices.org/EPaymentsDirectory/directories/%s?format=json" where %s is the list name
	logf          func(string, ...any)
}

type ClientOpts struct {
	HTTPClient                               *http.Client
	RoutingNumber, DownloadCode, DownloadURL string
	// Logf, if set, receives download status messages (HTTP status and FRB response code).
	Logf func(string, ...any)
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

	routingNum := firstNonEmpty(opts.RoutingNumber, os.Getenv("FRB_ROUTING_NUMBER"))
	downloadcd := firstNonEmpty(opts.DownloadCode, os.Getenv("FRB_DOWNLOAD_CODE"))
	downloadurltemp := firstNonEmpty(opts.DownloadURL, os.Getenv("FRB_DOWNLOAD_URL_TEMPLATE"))

	if downloadurltemp == "" {
		if routingNum == "" {
			return nil, fmt.Errorf("%w: %w", ErrMissingConfigValue, ErrMissingRoutingNumber)
		}

		if downloadcd == "" {
			return nil, fmt.Errorf("%w: %w", ErrMissingConfigValue, ErrMissingDownloadCD)
		}

		downloadurltemp = DefaultFRBDownloadURLTemplate
	}

	return &Client{
		httpClient:    opts.HTTPClient,
		routingNumber: routingNum,
		downloadCode:  downloadcd,
		downloadURL:   downloadurltemp,
		logf:          opts.Logf,
	}, nil
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
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
	if c.logf != nil {
		c.logf("FRB %s download HTTP %d", listName, resp.StatusCode)
	}

	var out bytes.Buffer
	if n, err := io.Copy(&out, resp.Body); n == 0 || err != nil {
		return nil, fmt.Errorf("copying n=%d: %v", n, err)
	}
	if out.Len() == 0 {
		return nil, nil
	}

	code, count, err := inspectFRBList(listName, out.Bytes())
	if c.logf != nil && code != 0 {
		c.logf("FRB %s directory response code %d (%d participants)", listName, code, count)
	}
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// frbListEnvelope is the official FRB eServices JSON wrapper for directory downloads.
type frbListEnvelope struct {
	FedACHParticipants *struct {
		Response struct {
			Code int `json:"code"`
		} `json:"response"`
		FedACHParticipants []json.RawMessage `json:"fedACHParticipants"`
	} `json:"fedACHParticipants"`
	FedwireParticipants *struct {
		Response struct {
			Code int `json:"code"`
		} `json:"response"`
		FedwireParticipants []json.RawMessage `json:"fedwireParticipants"`
	} `json:"fedwireParticipants"`
}

// inspectFRBList reads the FRB application-level response code from a JSON
// directory download. Plaintext files are ignored. A non-success code or an
// empty participant list is returned as an error; HTTP 200 is not sufficient
// because FRB still returns 200 when the download code is rejected.
func inspectFRBList(listName string, body []byte) (code, count int, err error) {
	if !json.Valid(body) {
		return 0, 0, nil
	}

	var env frbListEnvelope
	if err := json.Unmarshal(body, &env); err != nil {
		return 0, 0, nil
	}

	switch {
	case env.FedACHParticipants != nil:
		code = env.FedACHParticipants.Response.Code
		count = len(env.FedACHParticipants.FedACHParticipants)
	case env.FedwireParticipants != nil:
		code = env.FedwireParticipants.Response.Code
		count = len(env.FedwireParticipants.FedwireParticipants)
	default:
		return 0, 0, nil
	}

	if code != 0 && code != FRBSuccessResponseCode {
		return code, count, fmt.Errorf("%w: FRB %s response code %d with %d participants; check FRB_ROUTING_NUMBER and FRB_DOWNLOAD_CODE", ErrEmptyDirectory, listName, code, count)
	}
	if count == 0 {
		return code, count, fmt.Errorf("%w: FRB %s response code %d returned no participants; check FRB_ROUTING_NUMBER and FRB_DOWNLOAD_CODE", ErrEmptyDirectory, listName, code)
	}
	return code, count, nil
}
