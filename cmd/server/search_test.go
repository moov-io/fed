// Copyright 2019 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"net/url"
	"strings"
	"testing"
)

// TestSearch__fedachSearchRequest
func TestSearch__fedachSearchRequest(t *testing.T) {
	u, _ := url.Parse("https://moov.io/fed/ach/search?name=Farmers&routingNumber=044112187&city=CALDWELL&state=OH&postalCode=43724")
	req := readFEDACHSearchRequest(u)
	if req.Name != "FARMERS" {
		t.Errorf("req.Name=%s", req.Name)
	}
	if req.RoutingNumber != "044112187" {
		t.Errorf("req.RoutingNUmber=%s", req.RoutingNumber)
	}

	if req.City != "CALDWELL" {
		t.Errorf("req.City=%s", req.City)
	}
	if req.State != "OH" {
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
	req.Name = "FARMERS"
	if req.empty() {
		t.Error("req is not empty now")
	}
}

// TestSearch__fedachNameOnlySearchRequest by name only
func TestSearch__fedachNameOnlySearchRequest(t *testing.T) {
	u, _ := url.Parse("https://moov.io/fed/ach/search?name=Farmers")
	req := readFEDACHSearchRequest(u)
	if req.Name != "FARMERS" {
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
	if req.State != "OH" {
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
	if req.City != "CALDWELL" {
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

func TestSearcher_ACHFindNameOnly(t *testing.T) {
	s := searcher{}
	if err := s.readFEDACHData(); err != nil {
		t.Fatal(err)
	}

	achP, err := s.ACHFindNameOnly("Farmers")
	if err != nil {
		t.Fatal(err)
	}

	if len(achP) == 0 {
		t.Fatalf("%s", "No matches found for name")
	}

	for _, p := range achP {
		if !strings.Contains(p.CustomerName, strings.ToUpper("Farmer")) {
			t.Errorf("Name=%s", p.CustomerName)
		}
	}
}

func TestSearcher_ACHFindRoutingNumberOnly(t *testing.T) {
	s := searcher{}
	if err := s.readFEDACHData(); err != nil {
		t.Fatal(err)
	}

	achP, err := s.ACHFindRoutingNumberOnly("044112187")
	if err != nil {
		t.Fatal(err)
	}

	if len(achP) == 0 {
		t.Fatalf("%s", "No matches found for routing number")
	}

	for _, p := range achP {
		if !strings.Contains(p.RoutingNumber, "044112187") {
			t.Errorf("Routing Number=%s", p.RoutingNumber)
		}
	}
}

func TestSearcher_ACHFindCityOnly(t *testing.T) {
	s := searcher{}
	if err := s.readFEDACHData(); err != nil {
		t.Fatal(err)
	}

	achP := s.ACHFindCityOnly("CALDWELL")

	if len(achP) == 0 {
		t.Fatalf("%s", "No matches found for city")
	}

	for _, p := range achP {
		if !strings.Contains(p.City, strings.ToUpper("CALDWELL")) {
			t.Errorf("City=%s", p.City)
		}
	}
}

func TestSearcher_ACHFindStateOnly(t *testing.T) {
	s := searcher{}
	if err := s.readFEDACHData(); err != nil {
		t.Fatal(err)
	}

	achP := s.ACHFindStateOnly("OH")

	if len(achP) == 0 {
		t.Fatalf("%s", "No matches found for state")
	}

	for _, p := range achP {
		if !strings.Contains(p.State, "OH") {
			t.Errorf("State=%s", p.State)
		}
	}
}

func TestSearcher_ACHFindPostalCodeOnly(t *testing.T) {
	s := searcher{}
	if err := s.readFEDACHData(); err != nil {
		t.Fatal(err)
	}

	achP := s.ACHFindPostalCodeOnly("43724")

	if len(achP) == 0 {
		t.Fatalf("%s", "No matches found for postal code")
	}

	for _, p := range achP {
		if !strings.Contains(p.PostalCode, "43724") {
			t.Errorf("Postal Code=%s", p.PostalCode)
		}
	}
}

func TestSearcher_ACHFind(t *testing.T) {
	s := searcher{}
	if err := s.readFEDACHData(); err != nil {
		t.Fatal(err)
	}

	req := fedachSearchRequest{
		Name:          "Farmers",
		RoutingNumber: "044112187",
		City:          "Caldwell",
		State:         "OH",
		PostalCode:    "43724",
	}

	achP, err := s.FindFEDACH(req)

	if err != nil {
		t.Fatal(err)
	}

	if len(achP) == 0 {
		t.Fatalf("%s", "No matches found for postal code")
	}

	for _, p := range achP {
		if !strings.Contains(p.CustomerName, strings.ToUpper("Farmer")) {
			t.Errorf("Name=%s", p.CustomerName)
		}

		if !strings.Contains(p.RoutingNumber, "044112187") {
			t.Errorf("Routing Number=%s", p.RoutingNumber)
		}

		if !strings.Contains(p.City, strings.ToUpper("CALDWELL")) {
			t.Errorf("City=%s", p.City)
		}
		if !strings.Contains(p.State, "OH") {
			t.Errorf("State=%s", p.State)
		}

		if !strings.Contains(p.PostalCode, "43724") {
			t.Errorf("Postal Code=%s", p.PostalCode)
		}
	}
}
