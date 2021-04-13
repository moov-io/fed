// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"path/filepath"
	"testing"

	"github.com/moov-io/base/log"
)

func TestReader__readFEDACHData(t *testing.T) {
	s := &searcher{logger: log.NewNopLogger()}
	if err := s.readFEDACHData(filepath.Join("..", "..", "data", "FedACHdir.txt")); err != nil {
		t.Fatal(err)
	}
	if len(s.ACHDictionary.ACHParticipants) == 0 {
		t.Error("no ACH entries parsed")
	}

	// bad path
	if err := s.readFEDACHData("empty.txt"); err == nil {
		t.Error("expected error")
	}
}

func TestReader__readFEDWIREData(t *testing.T) {
	s := &searcher{logger: log.NewNopLogger()}
	if err := s.readFEDWIREData(filepath.Join("..", "..", "data", "fpddir.txt")); err != nil {
		t.Fatal(err)
	}
	if len(s.WIREDictionary.WIREParticipants) == 0 {
		t.Error("no Wire entries parsed")
	}

	// bad path
	if err := s.readFEDWIREData("empty.txt"); err == nil {
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
