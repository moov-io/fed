// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/moov-io/base/log"
)

func TestReader__readFEDACHData(t *testing.T) {
	s := &searcher{logger: log.NewNopLogger()}

	achFile, err := os.Open(filepath.Join("..", "..", "data", "FedACHdir.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if err := s.readFEDACHData(achFile); err != nil {
		t.Fatal(err)
	}
	if len(s.ACHDictionary.ACHParticipants) == 0 {
		t.Error("no ACH entries parsed")
	}

	// bad path
	achFile, err = os.Open("reader_test.go") // invalid fedach file
	if err != nil {
		t.Fatal(err)
	}
	defer achFile.Close()
	if err := s.readFEDACHData(achFile); err == nil {
		t.Error("expected error")
	}
}

func TestReader__readFEDWIREData(t *testing.T) {
	s := &searcher{logger: log.NewNopLogger()}

	wireFile, err := os.Open(filepath.Join("..", "..", "data", "fpddir.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if err := s.readFEDWIREData(wireFile); err != nil {
		t.Fatal(err)
	}
	if len(s.WIREDictionary.WIREParticipants) == 0 {
		t.Error("no Wire entries parsed")
	}

	// bad path
	wireFile, err = os.Open("reader_test.go")
	if err != nil {
		t.Fatal(err)
	}
	defer wireFile.Close()
	if err := s.readFEDWIREData(wireFile); err == nil {
		t.Error("expected error")
	}
}

func TestReader__readDataFilepath(t *testing.T) {
	if v := readDataFilepath("MISSING", "value"); v != "value" {
		t.Errorf("got %q", v)
	}
	if v := readDataFilepath("PATH", "value"); v == "" || v == "value" {
		t.Errorf("got %q", v)
	}
}
