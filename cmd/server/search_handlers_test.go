// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/moov-io/base/log"
	"github.com/moov-io/fed"
	"github.com/moov-io/fed/pkg/logos"

	"github.com/gorilla/mux"
)

func TestSearch__ACHName(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/fed/ach/search?name=Farmers", nil)

	s := searcher{}
	if err := s.helperLoadFEDACHFile(t); err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()
	addSearchRoutes(log.NewNopLogger(), router, &s, logos.NewNopGrabber())
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusOK {
		t.Errorf("incorrect status code: %d", w.Code)
	}

	//fmt.Printf("%s", w.Body.String())

	var wrapper struct {
		ACHParticipants  []*fed.ACHParticipant  `json:"achParticipants"`
		WIREParticipants []*fed.WIREParticipant `json:"wireParticipants"`
	}

	if err := json.NewDecoder(w.Body).Decode(&wrapper); err != nil {
		t.Fatal(err)
	}

	if wrapper.WIREParticipants != nil {
		t.Errorf("Wire participants should be nil")
	}

	for _, p := range wrapper.ACHParticipants {
		if !strings.Contains(p.CustomerName, strings.ToUpper("FARM")) {
			t.Errorf("Name=%s", p.CustomerName)
		}
	}
}

func TestSearch__ACHRoutingNumber(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/fed/ach/search?routingNumber=044112187&limit=3", nil)

	s := searcher{}
	if err := s.helperLoadFEDACHFile(t); err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()
	addSearchRoutes(log.NewNopLogger(), router, &s, logos.NewNopGrabber())
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusOK {
		t.Errorf("incorrect status code: %d", w.Code)
	}

	var wrapper struct {
		ACHParticipants  []*fed.ACHParticipant  `json:"achParticipants"`
		WIREParticipants []*fed.WIREParticipant `json:"wireParticipants"`
	}

	if err := json.NewDecoder(w.Body).Decode(&wrapper); err != nil {
		t.Fatal(err)
	}

	if wrapper.WIREParticipants != nil {
		t.Errorf("Wire participants should be nil")
	}

	for _, p := range wrapper.ACHParticipants {
		if strings.HasPrefix(p.RoutingNumber, "044") {
			continue
		}
		t.Errorf("Routing Number=%s", p.RoutingNumber)
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
	addSearchRoutes(log.NewNopLogger(), router, &s, logos.NewNopGrabber())
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusOK {
		t.Errorf("incorrect status code: %d", w.Code)
	}

	var wrapper struct {
		ACHParticipants  []*fed.ACHParticipant  `json:"achParticipants"`
		WIREParticipants []*fed.WIREParticipant `json:"wireParticipants"`
	}

	if err := json.NewDecoder(w.Body).Decode(&wrapper); err != nil {
		t.Fatal(err)
	}

	if wrapper.WIREParticipants != nil {
		t.Errorf("Wire participants should be nil")
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
	addSearchRoutes(log.NewNopLogger(), router, &s, logos.NewNopGrabber())
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusOK {
		t.Errorf("incorrect status code: %d", w.Code)
	}

	var wrapper struct {
		ACHParticipants  []*fed.ACHParticipant  `json:"achParticipants"`
		WIREParticipants []*fed.WIREParticipant `json:"wireParticipants"`
	}

	if err := json.NewDecoder(w.Body).Decode(&wrapper); err != nil {
		t.Fatal(err)
	}

	if wrapper.WIREParticipants != nil {
		t.Errorf("Wire participants should be nil")
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
	addSearchRoutes(log.NewNopLogger(), router, &s, logos.NewNopGrabber())
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusOK {
		t.Errorf("incorrect status code: %d", w.Code)
	}

	var wrapper struct {
		ACHParticipants  []*fed.ACHParticipant  `json:"achParticipants"`
		WIREParticipants []*fed.WIREParticipant `json:"wireParticipants"`
	}

	if wrapper.WIREParticipants != nil {
		t.Errorf("Wire participants should be nil")
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

func TestSearch__ACH(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/fed/ach/search?name=Farmers&routingNumber=044112187&city=CALDWELL&state=OH&postalCode=43724", nil)

	s := searcher{}
	if err := s.helperLoadFEDACHFile(t); err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()
	addSearchRoutes(log.NewNopLogger(), router, &s, logos.NewNopGrabber())
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusOK {
		t.Errorf("incorrect status code: %d", w.Code)
	}

	var wrapper struct {
		ACHParticipants  []*fed.ACHParticipant  `json:"achParticipants"`
		WIREParticipants []*fed.WIREParticipant `json:"wireParticipants"`
	}

	if err := json.NewDecoder(w.Body).Decode(&wrapper); err != nil {
		t.Fatal(err)
	}

	if wrapper.WIREParticipants != nil {
		t.Errorf("Wire participants should be nil")
	}

	for _, p := range wrapper.ACHParticipants {
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
			t.Errorf("PostalCode=%s", p.PostalCode)
		}
	}
}

func TestSearch__ACHEmpty(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/fed/ach/search", nil)

	s := searcher{}
	if err := s.helperLoadFEDACHFile(t); err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()
	addSearchRoutes(log.NewNopLogger(), router, &s, logos.NewNopGrabber())
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
	addSearchRoutes(log.NewNopLogger(), router, &s, logos.NewNopGrabber())
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusOK {
		t.Errorf("incorrect status code: %d", w.Code)
	}

	var wrapper struct {
		ACHParticipants  []*fed.ACHParticipant  `json:"achParticipants"`
		WIREParticipants []*fed.WIREParticipant `json:"wireParticipants"`
	}

	if err := json.NewDecoder(w.Body).Decode(&wrapper); err != nil {
		t.Fatal(err)
	}

	if wrapper.ACHParticipants != nil {
		t.Errorf("ACH participants should be nil")
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
	addSearchRoutes(log.NewNopLogger(), router, &s, logos.NewNopGrabber())
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusOK {
		t.Errorf("incorrect status code: %d", w.Code)
	}

	var wrapper struct {
		ACHParticipants  []*fed.ACHParticipant  `json:"achParticipants"`
		WIREParticipants []*fed.WIREParticipant `json:"wireParticipants"`
	}

	if err := json.NewDecoder(w.Body).Decode(&wrapper); err != nil {
		t.Fatal(err)
	}

	if wrapper.ACHParticipants != nil {
		t.Errorf("ACH participants should be nil")
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
	addSearchRoutes(log.NewNopLogger(), router, &s, logos.NewNopGrabber())
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusOK {
		t.Errorf("incorrect status code: %d", w.Code)
	}

	var wrapper struct {
		ACHParticipants  []*fed.ACHParticipant  `json:"achParticipants"`
		WIREParticipants []*fed.WIREParticipant `json:"wireParticipants"`
	}

	if err := json.NewDecoder(w.Body).Decode(&wrapper); err != nil {
		t.Fatal(err)
	}

	if wrapper.ACHParticipants != nil {
		t.Errorf("ACH participants should be nil")
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
	addSearchRoutes(log.NewNopLogger(), router, &s, logos.NewNopGrabber())
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusOK {
		t.Errorf("incorrect status code: %d", w.Code)
	}

	var wrapper struct {
		ACHParticipants  []*fed.ACHParticipant  `json:"achParticipants"`
		WIREParticipants []*fed.WIREParticipant `json:"wireParticipants"`
	}

	if err := json.NewDecoder(w.Body).Decode(&wrapper); err != nil {
		t.Fatal(err)
	}

	if wrapper.ACHParticipants != nil {
		t.Errorf("ACH participants should be nil")
	}

	for _, p := range wrapper.WIREParticipants {
		if !strings.Contains(p.City, strings.ToUpper("IOWA CITY")) {
			t.Errorf("City=%s", p.City)
		}
	}
}

func TestSearch__WIRE(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/fed/wire/search?name=MIDWEST&routingNumber=091905114&state=IA&city=IOWA+CITY", nil)

	s := searcher{}
	if err := s.helperLoadFEDWIREFile(t); err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()
	addSearchRoutes(log.NewNopLogger(), router, &s, logos.NewNopGrabber())
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusOK {
		t.Errorf("incorrect status code: %d", w.Code)
	}

	var wrapper struct {
		ACHParticipants  []*fed.ACHParticipant  `json:"achParticipants"`
		WIREParticipants []*fed.WIREParticipant `json:"wireParticipants"`
	}

	if err := json.NewDecoder(w.Body).Decode(&wrapper); err != nil {
		t.Fatal(err)
	}

	if wrapper.ACHParticipants != nil {
		t.Errorf("ACH participants should be nil")
	}

	for _, p := range wrapper.WIREParticipants {
		if !strings.Contains(p.CustomerName, strings.ToUpper("Midwest")) {
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

func TestSearch__ACHRoutingNumber1Digit(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/fed/ach/search?name=Farmers&routingNumber=0&city=CALDWELL&state=OH&postalCode=43724", nil)

	s := searcher{}
	if err := s.helperLoadFEDACHFile(t); err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()
	addSearchRoutes(log.NewNopLogger(), router, &s, logos.NewNopGrabber())
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusBadRequest {
		t.Errorf("incorrect status code: %d", w.Code)
	}
}

func TestSearch__ACHRoutingNumberOnly1Digit(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/fed/ach/search?routingNumber=0", nil)

	s := searcher{}
	if err := s.helperLoadFEDACHFile(t); err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()
	addSearchRoutes(log.NewNopLogger(), router, &s, logos.NewNopGrabber())
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusBadRequest {
		t.Errorf("incorrect status code: %d", w.Code)
	}
}

func TestSearch__WIRERoutingNumber1Digit(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/fed/wire/search?name=MIDWEST&routingNumber=0&state=IA&city=IOWA+CITY", nil)

	s := searcher{}
	if err := s.helperLoadFEDWIREFile(t); err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()
	addSearchRoutes(log.NewNopLogger(), router, &s, logos.NewNopGrabber())
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusBadRequest {
		t.Errorf("incorrect status code: %d", w.Code)
	}
}

func TestSearch__WIRERoutingNumberOnly1Digit(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/fed/wire/search?routingNumber=0", nil)

	s := searcher{}
	if err := s.helperLoadFEDWIREFile(t); err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()
	addSearchRoutes(log.NewNopLogger(), router, &s, logos.NewNopGrabber())
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusBadRequest {
		t.Errorf("incorrect status code: %d", w.Code)
	}
}

func TestSearch__WIREStateSoftLimit(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/fed/wire/search?state=IA", nil)

	s := searcher{}
	if err := s.helperLoadFEDWIREFile(t); err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()
	addSearchRoutes(log.NewNopLogger(), router, &s, logos.NewNopGrabber())
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusOK {
		t.Errorf("incorrect status code: %d", w.Code)
	}

	var wrapper struct {
		ACHParticipants  []*fed.ACHParticipant  `json:"achParticipants"`
		WIREParticipants []*fed.WIREParticipant `json:"wireParticipants"`
	}

	if err := json.NewDecoder(w.Body).Decode(&wrapper); err != nil {
		t.Fatal(err)
	}

	if wrapper.ACHParticipants != nil {
		t.Errorf("ACH participants should be nil")
	}

	if len(wrapper.WIREParticipants) != 100 {
		t.Errorf("exceeded the limit: %d", len(wrapper.WIREParticipants))
	}
}

func TestSearch__WIREStateLimit(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/fed/wire/search?state=PA&limit=110", nil)

	s := searcher{}
	if err := s.helperLoadFEDWIREFile(t); err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()
	addSearchRoutes(log.NewNopLogger(), router, &s, logos.NewNopGrabber())
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusOK {
		t.Errorf("incorrect status code: %d", w.Code)
	}

	var wrapper struct {
		ACHParticipants  []*fed.ACHParticipant  `json:"achParticipants"`
		WIREParticipants []*fed.WIREParticipant `json:"wireParticipants"`
	}

	if err := json.NewDecoder(w.Body).Decode(&wrapper); err != nil {
		t.Fatal(err)
	}

	if wrapper.ACHParticipants != nil {
		t.Errorf("ACH participants should be nil")
	}

	if len(wrapper.WIREParticipants) != 110 {
		t.Errorf("exceeded the limit: %d", len(wrapper.WIREParticipants))
	}
}

func TestSearch__WIREStateHardLimit(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/fed/wire/search?state=CA&limit=550", nil)

	s := searcher{}
	if err := s.helperLoadFEDWIREFile(t); err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()
	addSearchRoutes(log.NewNopLogger(), router, &s, logos.NewNopGrabber())
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusOK {
		t.Errorf("incorrect status code: %d", w.Code)
	}

	var wrapper struct {
		ACHParticipants  []*fed.ACHParticipant  `json:"achParticipants"`
		WIREParticipants []*fed.WIREParticipant `json:"wireParticipants"`
	}

	if err := json.NewDecoder(w.Body).Decode(&wrapper); err != nil {
		t.Fatal(err)
	}

	if wrapper.ACHParticipants != nil {
		t.Errorf("ACH participants should be nil")
	}

	if len(wrapper.WIREParticipants) != 500 {
		t.Errorf("exceeded the limit: %d", len(wrapper.WIREParticipants))
	}
}

func TestSearch__WireEmpty(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/fed/wire/search", nil)

	s := searcher{}
	if err := s.helperLoadFEDWIREFile(t); err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()
	addSearchRoutes(log.NewNopLogger(), router, &s, logos.NewNopGrabber())
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusBadRequest {
		t.Errorf("incorrect status code: %d", w.Code)
	}
}
