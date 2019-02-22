// Copyright 2019 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/moov-io/fed"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func (s *searcher) helperLogger(t *testing.T) {
}

func TestSearch__ACHName(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/fed/ach/search?name=Farmers", nil)

	s := searcher{}
	if err := s.helperLoadFEDACHFile(t); err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()
	addSearchRoutes(nil, router, &s)
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusOK {
		t.Errorf("incorrect status code: %d", w.Code)
	}

	//fmt.Printf("%s", w.Body.String())

	var wrapper struct {
		ACHParticipants []*fed.ACHParticipant `json:"achParticipants"`
	}

	if err := json.NewDecoder(w.Body).Decode(&wrapper); err != nil {
		t.Fatal(err)
	}

	for _, p := range wrapper.ACHParticipants {
		if !strings.Contains(p.CustomerName, strings.ToUpper("Farmer")) {
			t.Errorf("Name=%s", p.CustomerName)
		}
	}
}

func TestSearch__ACHRoutingNumber(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/fed/ach/search?routingNumber=044112187", nil)

	s := searcher{}
	if err := s.helperLoadFEDACHFile(t); err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()
	addSearchRoutes(nil, router, &s)
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusOK {
		t.Errorf("incorrect status code: %d", w.Code)
	}

	var wrapper struct {
		ACHParticipants []*fed.ACHParticipant `json:"achParticipants"`
	}

	if err := json.NewDecoder(w.Body).Decode(&wrapper); err != nil {
		t.Fatal(err)
	}

	for _, p := range wrapper.ACHParticipants {
		if !strings.Contains(p.RoutingNumber, "044112187") {
			t.Errorf("RoutingNumber=%s", p.RoutingNumber)
		}
	}
}

func TestSearch__ACHCity(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/fed/ach/search?city=CALDWELL", nil)

	s := searcher{}
	if err := s.helperLoadFEDACHFile(t); err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()
	addSearchRoutes(nil, router, &s)
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusOK {
		t.Errorf("incorrect status code: %d", w.Code)
	}

	var wrapper struct {
		ACHParticipants []*fed.ACHParticipant `json:"achParticipants"`
	}

	if err := json.NewDecoder(w.Body).Decode(&wrapper); err != nil {
		t.Fatal(err)
	}

	for _, p := range wrapper.ACHParticipants {
		if !strings.Contains(p.City, strings.ToUpper("CALDWELL")) {
			t.Errorf("City=%s", p.City)
		}
	}
}

func TestSearch__ACHState(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/fed/ach/search?state=OH", nil)

	s := searcher{}
	if err := s.helperLoadFEDACHFile(t); err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()
	addSearchRoutes(nil, router, &s)
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusOK {
		t.Errorf("incorrect status code: %d", w.Code)
	}

	var wrapper struct {
		ACHParticipants []*fed.ACHParticipant `json:"achParticipants"`
	}

	if err := json.NewDecoder(w.Body).Decode(&wrapper); err != nil {
		t.Fatal(err)
	}

	for _, p := range wrapper.ACHParticipants {
		if !strings.Contains(p.State, strings.ToUpper("OH")) {
			t.Errorf("State=%s", p.State)
		}
	}
}

func TestSearch__ACHPostalCode(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/fed/ach/search?postalCode=43724", nil)

	s := searcher{}
	if err := s.helperLoadFEDACHFile(t); err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()
	addSearchRoutes(nil, router, &s)
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusOK {
		t.Errorf("incorrect status code: %d", w.Code)
	}

	var wrapper struct {
		ACHParticipants []*fed.ACHParticipant `json:"achParticipants"`
	}

	if err := json.NewDecoder(w.Body).Decode(&wrapper); err != nil {
		t.Fatal(err)
	}

	for _, p := range wrapper.ACHParticipants {
		if !strings.Contains(p.PostalCode, "43724") {
			t.Errorf("PostalCode=%s", p.PostalCode)
		}
	}
}

func TestSearch__Empty(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/fed/ach/search", nil)

	s := searcher{}
	if err := s.helperLoadFEDACHFile(t); err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()
	addSearchRoutes(nil, router, &s)
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusBadRequest {
		t.Errorf("incorrect status code: %d", w.Code)
	}
}

// WIRES

func TestSearch__WIREName(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/fed/wire/search?name=MIDWEST", nil)

	s := searcher{}
	if err := s.helperLoadFEDWIREFile(t); err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()
	addSearchRoutes(nil, router, &s)
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusOK {
		t.Errorf("incorrect status code: %d", w.Code)
	}

	var wrapper struct {
		WIREParticipants []*fed.WIREParticipant `json:"wireParticipants"`
	}

	if err := json.NewDecoder(w.Body).Decode(&wrapper); err != nil {
		t.Fatal(err)
	}

	for _, p := range wrapper.WIREParticipants {
		if !strings.Contains(p.CustomerName, strings.ToUpper("Midwest")) {
			t.Errorf("Name=%s", p.CustomerName)
		}
	}
}

func TestSearch__WIRERoutingNumber(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/fed/wire/search?routingNumber=091905114", nil)

	s := searcher{}
	if err := s.helperLoadFEDWIREFile(t); err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()
	addSearchRoutes(nil, router, &s)
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusOK {
		t.Errorf("incorrect status code: %d", w.Code)
	}

	var wrapper struct {
		WIREParticipants []*fed.WIREParticipant `json:"wireParticipants"`
	}

	if err := json.NewDecoder(w.Body).Decode(&wrapper); err != nil {
		t.Fatal(err)
	}

	for _, p := range wrapper.WIREParticipants {
		if !strings.Contains(p.RoutingNumber, "091905114") {
			t.Errorf("Routing Number=%s", p.RoutingNumber)
		}
	}
}

func TestSearch__WIREState(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/fed/wire/search?state=IA", nil)

	s := searcher{}
	if err := s.helperLoadFEDWIREFile(t); err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()
	addSearchRoutes(nil, router, &s)
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusOK {
		t.Errorf("incorrect status code: %d", w.Code)
	}

	var wrapper struct {
		WIREParticipants []*fed.WIREParticipant `json:"wireParticipants"`
	}

	if err := json.NewDecoder(w.Body).Decode(&wrapper); err != nil {
		t.Fatal(err)
	}

	for _, p := range wrapper.WIREParticipants {
		if !strings.Contains(p.State, strings.ToUpper("IA")) {
			t.Errorf("State=%s", p.State)
		}
	}
}

func TestSearch__WIRECity(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/fed/wire/search?city=IOWA+CITY", nil)

	s := searcher{}
	if err := s.helperLoadFEDWIREFile(t); err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()
	addSearchRoutes(nil, router, &s)
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusOK {
		t.Errorf("incorrect status code: %d", w.Code)
	}

	var wrapper struct {
		WIREParticipants []*fed.WIREParticipant `json:"wireParticipants"`
	}

	if err := json.NewDecoder(w.Body).Decode(&wrapper); err != nil {
		t.Fatal(err)
	}

	for _, p := range wrapper.WIREParticipants {
		if !strings.Contains(p.City, strings.ToUpper("IOWA CITY")) {
			t.Errorf("City=%s", p.City)
		}
	}
}
