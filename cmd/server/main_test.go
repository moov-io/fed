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

func TestSearcher__setup(t *testing.T) {
	s := &searcher{logger: log.NewNopLogger()}

	logger := log.NewNopLogger()
	achFile, err := os.Open(filepath.Join("..", "..", "data", "FedACHdir.txt"))
	if err != nil {
		t.Fatal(err)
	}
	wireFile, err := os.Open(filepath.Join("..", "..", "data", "fpddir.txt"))
	if err != nil {
		t.Fatal(err)
	}

	if err := setupSearcher(logger, s, achFile, wireFile); err != nil {
		t.Fatal(err)
	}
	if err := setupSearcher(logger, s, achFile, nil); err == nil {
		t.Errorf("expected error")
	}
	if err := setupSearcher(logger, s, nil, nil); err == nil {
		t.Errorf("expected error")
	}
}
