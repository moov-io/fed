// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/moov-io/fed"
)

func (s *searcher) helperLoadFEDACHFile(t *testing.T) error {
	f, err := os.Open("../.././data/FedACHdir.txt")
	if err != nil {
		return fmt.Errorf("ERROR: opening FedACHdir.txt %v", err)
	}
	defer f.Close()

	s.ACHDictionary = fed.NewACHDictionary()
	if err := s.ACHDictionary.Read(f); err != nil {
		return fmt.Errorf("ERROR: reading FedACHdir.txt %v", err)
	}
	return nil
}

func (s *searcher) helperLoadFEDWIREFile(t *testing.T) error {
	f, err := os.Open("../.././data/fpddir.txt")
	if err != nil {
		return fmt.Errorf("ERROR: opening fpddir.txt %v", err)
	}
	defer f.Close()

	s.WIREDictionary = fed.NewWIREDictionary()
	if err := s.WIREDictionary.Read(f); err != nil {
		return fmt.Errorf("ERROR: reading fpddir.txt %v", err)
	}
	return nil
}

// TestSearch__fedachSearchRequest
func TestSearch__fedachSearchRequest(t *testing.T) {
	u, _ := url.Parse("https://moov.io/fed/ach/search?name=Farmers&routingNumber=044112187&city=CALDWELL&state=OH&postalCode=43724")
	req := readFEDSearchRequest(u)
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
	req = fedSearchRequest{}
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
	req := readFEDSearchRequest(u)
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
	req := readFEDSearchRequest(u)
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
	req := readFEDSearchRequest(u)
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
	req := readFEDSearchRequest(u)
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
	req := readFEDSearchRequest(u)
	if req.PostalCode != "43724" {
		t.Errorf("req.Zip=%s", req.PostalCode)
	}
	if !req.postalCodeOnly() {
		t.Errorf("req is not postal code only")
	}
}

func TestSearcher_ACHFindNameOnly(t *testing.T) {
	s := searcher{}
	if err := s.helperLoadFEDACHFile(t); err != nil {
		t.Fatal(err)
	}

	achP := s.ACHFindNameOnly(hardResultsLimit, "Farmers")

	if len(achP) == 0 {
		t.Fatalf("%s", "No matches found for name")
	}

	for _, p := range achP {
		if !strings.Contains(p.CustomerName, strings.ToUpper("Farm")) {
			t.Errorf("Name=%s", p.CustomerName)
		}
	}
}

func TestSearcher_ACHFindRoutingNumberOnly(t *testing.T) {
	s := searcher{}
	if err := s.helperLoadFEDACHFile(t); err != nil {
		t.Fatal(err)
	}

	achP, err := s.ACHFindRoutingNumberOnly(10, "044112187")
	if err != nil {
		t.Fatal(err)
	}

	if len(achP) == 0 {
		t.Fatalf("%s", "No matches found for routing number")
	}

	for _, p := range achP {
		if strings.HasPrefix(p.RoutingNumber, "041") {
			continue
		}
		if strings.HasPrefix(p.RoutingNumber, "044") {
			continue
		}
		t.Errorf("Routing Number=%s", p.RoutingNumber)
	}
}

func TestSearcher_ACHFindCityOnly(t *testing.T) {
	s := searcher{}
	if err := s.helperLoadFEDACHFile(t); err != nil {
		t.Fatal(err)
	}

	achP := s.ACHFindCityOnly(hardResultsLimit, "CALDWELL")

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
	if err := s.helperLoadFEDACHFile(t); err != nil {
		t.Fatal(err)
	}

	achP := s.ACHFindStateOnly(hardResultsLimit, "OH")

	if len(achP) == 0 {
		t.Fatalf("%s", "No matches found for state")
	}

	for _, p := range achP {
		if !strings.Contains(p.State, strings.ToUpper("OH")) {
			t.Errorf("State=%s", p.State)
		}
	}
}

