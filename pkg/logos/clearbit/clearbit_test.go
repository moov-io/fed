// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package clearbit

import (
	"os"
	"strings"
	"testing"
)

func setupTestClient(t *testing.T) *Client {
	apiKey := os.Getenv("CLEARBIT_API_KEY")
	if apiKey == "" {
		t.Skip("missing CLEARBIT_API_KEY")
	}
	return New(apiKey)
}

func TestClient(t *testing.T) {
	client := setupTestClient(t)

	logo, err := client.Lookup("Veridian Credit Union")
	if err != nil {
		t.Error(err)
	}

	if testing.Verbose() {
		t.Logf("logo=%#v", logo)
	}

	if !strings.EqualFold(logo.Name, "Veridian Credit Union") {
		t.Errorf("got %q", logo.Name)
	}
	if !strings.Contains(logo.URL, "veridiancu.org") {
		t.Errorf("got %q", logo.URL)
	}
}
