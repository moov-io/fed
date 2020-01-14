// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"path/filepath"
	"testing"

	"github.com/go-kit/kit/log"
)

func TestSearcher__setup(t *testing.T) {
	s := &searcher{logger: log.NewNopLogger()}

	logger := log.NewNopLogger()
	achPath := filepath.Join("..", "..", "data", "FedACHdir.txt")
	wirePath := filepath.Join("..", "..", "data", "fpddir.txt")

	if err := setupSearcher(logger, s, achPath, wirePath); err != nil {
		t.Fatal(err)
	}
	if err := setupSearcher(logger, s, achPath, "empty2.txt"); err == nil {
		t.Errorf("expected error")
	}
	if err := setupSearcher(logger, s, "empty.txt", "empty2.txt"); err == nil {
		t.Errorf("expected error")
	}
}