func TestSearcher_ACHFindPostalCodeOnly(t *testing.T) {
	s := searcher{}
	if err := s.helperLoadFEDACHFile(t); err != nil {
		t.Fatal(err)
	}

	achP := s.ACHFindPostalCodeOnly(hardResultsLimit, "43724")

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
	if err := s.helperLoadFEDACHFile(t); err != nil {
		t.Fatal(err)
	}

	req := fedSearchRequest{
		Name:          "Farmers",
		RoutingNumber: "044112187",
		City:          "Caldwell",
		State:         "OH",
		PostalCode:    "43724",
	}

	achP, err := s.ACHFind(hardResultsLimit, req)

	if err != nil {
		t.Fatal(err)
	}

	if len(achP) == 0 {
		t.Fatalf("%s", "No matches found for search")
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

// TestSearch__fedwireSearchRequest
func TestSearch__fedwireSearchRequest(t *testing.T) {
	u, _ := url.Parse("https://moov.io/fed/wire/search?name=MIDWest&routingNumber=091905114&state=IA&city=IOWA City")
	req := readFEDSearchRequest(u)
	if req.Name != "MIDWEST" {
		t.Errorf("req.Name=%s", req.Name)
	}
	if req.RoutingNumber != "091905114" {
		t.Errorf("req.RoutingNUmber=%s", req.RoutingNumber)
	}

	if req.City != "IOWA CITY" {
		t.Errorf("req.City=%s", req.City)
	}
	if req.State != "IA" {
		t.Errorf("req.State=%s", req.State)
	}
	if req.empty() {
		t.Error("req is not empty")
	}
	req = fedSearchRequest{}
	if !req.empty() {
		t.Error("req is empty now")
	}
	req.Name = "MIDWEST"
	if req.empty() {
		t.Error("req is not empty now")
	}
}

// TestSearch__fedwireNameOnlySearchRequest by name only
func TestSearch__fedwireNameOnlySearchRequest(t *testing.T) {
	u, _ := url.Parse("https://moov.io/fed/wire/search?name=MIDWest")
	req := readFEDSearchRequest(u)
	if req.Name != "MIDWEST" {
		t.Errorf("req.Name=%s", req.Name)
	}
	if !req.nameOnly() {
		t.Error("req is not name only")
	}
}

// TestSearch__fedwireRoutingNumberOnlySearchRequest by Routing Number Only
func TestSearch__fedwireRoutingNumberOnlySearchRequest(t *testing.T) {
	u, _ := url.Parse("https://moov.io/fed/wire/search?routingNumber=091905114")
	req := readFEDSearchRequest(u)
	if req.RoutingNumber != "091905114" {
		t.Errorf("req.RoutingNUmber=%s", req.RoutingNumber)
	}
	if !req.routingNumberOnly() {
		t.Errorf("req is not routing number only")
	}
}

// TestSearch__fedwireSearchStateOnlyRequest by state only
func TestSearch__fedwireSearchStateOnlyRequest(t *testing.T) {
	u, _ := url.Parse("https://moov.io/fed/wire/search?state=IA")
	req := readFEDSearchRequest(u)
	if req.State != "IA" {
		t.Errorf("req.State=%s", req.State)
	}
	if !req.stateOnly() {
		t.Errorf("req is not state only")
	}
}

// TestSearch__fedwireCityOnlySearchRequest by city only
func TestSearch__fedwireCityOnlySearchRequest(t *testing.T) {
	u, _ := url.Parse("https://moov.io/fed/wire/search?city=IOWA City")
	req := readFEDSearchRequest(u)
	if req.City != "IOWA CITY" {
		t.Errorf("req.City=%s", req.City)
	}
	if !req.cityOnly() {
		t.Errorf("req is not city only")
	}
}

func TestSearcher_WIREFindNameOnly(t *testing.T) {
	s := searcher{}
	if err := s.helperLoadFEDWIREFile(t); err != nil {
		t.Fatal(err)
	}

	wireP := s.WIREFindNameOnly(hardResultsLimit, "MIDWEST")

	if len(wireP) == 0 {
		t.Fatalf("%s", "No matches found for name")
	}

	for _, p := range wireP {
		if !strings.Contains(p.CustomerName, strings.ToUpper("MIDWEST")) {
			t.Errorf("Name=%s", p.CustomerName)
		}
	}
}

func TestSearcher_WIREFindRoutingNumberOnly(t *testing.T) {
	s := searcher{}
	if err := s.helperLoadFEDWIREFile(t); err != nil {
		t.Fatal(err)
	}

	wireP, err := s.WIREFindRoutingNumberOnly(hardResultsLimit, "091905114")
	if err != nil {
		t.Fatal(err)
	}

	if len(wireP) == 0 {
		t.Fatalf("%s", "No matches found for routing number")
	}

	for _, p := range wireP {
		if !strings.Contains(p.RoutingNumber, "091905114") {
			t.Errorf("Routing Number=%s", p.RoutingNumber)
		}
	}
}

func TestSearcher_WIREFindCityOnly(t *testing.T) {
	s := searcher{}
	if err := s.helperLoadFEDWIREFile(t); err != nil {
		t.Fatal(err)
	}

	wireP := s.WIREFindCityOnly(hardResultsLimit, "IOWA CITY")

	if len(wireP) == 0 {
		t.Fatalf("%s", "No matches found for city")
	}

	for _, p := range wireP {
		if !strings.Contains(p.City, strings.ToUpper("IOWA CITY")) {
			t.Errorf("City=%s", p.City)
		}
	}
}

func TestSearcher_WIREFindStateOnly(t *testing.T) {
	s := searcher{}
	if err := s.helperLoadFEDWIREFile(t); err != nil {
		t.Fatal(err)
	}
	wireP := s.WIREFindStateOnly(hardResultsLimit, "IA")

	if len(wireP) == 0 {
		t.Fatalf("%s", "No matches found for state")
	}

	for _, p := range wireP {
		if !strings.Contains(p.State, strings.ToUpper("IA")) {
			t.Errorf("State=%s", p.State)
		}
	}
}

func TestSearcher_WIREFind(t *testing.T) {
	s := searcher{}
	if err := s.helperLoadFEDWIREFile(t); err != nil {
		t.Fatal(err)
	}

	req := fedSearchRequest{
		Name:          "MIDWest",
		RoutingNumber: "091905114",
		City:          "IOWA CITY",
		State:         "IA",
	}

	wireP, err := s.WIREFind(hardResultsLimit, req)

	if err != nil {
		t.Fatal(err)
	}

	if len(wireP) == 0 {
		t.Fatalf("%s", "No matches found for search")
	}

	for _, p := range wireP {
		if !strings.Contains(p.CustomerName, strings.ToUpper("MIDWest")) {
			t.Errorf("Name=%s", p.CustomerName)
		}
		if !strings.Contains(p.RoutingNumber, "091905114") {
			t.Errorf("Routing Number=%s", p.RoutingNumber)
		}

		if !strings.Contains(p.City, strings.ToUpper("IOWA City")) {
			t.Errorf("City=%s", p.City)
		}
		if !strings.Contains(p.State, "IA") {
			t.Errorf("State=%s", p.State)
		}
	}
}

func TestSearch__extractSearchLimit(t *testing.T) {
	// Too high, fallback to hard max
	req := httptest.NewRequest("GET", "/?limit=1000", nil)
	if limit := extractSearchLimit(req); limit != hardResultsLimit {
		t.Errorf("got limit of %d", limit)
	}

	// No limit, use default
	req = httptest.NewRequest("GET", "/", nil)
	if limit := extractSearchLimit(req); limit != softResultsLimit {
		t.Errorf("got limit of %d", limit)
	}

	// Between soft and hard max
	req = httptest.NewRequest("GET", "/?limit=25", nil)
	if limit := extractSearchLimit(req); limit != 25 {
		t.Errorf("got limit of %d", limit)
	}

	// Lower than soft max
	req = httptest.NewRequest("GET", "/?limit=1", nil)
	if limit := extractSearchLimit(req); limit != 1 {
		t.Errorf("got limit of %d", limit)
	}
}
