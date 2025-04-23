// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/moov-io/base/log"

	"github.com/stretchr/testify/require"
)

func TestReader__fedACHDataFile(t *testing.T) {
	t.Setenv("FRB_ROUTING_NUMBER", "")
	t.Setenv("FRB_DOWNLOAD_CODE", "")

	r, err := fedACHDataFile(log.NewTestLogger())
	require.Nil(t, r)
	require.ErrorContains(t, err, "no such file or directory")
}

func TestReader__fedWireDataFile(t *testing.T) {
	t.Setenv("FRB_ROUTING_NUMBER", "")
	t.Setenv("FRB_DOWNLOAD_CODE", "")

	r, err := fedWireDataFile(log.NewTestLogger())
	require.Nil(t, r)
	require.ErrorContains(t, err, "no such file or directory")
}

func TestReader_inspectInitialDataDirectory(t *testing.T) {
	logger := log.NewNopLogger()

	dir := t.TempDir()

	err := os.WriteFile(filepath.Join(dir, "fedach.txt"), nil, 0600)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(dir, "fedwire.txt"), nil, 0600)
	require.NoError(t, err)

	// FedACH files
	fd, err := inspectInitialDataDirectory(logger, dir, fedachFilenames)
	require.NoError(t, err)

	file, ok := fd.(*os.File)
	require.True(t, ok)
	require.Equal(t, filepath.Join(dir, "fedach.txt"), file.Name())

	// FedWire files
	fd, err = inspectInitialDataDirectory(logger, dir, fedwireFilenames)
	require.NoError(t, err)

	file, ok = fd.(*os.File)
	require.True(t, ok)
	require.Equal(t, filepath.Join(dir, "fedwire.txt"), file.Name())
}

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
