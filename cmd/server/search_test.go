// Copyright 2019 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"net/url"
	"testing"
)

// TestSearch__fedachSearchRequest
func TestSearch__fedachSearchRequest(t *testing.T) {
	u, _ := url.Parse("https://moov.io/fed/ach/search?name=Farmers&routingNumber=044112187&city=CALDWELL&state=OH&postalCode=43724")
	req := readFEDACHSearchRequest(u)
	if req.Name != "farmers" {
		t.Errorf("req.Name=%s", req.Name)
	}
	if req.RoutingNumber != "044112187" {
		t.Errorf("req.RoutingNUmber=%s", req.RoutingNumber)
	}

	if req.City != "caldwell" {
		t.Errorf("req.City=%s", req.City)
	}
	if req.State != "oh" {
		t.Errorf("req.State=%s", req.State)
	}
	if req.PostalCode != "43724" {
		t.Errorf("req.Zip=%s", req.PostalCode)
	}
	if req.empty() {
		t.Error("req is not empty")
	}
	req = fedachSearchRequest{}
	if !req.empty() {
		t.Error("req is empty now")
	}
	req.Name = "farmers"
	if req.empty() {
		t.Error("req is not empty now")
	}
}

// TestSearch__fedachNameOnlySearchRequest by name only
func TestSearch__fedachNameOnlySearchRequest(t *testing.T) {
	u, _ := url.Parse("https://moov.io/fed/ach/search?name=Farmers")
	req := readFEDACHSearchRequest(u)
	if req.Name != "farmers" {
		t.Errorf("req.Name=%s", req.Name)
	}
	if !req.nameOnly() {
		t.Error("req is not name only")
	}
}

// TestSearch__fedachRoutingNumberOnlySearchRequest by Routing Number Only
func TestSearch__fedachRoutingNumberOnlySearchRequest(t *testing.T) {
	u, _ := url.Parse("https://moov.io/fed/ach/search?routingNumber=044112187")
	req := readFEDACHSearchRequest(u)
	if req.RoutingNumber != "044112187" {
		t.Errorf("req.RoutingNUmber=%s", req.RoutingNumber)
	}
	if !req.routingNumberOnly() {
		t.Errorf("req is not routing number only")
	}
}

// TestSearch__fedachSearchStateOnlyRequest by state only
func TestSearch__fedachSearchStateOnlyRequest(t *testing.T) {
	u, _ := url.Parse("https://moov.io/fed/ach/search?state=OH")
	req := readFEDACHSearchRequest(u)
	if req.State != "oh" {
		t.Errorf("req.State=%s", req.State)
	}
	if !req.stateOnly() {
		t.Errorf("req is not state only")
	}
}

// TestSearch__fedachCityOnlySearchRequest by city only
func TestSearch__fedachCityOnlySearchRequest(t *testing.T) {
	u, _ := url.Parse("https://moov.io/fed/ach/search?city=CALDWELL")
	req := readFEDACHSearchRequest(u)
	if req.City != "caldwell" {
		t.Errorf("req.City=%s", req.City)
	}
	if !req.cityOnly() {
		t.Errorf("req is not city only")
	}
}

// TestSearch__fedachPostalCodeOnlySearchRequest by postal code only
func TestSearch__fedachPostalCodeOnlySearchRequest(t *testing.T) {
	u, _ := url.Parse("https://moov.io/fed/ach/search?postalCode=43724")
	req := readFEDACHSearchRequest(u)
	if req.PostalCode != "43724" {
		t.Errorf("req.Zip=%s", req.PostalCode)
	}
	if !req.postalCodeOnly() {
		t.Errorf("req is not postal code only")
	}
}
